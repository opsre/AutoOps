package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	tmplhtml "html/template"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"gopkg.in/gomail.v2"

	alertDao "dodevops-api/api/monitor/dao"
	alertModel "dodevops-api/api/monitor/model"
)

type AlertServiceInterface interface {
	CreateTemplate(tpl *alertModel.PrometheusAlertDB) error
	DeleteTemplate(id int) error
	UpdateTemplate(tpl *alertModel.PrometheusAlertDB) error
	GetTemplateList() ([]*alertModel.PrometheusAlertDB, error)
	GetTemplateById(id int) (*alertModel.PrometheusAlertDB, error)

	PrometheusAlertHandle(pMsg alertModel.PrometheusAlertMsg, pJson interface{}, pAlertManagerJson map[string]interface{}, GlobalAlertRouter []*alertModel.AlertRouter) (string, error)
	GetAllAlertRouter() ([]*alertModel.AlertRouter, error)
}

type alertService struct {
	alertDao alertDao.AlertDao
}

func NewAlertService() AlertServiceInterface {
	return &alertService{
		alertDao: alertDao.NewAlertDao(),
	}
}

func (s *alertService) CreateTemplate(tpl *alertModel.PrometheusAlertDB) error {
	return s.alertDao.CreateTemplate(tpl)
}
func (s *alertService) DeleteTemplate(id int) error { return s.alertDao.DeleteTemplate(id) }
func (s *alertService) UpdateTemplate(tpl *alertModel.PrometheusAlertDB) error {
	return s.alertDao.UpdateTemplate(tpl)
}
func (s *alertService) GetTemplateList() ([]*alertModel.PrometheusAlertDB, error) {
	return s.alertDao.GetTemplateList()
}
func (s *alertService) GetTemplateById(id int) (*alertModel.PrometheusAlertDB, error) {
	return s.alertDao.GetTemplateById(id)
}
func (s *alertService) GetAllAlertRouter() ([]*alertModel.AlertRouter, error) {
	return s.alertDao.GetAllAlertRouter()
}

func (s *alertService) PrometheusAlertHandle(pMsg alertModel.PrometheusAlertMsg, pJson interface{}, pAlertManagerJson map[string]interface{}, GlobalAlertRouter []*alertModel.AlertRouter) (string, error) {
	logsign := "[AlertHandle]"

	templates, err := s.alertDao.GetTemplateList()
	if err != nil {
		return "", fmt.Errorf("读取模板列表失败: %v", err)
	}

	var PrometheusAlertTpl *alertModel.PrometheusAlertDB
	for _, Tpl := range templates {
		if Tpl.Tplname == pMsg.Tpl {
			PrometheusAlertTpl = Tpl
			break
		}
	}

	var message string
	if pMsg.Type != "" && PrometheusAlertTpl != nil {
		fmt.Printf("[AlertHandle] 匹配到模板: %s, 模式 Split: %s, Tpluse: %s\n", PrometheusAlertTpl.Tplname, pMsg.Split, PrometheusAlertTpl.Tpluse)
		if pMsg.Split != "false" && PrometheusAlertTpl.Tpluse == "Prometheus" {
			Alerts_Value, ok := pAlertManagerJson["alerts"].([]interface{})
			if ok {
				fmt.Printf("[AlertHandle] Alerts 数组长度: %d\n", len(Alerts_Value))
				for _, AlertValue := range Alerts_Value {
					pAlertManagerJson["alerts"] = Alerts_Value[0:0]
					pAlertManagerJson["alerts"] = append(pAlertManagerJson["alerts"].([]interface{}), AlertValue)
					go s.SetRecord(AlertValue)

					xalert := AlertValue.(map[string]interface{})
					Return_pMsgs := s.AlertRouterSet(xalert, pMsg, PrometheusAlertTpl.Tpl, GlobalAlertRouter)
					fmt.Printf("[AlertHandle] AlertRouterSet 路由匹配数量: %d\n", len(Return_pMsgs))
					for _, Return_pMsg := range Return_pMsgs {
						err, msg := s.TransformAlertMessage(pAlertManagerJson, Return_pMsg.Tpl)
						if err != nil {
							fmt.Printf("[AlertHandle] 模板渲染失败: %v\n", err)
							message += err.Error()
						} else {
							fmt.Printf("[AlertHandle] 模板渲染成功，准备发送 %s 消息\n", Return_pMsg.Type)
							message += s.SendMessagePrometheusAlert(msg, &Return_pMsg, logsign)
						}
					}
				}
			} else {
				fmt.Printf("[AlertHandle] pAlertManagerJson 中没有找到 alerts 数组\n")
			}
		} else {
			err, msg := s.TransformAlertMessage(pJson, PrometheusAlertTpl.Tpl)
			if err != nil {
				fmt.Printf("[AlertHandle] 模板渲染失败(非Split模式): %v\n", err)
				message = err.Error()
			} else {
				fmt.Printf("[AlertHandle] 模板渲染成功(非Split模式)，准备发送 %s 消息\n", pMsg.Type)
				message = s.SendMessagePrometheusAlert(msg, &pMsg, logsign)
			}
		}
	} else {
		fmt.Printf("[AlertHandle] 未匹配到模板或未指定 Type. 收到 Type: %s, Tpl: %v\n", pMsg.Type, (PrometheusAlertTpl != nil))
		message = "自定义模板接口参数异常(可能是不存在或找不到该模板名字)！"
	}
	return message, nil
}

func (s *alertService) AlertRouterSet(xalert map[string]interface{}, PMsg alertModel.PrometheusAlertMsg, Tpl string, GlobalAlertRouter []*alertModel.AlertRouter) []alertModel.PrometheusAlertMsg {
	return_Msgs := []alertModel.PrometheusAlertMsg{}
	PMsg.Tpl = Tpl
	return_Msgs = append(return_Msgs, PMsg)

	for _, router_value := range GlobalAlertRouter {
		LabelMap := []alertModel.LabelMap{}
		json.Unmarshal([]byte(router_value.Rules), &LabelMap)
		rules_num := len(LabelMap)
		rules_num_match := 0

		if xalert["status"] == "resolved" && !router_value.SendResolved {
			labelsObj, exists := xalert["labels"].(map[string]interface{})
			if exists {
				if alertName, ok := labelsObj["alertname"].(string); ok {
					fmt.Printf("告警名称：%s 路由规则：%s 路由恢复告警：%v\n", alertName, router_value.Name, router_value.SendResolved)
				}
			}
			continue
		}

		for _, rule := range LabelMap {
			if labelsObj, exists := xalert["labels"].(map[string]interface{}); exists {
				for label_key, label_value := range labelsObj {
					if rule.Regex {
						if rule.Name == label_key {
							if labelStr, ok := label_value.(string); ok {
								tz := regexp.MustCompile(rule.Value)
								if len(tz.FindAllString(labelStr, -1)) > 0 {
									rules_num_match += 1
								}
							}
						}
					} else {
						if rule.Name == label_key {
							if labelStr, ok := label_value.(string); ok && rule.Value == labelStr {
								rules_num_match += 1
							}
						}
					}
				}
			}
		}

		if rules_num == rules_num_match && rules_num > 0 {
			if router_value.Tpl != nil {
				PMsg.Type = router_value.Tpl.Tpltype
				PMsg.Tpl = router_value.Tpl.Tpl
			}
			atSomeOne := router_value.AtSomeOne
			if router_value.AtSomeOneRR {
				openIds := strings.Split(router_value.AtSomeOne, ",")
				if len(openIds) > 1 {
					duration := time.Since(time.Unix(0, 0))
					days := duration.Hours() / 24
					i := int(days) % len(openIds)
					atSomeOne = openIds[i]
				}
			}

			tplType := ""
			if router_value.Tpl != nil {
				tplType = router_value.Tpl.Tpltype
			}

			switch tplType {
			case "wx":
				PMsg.Wxurl = router_value.UrlOrPhone
				PMsg.AtSomeOne = atSomeOne
			case "dd":
				PMsg.Ddurl = router_value.UrlOrPhone
				PMsg.AtSomeOne = atSomeOne
			case "fs":
				PMsg.Fsurl = router_value.UrlOrPhone
				PMsg.AtSomeOne = atSomeOne
			case "webhook":
				PMsg.WebHookUrl = router_value.UrlOrPhone
			case "email":
				PMsg.Email = router_value.UrlOrPhone
			case "rl":
				PMsg.GroupId = router_value.UrlOrPhone
			case "txdx", "hwdx", "bddx", "alydx", "txdh", "alydh", "rlydh", "7moordx", "7moordh":
				PMsg.Phone = router_value.UrlOrPhone
			}
			return_Msgs = append(return_Msgs, PMsg)
		}
	}
	return return_Msgs
}

func (s *alertService) SetRecord(AlertValue interface{}) {
	var Alertname, Status, Level, Labels, Instance, Summary, Description, StartAt, EndAt string
	xalert := AlertValue.(map[string]interface{})

	if startsAt, ok := xalert["startsAt"].(string); ok {
		StartAt = startsAt
		// PCstTime handling omitted/mocked
	}
	if endsAt, ok := xalert["endsAt"].(string); ok {
		EndAt = endsAt
	}
	if status, ok := xalert["status"].(string); ok {
		Status = status
	}

	if labelsObj, ok := xalert["labels"].(map[string]interface{}); ok {
		if alertname, exists := labelsObj["alertname"]; exists {
			Alertname = alertname.(string)
		}
		if level, exists := labelsObj["level"]; exists {
			Level = level.(string)
		}
		if instance, exists := labelsObj["instance"]; exists {
			Instance = instance.(string)
		}
		if labelsJsonStr, err := json.Marshal(labelsObj); err == nil {
			Labels = string(labelsJsonStr)
		}
	}

	if annotationsObj, ok := xalert["annotations"].(map[string]interface{}); ok {
		if description, exists := annotationsObj["description"]; exists {
			Description = description.(string)
		}
		if summary, exists := annotationsObj["summary"]; exists {
			Summary = summary.(string)
		}
	}

	if !s.alertDao.GetRecordExist(Alertname, Level, Labels, Instance, StartAt, EndAt, Summary, Description, Status) {
		record := &alertModel.AlertRecord{
			Alertname:   Alertname,
			AlertLevel:  Level,
			Labels:      Labels,
			Instance:    Instance,
			StartsAt:    StartAt,
			EndsAt:      EndAt,
			Summary:     Summary,
			Description: Description,
			AlertStatus: Status,
			CreatedTime: time.Now(),
		}
		s.alertDao.AddAlertRecord(record)
	}

	// ES writing ignored for now based on legacy logic
}

func GetTimeDuration(start, end string) string { return "duration" }
func GetCSTtime(ts string) string              { return time.Now().Format(time.RFC3339) }
func TimeFormat(ts string) string              { return time.Now().Format("2006-01-02 15:04:05") }
func GetTime() string                          { return time.Now().Format("2006-01-02 15:04:05") }

func (s *alertService) TransformAlertMessage(p_json interface{}, tpltext string) (error error, msg string) {
	funcMap := template.FuncMap{
		"GetTimeDuration": GetTimeDuration,
		"GetCSTtime":      GetCSTtime,
		"TimeFormat":      TimeFormat,
		"GetTime":         GetTime,
		"toUpper":         strings.ToUpper,
		"toLower":         strings.ToLower,
		"title":           strings.Title,
		"join": func(sep string, arr []string) string {
			return strings.Join(arr, sep)
		},
		"match": regexp.MatchString,
		"safeHtml": func(text string) tmplhtml.HTML {
			return tmplhtml.HTML(text)
		},
		"reReplaceAll": func(pattern, repl, text string) string {
			re := regexp.MustCompile(pattern)
			return re.ReplaceAllString(text, repl)
		},
		"stringSlice": func(arr ...string) []string {
			return arr
		},
		"SplitString": func(pstring string, start int, stop int) string {
			if stop < 0 && len(pstring)+stop >= start {
				return pstring[start : len(pstring)+stop]
			}
			if stop > len(pstring) {
				stop = len(pstring)
			}
			if start >= 0 && start <= stop {
				return pstring[start:stop]
			}
			return pstring
		},
	}

	buf := new(bytes.Buffer)
	tpl, err := template.New("").Funcs(funcMap).Parse(tpltext)
	if err != nil {
		return err, ""
	}

	err = tpl.Execute(buf, p_json)
	if err != nil {
		return err, ""
	}

	return nil, buf.String()
}

// 模拟遗留的方法调用
func DoBalance(urls []string) string {
	if len(urls) > 0 {
		return urls[0]
	}
	return ""
}

type WXMessage struct {
	Msgtype  string `json:"msgtype"`
	Markdown Mark   `json:"markdown"`
}

type Mark struct {
	Content string `json:"content"`
}

func PostToWeiXin(message, url, atSomeone, logsign string) string {
	// (移除了beego配置和counter的依赖，仅保留核心发送逻辑)
	SendContent := message
	if atSomeone != "" {
		userid := strings.Split(atSomeone, ",")
		idtext := ""
		for _, id := range userid {
			idtext += "<@" + id + ">"
		}
		SendContent += "\n" + idtext
	}

	u := WXMessage{
		Msgtype:  "markdown",
		Markdown: Mark{Content: SendContent},
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)

	res, err := http.Post(url, "application/json", b)
	if err != nil {
		fmt.Printf("%s [weixin] http post error: %v\n", logsign, err)
		return err.Error()
	}
	defer res.Body.Close()

	result, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("%s [weixin] read body error: %v\n", logsign, err)
		return err.Error()
	}

	fmt.Printf("%s [weixin] response: %s\n", logsign, string(result))
	return string(result)
}

type DDMessage struct {
	Msgtype  string `json:"msgtype"`
	Markdown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown"`
	At struct {
		AtMobiles []string `json:"atMobiles"`
		IsAtAll   bool     `json:"isAtAll"`
	} `json:"at"`
}

func dingdingSign(ddurl string, logsign string) string {
	timestamp := time.Now()
	timestampMs := timestamp.UnixNano() / int64(time.Millisecond)
	tsMsStr := strconv.FormatInt(timestampMs, 10)

	u, err := url.Parse(ddurl)
	if err != nil {
		fmt.Printf("%s [dingdingSign] url 解析失败: %v\n", logsign, err)
		return ddurl
	}

	queryParams := u.Query()
	secret := queryParams.Get("secret")
	if len(secret) == 0 {
		return ddurl
	}

	signStr := tsMsStr + "\n" + secret
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(signStr))
	signature := h.Sum(nil)

	sign := base64.StdEncoding.EncodeToString(signature)

	delete(queryParams, "secret")
	queryParams.Add("timestamp", tsMsStr)
	queryParams.Add("sign", sign)
	u.RawQuery = queryParams.Encode()
	return u.String()
}

func PostToDingDing(title, message, urlString, atSomeone, logsign string) string {
	// 判断如果 url 里带了 secret，则进行加签
	if strings.Contains(urlString, "secret=") {
		urlString = dingdingSign(urlString, logsign)
	}

	Atall := true
	atMobile := []string{}
	SendText := message
	if atSomeone != "" {
		atMobile = strings.Split(atSomeone, ",")
		AtText := ""
		for _, phoneN := range atMobile {
			AtText += " @" + phoneN
		}
		SendText += AtText
		Atall = false
	}

	u := DDMessage{
		Msgtype: "markdown",
	}
	u.Markdown.Title = title
	u.Markdown.Text = SendText
	u.At.AtMobiles = atMobile
	u.At.IsAtAll = Atall

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Post(urlString, "application/json", b)
	if err != nil {
		fmt.Printf("%s [dingding] http post error: %v\n", logsign, err)
		return err.Error()
	}
	defer res.Body.Close()

	result, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("%s [dingding] read body error: %v\n", logsign, err)
		return err.Error()
	}

	fmt.Printf("%s [dingding] response: %s\n", logsign, string(result))
	return string(result)
}

type FSMessage struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type Conf struct {
	WideScreenMode bool `json:"wide_screen_mode"`
	EnableForward  bool `json:"enable_forward"`
}

type Te struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type Element struct {
	Tag      string    `json:"tag"`
	Text     Te        `json:"text"`
	Content  string    `json:"content"`
	Elements []Element `json:"elements"`
}

type Titles struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type Headers struct {
	Title    Titles `json:"title"`
	Template string `json:"template"`
}

type Cards struct {
	Config   Conf      `json:"config"`
	Elements []Element `json:"elements"`
	Header   Headers   `json:"header"`
}

type FSMessagev2 struct {
	MsgType string `json:"msg_type"`
	Email   string `json:"email"`
	Card    Cards  `json:"card"`
}

func PostToFeiShu(title, text, Fsurl, logsign string) string {
	u := FSMessage{Title: title, Text: text}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Post(Fsurl, "application/json", b)
	if err != nil {
		fmt.Printf("%s [feishu] http post error: %v\n", logsign, err)
		return err.Error()
	}
	defer res.Body.Close()

	result, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("%s [feishu] read body error: %v\n", logsign, err)
		return err.Error()
	}

	fmt.Printf("%s [feishu] response: %s\n", logsign, string(result))
	return string(result)
}

func PostToFeiShuv2(title, text, Fsurl, userOpenId, logsign string) string {
	var color string
	if strings.Count(text, "resolved") > 0 && strings.Count(text, "firing") > 0 {
		color = "orange"
	} else if strings.Count(text, "resolved") > 0 {
		color = "green"
	} else {
		color = "red"
	}

	SendContent := text
	if userOpenId != "" {
		OpenIds := strings.Split(userOpenId, ",")
		OpenIdtext := ""
		for _, OpenId := range OpenIds {
			OpenIdtext += "<at user_id=\"" + OpenId + "\" id=\"" + OpenId + "\" email=\"" + OpenId + "\"></at>"
		}
		SendContent += "\n" + OpenIdtext
	}

	u := FSMessagev2{
		MsgType: "interactive",
		Email:   "",
		Card: Cards{
			Config: Conf{
				WideScreenMode: true,
				EnableForward:  true,
			},
			Header: Headers{
				Title: Titles{
					Content: title,
					Tag:     "plain_text",
				},
				Template: color,
			},
			Elements: []Element{
				{
					Tag: "div",
					Text: Te{
						Content: SendContent,
						Tag:     "lark_md",
					},
				},
				{
					Tag: "hr",
				},
				{
					Tag: "note",
					Elements: []Element{
						{
							Content: title,
							Tag:     "lark_md",
						},
					},
				},
			},
		},
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Post(Fsurl, "application/json", b)
	if err != nil {
		fmt.Printf("%s [feishuv2] %s http post error: %v\n", logsign, title, err)
		return err.Error()
	}
	defer res.Body.Close()

	result, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("%s [feishuv2] %s read body error: %v\n", logsign, title, err)
		return err.Error()
	}

	fmt.Printf("%s [feishuv2] %s response: %s\n", logsign, title, string(result))
	return string(result)
}

func PostToFS(title, message, url, atSomeone, logsign string) string {
	if strings.Contains(url, "/v2/") {
		return PostToFeiShuv2(title, message, url, atSomeone, logsign)
	}
	return PostToFeiShu(title, message, url, logsign)
}
func PostToWebhook(text, WebhookUrl, logsign string, contentType string) string {
	fmt.Printf("%s [Webhook] Sending payload to %s: %s\n", logsign, WebhookUrl, text)

	JsonMsg := bytes.NewReader([]byte(text))
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	if contentType == "" {
		contentType = "application/json"
	}

	res, err := client.Post(WebhookUrl, contentType, JsonMsg)
	if err != nil {
		fmt.Printf("%s [Webhook] http post error: %v\n", logsign, err)
		return err.Error()
	}
	defer res.Body.Close()

	result, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("%s [Webhook] read body error: %v\n", logsign, err)
		return err.Error()
	}

	fmt.Printf("%s [Webhook] response: %s\n", logsign, string(result))
	return string(result)
}

func PostTXmessage(message, phone, logsign string) string      { return "txdx" }
func PostHWmessage(message, phone, logsign string) string      { return "hwdx" }
func PostBDYmessage(message, phone, logsign string) string     { return "bddx" }
func PostALYmessage(message, phone, logsign string) string     { return "alydx" }
func PostTXphonecall(message, phone, logsign string) string    { return "txdh" }
func PostALYphonecall(message, phone, logsign string) string   { return "alydh" }
func PostRLYphonecall(message, phone, logsign string) string   { return "rlydh" }
func Post7MOORmessage(message, phone, logsign string) string   { return "7moordx" }
func Post7MOORphonecall(message, phone, logsign string) string { return "7moordh" }

func (s *alertService) SendEmail(EmailBody, Emails, EmailTitle, logsign string) string {
        openEmail := s.alertDao.GetAlertConfig("open-email")
        if openEmail != "1" {
                fmt.Printf("%s [email] email未配置未开启状态,请先开启\n", logsign)
                return "email未配置未开启状态,请先开启"
        }

        serverHost := s.alertDao.GetAlertConfig("Email_host")
        serverPort, _ := strconv.Atoi(s.alertDao.GetAlertConfig("Email_port"))
        fromEmail := s.alertDao.GetAlertConfig("Email_user")
        Passwd := s.alertDao.GetAlertConfig("Email_password")
        if EmailTitle == "" {
                EmailTitle = s.alertDao.GetAlertConfig("Email_title")
        }

        SendToEmails := []string{}
        m := gomail.NewMessage()
        if len(Emails) == 0 { return "收件人不能为空" }

	for _, Email := range strings.Split(Emails, ",") {
		SendToEmails = append(SendToEmails, strings.TrimSpace(Email))
	}

	// 收件人,...代表打散列表填充不定参数
	m.SetHeader("To", SendToEmails...)
	// 发件人
	m.SetAddressHeader("From", fromEmail, EmailTitle)
	// 主题
	m.SetHeader("Subject", EmailTitle)
	// 正文
	m.SetBody("text/html", EmailBody)

	d := gomail.NewDialer(serverHost, serverPort, fromEmail, Passwd)
	// 忽略证书
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// 发送
	err := d.DialAndSend(m)
	if err != nil {
		fmt.Printf("%s [email] http post/send error: %v\n", logsign, err)
		return err.Error()
	}

	fmt.Printf("%s [email] email send ok to %s\n", logsign, Emails)
	return "email send ok to " + Emails
}

func SendTG(message, logsign string) string                                 { return "tg" }
func SendWorkWechat(touser, toparty, totag, message, logsign string) string { return "workwechat" }
func PostToRuLiu(groupId, message, url, logsign string) string              { return "ruliu" }
func SendBark(message, logsign string) string                               { return "bark" }
func SendVoice(message, logsign string) string                              { return "voice" }
func PostToFeiShuApp(title, message, atsomeone, logsign string) string      { return "fsapp" }
func SendKafka(message, logsign string) string                              { return "kafka" }

func (s *alertService) SendMessagePrometheusAlert(message string, pmsg *alertModel.PrometheusAlertMsg, logsign string) string {
	Title := "AutoOps Alert"
	var ReturnMsg string
        if pmsg.Wxurl == "" { pmsg.Wxurl = s.alertDao.GetAlertConfig("wxurl") }
        if pmsg.Ddurl == "" { pmsg.Ddurl = s.alertDao.GetAlertConfig("ddurl") }
        if pmsg.Fsurl == "" { pmsg.Fsurl = s.alertDao.GetAlertConfig("fsurl") }
        if pmsg.Email == "" { pmsg.Email = s.alertDao.GetAlertConfig("Default_emails") }

	switch pmsg.Type {
	case "wx":
		Wxurl := strings.Split(pmsg.Wxurl, ",")
		if pmsg.RoundRobin == "true" {
			ReturnMsg += PostToWeiXin(message, DoBalance(Wxurl), pmsg.AtSomeOne, logsign)
		} else {
			for _, url := range Wxurl {
				ReturnMsg += PostToWeiXin(message, url, pmsg.AtSomeOne, logsign)
			}
		}
	case "dd":
		Ddurl := strings.Split(pmsg.Ddurl, ",")
		if pmsg.RoundRobin == "true" {
			ReturnMsg += PostToDingDing(Title, message, DoBalance(Ddurl), pmsg.AtSomeOne, logsign)
		} else {
			for _, url := range Ddurl {
				ReturnMsg += PostToDingDing(Title, message, url, pmsg.AtSomeOne, logsign)
			}
		}
	case "fs":
		Fsurl := strings.Split(pmsg.Fsurl, ",")
		if pmsg.RoundRobin == "true" {
			ReturnMsg += PostToFS(Title, message, DoBalance(Fsurl), pmsg.AtSomeOne, logsign)
		} else {
			for _, url := range Fsurl {
				ReturnMsg += PostToFS(Title, message, url, pmsg.AtSomeOne, logsign)
			}
		}
	case "webhook":
		Fwebhookurl := strings.Split(pmsg.WebHookUrl, ",")
		if pmsg.RoundRobin == "true" {
			ReturnMsg += PostToWebhook(message, DoBalance(Fwebhookurl), logsign, pmsg.WebhookContentType)
		} else {
			for _, url := range Fwebhookurl {
				ReturnMsg += PostToWebhook(message, url, logsign, pmsg.WebhookContentType)
			}
		}
	case "txdx":
		ReturnMsg += PostTXmessage(message, pmsg.Phone, logsign)
	case "hwdx":
		ReturnMsg += PostHWmessage(message, pmsg.Phone, logsign)
	case "bddx":
		ReturnMsg += PostBDYmessage(message, pmsg.Phone, logsign)
	case "alydx":
		ReturnMsg += PostALYmessage(message, pmsg.Phone, logsign)
	case "txdh":
		ReturnMsg += PostTXphonecall(message, pmsg.Phone, logsign)
	case "alydh":
		ReturnMsg += PostALYphonecall(message, pmsg.Phone, logsign)
	case "rlydh":
		ReturnMsg += PostRLYphonecall(message, pmsg.Phone, logsign)
	case "7moordx":
		ReturnMsg += Post7MOORmessage(message, pmsg.Phone, logsign)
	case "7moordh":
		ReturnMsg += Post7MOORphonecall(message, pmsg.Phone, logsign)
	case "email":
		ReturnMsg += s.SendEmail(message, pmsg.Email, pmsg.EmailTitle, logsign)
	case "tg":
		ReturnMsg += SendTG(message, logsign)
	case "workwechat":
		ReturnMsg += SendWorkWechat(pmsg.ToUser, pmsg.ToParty, pmsg.ToTag, message, logsign)
	case "rl":
		ReturnMsg += PostToRuLiu(pmsg.GroupId, message, "mock_rl_url", logsign)
	case "bark":
		ReturnMsg += SendBark(message, logsign)
	case "voice":
		ReturnMsg += SendVoice(message, logsign)
	case "fsapp":
		ReturnMsg += PostToFeiShuApp(Title, message, pmsg.AtSomeOne, logsign)
	case "kafka":
		ReturnMsg += SendKafka(message, logsign)
	default:
		ReturnMsg = "参数错误"
	}
	return ReturnMsg
}
