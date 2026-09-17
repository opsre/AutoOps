package service

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "time"

    "dodevops-api/api/cmdb/dao"
    "dodevops-api/api/monitor/model"
    "dodevops-api/common/config"
)

type WindowsMonitorService struct {
    deviceDao      *dao.CmdbSnmpDeviceDao
    prometheusURL  string
}

func NewWindowsMonitorService() *WindowsMonitorService {
    return &WindowsMonitorService{
        deviceDao:     dao.NewCmdbSnmpDeviceDao(),
        prometheusURL: config.Config.Monitor.Prometheus.URL,
    }
}

// WindowsMonitorResponse 监控返回结构
type WindowsMonitorResponse struct {
    DeviceID   uint                   `json:"deviceId"`
    DeviceName string                 `json:"deviceName"`
    IPAddress  string                 `json:"ipAddress"`
    Ranges     map[string]interface{} `json:"ranges"` // 30m/1d/7d -> 各指标时序
    UptimeSec  float64                `json:"uptimeSec"`
}

// queryRange 通用的 Prometheus range 查询
func (s *WindowsMonitorService) queryRange(expr string, start, end time.Time, step string) (*model.PrometheusQueryResult, error) {
    if s.prometheusURL == "" {
        return nil, fmt.Errorf("Prometheus URL not configured")
    }
    req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/query_range", s.prometheusURL), nil)
    if err != nil { return nil, err }
    q := req.URL.Query()
    q.Set("query", expr)
    q.Set("start", fmt.Sprintf("%d", start.Unix()))
    q.Set("end", fmt.Sprintf("%d", end.Unix()))
    q.Set("step", step)
    req.URL.RawQuery = q.Encode()
    req.Header.Set("Cache-Control", "no-cache")
    client := &http.Client{ Timeout: 30 * time.Second }
    resp, err := client.Do(req)
    if err != nil { return nil, err }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil { return nil, err }
    var r model.PrometheusQueryResult
    if err := json.Unmarshal(body, &r); err != nil {
        return nil, fmt.Errorf("parse prometheus response failed: %v", err)
    }
    if r.Status != "success" { return nil, fmt.Errorf("prometheus query failed: %s", string(body)) }
    return &r, nil
}

// queryInstant 即时查询（用于 uptime 等）
func (s *WindowsMonitorService) queryInstant(expr string) (float64, error) {
    if s.prometheusURL == "" { return 0, fmt.Errorf("Prometheus URL not configured") }
    u := fmt.Sprintf("%s/api/v1/query?query=%s&time=%d", s.prometheusURL, url.QueryEscape(expr), time.Now().Unix())
    resp, err := http.Get(u)
    if err != nil { return 0, err }
    defer resp.Body.Close()
    b, _ := ioutil.ReadAll(resp.Body)
    var r model.PrometheusQueryResult
    if err := json.Unmarshal(b, &r); err != nil { return 0, err }
    if r.Status != "success" || len(r.Data.Result) == 0 || len(r.Data.Result[0].Value) < 2 { return 0, fmt.Errorf("no data") }
    vStr, _ := r.Data.Result[0].Value[1].(string)
    var v float64
    fmt.Sscanf(vStr, "%f", &v)
    return v, nil
}

// GetWindowsMonitorMetrics 获取Windows主机关键指标（30m/1d/7d）
func (s *WindowsMonitorService) GetWindowsMonitorMetrics(deviceID uint) (*WindowsMonitorResponse, error) {
    dev, err := s.deviceDao.GetSnmpDeviceById(deviceID)
    if err != nil { return nil, fmt.Errorf("设备不存在: %v", err) }
    if dev.DeviceType != "windows" { return nil, fmt.Errorf("设备类型必须为 windows") }

    // 时间范围定义
    now := time.Now()
    ranges := map[string]struct{ dur time.Duration; step string }{
        "30m": {30 * time.Minute, "30s"},
        "1d":  {24 * time.Hour, "5m"},
        "7d":  {7 * 24 * time.Hour, "15m"},
    }

    // 标签：安装脚本使用 job=<DeviceName>
    job := dev.DeviceName

    resp := &WindowsMonitorResponse{
        DeviceID:   dev.ID,
        DeviceName: dev.DeviceName,
        IPAddress:  dev.IPAddress,
        Ranges:     map[string]interface{}{},
        UptimeSec:  0,
    }

    // uptime 优先 windows_system_system_up_time 或 windows_system_system_up_time_seconds
    if v, err := s.queryInstant(fmt.Sprintf("windows_system_system_up_time{job=\"%s\"}", job)); err == nil {
        resp.UptimeSec = v
    } else {
        if v1, err1 := s.queryInstant(fmt.Sprintf("windows_system_system_up_time_seconds{job=\"%s\"}", job)); err1 == nil {
            resp.UptimeSec = v1
        } else if v2, err2 := s.queryInstant(fmt.Sprintf("time() - min(windows_process_start_time_seconds{job=\"%s\"})", job)); err2 == nil {
            resp.UptimeSec = v2
        }
    }

    // 若仍未取到，使用区间查询回退（取最近点）
    if resp.UptimeSec <= 0 {
        // 30分钟区间内的最大值作为近似（避免瞬时空窗）
        if r, err := s.queryRange(fmt.Sprintf("max_over_time(windows_system_system_up_time_seconds{job=\"%s\"}[15m])", job), now.Add(-30*time.Minute), now, "1m"); err == nil {
            if len(r.Data.Result) > 0 {
                vals := r.Data.Result[0].Values
                if len(vals) > 0 {
                    if str, ok := vals[len(vals)-1][1].(string); ok {
                        var v float64
                        fmt.Sscanf(str, "%f", &v)
                        if v > 0 { resp.UptimeSec = v }
                    }
                }
            }
        }
    }
    if resp.UptimeSec <= 0 {
        if r, err := s.queryRange(fmt.Sprintf("max_over_time(windows_system_system_up_time{job=\"%s\"}[15m])", job), now.Add(-30*time.Minute), now, "1m"); err == nil {
            if len(r.Data.Result) > 0 {
                vals := r.Data.Result[0].Values
                if len(vals) > 0 {
                    if str, ok := vals[len(vals)-1][1].(string); ok {
                        var v float64
                        fmt.Sscanf(str, "%f", &v)
                        if v > 0 { resp.UptimeSec = v }
                    }
                }
            }
        }
    }

    // 逐范围查询
    for key, cfg := range ranges {
        start := now.Add(-cfg.dur)
        end := now
        step := cfg.step

        // CPU 使用率：100 - avg(rate(idle[5m]))*100，限制0~100（Windows 推荐口径）
        cpuExpr := fmt.Sprintf("clamp_min(clamp_max(100 - (avg by (job) (rate(windows_cpu_time_total{job=\"%s\",mode=\"idle\"}[5m])) * 100), 0), 100)", job)
        // 内存使用率：100*(1 - avail/total)
        memExpr := fmt.Sprintf("100 * (1 - (windows_memory_available_bytes{job=\"%s\"} / windows_memory_physical_total_bytes{job=\"%s\"}))", job, job)
        // 磁盘容量使用率（C盘）
        diskExpr := fmt.Sprintf("100 * (1 - (windows_logical_disk_free_bytes{job=\"%s\",volume=\"C:\"} / windows_logical_disk_size_bytes{job=\"%s\",volume=\"C:\"}))", job, job)
        // 磁盘读/写平均延迟（C盘）
        readLatExpr := fmt.Sprintf("rate(windows_logical_disk_read_latency_seconds_total{job=\"%s\",volume=\"C:\"}[%s]) / rate(windows_logical_disk_reads_total{job=\"%s\",volume=\"C:\"}[%s])", job, step, job, step)
        writeLatExpr := fmt.Sprintf("rate(windows_logical_disk_write_latency_seconds_total{job=\"%s\",volume=\"C:\"}[%s]) / rate(windows_logical_disk_writes_total{job=\"%s\",volume=\"C:\"}[%s])", job, step, job, step)
        // 磁盘吞吐（字节/秒）
        readBpsExpr := fmt.Sprintf("rate(windows_logical_disk_read_bytes_total{job=\"%s\",volume=\"C:\"}[%s])", job, step)
        writeBpsExpr := fmt.Sprintf("rate(windows_logical_disk_write_bytes_total{job=\"%s\",volume=\"C:\"}[%s])", job, step)
        // 网络带宽（字节/秒，汇总）
        netInExpr := fmt.Sprintf("sum(rate(windows_net_bytes_received_total{job=\"%s\"}[%s]))", job, step)
        netOutExpr := fmt.Sprintf("sum(rate(windows_net_bytes_sent_total{job=\"%s\"}[%s]))", job, step)
        // 进程数：优先 windows_system_processes，其次多重回退
        procCountExpr := fmt.Sprintf(
            "sum by (job) (windows_system_processes{job=\"%s\"}) or "+
                "sum by (job) (windows_os_processes{job=\"%s\"}) or "+
                "sum by (job) (windows_cs_processes{job=\"%s\"}) or "+
                // 以 pid 去重，避免同一进程的多条时间序列导致放大
                "count by (job) (count by (job,pid) (windows_process_start_time_seconds{job=\"%s\"}))",
            job, job, job, job)
        // 线程数
        threadsExpr := fmt.Sprintf("windows_system_threads{job=\"%s\"}", job)

        // 覆盖CPU表达式为 per-core 归一化后聚合
        cpuExpr = fmt.Sprintf(
            "clamp_min(clamp_max( avg by (job) ( sum by (job, instance, core) (rate(windows_cpu_time_total{job=\"%s\",mode!=\"idle\"}[5m])) / ( sum by (job, instance, core) (rate(windows_cpu_time_total{job=\"%s\",mode!=\"idle\"}[5m])) + sum by (job, instance, core) (rate(windows_cpu_time_total{job=\"%s\",mode=\"idle\"}[5m])) ) ), 1), 0) * 100",
            job, job, job,
        )

        // 执行查询
        data := map[string]*model.PrometheusQueryResult{}
        if r, err := s.queryRange(cpuExpr, start, end, step); err == nil { data["cpu"] = r }
        if r, err := s.queryRange(memExpr, start, end, step); err == nil { data["memory"] = r }
        if r, err := s.queryRange(diskExpr, start, end, step); err == nil { data["disk_usage"] = r }
        if r, err := s.queryRange(readLatExpr, start, end, step); err == nil { data["disk_read_latency"] = r }
        if r, err := s.queryRange(writeLatExpr, start, end, step); err == nil { data["disk_write_latency"] = r }
        if r, err := s.queryRange(readBpsExpr, start, end, step); err == nil { data["disk_read_bps"] = r }
        if r, err := s.queryRange(writeBpsExpr, start, end, step); err == nil { data["disk_write_bps"] = r }
        if r, err := s.queryRange(netInExpr, start, end, step); err == nil { data["net_in_bps"] = r }
        if r, err := s.queryRange(netOutExpr, start, end, step); err == nil { data["net_out_bps"] = r }
        if r, err := s.queryRange(procCountExpr, start, end, step); err == nil { data["process_count"] = r }
        if r, err := s.queryRange(threadsExpr, start, end, step); err == nil { data["threads"] = r }

        resp.Ranges[key] = data
    }

    return resp, nil
}

// WindowsRealtimeMetrics 实时指标返回结构
type WindowsRealtimeMetrics struct {
    DeviceID   uint    `json:"deviceId"`
    DeviceName string  `json:"deviceName"`
    IPAddress  string  `json:"ipAddress"`
    CPUPercent float64 `json:"cpuPercent"`   // CPU使用率
    MemPercent float64 `json:"memPercent"`   // 内存使用率
    Timestamp  int64   `json:"timestamp"`    // 查询时间戳
}

// GetWindowsRealtimeMetrics 获取Windows主机实时CPU和内存使用率
func (s *WindowsMonitorService) GetWindowsRealtimeMetrics(deviceID uint) (*WindowsRealtimeMetrics, error) {
    dev, err := s.deviceDao.GetSnmpDeviceById(deviceID)
    if err != nil { return nil, fmt.Errorf("设备不存在: %v", err) }
    if dev.DeviceType != "windows" { return nil, fmt.Errorf("设备类型必须为 windows") }

    job := dev.DeviceName
    now := time.Now()

    resp := &WindowsRealtimeMetrics{
        DeviceID:   dev.ID,
        DeviceName: dev.DeviceName,
        IPAddress:  dev.IPAddress,
        Timestamp:  now.Unix(),
    }

    // CPU使用率（实时）
    cpuExpr := fmt.Sprintf(
        "clamp_min(clamp_max( avg by (job) ( sum by (job, instance, core) (rate(windows_cpu_time_total{job=\"%s\",mode!=\"idle\"}[1m])) / ( sum by (job, instance, core) (rate(windows_cpu_time_total{job=\"%s\",mode!=\"idle\"}[1m])) + sum by (job, instance, core) (rate(windows_cpu_time_total{job=\"%s\",mode=\"idle\"}[1m])) ) ), 1), 0) * 100",
        job, job, job,
    )
    if cpu, err := s.queryInstant(cpuExpr); err == nil {
        resp.CPUPercent = cpu
    }

    // 内存使用率（实时）
    memExpr := fmt.Sprintf("100 * (1 - (windows_memory_available_bytes{job=\"%s\"} / windows_memory_physical_total_bytes{job=\"%s\"}))", job, job)
    if mem, err := s.queryInstant(memExpr); err == nil {
        resp.MemPercent = mem
    }

    return resp, nil
}

// GetBatchWindowsRealtimeMetrics 批量获取Windows主机实时CPU和内存使用率
func (s *WindowsMonitorService) GetBatchWindowsRealtimeMetrics(deviceIDs []uint) ([]*WindowsRealtimeMetrics, error) {
    if len(deviceIDs) == 0 {
        return []*WindowsRealtimeMetrics{}, nil
    }

    results := make([]*WindowsRealtimeMetrics, 0, len(deviceIDs))
    for _, deviceID := range deviceIDs {
        if metrics, err := s.GetWindowsRealtimeMetrics(deviceID); err == nil {
            results = append(results, metrics)
        }
        // 忽略单个设备查询失败的情况，继续查询其他设备
    }

    return results, nil
}