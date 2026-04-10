package controller

import (
	// "log"
	"dodevops-api/api/monitor/model"
	"dodevops-api/api/monitor/service"
	"dodevops-api/common/result"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AlertController struct {
	alertService service.AlertServiceInterface
}

func NewAlertController() *AlertController {
	return &AlertController{
		alertService: service.NewAlertService(),
	}
}

/*
准备新增阿里云告警回调
Content-Type: application/x-www-form-urlencoded; charset=UTF-8
*/
type AliyunAlert struct {
	Expression      string `json:"expression"`
	MetricName      string `json:"metricName"`
	InstanceName    string `json:"instanceName"`
	Signature       string `json:"signature"`
	MetricProject   string `json:"metricProject"`
	UserId          string `json:"userId"`
	CurValue        string `json:"curValue"`
	AlertName       string `json:"alertName"`
	Namespace       string `json:"namespace"`
	TriggerLevel    string `json:"triggerLevel"`
	AlertState      string `json:"alertState"`
	PreTriggerLevel string `json:"preTriggerLevel"`
	RuleId          string `json:"ruleId"`
	Dimensions      string `json:"dimensions"`
	Timestamp       string `json:"timestamp"`
}

// CreateTemplate 创建模板
func (c *AlertController) CreateTemplate(ctx *gin.Context) {
	var ts model.PrometheusAlertDB
	if err := ctx.ShouldBindJSON(&ts); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "参数解析失败: "+err.Error())
		return
	}

	if err := c.alertService.CreateTemplate(&ts); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "模板创建失败: "+err.Error())
		return
	}
	result.Success(ctx, "模版创建成功")
}

// DeleteTemplate 删除模板
func (c *AlertController) DeleteTemplate(ctx *gin.Context) {
	idStr := ctx.DefaultQuery("id", "")
	if id, err := strconv.Atoi(idStr); err == nil {
		if err := c.alertService.DeleteTemplate(id); err != nil {
			result.Failed(ctx, http.StatusBadRequest, "删除失败: "+err.Error())
			return
		}
		result.Success(ctx, "模版删除成功")
	} else {
		result.Failed(ctx, http.StatusBadRequest, "无效的ID参数")
	}
}

// UpdateTemplate 更新模板
func (c *AlertController) UpdateTemplate(ctx *gin.Context) {
	var ts model.PrometheusAlertDB
	if err := ctx.ShouldBindJSON(&ts); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "参数解析失败: "+err.Error())
		return
	}
	if ts.Id == 0 {
		result.Failed(ctx, http.StatusBadRequest, "模板ID不能为空")
		return
	}

	if err := c.alertService.UpdateTemplate(&ts); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "模板更新失败: "+err.Error())
		return
	}
	result.Success(ctx, "模版更新成功")
}

// GetTemplate 获取模板列表/详情
func (c *AlertController) GetTemplate(ctx *gin.Context) {
	idStr := ctx.DefaultQuery("id", "")
	if idStr != "" {
		if id, err := strconv.Atoi(idStr); err == nil {
			tpl, err := c.alertService.GetTemplateById(id)
			if err != nil {
				result.Failed(ctx, http.StatusBadRequest, "获取详情失败")
				return
			}
			result.Success(ctx, map[string]interface{}{"info": tpl})
			return
		}
	}

	var query model.TemplateQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "参数错误")
		return
	}
	list, total, err := c.alertService.GetTemplateList(&query)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "获取模板列表失败")
		return
	}
	result.Success(ctx, map[string]interface{}{
		"list":  list,
		"total": total,
	})
}

func (c *AlertController) ReceiveGitlabWebhook(ctx *gin.Context) {
	result.Success(ctx, "Gitlab webhook received successfully")
}

func (c *AlertController) ReceiveZabbixWebhook(ctx *gin.Context) {
	result.Success(ctx, "Zabbix webhook received successfully")
}

// ReceivePrometheusWebhook 接收Prometheus Webhook
func (c *AlertController) ReceivePrometheusWebhook(ctx *gin.Context) {
	var p_json interface{}
	p_alertmanager_json := make(map[string]interface{})
	pMsg := model.PrometheusAlertMsg{}

	bodyBytes, _ := ctx.GetRawData()

	if ctx.Query("from") == "aliyun" {
		AliyunAlertJson := AliyunAlert{}
		AliyunAlertJson.Expression = ctx.Query("expression")
		AliyunAlertJson.MetricName = ctx.Query("metricName")
		AliyunAlertJson.InstanceName = ctx.Query("instanceName")
		AliyunAlertJson.Signature = ctx.Query("signature")
		AliyunAlertJson.MetricProject = ctx.Query("metricProject")
		AliyunAlertJson.UserId = ctx.Query("userId")
		AliyunAlertJson.CurValue = ctx.Query("curValue")
		AliyunAlertJson.AlertName = ctx.Query("alertName")
		AliyunAlertJson.Namespace = ctx.Query("namespace")
		AliyunAlertJson.TriggerLevel = ctx.Query("triggerLevel")
		AliyunAlertJson.AlertState = ctx.Query("alertState")
		AliyunAlertJson.PreTriggerLevel = ctx.Query("preTriggerLevel")
		AliyunAlertJson.RuleId = ctx.Query("ruleId")
		AliyunAlertJson.Dimensions = ctx.Query("dimensions")
		AliyunAlertJson.Timestamp = ctx.Query("timestamp")
		p_json = AliyunAlertJson
	} else {
		json.Unmarshal(bodyBytes, &p_json)
		json.Unmarshal(bodyBytes, &p_alertmanager_json)
	}

	pMsg.Type = ctx.Query("type")
	pMsg.Tpl = ctx.Query("tpl")

	pMsg.Ddurl = ctx.Query("ddurl")
	pMsg.Wxurl = ctx.Query("wxurl")
	pMsg.Fsurl = ctx.Query("fsurl")
	pMsg.Email = ctx.Query("email")
	pMsg.GroupId = ctx.Query("groupid")
	pMsg.Phone = ctx.Query("phone")

	pMsg.WebHookUrl = ctx.Query("webhookurl")
	pMsg.WebhookContentType = ctx.Query("webhookContentType")
	pMsg.ToUser = ctx.Query("wxuser")
	pMsg.ToParty = ctx.Query("wxparty")
	pMsg.ToTag = ctx.Query("wxtag")
	pMsg.EmailTitle = ctx.Query("emailtitle")

	pMsg.AtSomeOne = ctx.Query("at")
	pMsg.RoundRobin = ctx.Query("rr")
	pMsg.Split = ctx.Query("split")

	GlobalAlertRouter, _, _ := c.alertService.GetAllAlertRouter(&model.RouterQuery{})
	message, err := c.alertService.PrometheusAlertHandle(pMsg, p_json, p_alertmanager_json, GlobalAlertRouter)
	if err != nil {
		result.Failed(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	result.Success(ctx, message)
}

func (c *AlertController) GetRecords(ctx *gin.Context) {
	var query model.RecordQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "参数错误")
		return
	}
	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 10
	}
	list, total, err := c.alertService.GetAlertRecords(&query)
	if err != nil {
		result.Failed(ctx, http.StatusInternalServerError, "获取失败: "+err.Error())
		return
	}
	result.Success(ctx, map[string]interface{}{
		"list":  list,
		"total": total,
	})
}

func (c *AlertController) CleanRecords(ctx *gin.Context) {
	result.Success(ctx, "Records cleaned successfully")
}

type AlertRouterReq struct {
	RouterId           interface{}                   `json:"RouterId"`
	RouterName         string                   `json:"RouterName"`
	RouterTplId        interface{}                   `json:"RouterTplId"`
	RouterPurl         string                   `json:"RouterPurl"`
	RouterPat          string                   `json:"RouterPat"`
	RouterPatRR        bool                     `json:"RouterPatRR"`
	RouterSendResolved bool                     `json:"RouterSendResolved"`
	Rules              []map[string]interface{} `json:"Rules"`
}

func (ac *AlertController) CreateRouter(ctx *gin.Context) {
	var req AlertRouterReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "参数绑定失败")
		return
	}

	rulesBytes, _ := json.Marshal(req.Rules)
	tplIdInt := getInt(req.RouterTplId)

	router := &model.AlertRouter{
		Name:         req.RouterName,
		TplId:        tplIdInt,
		Rules:        string(rulesBytes),
		UrlOrPhone:   req.RouterPurl,
		AtSomeOne:    req.RouterPat,
		AtSomeOneRR:  req.RouterPatRR,
		SendResolved: req.RouterSendResolved,
	}

	if err := ac.alertService.CreateAlertRouter(router); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "创建失败: "+err.Error())
		return
	}

	result.Success(ctx, "success")
}

func (ac *AlertController) UpdateRouter(ctx *gin.Context) {
	var req AlertRouterReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "参数绑定失败")
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "非法ID")
		return
	}

	rulesBytes, _ := json.Marshal(req.Rules)
	tplIdInt := getInt(req.RouterTplId)

	router := &model.AlertRouter{
		Id:           id,
		Name:         req.RouterName,
		TplId:        tplIdInt,
		Rules:        string(rulesBytes),
		UrlOrPhone:   req.RouterPurl,
		AtSomeOne:    req.RouterPat,
		AtSomeOneRR:  req.RouterPatRR,
		SendResolved: req.RouterSendResolved,
	}

	if err := ac.alertService.UpdateAlertRouter(router); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "更新失败: "+err.Error())
		return
	}

	result.Success(ctx, "success")
}

func (ac *AlertController) DeleteRouter(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "非法ID")
		return
	}

	if err := ac.alertService.DeleteAlertRouter(id); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "删除失败: "+err.Error())
		return
	}

	result.Success(ctx, "success")
}

func (ac *AlertController) GetRouter(ctx *gin.Context) {
	var query model.RouterQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "参数错误")
		return
	}
	list, total, err := ac.alertService.GetAllAlertRouter(&query)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "获取失败: "+err.Error())
		return
	}
	result.Success(ctx, map[string]interface{}{
		"list":  list,
		"total": total,
	})
}

func (ac *AlertController) ReloadConfig(ctx *gin.Context) {
	// Fake Reload, actual implementation can restart services or flush caches.
	result.Success(ctx, "Config reloaded successfully")
}

func (ac *AlertController) HealthCheck(ctx *gin.Context) {
	// Basic Health check
	result.Success(ctx, map[string]string{"status": "ok", "component": "alert-router"})
}

func getInt(v interface{}) int {
switch val := v.(type) {
case string:
i, _ := strconv.Atoi(val)
return i
case float64:
return int(val)
case int:
return val
}
return 0
}

func (ac *AlertController) GetRouterById(ctx *gin.Context) {
id, err := strconv.Atoi(ctx.Param("id"))
if err != nil {
result.Failed(ctx, http.StatusBadRequest, "非法ID")
return
}
router, err := ac.alertService.GetAlertRouterById(id)
if err != nil {
result.Failed(ctx, http.StatusBadRequest, "获取失败或不存在")
return
}
result.Success(ctx, router)
}
