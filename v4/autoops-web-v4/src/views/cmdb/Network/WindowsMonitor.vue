<template>
  <div class="windows-monitor-page">
    <el-card shadow="hover" class="monitor-card">
      <template #header>
        <div class="card-header">
          <div class="title-wrapper">
            <span class="title">Windows 监控</span>
          </div>
          <div class="header-actions">
            <el-radio-group v-model="currentRange" size="small" @change="renderCharts">
              <el-radio-button label="30m">30分钟</el-radio-button>
              <el-radio-button label="1d">1天</el-radio-button>
              <el-radio-button label="7d">7天</el-radio-button>
            </el-radio-group>
            <el-button size="small" @click="goBack" class="back-btn" style="margin-left: 8px;">
              <el-icon><ArrowLeft /></el-icon>
              返回
            </el-button>
          </div>
        </div>
      </template>

      <!-- 设备基础信息 -->
      <el-card class="info-section" shadow="never">
        <template #header>
          <div class="section-header-with-line">
            <span class="section-title">设备信息</span>
            <div class="title-line"></div>
          </div>
        </template>
        <el-descriptions :column="3" border>
          <el-descriptions-item label="设备名称">{{ deviceInfo.deviceName || '-' }}</el-descriptions-item>
          <el-descriptions-item label="IP地址">{{ deviceInfo.ipAddress || '-' }}</el-descriptions-item>
          <el-descriptions-item label="运行时长">{{ formatUptime(deviceInfo.uptimeSec) }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 实时指标概览 -->
      <el-card class="metrics-section" shadow="never">
        <template #header>
          <div class="section-header-with-line">
            <span class="section-title">实时指标</span>
            <div class="title-line"></div>
          </div>
        </template>
        <el-row :gutter="15" class="metrics-row">
          <el-col class="metric-col">
            <div class="metric-item">
              <div class="metric-content">
                <div class="metric-icon cpu-icon">
                  <el-icon><Cpu /></el-icon>
                </div>
                <div class="metric-info">
                  <div class="metric-label">CPU使用率</div>
                  <div class="metric-value" :class="getMetricClass(latest('cpu'))">
                    {{ formatPercent(latest('cpu')) }}
                  </div>
                </div>
              </div>
            </div>
          </el-col>
          <el-col class="metric-col">
            <div class="metric-item">
              <div class="metric-content">
                <div class="metric-icon memory-icon">
                  <el-icon><MemoryCard /></el-icon>
                </div>
                <div class="metric-info">
                  <div class="metric-label">内存使用率</div>
                  <div class="metric-value" :class="getMetricClass(latest('memory'))">
                    {{ formatPercent(latest('memory')) }}
                  </div>
                </div>
              </div>
            </div>
          </el-col>
          <el-col class="metric-col">
            <div class="metric-item">
              <div class="metric-content">
                <div class="metric-icon disk-icon">
                  <el-icon><Odometer /></el-icon>
                </div>
                <div class="metric-info">
                  <div class="metric-label">C盘使用率</div>
                  <div class="metric-value" :class="getMetricClass(latest('disk_usage'))">
                    {{ formatPercent(latest('disk_usage')) }}
                  </div>
                </div>
              </div>
            </div>
          </el-col>
          <el-col class="metric-col">
            <div class="metric-item">
              <div class="metric-content">
                <div class="metric-icon process-icon">
                  <el-icon><Monitor /></el-icon>
                </div>
                <div class="metric-info">
                  <div class="metric-label">进程数</div>
                  <div class="metric-value normal">
                    {{ formatNumber(latest('process_count')) }}
                  </div>
                </div>
              </div>
            </div>
          </el-col>
        </el-row>
      </el-card>

      <!-- 趋势图表 -->
      <el-card class="charts-section" shadow="never">
        <template #header>
          <span class="section-title">{{ rangeLabel }}趋势</span>
        </template>
        <el-row :gutter="20">
          <el-col :span="12">
            <div class="chart-card"><div class="chart-title">CPU使用率(%)</div><div ref="cpuChart" class="chart"></div></div>
          </el-col>
          <el-col :span="12">
            <div class="chart-card"><div class="chart-title">内存使用率(%)</div><div ref="memoryChart" class="chart"></div></div>
          </el-col>
          <el-col :span="12">
            <div class="chart-card"><div class="chart-title">C盘使用率(%)</div><div ref="diskUsageChart" class="chart"></div></div>
          </el-col>
          <el-col :span="12">
            <div class="chart-card"><div class="chart-title">进程数</div><div ref="processChart" class="chart"></div></div>
          </el-col>
          <el-col :span="12">
            <div class="chart-card"><div class="chart-title">磁盘读/写延迟(ms/次)</div><div ref="diskLatencyChart" class="chart"></div></div>
          </el-col>
          <el-col :span="12">
            <div class="chart-card"><div class="chart-title">磁盘读/写吞吐(MB/s)</div><div ref="diskBpsChart" class="chart"></div></div>
          </el-col>
          <el-col :span="12">
            <div class="chart-card"><div class="chart-title">网络入/出带宽(Mbps)</div><div ref="netChart" class="chart"></div></div>
          </el-col>
        </el-row>
      </el-card>
    </el-card>
  </div>
  </template>

<script>
import * as echarts from 'echarts'
import API from '@/api/cmdb'
import { ArrowLeft, Cpu, MemoryCard, Monitor, Odometer } from '@element-plus/icons-vue'

export default {
  name: 'WindowsMonitor',
  components: { ArrowLeft, Cpu, MemoryCard, Monitor, Odometer },
  data() {
    return {
      deviceId: null,
      deviceInfo: { deviceName: '', ipAddress: '', uptimeSec: 0 },
      ranges: { '30m': {}, '1d': {}, '7d': {} },
      currentRange: '30m',
      rangeOptions: [],
      chartInstances: { cpu: null, memory: null, diskUsage: null, process: null, diskLatency: null, diskBps: null, net: null },
      refreshInterval: null
    }
  },
  computed: {
    rangeLabel() {
      return this.currentRange === '30m' ? '近30分钟' : this.currentRange === '1d' ? '近1天' : '近7天'
    }
  },
  methods: {
    async loadDeviceInfo() {
      try {
        const { data: res } = await API.getSnmpDeviceById(this.deviceId)
        if (res.code === 200) {
          this.deviceInfo.deviceName = res.data.deviceName || ''
          this.deviceInfo.ipAddress = res.data.ipAddress || ''
        }
      } catch (e) { /* 静默失败 */ }
    },
    async loadMonitorData() {
      try {
        const { data: res } = await API.getWindowsDeviceMonitor(this.deviceId)
        if (res.code !== 200) throw new Error(res.message || '获取监控数据失败')
        const data = res.data || {}
        this.deviceInfo.uptimeSec = Number(data.uptimeSec || 0)

        // 解析各时间范围
        const parseRange = (rangeObj) => ({
          cpu: this.parsePromResult(rangeObj?.cpu),
          memory: this.parsePromResult(rangeObj?.memory),
          disk_usage: this.parsePromResult(rangeObj?.disk_usage),
          disk_read_latency: this.parsePromResult(rangeObj?.disk_read_latency),
          disk_write_latency: this.parsePromResult(rangeObj?.disk_write_latency),
          disk_read_bps: this.parsePromResult(rangeObj?.disk_read_bps),
          disk_write_bps: this.parsePromResult(rangeObj?.disk_write_bps),
          net_in_bps: this.parsePromResult(rangeObj?.net_in_bps),
          net_out_bps: this.parsePromResult(rangeObj?.net_out_bps),
          process_count: this.parsePromResult(rangeObj?.process_count)
        })

        const ranges = data.ranges || {}
        this.ranges['30m'] = parseRange(ranges['30m'] || {})
        this.ranges['1d'] = parseRange(ranges['1d'] || {})
        this.ranges['7d'] = parseRange(ranges['7d'] || {})

        // 若后端未返回运行时长，基于最早时序点估算(当前时间 - 最早点)
        if (!this.deviceInfo.uptimeSec || this.deviceInfo.uptimeSec <= 0) {
          const earliest = this.getEarliestTimestamp()
          if (earliest) {
            const nowSec = Math.floor(Date.now() / 1000)
            const est = Math.max(0, nowSec - earliest)
            this.deviceInfo.uptimeSec = est
          }
        }

        this.$nextTick(this.renderCharts)
      } catch (error) {
        console.error('加载监控失败:', error)
        this.$message.error(error?.response?.data?.message || error?.message || '获取监控数据失败')
      }
    },
    getEarliestTimestamp() {
      const all = []
      const pushSeries = (obj) => {
        if (!obj) return
        Object.values(obj).forEach(arr => {
          if (Array.isArray(arr)) {
            arr.forEach(p => { if (p && Number.isFinite(p.timestamp)) all.push(Number(p.timestamp)) })
          }
        })
      }
      pushSeries(this.ranges['7d'])
      pushSeries(this.ranges['1d'])
      pushSeries(this.ranges['30m'])
      if (!all.length) return 0
      return Math.min(...all)
    },
    parsePromResult(result) {
      if (!result || result.status !== 'success' || !result.data) return []
      const coerce = (t, v) => ({ timestamp: Number(t), value: Number(v) })
      let points = []
      if (result.data.resultType === 'matrix') {
        const values = result.data.result?.[0]?.values || []
        points = values.map(([t, v]) => coerce(t, v))
      } else if (result.data.resultType === 'vector') {
        const v = result.data.result?.[0]?.value
        if (Array.isArray(v)) {
          const [t, val] = v
          points = [coerce(t, val)]
        }
      }
      // 过滤 NaN/Infinity
      return points.filter(p => Number.isFinite(p.value))
    },
    getSeries(key) {
      // 针对不同指标清洗数据: 百分比夹断到[0,100]，其余非负
      const raw = this.ranges[this.currentRange]?.[key] || []
      const sanitize = (p) => {
        let v = p.value
        if (['cpu','memory','disk_usage'].includes(key)) {
          if (!Number.isFinite(v)) return null
          v = Math.max(0, Math.min(100, v))
        } else {
          if (!Number.isFinite(v)) return null
          v = Math.max(0, v)
        }
        return { timestamp: p.timestamp, value: v }
      }
      return raw.map(sanitize).filter(Boolean)
    },
    latest(key) {
      // 优先当前范围，若无有效点，回退到1d、7d
      const pref = [this.currentRange, '1d', '7d']
      let series = []
      for (const r of pref) {
        series = (this.ranges[r]?.[key] || []).filter(p => Number.isFinite(p.value))
        if (series.length) break
      }
      if (!series.length) return undefined
      // 对CPU做平滑后取末值，降低异常波动影响
      if (key === 'cpu') {
        const smoothed = this.smoothSeries(series, 3)
        return smoothed.length ? smoothed[smoothed.length - 1].value : series[series.length - 1].value
      }
      return series[series.length - 1].value
    },
    smoothSeries(arr, windowSize = 3) {
      if (!arr || arr.length <= 1) return arr || []
      const out = []
      for (let i = 0; i < arr.length; i++) {
        let sum = 0, n = 0
        for (let j = Math.max(0, i - windowSize + 1); j <= i; j++) {
          const v = Number(arr[j].value)
          if (Number.isFinite(v)) { sum += v; n++ }
        }
        const avg = n ? sum / n : arr[i].value
        out.push({ timestamp: arr[i].timestamp, value: avg })
      }
      return out
    },
    formatPercent(v) {
      return v === undefined ? '-' : `${Number(v).toFixed(2)}%`
    },
    formatNumber(v) {
      return v === undefined ? '-' : Number(v).toFixed(0)
    },
    toMbps(bps) { return (Number(bps) * 8) / 1_000_000 },
    toMBps(bps) { return Number(bps) / 1024 / 1024 },
    toMs(sec) { return Number(sec) * 1000 },

    renderCharts() {
      this.renderLine('cpuChart', 'cpu', '%', '#00d8ff')
      this.renderLine('memoryChart', 'memory', '%', '#4fc3f7')
      this.renderLine('diskUsageChart', 'disk_usage', '%', '#7da6ff')
      this.renderLine('processChart', 'process_count', 'count', '#67c23a', { fixed: 0 })

      // 复合图：磁盘延迟 读/写(ms)
      this.renderDualLine('diskLatencyChart',
        this.transform(this.ranges[this.currentRange]?.disk_read_latency, (v) => this.toMs(v)),
        this.transform(this.ranges[this.currentRange]?.disk_write_latency, (v) => this.toMs(v)),
        '读延迟', '写延迟', 'ms/次', ['#e67e22', '#f56c6c'])

      // 复合图：磁盘吞吐(MB/s)
      this.renderDualLine('diskBpsChart',
        this.transform(this.ranges[this.currentRange]?.disk_read_bps, (v) => this.toMBps(v)),
        this.transform(this.ranges[this.currentRange]?.disk_write_bps, (v) => this.toMBps(v)),
        '读吞吐', '写吞吐', 'MB/s', ['#00bcd4', '#9c27b0'])

      // 复合图：网络(Mbps)
      this.renderDualLine('netChart',
        this.transform(this.ranges[this.currentRange]?.net_in_bps, (v) => this.toMbps(v)),
        this.transform(this.ranges[this.currentRange]?.net_out_bps, (v) => this.toMbps(v)),
        '入带宽', '出带宽', 'Mbps', ['#409eff', '#67c23a'])
    },
    transform(arr, mapFn) {
      return (arr || []).map(p => ({ timestamp: p.timestamp, value: mapFn(p.value) }))
    },
    renderLine(refName, key, unit, color, { fixed = 2 } = {}) {
      const el = this.$refs[refName]
      if (!el) return
      const chartKey = this._chartKeyByRef(refName)
      if (!this.chartInstances[chartKey]) this.chartInstances[chartKey] = echarts.init(el, 'dark')
      let data = this.getSeries(key)
      if (key === 'cpu') data = this.smoothSeries(data, 3)
      const xAxisData = data.map(p => this.formatTime(p.timestamp))
      const seriesData = data.map(p => Number(p.value).toFixed(fixed))
      const option = this.getChartOption(xAxisData, [{ name: key, data: seriesData, color }], unit)
      this.chartInstances[chartKey].setOption(option, true)
    },
    renderDualLine(refName, arr1, arr2, name1, name2, unit, colors) {
      const el = this.$refs[refName]
      if (!el) return
      const chartKey = this._chartKeyByRef(refName)
      if (!this.chartInstances[chartKey]) this.chartInstances[chartKey] = echarts.init(el, 'dark')
      const data = arr1?.length ? arr1 : arr2 || []
      const xAxisData = data.map(p => this.formatTime(p.timestamp))
      const s1 = (arr1 || []).map(p => Number(p.value).toFixed(2))
      const s2 = (arr2 || []).map(p => Number(p.value).toFixed(2))
      const option = this.getChartOption(xAxisData, [
        { name: name1, data: s1, color: colors[0] },
        { name: name2, data: s2, color: colors[1] }
      ], unit)
      this.chartInstances[chartKey].setOption(option, true)
    },
    _chartKeyByRef(refName) {
      const map = { cpuChart: 'cpu', memoryChart: 'memory', diskUsageChart: 'diskUsage', processChart: 'process', diskLatencyChart: 'diskLatency', diskBpsChart: 'diskBps', netChart: 'net' }
      return map[refName]
    },
    getChartOption(xAxisData, seriesDefs, unit) {
      return {
        backgroundColor: 'transparent',
        tooltip: {
          trigger: 'axis',
          confine: false,
          backgroundColor: 'rgba(10, 14, 28, 0.95)',
          borderColor: '#1e3a8a',
          borderWidth: 1,
          textStyle: { color: '#e5e7eb' },
          axisPointer: {
            type: 'cross',
            label: {
              backgroundColor: '#1e3a8a'
            }
          },
          formatter: (params) => {
            if (!params || !params.length) return ''
            const lines = [
              `<div style="font-weight:600;color:#93c5fd;margin-bottom:4px;">${params[0].axisValue}</div>`
            ]
            params.forEach(p => {
              lines.push(`<div style="color:${p.color}">• ${p.seriesName}: <b>${p.value}${unit}</b></div>`)
            })
            return lines.join('')
          }
        },
        grid: { left: '3%', right: '3%', bottom: '8%', containLabel: true },
        xAxis: { type: 'category', boundaryGap: false, data: xAxisData, axisLabel: { color: '#9ca3af' }, axisLine: { lineStyle: { color: '#374151' } } },
        yAxis: { type: 'value', name: unit, axisLabel: { color: '#9ca3af' }, axisLine: { lineStyle: { color: '#374151' } }, splitLine: { lineStyle: { color: '#1f2937' } } },
        series: seriesDefs.map(def => ({
          name: def.name,
          type: 'line',
          data: def.data,
          smooth: true,
          showSymbol: false,
          lineStyle: { color: def.color, width: 1 },
          areaStyle: { opacity: 0.12, color: def.color }
        }))
      }
    },
    formatTime(ts) {
      const d = new Date(Number(ts) * 1000)
      const hh = String(d.getHours()).padStart(2, '0')
      const mm = String(d.getMinutes()).padStart(2, '0')
      return `${hh}:${mm}`
    },
    getMetricClass(v) {
      if (v === undefined || v === null) return 'normal'
      if (v > 80) return 'danger'
      if (v > 60) return 'warning'
      return 'normal'
    },
    formatUptime(seconds) {
      if (!seconds) return '-'
      const days = Math.floor(seconds / 86400)
      const hours = Math.floor((seconds % 86400) / 3600)
      const minutes = Math.floor((seconds % 3600) / 60)
      const parts = []
      if (days > 0) parts.push(`${days}天`)
      if (hours > 0) parts.push(`${hours}小时`)
      if (minutes > 0) parts.push(`${minutes}分钟`)
      return parts.length > 0 ? parts.join(' ') : '0分钟'
    },
    startAutoRefresh() {
      this.refreshInterval = setInterval(() => { this.loadMonitorData() }, 30000)
    },
    stopAutoRefresh() {
      if (this.refreshInterval) { clearInterval(this.refreshInterval); this.refreshInterval = null }
    },
    handleResize() {
      Object.values(this.chartInstances).forEach(c => c && c.resize())
    },
    goBack() { this.$router.back() }
  },
  mounted() {
    this.deviceId = this.$route.params.id
    if (!this.deviceId) { this.$message.error('设备ID不存在'); this.$router.back(); return }
    this.loadDeviceInfo()
    this.loadMonitorData()
    this.startAutoRefresh()
    window.addEventListener('resize', this.handleResize)
  },
  beforeUnmount() {
    this.stopAutoRefresh()
    Object.values(this.chartInstances).forEach(c => c && c.dispose && c.dispose())
    window.removeEventListener('resize', this.handleResize)
  }
}
</script>

<style scoped>
.windows-monitor-page { padding: 20px; min-height: 100vh; background: #eef1f6; }
.monitor-card { background: rgba(255, 255, 255, 0.95); border-radius: 16px; box-shadow: 0 8px 32px rgba(0,0,0,0.1); border: 1px solid rgba(255, 255, 255, 0.2); }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.title { font-size: 20px; font-weight: 700; color: #2c3e50; }
.info-section, .metrics-section, .charts-section { margin-bottom: 20px; border-radius: 12px; }
.metrics-row { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px; }
/* 黑底蓝字：仅指标与图表卡片 */
.metric-item { background: rgba(103, 126, 234, 0.05); border: 1px solid rgba(103, 126, 234, 0.1); border-radius: 12px; padding: 12px 16px; }
.metric-content { display: flex; align-items: center; gap: 12px; }
.metric-icon { font-size: 18px; color: #409eff; }
.metric-label { color: #909399; font-size: 12px; }
.metric-value { font-size: 18px; font-weight: 700; color: #303133; }
.metric-value.warning { color: #fbbf24; }
.metric-value.danger { color: #f87171; }
.chart-card { background: #0b1220; border: 1px solid rgba(59,130,246,0.25); border-radius: 10px; padding: 8px; margin-bottom: 16px; }
.chart-title { font-weight: 700; color: #93c5fd; margin: 4px 8px 8px; }
.chart { width: 100%; height: 260px; }
.section-header-with-line { display: flex; align-items: center; gap: 10px; }
.section-title { font-weight: 700; color: #303133; }
.title-line { flex: 1; height: 1px; background: #ebeef5; }
.header-actions { display: flex; align-items: center; }
</style>
