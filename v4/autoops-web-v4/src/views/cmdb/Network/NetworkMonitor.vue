<template>
  <div class="network-monitor-page">
    <!-- 动态背景粒子效果 -->
    <div class="bg-animation">
      <div class="particle"></div>
      <div class="particle"></div>
      <div class="particle"></div>
      <div class="particle"></div>
      <div class="particle"></div>
    </div>

    <el-card shadow="hover" class="monitor-card">
      <template #header>
        <div class="card-header">
     内存使用率(%)
     <div class="title-wrapper">
            <span class="title">交换机监控</span>
          </div>
          <el-button size="small" @click="goBack" class="back-btn">
            <el-icon><ArrowLeft /></el-icon>
            返回
          </el-button>
        </div>
      </template>

      <!-- 设备基础信息 -->
      <el-card class="info-section animate-slide-in" shadow="never">
        <template #header>
          <div class="section-header-with-line">
            <span class="section-title">设备信息</span>
            <div class="title-line"></div>
          </div>
        </template>
        <el-descriptions :column="3" border>
          <el-descriptions-item label="设备名称">{{ monitorData.deviceName || '-' }}</el-descriptions-item>
          <el-descriptions-item label="IP地址">{{ monitorData.ipAddress || '-' }}</el-descriptions-item>
          <el-descriptions-item label="设备类型">{{ monitorData.deviceType || '-' }}</el-descriptions-item>
          <el-descriptions-item label="位置">{{ monitorData.location || '-' }}</el-descriptions-item>
          <el-descriptions-item label="负责人">{{ monitorData.owner || '-' }}</el-descriptions-item>
          <el-descriptions-item label="备注">{{ monitorData.remark || '-' }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 实时监控指标 -->
      <el-card class="metrics-section animate-slide-in" shadow="never" style="animation-delay: 0.1s">
        <template #header>
          <div class="section-header-with-line">
            <span class="section-title">实时监控指标</span>
            <div class="title-line"></div>
          </div>
        </template>
        <el-row :gutter="15" class="metrics-row">
          <el-col class="metric-col">
            <div class="metric-item pulse-animation">
              <div class="metric-content">
                <div class="metric-icon cpu-icon">
                  <el-icon><Cpu /></el-icon>
                </div>
                <div class="metric-info">
                  <div class="metric-label">CPU使用率</div>
                  <div class="metric-value" :class="getMetricClass(monitorData.cpuUsage)">
                    {{ monitorData.cpuUsage !== undefined ? monitorData.cpuUsage.toFixed(2) + '%' : '-' }}
                  </div>
                </div>
              </div>
            </div>
          </el-col>
          <el-col class="metric-col">
            <div class="metric-item pulse-animation" style="animation-delay: 0.1s">
              <div class="metric-content">
                <div class="metric-icon memory-icon">
                  <el-icon><MemoryCard /></el-icon>
                </div>
                <div class="metric-info">
                  <div class="metric-label">内存使用率</div>
                  <div class="metric-value" :class="getMetricClass(monitorData.memoryUsage)">
                    {{ monitorData.memoryUsage !== undefined ? monitorData.memoryUsage.toFixed(2) + '%' : '-' }}
                  </div>
                </div>
              </div>
            </div>
          </el-col>
          <el-col class="metric-col">
            <div class="metric-item">
              <div class="metric-content">
                <div class="metric-icon temp-icon">
                  <el-icon><Odometer /></el-icon>
                </div>
                <div class="metric-info">
                  <div class="metric-label">温度</div>
                  <div class="metric-value" :class="getTempClass(monitorData.temperature)">
                    {{ monitorData.temperature !== undefined ? monitorData.temperature.toFixed(1) + '°C' : '-' }}
                  </div>
                </div>
              </div>
            </div>
          </el-col>
          <el-col class="metric-col">
            <div class="metric-item">
              <div class="metric-content">
                <div class="metric-icon uptime-icon">
                  <el-icon><Clock /></el-icon>
                </div>
                <div class="metric-info">
                  <div class="metric-label">运行时间</div>
                  <div class="metric-value normal">
                    {{ formatUptime(monitorData.uptime) }}
                  </div>
                </div>
              </div>
            </div>
          </el-col>
          <el-col class="metric-col">
            <div class="metric-item port-stats-item">
              <div class="port-stats-content">
                <div class="metric-icon port-stats-icon port-stats-icon-large">
                  <el-icon><Connection /></el-icon>
                </div>
                <div class="port-stats-info">
                  <div class="port-stat-row">
                    <div class="metric-label">总端口</div>
                    <div class="metric-value" style="color: #409eff;">
                      {{ monitorData.totalPorts || 0 }}
                    </div>
                  </div>
                  <div class="port-stat-row">
                    <div class="metric-label">空闲端口</div>
                    <div class="metric-value" style="color: #67c23a;">
                      {{ monitorData.idlePorts || 0 }}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </el-col>
        </el-row>
      </el-card>

      <!-- 端口流量Top5 -->
      <el-row :gutter="20" class="traffic-section">
        <el-col :span="12">
          <el-card shadow="never" class="traffic-card">
            <template #header>
              <div class="section-header-with-line">
                <span class="section-title">入站流量 Top5</span>
                <div class="title-line"></div>
              </div>
            </template>
            <div class="traffic-list">
              <div
                v-for="(item, index) in monitorData.inboundTop5 || []"
                :key="index"
                class="traffic-item"
                :class="`traffic-item-${index + 1}`">
                <div class="traffic-rank">
                  <span class="rank-number">{{ index + 1 }}</span>
                </div>
                <div class="traffic-port">
                  <i class="port-icon"></i>
                  <span class="port-name">{{ item.portName || '-' }}</span>
                </div>
                <div class="traffic-value">
                  <span class="value-number">{{ item.trafficMbps !== undefined ? item.trafficMbps.toFixed(2) : '-' }}</span>
                  <span class="value-unit">Mbps</span>
                </div>
              </div>
              <div v-if="!monitorData.inboundTop5 || monitorData.inboundTop5.length === 0" class="no-data">
                暂无数据
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card shadow="never" class="traffic-card">
            <template #header>
              <div class="section-header-with-line">
                <span class="section-title">出站流量 Top5</span>
                <div class="title-line"></div>
              </div>
            </template>
            <div class="traffic-list">
              <div
                v-for="(item, index) in monitorData.outboundTop5 || []"
                :key="index"
                class="traffic-item"
                :class="`traffic-item-${index + 1}`">
                <div class="traffic-rank">
                  <span class="rank-number">{{ index + 1 }}</span>
                </div>
                <div class="traffic-port">
                  <i class="port-icon"></i>
                  <span class="port-name">{{ item.portName || '-' }}</span>
                </div>
                <div class="traffic-value">
                  <span class="value-number">{{ item.trafficMbps !== undefined ? item.trafficMbps.toFixed(2) : '-' }}</span>
                  <span class="value-unit">Mbps</span>
                </div>
              </div>
              <div v-if="!monitorData.outboundTop5 || monitorData.outboundTop5.length === 0" class="no-data">
                暂无数据
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 24小时趋势图表 -->
      <el-card class="charts-section" shadow="never">
        <template #header>
          <span class="section-title">24小时趋势</span>
        </template>
        <el-row :gutter="20">
          <el-col :span="8">
            <div class="chart-card">
              <div class="chart-card-header">
                <span class="chart-card-title">CPU使用率趋势</span>
              </div>
              <div ref="cpuChart" class="chart"></div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="chart-card">
              <div class="chart-card-header">
                <span class="chart-card-title">内存使用率趋势</span>
              </div>
              <div ref="memoryChart" class="chart"></div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="chart-card">
              <div class="chart-card-header">
                <span class="chart-card-title">温度趋势</span>
              </div>
              <div ref="tempChart" class="chart"></div>
            </div>
          </el-col>
        </el-row>
      </el-card>
    </el-card>
  </div>
</template>

<script>
import * as echarts from 'echarts'
import API from '@/api/cmdb'
import { ArrowLeft, Refresh, Cpu, MemoryCard, Odometer, Clock, Connection } from '@element-plus/icons-vue'

export default {
  name: 'NetworkMonitor',
  components: {
    ArrowLeft,
    Refresh,
    Cpu,
    MemoryCard,
    Odometer,
    Clock,
    Connection
  },
  data() {
    return {
      deviceId: null,
      loading: false,
      monitorData: {
        deviceName: '',
        ipAddress: '',
        deviceType: '',
        location: '',
        owner: '',
        remark: '',
        cpuUsage: 0,
        memoryUsage: 0,
        temperature: 0,
        uptime: 0,
        inboundTop5: [],
        outboundTop5: [],
        idlePorts: 0,
        totalPorts: 0,
        cpuTrend: [],
        memoryTrend: [],
        temperatureTrend: []
      },
      chartInstances: {
        cpu: null,
        memory: null,
        temp: null
      },
      refreshInterval: null
    }
  },
  methods: {
    async loadDeviceInfo() {
      try {
        const { data: res } = await API.getSnmpDeviceById(this.deviceId)

        if (res.code === 200) {
          // 更新设备基础信息
          this.monitorData.deviceName = res.data.deviceName || ''
          this.monitorData.ipAddress = res.data.ipAddress || ''
          this.monitorData.deviceType = res.data.deviceType || ''
          this.monitorData.location = res.data.location || ''
          this.monitorData.owner = res.data.owner || ''
          this.monitorData.remark = res.data.remark || ''
        } else {
          this.$message.error(res.message || '获取设备信息失败')
        }
      } catch (error) {
        console.error('获取设备信息失败:', error)
        this.$message.error(error?.response?.data?.message || error?.message || '获取设备信息失败')
      }
    },

    async loadMonitorData() {
      try {
        const { data: res } = await API.getSnmpDeviceMonitorData(this.deviceId)

        if (res.code === 200) {
          // 映射后端字段到前端字段
          const apiData = res.data

          this.monitorData = {
            ...this.monitorData,
            cpuUsage: apiData.cpuUsage || 0,
            memoryUsage: apiData.memoryUsage || 0,
            temperature: apiData.temperature || 0,
            uptime: apiData.uptimeDays ? Math.floor(apiData.uptimeDays) * 86400 : 0, // 转换天数为秒(取整)
            inboundTop5: (apiData.portInTop5 || []).map(item => ({
              portName: item.portName,
              trafficMbps: item.trafficRateMbps || 0
            })),
            outboundTop5: (apiData.portOutTop5 || []).map(item => ({
              portName: item.portName,
              trafficMbps: item.trafficRateMbps || 0
            })),
            idlePorts: apiData.idlePortCount || 0,
            totalPorts: apiData.totalPortCount || 0,
            cpuTrend: (apiData.cpuTrend || []).map(item => ({
              timestamp: new Date(item.timestamp).getTime() / 1000, // 转换为Unix时间戳(秒)
              value: item.value
            })),
            memoryTrend: (apiData.memoryTrend || []).map(item => ({
              timestamp: new Date(item.timestamp).getTime() / 1000,
              value: item.value
            })),
            temperatureTrend: (apiData.temperatureTrend || []).map(item => ({
              timestamp: new Date(item.timestamp).getTime() / 1000,
              value: item.value
            }))
          }

          // 渲染图表
          this.$nextTick(() => {
            this.renderCharts()
          })
        } else {
          this.$message.error(res.message || '获取监控数据失败')
        }
      } catch (error) {
        console.error('获取监控数据失败:', error)
        this.$message.error(error?.response?.data?.message || error?.message || '获取监控数据失败')
      }
    },

    renderCharts() {
      this.renderCpuChart()
      this.renderMemoryChart()
      this.renderTempChart()
    },

    renderCpuChart() {
      if (!this.$refs.cpuChart) return

      if (!this.chartInstances.cpu) {
        this.chartInstances.cpu = echarts.init(this.$refs.cpuChart, 'dark')
      }

      const data = this.monitorData.cpuTrend || []
      const xAxisData = data.map(item => {
        const date = new Date(item.timestamp * 1000)
        return `${date.getHours()}:${date.getMinutes().toString().padStart(2, '0')}`
      })
      const seriesData = data.map(item => item.value.toFixed(2))

      const option = this.getChartOption('CPU使用率', xAxisData, seriesData, '%', '#00d8ff')
      this.chartInstances.cpu.setOption(option, true)
    },

    renderMemoryChart() {
      if (!this.$refs.memoryChart) return

      if (!this.chartInstances.memory) {
        this.chartInstances.memory = echarts.init(this.$refs.memoryChart, 'dark')
      }

      const data = this.monitorData.memoryTrend || []
      const xAxisData = data.map(item => {
        const date = new Date(item.timestamp * 1000)
        return `${date.getHours()}:${date.getMinutes().toString().padStart(2, '0')}`
      })
      const seriesData = data.map(item => item.value.toFixed(2))

      const option = this.getChartOption('内存使用率', xAxisData, seriesData, '%', '#ff9800')
      this.chartInstances.memory.setOption(option, true)
    },

    renderTempChart() {
      if (!this.$refs.tempChart) return

      if (!this.chartInstances.temp) {
        this.chartInstances.temp = echarts.init(this.$refs.tempChart, 'dark')
      }

      const data = this.monitorData.temperatureTrend || []
      const xAxisData = data.map(item => {
        const date = new Date(item.timestamp * 1000)
        return `${date.getHours()}:${date.getMinutes().toString().padStart(2, '0')}`
      })
      const seriesData = data.map(item => item.value.toFixed(1))

      const option = this.getChartOption('温度', xAxisData, seriesData, '°C', '#f56c6c')
      this.chartInstances.temp.setOption(option, true)
    },

    getChartOption(title, xAxisData, seriesData, unit, color) {
      return {
        backgroundColor: 'transparent',
        tooltip: {
          show: true,
          trigger: 'axis',
          backgroundColor: 'rgba(0, 0, 0, 0.9)',
          borderColor: color,
          borderWidth: 1,
          confine: true,
          textStyle: {
            color: '#fff',
            fontSize: 12
          },
          axisPointer: {
            type: 'cross',
            crossStyle: {
              color: color,
              width: 1,
              type: 'dashed'
            },
            lineStyle: {
              color: color,
              width: 1,
              opacity: 0.5
            }
          },
          formatter: function(params) {
            if (params && params.length > 0) {
              return `<div style="color: #fff; font-weight: bold; margin-bottom: 5px;">${params[0].axisValue}</div>
                      <div style="color: ${color};">• ${params[0].seriesName}: <span style="font-weight: bold;">${params[0].value}${unit}</span></div>`
            }
            return ''
          }
        },
        grid: {
          left: '3%',
          right: '4%',
          bottom: '3%',
          top: '15%',
          containLabel: true
        },
        xAxis: {
          type: 'category',
          boundaryGap: false,
          data: xAxisData,
          axisLine: {
            lineStyle: {
              color: '#6E7079'
            }
          },
          axisLabel: {
            color: '#E4E7ED'
          }
        },
        yAxis: {
          type: 'value',
          name: unit,
          nameTextStyle: {
            color: '#E4E7ED'
          },
          axisLine: {
            lineStyle: {
              color: '#6E7079'
            }
          },
          axisLabel: {
            color: '#E4E7ED'
          },
          splitLine: {
            lineStyle: {
              color: ['#2D3A4B']
            }
          }
        },
        series: [{
          name: title,
          type: 'line',
          data: seriesData,
          smooth: true,
          symbol: 'circle',
          symbolSize: 6,
          showSymbol: false,
          lineStyle: {
            color: color,
            width: 1
          },
          itemStyle: {
            color: color,
            borderColor: '#fff',
            borderWidth: 2
          },
          emphasis: {
            focus: 'series',
            itemStyle: {
              borderWidth: 2,
              shadowBlur: 10,
              shadowColor: color
            }
          },
          areaStyle: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [
                {
                  offset: 0,
                  color: color === '#00d8ff' ? 'rgba(0, 216, 255, 0.3)' :
                         color === '#ff9800' ? 'rgba(255, 152, 0, 0.3)' :
                         'rgba(245, 108, 108, 0.3)'
                },
                {
                  offset: 1,
                  color: color === '#00d8ff' ? 'rgba(0, 216, 255, 0)' :
                         color === '#ff9800' ? 'rgba(255, 152, 0, 0)' :
                         'rgba(245, 108, 108, 0)'
                }
              ]
            }
          }
        }]
      }
    },

    getMetricClass(value) {
      if (value === undefined || value === null) return 'normal'
      if (value > 80) return 'danger'
      if (value > 60) return 'warning'
      return 'normal'
    },

    getTempClass(value) {
      if (value === undefined || value === null) return 'normal'
      if (value > 70) return 'danger'
      if (value > 50) return 'warning'
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

    refreshData() {
      this.loadMonitorData()
      this.$message.success('刷新成功')
    },

    goBack() {
      this.$router.back()
    },

    startAutoRefresh() {
      this.refreshInterval = setInterval(() => {
        this.loadMonitorData()
      }, 30000) // 30秒刷新一次
    },

    stopAutoRefresh() {
      if (this.refreshInterval) {
        clearInterval(this.refreshInterval)
        this.refreshInterval = null
      }
    },

    handleResize() {
      Object.values(this.chartInstances).forEach(chart => {
        chart && chart.resize()
      })
    }
  },
  mounted() {
    this.deviceId = this.$route.params.id

    if (!this.deviceId) {
      this.$message.error('设备ID不存在')
      this.$router.back()
      return
    }

    // 先加载设备基础信息
    this.loadDeviceInfo()
    // 再加载监控数据
    this.loadMonitorData()
    // 启动自动刷新
    this.startAutoRefresh()

    window.addEventListener('resize', this.handleResize)
  },
  beforeUnmount() {
    this.stopAutoRefresh()

    Object.values(this.chartInstances).forEach(chart => {
      if (chart) {
        chart.dispose()
      }
    })

    window.removeEventListener('resize', this.handleResize)
  }
}
</script>

<style scoped>
.network-monitor-page {
  padding: 20px;
  min-height: 100vh;
  background: #eef1f6;
}

.monitor-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  font-size: 20px;
  font-weight: 600;
  color: #2c3e50;
  background: linear-gradient(45deg, #667eea, #764ba2);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.info-section,
.metrics-section,
.port-stats-section,
.charts-section {
  margin-bottom: 20px;
  border-radius: 12px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
}

/* 5列平均分布布局 */
.metrics-row {
  display: flex;
  flex-wrap: nowrap;
}

.metric-col {
  flex: 1;
  min-width: 0;
  padding-left: 7.5px;
  padding-right: 7.5px;
}

.metric-col:first-child {
  padding-left: 0;
}

.metric-col:last-child {
  padding-right: 0;
}

.metric-item {
  padding: 12px 16px;
  background: rgba(103, 126, 234, 0.05);
  border-radius: 12px;
  border: 1px solid rgba(103, 126, 234, 0.1);
  position: relative;
  height: 80px;
  display: flex;
  align-items: center;
}

.metric-content {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
}

.metric-info {
  flex: 1;
}

.metric-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 12px;
  font-size: 20px;
  transition: all 0.3s ease;
  flex-shrink: 0;
}

.cpu-icon {
  background: linear-gradient(135deg, rgba(0, 216, 255, 0.15), rgba(0, 216, 255, 0.05));
  color: #00d8ff;
  box-shadow: 0 4px 12px rgba(0, 216, 255, 0.2);
}

.memory-icon {
  background: linear-gradient(135deg, rgba(255, 152, 0, 0.15), rgba(255, 152, 0, 0.05));
  color: #ff9800;
  box-shadow: 0 4px 12px rgba(255, 152, 0, 0.2);
}

.temp-icon {
  background: linear-gradient(135deg, rgba(245, 108, 108, 0.15), rgba(245, 108, 108, 0.05));
  color: #f56c6c;
  box-shadow: 0 4px 12px rgba(245, 108, 108, 0.2);
}

.uptime-icon {
  background: linear-gradient(135deg, rgba(103, 194, 58, 0.15), rgba(103, 194, 58, 0.05));
  color: #67c23a;
  box-shadow: 0 4px 12px rgba(103, 194, 58, 0.2);
}

.idle-port-icon {
  background: linear-gradient(135deg, rgba(103, 194, 58, 0.15), rgba(103, 194, 58, 0.05));
  color: #67c23a;
  box-shadow: 0 4px 12px rgba(103, 194, 58, 0.2);
}

.total-port-icon {
  background: linear-gradient(135deg, rgba(64, 158, 255, 0.15), rgba(64, 158, 255, 0.05));
  color: #409eff;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.2);
}

.metric-item:hover .metric-icon {
  transform: scale(1.1);
  box-shadow: 0 6px 16px rgba(103, 126, 234, 0.3);
}

.pulse-animation .metric-icon {
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.05);
  }
}

.metric-label {
  font-size: 12px;
  color: #909399;
  margin-bottom: 4px;
}

.metric-value {
  font-size: 20px;
  font-weight: 700;
}

.metric-value.normal {
  color: #67c23a;
}

.metric-value.warning {
  color: #e6a23c;
}

.metric-value.danger {
  color: #f56c6c;
}

/* 端口统计卡片样式 */
.port-stats-item {
  background: linear-gradient(135deg, rgba(103, 126, 234, 0.08), rgba(103, 126, 234, 0.02));
  border: 1px solid rgba(103, 126, 234, 0.15);
  padding: 12px 16px;
  height: 80px;
}

.port-stats-content {
  display: flex;
  align-items: center;
  gap: 12px;
  height: 100%;
}

.port-stats-icon {
  background: linear-gradient(135deg, rgba(103, 194, 58, 0.15), rgba(103, 194, 58, 0.05));
  color: #67c23a;
  box-shadow: 0 4px 12px rgba(103, 194, 58, 0.2);
  width: 40px;
  height: 40px;
  font-size: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  transition: all 0.3s ease;
  flex-shrink: 0;
}

.port-stats-icon-large {
  width: 56px !important;
  height: 56px !important;
  font-size: 28px !important;
}

.port-stats-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.port-stat-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.port-stat-row .metric-label {
  margin-bottom: 0;
  font-size: 12px;
}

.port-stat-row .metric-value {
  font-size: 18px;
}

.traffic-section {
  margin-bottom: 20px;
}

.traffic-card {
  border-radius: 12px;
  overflow: hidden;
}

.traffic-list {
  padding: 0;
}

.traffic-item {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  margin-bottom: 6px;
  border-radius: 6px;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.traffic-item::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 2px;
  transition: width 0.3s ease;
}

.traffic-item:hover {
  transform: translateX(3px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.traffic-item:hover::before {
  width: 3px;
}

/* 渐变色彩方案 - 每行不同颜色,从深到浅 */
.traffic-item-1 {
  background: linear-gradient(90deg, rgba(255, 77, 79, 0.18), rgba(255, 77, 79, 0.06));
  border-left: 2px solid #ff4d4f;
}

.traffic-item-1::before {
  background: #ff4d4f;
}

.traffic-item-2 {
  background: linear-gradient(90deg, rgba(255, 152, 0, 0.14), rgba(255, 152, 0, 0.05));
  border-left: 2px solid #ff9800;
}

.traffic-item-2::before {
  background: #ff9800;
}

.traffic-item-3 {
  background: linear-gradient(90deg, rgba(103, 194, 58, 0.12), rgba(103, 194, 58, 0.04));
  border-left: 2px solid #67c23a;
}

.traffic-item-3::before {
  background: #67c23a;
}

.traffic-item-4 {
  background: linear-gradient(90deg, rgba(64, 158, 255, 0.10), rgba(64, 158, 255, 0.03));
  border-left: 2px solid #409eff;
}

.traffic-item-4::before {
  background: #409eff;
}

.traffic-item-5 {
  background: linear-gradient(90deg, rgba(144, 147, 153, 0.08), rgba(144, 147, 153, 0.02));
  border-left: 2px solid #909399;
}

.traffic-item-5::before {
  background: #909399;
}

.traffic-rank {
  min-width: 28px;
  margin-right: 10px;
}

.rank-number {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: #fff;
  font-weight: 700;
  font-size: 12px;
  box-shadow: 0 2px 6px rgba(103, 126, 234, 0.25);
}

.traffic-port {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 8px;
}

.port-icon {
  display: inline-block;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea, #764ba2);
  box-shadow: 0 0 6px rgba(103, 126, 234, 0.4);
}

.port-name {
  font-size: 13px;
  font-weight: 600;
  color: #2c3e50;
  font-family: 'Monaco', 'Courier New', monospace;
}

.traffic-value {
  display: flex;
  align-items: baseline;
  gap: 4px;
  padding: 4px 10px;
  background: rgba(103, 126, 234, 0.06);
  border-radius: 4px;
  border: 1px solid rgba(103, 126, 234, 0.12);
}

.value-number {
  font-size: 15px;
  font-weight: 700;
  color: #409eff;
  font-family: 'Monaco', 'Courier New', monospace;
}

.value-unit {
  font-size: 11px;
  font-weight: 500;
  color: #909399;
}

.no-data {
  text-align: center;
  padding: 30px 0;
  color: #c0c4cc;
  font-size: 13px;
}

.stat-item {
  padding: 15px;
  background: rgba(103, 126, 234, 0.05);
  border-radius: 8px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stat-label {
  font-size: 14px;
  color: #606266;
  font-weight: 500;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #409eff;
}

/* 图表卡片 - 黑色背景蓝色线条高级感设计 */
.chart-card {
  background-color: #1a1a2e;
  border: 1px solid #00d8ff;
  border-radius: 12px;
  padding: 0;
  margin-bottom: 10px;
  box-shadow: 0 4px 16px rgba(0, 216, 255, 0.15);
  transition: all 0.3s ease;
}

.chart-card:hover {
  box-shadow: 0 6px 24px rgba(0, 216, 255, 0.25);
  transform: translateY(-2px);
}

.chart-card-header {
  padding: 15px 20px;
  border-bottom: 1px solid rgba(0, 216, 255, 0.3);
  background: linear-gradient(135deg, rgba(0, 216, 255, 0.05), rgba(0, 216, 255, 0.02));
}

.chart-card-title {
  font-size: 14px;
  font-weight: 700;
  color: #00d8ff;
  text-shadow: 0 0 10px rgba(0, 216, 255, 0.5);
}

.chart {
  width: 100%;
  height: 280px;
  padding: 10px;
}

/* 响应式布局 */
@media (max-width: 768px) {
  .network-monitor-page {
    padding: 15px;
  }

  .title {
    font-size: 18px;
  }

  .metric-value {
    font-size: 24px;
  }

  .chart {
    height: 250px;
  }
}
</style>
