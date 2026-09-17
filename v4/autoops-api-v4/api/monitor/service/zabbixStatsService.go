package service

import (
	"bytes"
	"dodevops-api/api/configcenter/dao"
	"dodevops-api/api/monitor/model"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"
)

// ZabbixStatsService Zabbix告警统计服务
type ZabbixStatsService struct {
	accountDao *dao.AccountAuthDao
}

// NewZabbixStatsService 创建Zabbix统计服务实例
func NewZabbixStatsService() *ZabbixStatsService {
	return &ZabbixStatsService{
		accountDao: dao.NewAccountAuthDao(),
	}
}

// ZabbixClient Zabbix客户端
type ZabbixClient struct {
	BaseURL    string
	AuthToken  string
	HTTPClient *http.Client
}

// zabbixRequest Zabbix API请求结构
type zabbixRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Auth    interface{} `json:"auth"`
	ID      int         `json:"id"`
}

// zabbixResponse Zabbix API响应结构
type zabbixResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    string `json:"data"`
	} `json:"error"`
	ID int `json:"id"`
}

// NewZabbixClient 创建Zabbix客户端
func NewZabbixClient(accountID uint, service *ZabbixStatsService) (*ZabbixClient, error) {
	account, err := service.accountDao.GetByID(accountID)
	if err != nil {
		return nil, fmt.Errorf("获取账号失败: %w", err)
	}

	if account.Type != 5 { // Zabbix类型
		return nil, fmt.Errorf("账号类型不是Zabbix")
	}

	password, err := account.DecryptPassword()
	if err != nil {
		return nil, fmt.Errorf("解密密码失败: %w", err)
	}

	baseURL := account.Host
	// 移除可能的 /zabbix.php 后缀,并确保使用 api_jsonrpc.php
	if len(baseURL) > 0 && baseURL[len(baseURL)-1] == '/' {
		baseURL = baseURL[:len(baseURL)-1]
	}
	if len(baseURL) > 11 && baseURL[len(baseURL)-11:] == "/zabbix.php" {
		baseURL = baseURL[:len(baseURL)-11]
	}

	client := &ZabbixClient{
		BaseURL: baseURL + "/api_jsonrpc.php",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// 登录获取token
	authToken, err := client.login(account.Name, password)
	if err != nil {
		return nil, fmt.Errorf("登录Zabbix失败: %w", err)
	}

	client.AuthToken = authToken
	return client, nil
}

// login Zabbix登录
func (c *ZabbixClient) login(username, password string) (string, error) {
	req := zabbixRequest{
		JSONRPC: "2.0",
		Method:  "user.login",
		Params: map[string]string{
			"user":     username,
			"password": password,
		},
		Auth: nil,
		ID:   1,
	}

	var resp zabbixResponse
	if err := c.doRequest(&req, &resp); err != nil {
		return "", err
	}

	if resp.Error != nil {
		return "", fmt.Errorf("Zabbix登录错误: %s", resp.Error.Data)
	}

	var token string
	if err := json.Unmarshal(resp.Result, &token); err != nil {
		return "", fmt.Errorf("解析登录响应失败: %w", err)
	}

	return token, nil
}

// doRequest 执行Zabbix API请求
func (c *ZabbixClient) doRequest(req *zabbixRequest, resp *zabbixResponse) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer httpResp.Body.Close()

	if err := json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	return nil
}

// getEvents 获取指定时间范围的事件
func (c *ZabbixClient) getEvents(startTime, endTime int64) ([]model.ZabbixEvent, error) {
	req := zabbixRequest{
		JSONRPC: "2.0",
		Method:  "event.get",
		Params: map[string]interface{}{
			"output":     []string{"eventid", "source", "object", "objectid", "clock", "value", "severity", "name"},
			"time_from":  startTime,
			"time_till":  endTime,
			"source":     0, // 触发器事件
			"object":     0,
			"sortfield":  []string{"clock"},
			"sortorder":  "DESC",
		},
		Auth: c.AuthToken,
		ID:   2,
	}

	var resp zabbixResponse
	if err := c.doRequest(&req, &resp); err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("获取事件错误: %s", resp.Error.Data)
	}

	var events []model.ZabbixEvent
	if err := json.Unmarshal(resp.Result, &events); err != nil {
		return nil, fmt.Errorf("解析事件数据失败: %w", err)
	}

	return events, nil
}

// getProblems 获取指定时间范围的问题(告警)
func (c *ZabbixClient) getProblems(startTime, endTime int64) ([]map[string]interface{}, error) {
	req := zabbixRequest{
		JSONRPC: "2.0",
		Method:  "problem.get",
		Params: map[string]interface{}{
			"output":        "extend",
			"selectHosts":   []string{"hostid", "host", "name"},
			"selectGroups":  []string{"groupid", "name"},
			"time_from":     startTime,
			"time_till":     endTime,
			"recent":        false,
			"sortfield":     []string{"eventid"},
			"sortorder":     "DESC",
		},
		Auth: c.AuthToken,
		ID:   3,
	}

	var resp zabbixResponse
	if err := c.doRequest(&req, &resp); err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("获取问题错误: %s", resp.Error.Data)
	}

	var problems []map[string]interface{}
	if err := json.Unmarshal(resp.Result, &problems); err != nil {
		return nil, fmt.Errorf("解析问题数据失败: %w", err)
	}

	return problems, nil
}

// GetAlertTopByName 获取告警TOP(按告警名称)
func (s *ZabbixStatsService) GetAlertTopByName(accountID uint, startTime, endTime int64, topN int, timeRange string) (*model.ZabbixAlertTopResult, error) {
	client, err := NewZabbixClient(accountID, s)
	if err != nil {
		return nil, err
	}

	problems, err := client.getProblems(startTime, endTime)
	if err != nil {
		return nil, err
	}

	// 统计每个告警名称的数量
	alertCount := make(map[string]int)
	for _, problem := range problems {
		name, ok := problem["name"].(string)
		if !ok {
			continue
		}
		alertCount[name]++
	}

	// 转换为切片并排序
	items := make([]model.ZabbixAlertTopItem, 0, len(alertCount))
	totalCount := len(problems)

	for name, count := range alertCount {
		percentage := 0.0
		if totalCount > 0 {
			percentage = float64(count) / float64(totalCount) * 100
		}
		items = append(items, model.ZabbixAlertTopItem{
			Name:       name,
			Count:      count,
			Percentage: percentage,
		})
	}

	// 按数量降序排序
	sort.Slice(items, func(i, j int) bool {
		return items[i].Count > items[j].Count
	})

	// 取TOP N
	if len(items) > topN {
		items = items[:topN]
	}

	return &model.ZabbixAlertTopResult{
		Items:     items,
		Total:     totalCount,
		StartTime: startTime,
		EndTime:   endTime,
		TimeRange: timeRange,
	}, nil
}

// GetAlertTopByHost 获取告警TOP(按主机)
func (s *ZabbixStatsService) GetAlertTopByHost(accountID uint, startTime, endTime int64, topN int, timeRange string) (*model.ZabbixAlertTopResult, error) {
	client, err := NewZabbixClient(accountID, s)
	if err != nil {
		return nil, err
	}

	problems, err := client.getProblems(startTime, endTime)
	if err != nil {
		return nil, err
	}

	// 统计每个主机的告警数量
	hostCount := make(map[string]int)
	for _, problem := range problems {
		hosts, ok := problem["hosts"].([]interface{})
		if !ok || len(hosts) == 0 {
			continue
		}

		for _, hostInterface := range hosts {
			host, ok := hostInterface.(map[string]interface{})
			if !ok {
				continue
			}
			hostName, ok := host["name"].(string)
			if !ok || hostName == "" {
				hostName, _ = host["host"].(string)
			}
			if hostName != "" {
				hostCount[hostName]++
			}
		}
	}

	// 转换为切片并排序
	items := make([]model.ZabbixAlertTopItem, 0, len(hostCount))
	totalCount := len(problems)

	for hostName, count := range hostCount {
		percentage := 0.0
		if totalCount > 0 {
			percentage = float64(count) / float64(totalCount) * 100
		}
		items = append(items, model.ZabbixAlertTopItem{
			Name:       hostName,
			Count:      count,
			Percentage: percentage,
		})
	}

	// 按数量降序排序
	sort.Slice(items, func(i, j int) bool {
		return items[i].Count > items[j].Count
	})

	// 取TOP N
	if len(items) > topN {
		items = items[:topN]
	}

	return &model.ZabbixAlertTopResult{
		Items:     items,
		Total:     totalCount,
		StartTime: startTime,
		EndTime:   endTime,
		TimeRange: timeRange,
	}, nil
}

// GetAlertTopByHostGroup 获取告警TOP(按主机群组)
func (s *ZabbixStatsService) GetAlertTopByHostGroup(accountID uint, startTime, endTime int64, topN int, timeRange string) (*model.ZabbixAlertTopResult, error) {
	client, err := NewZabbixClient(accountID, s)
	if err != nil {
		return nil, err
	}

	// 1. 获取问题列表
	req := zabbixRequest{
		JSONRPC: "2.0",
		Method:  "problem.get",
		Params: map[string]interface{}{
			"output":    "extend",
			"time_from": startTime,
			"time_till": endTime,
			"recent":    false,
			"sortfield": []string{"eventid"},
			"sortorder": "DESC",
		},
		Auth: client.AuthToken,
		ID:   3,
	}

	var resp zabbixResponse
	if err := client.doRequest(&req, &resp); err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("获取问题错误: %s", resp.Error.Data)
	}

	var problems []map[string]interface{}
	if err := json.Unmarshal(resp.Result, &problems); err != nil {
		return nil, fmt.Errorf("解析问题数据失败: %w", err)
	}

	if len(problems) == 0 {
		return &model.ZabbixAlertTopResult{
			Items:     []model.ZabbixAlertTopItem{},
			Total:     0,
			StartTime: startTime,
			EndTime:   endTime,
			TimeRange: timeRange,
		}, nil
	}

	// 2. 收集所有触发器ID
	triggerIDs := make(map[string]bool)
	problemTriggerMap := make(map[int]string) // problem index -> triggerid

	for idx, problem := range problems {
		objectid, ok := problem["objectid"].(string)
		if !ok || objectid == "" {
			continue
		}
		triggerIDs[objectid] = true
		problemTriggerMap[idx] = objectid
	}

	if len(triggerIDs) == 0 {
		return &model.ZabbixAlertTopResult{
			Items:     []model.ZabbixAlertTopItem{},
			Total:     len(problems),
			StartTime: startTime,
			EndTime:   endTime,
			TimeRange: timeRange,
		}, nil
	}

	// 3. 通过触发器获取主机信息
	triggerIDList := make([]string, 0, len(triggerIDs))
	for tid := range triggerIDs {
		triggerIDList = append(triggerIDList, tid)
	}

	triggerReq := zabbixRequest{
		JSONRPC: "2.0",
		Method:  "trigger.get",
		Params: map[string]interface{}{
			"output":      []string{"triggerid"},
			"selectHosts": []string{"hostid", "host", "name"},
			"triggerids":  triggerIDList,
		},
		Auth: client.AuthToken,
		ID:   4,
	}

	var triggerResp zabbixResponse
	if err := client.doRequest(&triggerReq, &triggerResp); err != nil {
		return nil, err
	}

	if triggerResp.Error != nil {
		return nil, fmt.Errorf("获取触发器信息错误: %s", triggerResp.Error.Data)
	}

	var triggers []map[string]interface{}
	if err := json.Unmarshal(triggerResp.Result, &triggers); err != nil {
		return nil, fmt.Errorf("解析触发器数据失败: %w", err)
	}

	// 4. 构建 triggerID -> hostIDs 映射
	triggerHostMap := make(map[string][]string)
	allHostIDs := make(map[string]bool)

	for _, trigger := range triggers {
		triggerID, ok := trigger["triggerid"].(string)
		if !ok {
			continue
		}

		hosts, ok := trigger["hosts"].([]interface{})
		if !ok || len(hosts) == 0 {
			continue
		}

		hostIDs := make([]string, 0)
		for _, hostInterface := range hosts {
			host, ok := hostInterface.(map[string]interface{})
			if !ok {
				continue
			}

			var hostID string
			if id, ok := host["hostid"].(string); ok {
				hostID = id
			} else if id, ok := host["hostid"].(float64); ok {
				hostID = fmt.Sprintf("%.0f", id)
			}

			if hostID != "" {
				hostIDs = append(hostIDs, hostID)
				allHostIDs[hostID] = true
			}
		}
		triggerHostMap[triggerID] = hostIDs
	}

	if len(allHostIDs) == 0 {
		return &model.ZabbixAlertTopResult{
			Items:     []model.ZabbixAlertTopItem{},
			Total:     len(problems),
			StartTime: startTime,
			EndTime:   endTime,
			TimeRange: timeRange,
		}, nil
	}

	// 5. 获取主机的群组信息
	hostIDList := make([]string, 0, len(allHostIDs))
	for hostID := range allHostIDs {
		hostIDList = append(hostIDList, hostID)
	}

	hostReq := zabbixRequest{
		JSONRPC: "2.0",
		Method:  "host.get",
		Params: map[string]interface{}{
			"output":       []string{"hostid", "host"},
			"selectGroups": []string{"groupid", "name"},
			"hostids":      hostIDList,
		},
		Auth: client.AuthToken,
		ID:   5,
	}

	var hostResp zabbixResponse
	if err := client.doRequest(&hostReq, &hostResp); err != nil {
		return nil, err
	}

	if hostResp.Error != nil {
		return nil, fmt.Errorf("获取主机群组信息错误: %s", hostResp.Error.Data)
	}

	var hosts []map[string]interface{}
	if err := json.Unmarshal(hostResp.Result, &hosts); err != nil {
		return nil, fmt.Errorf("解析主机群组数据失败: %w", err)
	}

	// 6. 构建 hostID -> groups 的映射
	hostGroupMap := make(map[string][]string)
	for _, host := range hosts {
		var hostID string
		if id, ok := host["hostid"].(string); ok {
			hostID = id
		} else if id, ok := host["hostid"].(float64); ok {
			hostID = fmt.Sprintf("%.0f", id)
		}

		if hostID == "" {
			continue
		}

		groups, ok := host["groups"].([]interface{})
		if !ok {
			continue
		}

		groupNames := make([]string, 0)
		for _, groupInterface := range groups {
			group, ok := groupInterface.(map[string]interface{})
			if !ok {
				continue
			}
			groupName, ok := group["name"].(string)
			if ok && groupName != "" {
				groupNames = append(groupNames, groupName)
			}
		}
		hostGroupMap[hostID] = groupNames
	}

	// 7. 统计每个主机群组的告警数量
	groupCount := make(map[string]int)
	for problemIdx, triggerID := range problemTriggerMap {
		if problemIdx >= len(problems) {
			continue
		}

		hostIDsForTrigger := triggerHostMap[triggerID]
		for _, hostID := range hostIDsForTrigger {
			groupNames := hostGroupMap[hostID]
			for _, groupName := range groupNames {
				groupCount[groupName]++
			}
		}
	}

	// 转换为切片并排序
	items := make([]model.ZabbixAlertTopItem, 0, len(groupCount))
	totalCount := len(problems)

	for groupName, count := range groupCount {
		percentage := 0.0
		if totalCount > 0 {
			percentage = float64(count) / float64(totalCount) * 100
		}
		items = append(items, model.ZabbixAlertTopItem{
			Name:       groupName,
			Count:      count,
			Percentage: percentage,
		})
	}

	// 按数量降序排序
	sort.Slice(items, func(i, j int) bool {
		return items[i].Count > items[j].Count
	})

	// 取TOP N
	if len(items) > topN {
		items = items[:topN]
	}

	return &model.ZabbixAlertTopResult{
		Items:     items,
		Total:     totalCount,
		StartTime: startTime,
		EndTime:   endTime,
		TimeRange: timeRange,
	}, nil
}

// GetNotificationStats 获取告警数量统计(按天)
func (s *ZabbixStatsService) GetNotificationStats(accountID uint, startTime, endTime int64, timeRange string) (*model.ZabbixNotificationStats, error) {
	client, err := NewZabbixClient(accountID, s)
	if err != nil {
		return nil, err
	}

	// 查询告警问题 (problem.get)
	req := zabbixRequest{
		JSONRPC: "2.0",
		Method:  "problem.get",
		Params: map[string]interface{}{
			"output":    []string{"eventid", "clock", "name", "severity"},
			"time_from": startTime,
			"time_till": endTime,
			"recent":    false,
			"sortfield": []string{"eventid"},
			"sortorder": "DESC",
		},
		Auth: client.AuthToken,
		ID:   4,
	}

	var resp zabbixResponse
	if err := client.doRequest(&req, &resp); err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("获取告警错误: %s", resp.Error.Data)
	}

	var problems []map[string]interface{}
	if err := json.Unmarshal(resp.Result, &problems); err != nil {
		return nil, fmt.Errorf("解析告警数据失败: %w", err)
	}

	// 按天统计
	dailyCount := make(map[string]int)
	totalCount := len(problems)

	for _, problem := range problems {
		clockStr, ok := problem["clock"].(string)
		if !ok {
			continue
		}

		// 解析时间戳
		var timestamp int64
		fmt.Sscanf(clockStr, "%d", &timestamp)

		// 转换为日期
		t := time.Unix(timestamp, 0)
		dateStr := t.Format("2006-01-02")

		dailyCount[dateStr]++
	}

	// 生成完整的日期范围(包含0数据的日期)
	items := make([]model.ZabbixNotificationItem, 0)
	start := time.Unix(startTime, 0)
	end := time.Unix(endTime, 0)

	// 将开始和结束时间都调整为当天的0点
	startDate := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	endDate := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())

	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		count := dailyCount[dateStr]
		percentage := 0.0
		if totalCount > 0 {
			percentage = float64(count) / float64(totalCount) * 100
		}

		items = append(items, model.ZabbixNotificationItem{
			Date:       dateStr,
			Timestamp:  d.Unix(),
			Count:      count,
			Percentage: percentage,
		})
	}

	return &model.ZabbixNotificationStats{
		Items:      items,
		TotalCount: totalCount,
		StartTime:  startTime,
		EndTime:    endTime,
		TimeRange:  timeRange,
	}, nil
}

// getTimeRange 获取时间范围
func getTimeRange(rangeType string) (startTime, endTime int64) {
	now := time.Now()
	endTime = now.Unix()

	switch rangeType {
	case "today":
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		startTime = today.Unix()
	case "yesterday":
		yesterday := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location())
		yesterdayEnd := yesterday.Add(24 * time.Hour)
		startTime = yesterday.Unix()
		endTime = yesterdayEnd.Unix()
	case "last7days":
		startTime = now.Add(-7 * 24 * time.Hour).Unix()
	case "last30days":
		startTime = now.Add(-30 * 24 * time.Hour).Unix()
	default:
		// 默认今天
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		startTime = today.Unix()
	}

	return startTime, endTime
}
