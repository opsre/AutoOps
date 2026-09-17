<template>
  <el-dialog
    title="进程监控"
    v-model:visible="dialogVisible"
    width="80%"
    top="5vh"
    @close="handleClose"
    custom-class="process-monitor-dialog"
    :modal="true"
    :append-to-body="true"
  >
    <div class="process-monitor-header">
      <div class="host-info">
        <span class="host-name">{{ processData.hostName || '未知主机' }}</span>
        <span class="update-time">最后更新: {{ formatTime(processData.updateTime) }}</span>
      </div>
      <el-button 
        type="primary" 
        size="mini" 
        icon="Refresh"
        :loading="loading"
        @click="fetchProcessData"
      >
        手动刷新
      </el-button>
    </div>

    <el-row :gutter="20" class="process-charts">
      <el-col :span="12">
        <el-card class="chart-card">
          <div class="chart-header">
            <h4>Top CPU 进程</h4>
            <el-tag size="mini" type="info">CPU使用率 %</el-tag>
          </div>
          <div class="chart-container" ref="cpuChart"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="chart-card">
          <div class="chart-header">
            <h4>Top 内存 进程</h4>
            <el-tag size="mini" type="warning">内存使用率 %</el-tag>
          </div>
          <div class="chart-container" ref="memoryChart"></div>
        </el-card>
      </el-col>
    </el-row>
  </el-dialog>
</template>

<script>
import * as echarts from 'echarts'

export default {
  name: 'ProcessMonitorDialog',
  model: {
    prop: 'visible',
    event: 'change'
  },
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    hostId: {
      type: [String, Number],
      required: true
    }
  },
  data() {
    return {
      loading: false,
      processData: {},
      cpuChart: null,
      memoryChart: null,
      refreshInterval: null
    }
  },
  computed: {
    dialogVisible: {
      get() {
        return this.visible
      },
      set(val) {
        this.$emit('change', val)
      }
    }
  },
  watch: {
    visible(newVal) {
      if (newVal) {
        this.startMonitoring()
      } else {
        this.stopMonitoring()
      }
    },
    hostId(newVal) {
      // 主机ID变化时重新加载数据
      if (newVal && this.visible) {
        this.startMonitoring()
      }
    }
  },
  methods: {
    async fetchProcessData() {
      if (!this.hostId) return
      
      this.loading = true
      try {
        const { data: res } = await this.$api.getHostTopProcesses(this.hostId)
        if (res.code === 200) {
          this.processData = res.data
          this.$nextTick(() => {
            this.renderCharts()
          })
        }
      } catch (error) {
        console.error('获取进程数据失败:', error)
        this.$message.error('获取进程数据失败')
      } finally {
        this.loading = false
      }
    },
    
    renderCharts() {
      this.renderCpuChart()
      this.renderMemoryChart()
    },
    
    renderCpuChart() {
      if (!this.processData.topCPU || !this.$refs.cpuChart) return
      
      if (!this.cpuChart) {
        this.cpuChart = echarts.init(this.$refs.cpuChart, 'dark')
      }
      
      const data = this.processData.topCPU.slice(0, 5).map(process => ({
        name: process.name,
        value: parseFloat(process.cpuPercent.toFixed(2)),
        pid: process.pid
      }))
      
      const option = {
        backgroundColor: 'transparent',
        tooltip: {
          trigger: 'item',
          backgroundColor: 'rgba(0, 0, 0, 0.9)',
          borderColor: '#00d8ff',
          borderWidth: 1,
          textStyle: {
            color: '#00d8ff',
            fontSize: 12
          },
          confine: true,
          enterable: true,
          hideDelay: 100,
          showDelay: 50,
          formatter: function(params) {
            return `
              <div style="color: #00d8ff; padding: 8px; font-size: 12px;">
                <div style="font-weight: bold; margin-bottom: 5px; font-size: 13px;">${params.data.name}</div>
                <div>PID: ${params.data.pid}</div>
                <div>CPU使用率: ${params.data.value}%</div>
              </div>
            `
          }
        },
        grid: {
          left: '3%',
          right: '4%',
          bottom: '15%',
          top: '10%',
          containLabel: true
        },
        xAxis: {
          type: 'category',
          data: data.map(item => item.name),
          axisLine: {
            lineStyle: {
              color: '#6E7079'
            }
          },
          axisLabel: {
            rotate: 0,
            fontSize: 14,
            color: '#E4E7ED',
            interval: 0,
            formatter: function(value) {
              // 缩短截断长度以避免显示问题
              return value.length > 6 ? value.substring(0, 6) + '...' : value
            }
          },
          axisTick: {
            show: true,
            lineStyle: {
              color: '#6E7079'
            }
          }
        },
        yAxis: {
          type: 'value',
          name: 'CPU %',
          nameTextStyle: {
            color: '#E4E7ED',
            fontSize: 14
          },
          axisLine: {
            lineStyle: {
              color: '#6E7079'
            }
          },
          axisLabel: {
            color: '#E4E7ED',
            fontSize: 13
          },
          splitLine: {
            lineStyle: {
              color: ['#2D3A4B'],
              type: 'dashed'
            }
          }
        },
        series: [{
          type: 'bar',
          data: data,
          barWidth: '40%',
          itemStyle: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [
                { offset: 0, color: '#00d8ff' },
                { offset: 0.5, color: '#3A4DE9' },
                { offset: 1, color: '#1890ff' }
              ]
            },
            borderRadius: [4, 4, 0, 0],
            shadowColor: 'rgba(0, 216, 255, 0.3)',
            shadowBlur: 10,
            shadowOffsetY: 2
          },
          emphasis: {
            itemStyle: {
              color: {
                type: 'linear',
                x: 0,
                y: 0,
                x2: 0,
                y2: 1,
                colorStops: [
                  { offset: 0, color: '#00f5ff' },
                  { offset: 0.5, color: '#4169E1' },
                  { offset: 1, color: '#1E90FF' }
                ]
              },
              shadowColor: 'rgba(0, 245, 255, 0.5)',
              shadowBlur: 15
            }
          },
          animationDuration: 1000,
          animationEasing: 'cubicOut'
        }]
      }
      
      this.cpuChart.setOption(option, true)
    },
    
    renderMemoryChart() {
      if (!this.processData.topMemory || !this.$refs.memoryChart) return
      
      if (!this.memoryChart) {
        this.memoryChart = echarts.init(this.$refs.memoryChart, 'dark')
      }
      
      const data = this.processData.topMemory.slice(0, 5).map(process => ({
        name: process.name,
        value: parseFloat(process.memPercent.toFixed(2)),
        pid: process.pid
      }))
      
      const option = {
        backgroundColor: 'transparent',
        tooltip: {
          trigger: 'item',
          backgroundColor: 'rgba(0, 0, 0, 0.9)',
          borderColor: '#ff9800',
          borderWidth: 1,
          textStyle: {
            color: '#ff9800',
            fontSize: 12
          },
          confine: true,
          enterable: true,
          hideDelay: 100,
          showDelay: 50,
          formatter: function(params) {
            return `
              <div style="color: #ff9800; padding: 8px; font-size: 12px;">
                <div style="font-weight: bold; margin-bottom: 5px; font-size: 13px;">${params.data.name}</div>
                <div>PID: ${params.data.pid}</div>
                <div>内存使用率: ${params.data.value}%</div>
              </div>
            `
          }
        },
        grid: {
          left: '3%',
          right: '4%',
          bottom: '15%',
          top: '10%',
          containLabel: true
        },
        xAxis: {
          type: 'category',
          data: data.map(item => item.name),
          axisLine: {
            lineStyle: {
              color: '#6E7079'
            }
          },
          axisLabel: {
            rotate: 0,
            fontSize: 14,
            color: '#E4E7ED',
            interval: 0,
            formatter: function(value) {
              // 缩短截断长度以避免显示问题
              return value.length > 6 ? value.substring(0, 6) + '...' : value
            }
          },
          axisTick: {
            show: true,
            lineStyle: {
              color: '#6E7079'
            }
          }
        },
        yAxis: {
          type: 'value',
          name: '内存 %',
          nameTextStyle: {
            color: '#E4E7ED',
            fontSize: 14
          },
          axisLine: {
            lineStyle: {
              color: '#6E7079'
            }
          },
          axisLabel: {
            color: '#E4E7ED',
            fontSize: 13
          },
          splitLine: {
            lineStyle: {
              color: ['#2D3A4B'],
              type: 'dashed'
            }
          }
        },
        series: [{
          type: 'bar',
          data: data,
          barWidth: '40%',
          itemStyle: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [
                { offset: 0, color: '#ffb74d' },
                { offset: 0.5, color: '#ff9800' },
                { offset: 1, color: '#f57c00' }
              ]
            },
            borderRadius: [4, 4, 0, 0],
            shadowColor: 'rgba(255, 152, 0, 0.3)',
            shadowBlur: 10,
            shadowOffsetY: 2
          },
          emphasis: {
            itemStyle: {
              color: {
                type: 'linear',
                x: 0,
                y: 0,
                x2: 0,
                y2: 1,
                colorStops: [
                  { offset: 0, color: '#ffc947' },
                  { offset: 0.5, color: '#ffab00' },
                  { offset: 1, color: '#ff6f00' }
                ]
              },
              shadowColor: 'rgba(255, 171, 0, 0.5)',
              shadowBlur: 15
            }
          },
          animationDuration: 1000,
          animationEasing: 'cubicOut'
        }]
      }
      
      this.memoryChart.setOption(option, true)
    },
    
    formatTime(timestamp) {
      if (!timestamp) return '未知时间'
      const date = new Date(timestamp * 1000)
      return date.toLocaleString()
    },
    
    startMonitoring() {
      console.log('=== startMonitoring 开始 ===', {
        hostId: this.hostId,
        visible: this.visible
      })
      
      // 检查必要参数
      if (!this.hostId) {
        console.error('hostId 未设置，无法加载数据')
        return
      }
      
      // 停止之前的刷新
      this.stopMonitoring()
      
      // 立即加载数据
      console.log('开始加载初始数据...')
      this.fetchProcessData().then(() => {
        console.log('初始数据加载完成')
      }).catch(error => {
        console.error('初始数据加载失败:', error)
      })
      
      // 设置定时器，10秒刷新一次
      this.refreshInterval = setInterval(() => {
        console.log('定时刷新进程数据...')
        this.fetchProcessData().catch(error => {
          console.error('定时刷新失败:', error)
        })
      }, 10000)
      
      console.log('定时器已设置，间雔10秒')
      console.log('=== startMonitoring 结束 ===')
    },
    
    stopMonitoring() {
      if (this.refreshInterval) {
        clearInterval(this.refreshInterval)
        this.refreshInterval = null
      }
    },
    
    handleClose() {
      this.stopMonitoring()
      this.dialogVisible = false
    },
    
    handleResize() {
      this.cpuChart && this.cpuChart.resize()
      this.memoryChart && this.memoryChart.resize()
    }
  },
  
  mounted() {
    console.log('ProcessMonitorDialog mounted, visible:', this.visible, 'hostId:', this.hostId)
    window.addEventListener('resize', this.handleResize)
    // 利用v-if组件创建即可见的特性，立即加载数据
    if (this.hostId) {
      console.log('Component mounted with hostId, immediately starting monitoring')
      this.$nextTick(() => {
        setTimeout(() => {
          this.startMonitoring()
        }, 100) // 稍微延迟确保 DOM 完全渲染
      })
    }
  },
  
  beforeUnmount() {
    window.removeEventListener('resize', this.handleResize)
    this.stopMonitoring()
    this.cpuChart && this.cpuChart.dispose()
    this.memoryChart && this.memoryChart.dispose()
  }
}
</script>

<style scoped>
.process-monitor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 15px 20px;
  background: #1e1e1e; /* 改为黑色背景 */
  border-radius: 8px;
  border: 1px solid #2D3A4B;
}

.host-info {
  display: flex;
  flex-direction: column;
}

.host-name {
  font-size: 16px;
  font-weight: bold;
  color: #00d8ff; /* 使用主机监控的蓝色 */
  margin-bottom: 5px;
}

.update-time {
  font-size: 12px;
  color: #6E7079; /* 使用暗色主题颜色 */
}

.process-charts {
  margin-bottom: 20px;
}

.chart-card {
  height: 400px;
  background: #1e1e1e; /* 改为黑色背景 */
  border: 1px solid #2D3A4B;
}

.chart-card .el-card__body {
  background: #1e1e1e; /* 卡片内容区域也是黑色 */
  color: #E4E7ED;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.chart-header h4 {
  margin: 0;
  color: #E4E7ED; /* 使用暗色主题的文字颜色 */
}

.chart-container {
  height: 320px;
  width: 100%;
  background: #1e1e1e; /* 图表容器也是黑色 */
}

/* 覆盖Element UI对话框的默认样式 */
.process-monitor-dialog .el-dialog {
  background-color: #1e1e1e !important;
}

.process-monitor-dialog .el-dialog__header {
  background-color: #1e1e1e !important;
  border-bottom: 1px solid #2D3A4B !important;
}

.process-monitor-dialog .el-dialog__title {
  color: #E4E7ED !important;
}

.process-monitor-dialog .el-dialog__body {
  background-color: #1e1e1e !important;
  color: #E4E7ED !important;
}

/* 覆盖卡片的默认样式 */
.chart-card .el-card {
  background-color: #1e1e1e !important;
  border: 1px solid #2D3A4B !important;
}

.chart-card .el-card__header {
  background-color: #1e1e1e !important;
  border-bottom: 1px solid #2D3A4B !important;
  color: #E4E7ED !important;
}

.chart-card .el-card__body {
  background-color: #1e1e1e !important;
  color: #E4E7ED !important;
}
</style>