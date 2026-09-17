<template>
    <el-dialog
    title="主机监控"
    v-model:visible="dialogVisible"
    width="80%"
    top="5vh"
    @close="handleClose"
    custom-class="monitor-dialog"
    :modal="true"
    :append-to-body="true"
    :lock-scroll="false"
    style="z-index: 2001"
  >
    <div class="time-selector">
      <el-radio-group v-model="timeRangeType" size="small" @change="handleTimeRangeChange">
        <el-radio-button label="preset">预设时间</el-radio-button>
        <el-radio-button label="custom">自定义时间</el-radio-button>
      </el-radio-group>

      <div v-if="timeRangeType === 'preset'" class="preset-buttons">
        <el-button 
          v-for="duration in presetDurations" 
          :key="duration.value"
          size="mini"
          :type="activeDuration === duration.value ? 'primary' : ''"
          @click="selectPresetDuration(duration.value)"
        >
          {{ duration.label }}
        </el-button>
      </div>

      <div v-else class="custom-time">
        <el-date-picker
          v-model="customTimeRange"
          type="datetimerange"
          range-separator="至"
          start-placeholder="开始时间"
          end-placeholder="结束时间"
          :value-format="null"
          @input="handleRawDateChange"
        />
        <el-button 
          size="mini" 
          type="primary" 
          style="margin-left: 10px"
          @click="handleQueryClick"
          :disabled="!customTimeRange || customTimeRange.length !== 2"
        >
          查询
        </el-button>
      </div>
    </div>

    <el-row :gutter="20" class="monitor-row">
      <!-- 普通指标卡片 -->
      <el-col :span="8" v-for="(metric, index) in nonNetworkIOMetrics" :key="metric.key">
        <el-card class="metric-card">
          <div class="metric-header">
            <span class="metric-title">{{ metric.title }}</span>
            <el-tag size="mini" :type="getMetricStatus(metric.value)">
              {{ metric.value }} {{ metric.unit }}
            </el-tag>
          </div>
          <div class="metric-chart" :ref="'chart-'+metric.key"></div>
        </el-card>
      </el-col>
      
      <!-- 磁盘I/O合并卡片 -->
      <el-col :span="8">
        <el-card class="metric-card">
          <div class="metric-header">
            <span class="metric-title">磁盘I/O速率</span>
            <div class="disk-io-tags">
              <el-tag size="mini" type="success" style="margin-right: 5px;">
                读: {{ diskReadMetric.value }} {{ diskReadMetric.unit }}
              </el-tag>
              <el-tag size="mini" type="warning">
                写: {{ diskWriteMetric.value }} {{ diskWriteMetric.unit }}
              </el-tag>
            </div>
          </div>
          <div class="metric-chart" ref="chart-diskIO"></div>
        </el-card>
      </el-col>
      
      <!-- 网络I/O合并卡片 -->
      <el-col :span="8">
        <el-card class="metric-card">
          <div class="metric-header">
            <span class="metric-title">网络I/O速率</span>
            <div class="network-io-tags">
              <el-tag size="mini" type="info" style="margin-right: 5px;">
                接收: {{ networkReceiveMetric.value }} {{ networkReceiveMetric.unit }}
              </el-tag>
              <el-tag size="mini" type="primary">
                发送: {{ networkSendMetric.value }} {{ networkSendMetric.unit }}
              </el-tag>
            </div>
          </div>
          <div class="metric-chart" ref="chart-networkIO"></div>
        </el-card>
      </el-col>
    </el-row>
  </el-dialog>
</template>

<script>
import * as echarts from 'echarts';
import moment from 'moment';

export default {
  name: 'MonitorDialog',
  components: {},
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
      type: String,
      required: true
    }
  },
  data() {
    return {
      timeRangeType: 'preset',
      activeDuration: '30m',
      presetDurations: [
        { label: '30分钟', value: '30m' },
        { label: '1小时', value: '1h' },
        { label: '3小时', value: '3h' },
        { label: '6小时', value: '6h' },
        { label: '12小时', value: '12h' },
        { label: '24小时', value: '24h' }
      ],
      customTimeRange: [],
      chartInstances: [],
      metrics: [
        { 
          title: 'CPU使用率', 
          key: 'cpu', 
          unit: '%',
          value: 0,
          chartData: {
            xAxis: [],
            series: []
          }
        },
        { 
          title: '内存使用率', 
          key: 'memory', 
          unit: '%',
          value: 0,
          chartData: {
            xAxis: [],
            series: []
          }
        },
        { 
          title: '磁盘使用率', 
          key: 'disk', 
          unit: '%',
          value: 0,
          chartData: {
            xAxis: [],
            series: []
          }
        },
        { 
          title: '磁盘读取速率', 
          key: 'diskReadKB', 
          unit: 'KB/s',
          value: 0,
          chartData: {
            xAxis: [],
            series: []
          }
        },
        { 
          title: '磁盘写入速率', 
          key: 'diskWriteKB', 
          unit: 'KB/s',
          value: 0,
          chartData: {
            xAxis: [],
            series: []
          }
        },
        { 
          title: '网络接收', 
          key: 'networkReceive', 
          unit: 'KB/s',
          value: 0,
          chartData: {
            xAxis: [],
            series: []
          }
        },
        { 
          title: '网络发送', 
          key: 'networkSend', 
          unit: 'KB/s',
          value: 0,
          chartData: {
            xAxis: [],
            series: []
          }
        },
        { 
          title: '1分钟负载', 
          key: 'load1min', 
          unit: '',
          value: 0,
          chartData: {
            xAxis: [],
            series: []
          }
        },
        { 
          title: '5分钟负载', 
          key: 'load5min', 
          unit: '',
          value: 0,
          chartData: {
            xAxis: [],
            series: []
          }
        },
        { 
          title: '15分钟负载', 
          key: 'load15min', 
          unit: '',
          value: 0,
          chartData: {
            xAxis: [],
            series: []
          }
        },
        { 
          title: '进程总数', 
          key: 'totalProcesses', 
          unit: '',
          value: 0,
          chartData: {
            xAxis: [],
            series: []
          }
        }
      ],
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
    },
    // 非磁盘I/O和网络I/O指标（排除磁盘读写和网络收发速率）
    nonNetworkIOMetrics() {
      return this.metrics.filter(metric => 
        metric.key !== 'diskReadKB' && 
        metric.key !== 'diskWriteKB' && 
        metric.key !== 'networkReceive' && 
        metric.key !== 'networkSend'
      )
    },
    // 磁盘读取指标
    diskReadMetric() {
      return this.metrics.find(metric => metric.key === 'diskReadKB') || { value: 0, unit: 'KB/s' }
    },
    // 磁盘写入指标
    diskWriteMetric() {
      return this.metrics.find(metric => metric.key === 'diskWriteKB') || { value: 0, unit: 'KB/s' }
    },
    // 网络接收指标
    networkReceiveMetric() {
      return this.metrics.find(metric => metric.key === 'networkReceive') || { value: 0, unit: 'KB/s' }
    },
    // 网络发送指标
    networkSendMetric() {
      return this.metrics.find(metric => metric.key === 'networkSend') || { value: 0, unit: 'KB/s' }
    }
  },
  methods: {
    getMetricStatus(value) {
      if (value > 80) return 'danger'
      if (value > 60) return 'warning'
      return 'success'
    },
    async fetchMetrics() {
      try {
        console.log('开始获取监控数据...')
        const params = {
          id: this.hostId
        }
        
        console.log('初始请求参数:', JSON.stringify(params, null, 2))
        
        if (this.timeRangeType === 'preset') {
          console.log('使用预设时间范围:', this.activeDuration)
          params.duration = this.activeDuration
        } else if (this.customTimeRange && this.customTimeRange.length === 2) {
          console.log('使用自定义时间范围:', this.customTimeRange)
          console.log('原始时间值:', {
            start: this.customTimeRange[0],
            end: this.customTimeRange[1]
          })
          
          // 严格格式化时间 - 直接返回加号分隔格式
          const formatTime = (timeValue) => {
            if (!timeValue) return timeValue
            
            // 确保我们处理的是Date对象
            const date = timeValue instanceof Date ? timeValue : new Date(timeValue)
            
            if (isNaN(date.getTime())) {
              console.error('无效的日期:', timeValue)
              return timeValue
            }
            
            const year = date.getFullYear()
            const month = String(date.getMonth() + 1).padStart(2, '0')
            const day = String(date.getDate()).padStart(2, '0')
            const hours = String(date.getHours()).padStart(2, '0')
            const minutes = String(date.getMinutes()).padStart(2, '0')
            const seconds = String(date.getSeconds()).padStart(2, '0')
            
            return `${year}-${month}-${day}+${hours}:${minutes}:${seconds}`
          }
          
          // 返回未编码的加号格式，让请求库处理编码
          const encodeOnce = (str) => {
            return str // 返回原始字符串，不进行任何编码
          }
          
          // 直接从el-date-picker获取Date对象
          const startDate = this.customTimeRange[0] instanceof Date ? 
            this.customTimeRange[0] : new Date(this.customTimeRange[0])
          const endDate = this.customTimeRange[1] instanceof Date ? 
            this.customTimeRange[1] : new Date(this.customTimeRange[1])
            
          params.start = encodeOnce(formatTime(startDate))
          params.end = encodeOnce(formatTime(endDate))
          
          console.log('最终请求参数:', {
            start: params.start,
            end: params.end,
            fullURL: `start=${params.start}&end=${params.end}`
          })
          
          console.log('原始时间对象:', {
            start: startDate,
            end: endDate
          })
          
          console.log('格式化后的时间:', {
            start: formatTime(startDate),
            end: formatTime(endDate)
          })
          
          console.log('编码后的参数:', {
            start: params.start,
            end: params.end
          })
          
          console.log('完整请求URL:', `start=${params.start}&end=${params.end}`)
          
          console.log('原始时间输入:', {
            start: this.customTimeRange[0], 
            end: this.customTimeRange[1]
          })
          
          console.log('验证时间格式:', {
            start: new Date(this.customTimeRange[0]),
            end: new Date(this.customTimeRange[1])
          })
        }

        console.log('请求参数:', JSON.stringify(params, null, 2))
        console.log('编码后URL:', `start=${params.start}&end=${params.end}`)
        console.log('最终请求参数:', JSON.stringify(params, null, 2))
        const response = await this.$api.getHostMonitorData(params)
        console.log('API响应:', {
          status: response.status,
          data: response.data
        })
        console.log('完整响应:', response)
        console.log('响应数据:', JSON.stringify(response.data, null, 2))
        
        if (response.data.code === 200) {
          await this.updateMetrics(response.data.data)
        }
      } catch (error) {
        console.error('获取监控指标失败:', error)
      }
    },
    async updateMetrics(metricsData) {
      console.log('更新指标数据 - 开始处理')
      if (!metricsData) {
        console.error('metricsData为空')
        return
      }
      
      // 确保DOM已渲染
      await this.$nextTick()
      
      // 处理普通指标
      for (let index = 0; index < this.metrics.length; index++) {
        const metric = this.metrics[index]
        
        // 跳过磁盘I/O和网络I/O指标，它们将在单独的方法中处理
        if (metric.key === 'diskReadKB' || metric.key === 'diskWriteKB' || 
            metric.key === 'networkReceive' || metric.key === 'networkSend') {
          // 仍然更新数值显示
          const dataPoints = metricsData[metric.key] || []
          const validDataPoints = dataPoints.filter(point => 
            point && 
            typeof point.value === 'number' && 
            !isNaN(point.value) && 
            typeof point.timestamp === 'number'
          )
          
          if (validDataPoints.length > 0) {
            const sum = validDataPoints.reduce((acc, point) => acc + point.value, 0)
            const avgValue = sum / validDataPoints.length
            metric.value = Math.max(0, avgValue).toFixed(2)
          } else {
            metric.value = '0.00'
          }
          continue
        }
        
        try {
          const dataPoints = metricsData[metric.key] || []
          
          // 验证和清理数据
          const validDataPoints = dataPoints.filter(point => 
            point && 
            typeof point.value === 'number' && 
            !isNaN(point.value) && 
            typeof point.timestamp === 'number'
          )
          
          // 如果没有有效数据，使用默认值
          if (validDataPoints.length === 0) {
            metric.value = '0.00'
            // 创建空图表
            this.createEmptyChart(metric.key, metric)
            continue
          }
          
          // 计算当前值（平均值）
          const sum = validDataPoints.reduce((acc, point) => acc + point.value, 0)
          let avgValue = sum / validDataPoints.length
          
          // 处理百分比指标
          if (metric.key === 'cpu' || metric.key === 'memory' || metric.key === 'disk') {
            // 如果原始值小于1，则认为是小数形式，需要乘以100
            avgValue = avgValue < 1 ? avgValue * 100 : avgValue
            metric.value = Math.max(0, Math.min(100, avgValue)).toFixed(2)
          } else {
            metric.value = Math.max(0, avgValue).toFixed(2)
          }
          
          // 格式化时间戳和数据
          const xAxisData = []
          const seriesData = []
          
          validDataPoints.forEach(point => {
            const date = new Date(point.timestamp * 1000)
            const timeStr = `${date.getHours()}:${date.getMinutes().toString().padStart(2, '0')}`
            xAxisData.push(timeStr)
            
            let value = point.value
            // 根据指标类型处理数据
            if (metric.key === 'cpu' || metric.key === 'memory' || metric.key === 'disk') {
              value = value < 1 ? value * 100 : value
              value = Math.max(0, Math.min(100, value))
            } else {
              value = Math.max(0, value)
            }
            // 确保数据为有限数字
            const numValue = parseFloat(value.toFixed(2))
            seriesData.push(isFinite(numValue) ? numValue : 0)
          })
          
          console.log(`准备渲染图表 ${metric.key} - 数据点: ${seriesData.length}`)
          
          // 获取图表容器
          const chartRef = this.$refs[`chart-${metric.key}`]
          if (!chartRef || !chartRef[0]) {
            console.error(`图表容器chart-${metric.key}未找到`, chartRef)
            continue
          }
          
          this.renderSingleChart(chartRef[0], metric, xAxisData, seriesData)
        } catch (error) {
          console.error(`更新图表 ${metric.key} 失败:`, error)
        }
      }
      
      // 处理磁盘I/O合并图表
      this.renderDiskIOChart(metricsData)
      
      // 处理网络I/O合并图表
      this.renderNetworkIOChart(metricsData)
      
      console.log('更新指标数据 - 完成')
    },
    
    // 渲染单个指标图表
    renderSingleChart(chartDom, metric, xAxisData, seriesData) {
      try {
        console.log(`渲染图表 ${metric.key} - DOM:`, chartDom)
        
        // 获取或创建图表实例
        let chartInstance = this.chartInstances[metric.key]
        if (!chartInstance) {
          chartInstance = echarts.init(chartDom, 'dark')
          this.chartInstances[metric.key] = chartInstance
        }
        
        // 图表配置
        const option = {
          backgroundColor: 'transparent',
          tooltip: {
            trigger: 'axis',
            backgroundColor: 'rgba(0, 0, 0, 0.9)',
            borderColor: '#00d8ff',
            borderWidth: 1,
            textStyle: {
              color: '#00d8ff'
            },
            axisPointer: {
              type: 'line',
              lineStyle: {
                color: '#00d8ff',
                width: 1,
                opacity: 0.8
              }
            },
            formatter: function(params) {
              if (params && params.length > 0) {
                let result = `<div style="color: #00d8ff; font-weight: bold;">${params[0].axisValue}</div>`
                result += `<div style="color: #00d8ff; margin-top: 5px;">• ${params[0].seriesName}: ${params[0].value} ${metric.unit}</div>`
                return result
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
            name: metric.unit,
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
            name: metric.title,
            type: 'line',
            data: seriesData,
            smooth: true,
            symbol: 'none', // 去掉小圆点
            lineStyle: {
              color: '#00d8ff',
              width: 1 // 细化线条
            },
            itemStyle: {
              color: '#00d8ff'
            },
            areaStyle: {
              color: {
                type: 'linear',
                x: 0,
                y: 0,
                x2: 0,
                y2: 1,
                colorStops: [
                  { offset: 0, color: 'rgba(0, 216, 255, 0.2)' },
                  { offset: 1, color: 'rgba(0, 216, 255, 0)' }
                ]
              }
            }
          }]
        }
        
        chartInstance.setOption(option, true)
        console.log(`图表 ${metric.key} 渲染成功`)
      } catch (error) {
        console.error(`渲染图表 ${metric.key} 失败:`, error)
      }
    },
    
    // 渲染磁盘I/O合并图表
    renderDiskIOChart(metricsData) {
      try {
        const chartRef = this.$refs['chart-diskIO']
        if (!chartRef) {
          console.error('磁盘I/O图表容器未找到')
          return
        }
        
        // 获取磁盘读写数据
        const readData = metricsData['diskReadKB'] || []
        const writeData = metricsData['diskWriteKB'] || []
        
        // 验证和清理数据
        const validReadData = readData.filter(point => 
          point && 
          typeof point.value === 'number' && 
          !isNaN(point.value) && 
          typeof point.timestamp === 'number'
        )
        
        const validWriteData = writeData.filter(point => 
          point && 
          typeof point.value === 'number' && 
          !isNaN(point.value) && 
          typeof point.timestamp === 'number'
        )
        
        // 如果没有有效数据，创建空图表
        if (validReadData.length === 0 && validWriteData.length === 0) {
          this.createEmptyChart('diskIO', { key: 'diskIO', title: '磁盘I/O速率' })
          return
        }
        
        // 获取或创建图表实例
        let chartInstance = this.chartInstances['diskIO']
        if (!chartInstance) {
          chartInstance = echarts.init(chartRef, 'dark')
          this.chartInstances['diskIO'] = chartInstance
        }
        
        // 处理时间轴数据（使用读取数据的时间轴，如果没有则使用写入数据的时间轴）
        const timeData = validReadData.length > 0 ? validReadData : validWriteData
        const xAxisData = timeData.map(point => {
          const date = new Date(point.timestamp * 1000)
          return `${date.getHours()}:${date.getMinutes().toString().padStart(2, '0')}`
        })
        
        // 处理读取数据
        const readSeriesData = validReadData.map(point => {
          const value = Math.max(0, point.value)
          return isFinite(value) ? parseFloat(value.toFixed(2)) : 0
        })
        
        // 处理写入数据
        const writeSeriesData = validWriteData.map(point => {
          const value = Math.max(0, point.value)
          return isFinite(value) ? parseFloat(value.toFixed(2)) : 0
        })
        
        // 图表配置
        const option = {
          backgroundColor: 'transparent',
          tooltip: {
            trigger: 'axis',
            backgroundColor: 'rgba(0, 0, 0, 0.9)',
            borderColor: '#00d8ff',
            borderWidth: 1,
            textStyle: {
              color: '#00d8ff'
            },
            axisPointer: {
              type: 'line',
              lineStyle: {
                color: '#00d8ff',
                width: 1,
                opacity: 0.8
              }
            },
            formatter: function(params) {
              if (!params || params.length === 0) return ''
              let result = `<div style="color: #00d8ff; font-weight: bold;">${params[0].axisValue}</div>`
              params.forEach(param => {
                const color = param.seriesName === '磁盘读取' ? '#00d8ff' : '#ff9800'
                result += `<div style="color: ${color}; margin-top: 3px;">• ${param.seriesName}: ${param.value} KB/s</div>`
              })
              return result
            }
          },
          legend: {
            data: ['磁盘读取', '磁盘写入'],
            textStyle: {
              color: '#E4E7ED'
            },
            top: 0
          },
          grid: {
            left: '3%',
            right: '4%',
            bottom: '3%',
            top: '20%',
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
            name: 'KB/s',
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
          series: [
            {
              name: '磁盘读取',
              type: 'line',
              data: readSeriesData,
              smooth: true,
              symbol: 'none', // 去掉小圆点
              lineStyle: {
                color: '#00d8ff',
                width: 1 // 细化线条
              },
              itemStyle: {
                color: '#00d8ff'
              },
              areaStyle: {
                color: {
                  type: 'linear',
                  x: 0,
                  y: 0,
                  x2: 0,
                  y2: 1,
                  colorStops: [
                    { offset: 0, color: 'rgba(0, 216, 255, 0.2)' },
                    { offset: 1, color: 'rgba(0, 216, 255, 0)' }
                  ]
                }
              }
            },
            {
              name: '磁盘写入',
              type: 'line',
              data: writeSeriesData,
              smooth: true,
              symbol: 'none', // 去掉小圆点
              lineStyle: {
                color: '#ff9800',
                width: 1 // 细化线条
              },
              itemStyle: {
                color: '#ff9800'
              },
              areaStyle: {
                color: {
                  type: 'linear',
                  x: 0,
                  y: 0,
                  x2: 0,
                  y2: 1,
                  colorStops: [
                    { offset: 0, color: 'rgba(255, 152, 0, 0.2)' },
                    { offset: 1, color: 'rgba(255, 152, 0, 0)' }
                  ]
                }
              }
            }
          ]
        }
        
        chartInstance.setOption(option, true)
        console.log('磁盘I/O合并图表渲染成功')
      } catch (error) {
        console.error('渲染磁盘I/O图表失败:', error)
      }
    },
    
    // 渲染网络I/O合并图表
    renderNetworkIOChart(metricsData) {
      try {
        const chartRef = this.$refs['chart-networkIO']
        if (!chartRef) {
          console.error('网络I/O图表容器未找到')
          return
        }
        
        // 获取网络收发数据
        const receiveData = metricsData['networkReceive'] || []
        const sendData = metricsData['networkSend'] || []
        
        // 验证和清理数据
        const validReceiveData = receiveData.filter(point => 
          point && 
          typeof point.value === 'number' && 
          !isNaN(point.value) && 
          typeof point.timestamp === 'number'
        )
        
        const validSendData = sendData.filter(point => 
          point && 
          typeof point.value === 'number' && 
          !isNaN(point.value) && 
          typeof point.timestamp === 'number'
        )
        
        // 如果没有有效数据，创建空图表
        if (validReceiveData.length === 0 && validSendData.length === 0) {
          this.createEmptyChart('networkIO', { key: 'networkIO', title: '网络I/O速率' })
          return
        }
        
        // 获取或创建图表实例
        let chartInstance = this.chartInstances['networkIO']
        if (!chartInstance) {
          chartInstance = echarts.init(chartRef, 'dark')
          this.chartInstances['networkIO'] = chartInstance
        }
        
        // 处理时间轴数据
        const timeData = validReceiveData.length > 0 ? validReceiveData : validSendData
        const xAxisData = timeData.map(point => {
          const date = new Date(point.timestamp * 1000)
          return `${date.getHours()}:${date.getMinutes().toString().padStart(2, '0')}`
        })
        
        // 处理接收数据
        const receiveSeriesData = validReceiveData.map(point => {
          const value = Math.max(0, point.value)
          return isFinite(value) ? parseFloat(value.toFixed(2)) : 0
        })
        
        // 处理发送数据
        const sendSeriesData = validSendData.map(point => {
          const value = Math.max(0, point.value)
          return isFinite(value) ? parseFloat(value.toFixed(2)) : 0
        })
        
        // 图表配置
        const option = {
          backgroundColor: 'transparent',
          tooltip: {
            trigger: 'axis',
            backgroundColor: 'rgba(0, 0, 0, 0.9)',
            borderColor: '#00d8ff',
            borderWidth: 1,
            textStyle: {
              color: '#00d8ff'
            },
            axisPointer: {
              type: 'line',
              lineStyle: {
                color: '#00d8ff',
                width: 1,
                opacity: 0.8
              }
            },
            formatter: function(params) {
              if (!params || params.length === 0) return ''
              let result = `<div style="color: #00d8ff; font-weight: bold;">${params[0].axisValue}</div>`
              params.forEach(param => {
                const color = param.seriesName === '网络接收' ? '#17a2b8' : '#6610f2'
                result += `<div style="color: ${color}; margin-top: 3px;">• ${param.seriesName}: ${param.value} KB/s</div>`
              })
              return result
            }
          },
          legend: {
            data: ['网络接收', '网络发送'],
            textStyle: {
              color: '#E4E7ED'
            },
            top: 0
          },
          grid: {
            left: '3%',
            right: '4%',
            bottom: '3%',
            top: '20%',
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
            name: 'KB/s',
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
          series: [
            {
              name: '网络接收',
              type: 'line',
              data: receiveSeriesData,
              smooth: true,
              symbol: 'none',
              lineStyle: {
                color: '#17a2b8',
                width: 1
              },
              itemStyle: {
                color: '#17a2b8'
              },
              areaStyle: {
                color: {
                  type: 'linear',
                  x: 0,
                  y: 0,
                  x2: 0,
                  y2: 1,
                  colorStops: [
                    { offset: 0, color: 'rgba(23, 162, 184, 0.2)' },
                    { offset: 1, color: 'rgba(23, 162, 184, 0)' }
                  ]
                }
              }
            },
            {
              name: '网络发送',
              type: 'line',
              data: sendSeriesData,
              smooth: true,
              symbol: 'none',
              lineStyle: {
                color: '#6610f2',
                width: 1
              },
              itemStyle: {
                color: '#6610f2'
              },
              areaStyle: {
                color: {
                  type: 'linear',
                  x: 0,
                  y: 0,
                  x2: 0,
                  y2: 1,
                  colorStops: [
                    { offset: 0, color: 'rgba(102, 16, 242, 0.2)' },
                    { offset: 1, color: 'rgba(102, 16, 242, 0)' }
                  ]
                }
              }
            }
          ]
        }
        
        chartInstance.setOption(option, true)
        console.log('网络I/O合并图表渲染成功')
      } catch (error) {
        console.error('渲染网络I/O图表失败:', error)
      }
    },
    createEmptyChart(chartKey, metric) {
      try {
        // 获取图表容器
        const chartRef = this.$refs[`chart-${chartKey}`]
        if (!chartRef && !chartRef[0]) {
          console.error(`图表容器chart-${chartKey}未找到`)
          return
        }
        
        // 对于磁盘I/O图表，直接使用chartRef；对于其他图表，使用chartRef[0]
        const chartDom = chartKey === 'diskIO' ? chartRef : (chartRef[0] || chartRef)
        
        // 初始化图表
        if (!this.chartInstances[chartKey]) {
          this.chartInstances[chartKey] = echarts.init(chartDom, 'dark')
        }
        
        // 空数据图表配置
        const emptyOption = {
          backgroundColor: 'transparent',
          graphic: {
            type: 'text',
            left: 'center',
            top: 'middle',
            style: {
              text: '暂无数据',
              fontSize: 14,
              fill: '#6E7079'
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
            data: [],
            axisLine: {
              lineStyle: {
                color: '#6E7079'
              }
            }
          },
          yAxis: {
            type: 'value',
            axisLine: {
              lineStyle: {
                color: '#6E7079'
              }
            },
            splitLine: {
              lineStyle: {
                color: ['#2D3A4B']
              }
            }
          },
          series: [{
            name: metric.title,
            type: 'line',
            data: []
          }]
        }
        
        this.chartInstances[chartKey].setOption(emptyOption, true)
        console.log(`创建空图表 ${metric.key} 成功`)
      } catch (error) {
        console.error(`创建空图表 ${metric.key} 失败:`, error)
      }
    },
    startRefresh() {
      console.log('=== startRefresh 开始 ===', {
        hostId: this.hostId,
        currentVisible: this.visible
      })
      
      // 检查必要参数
      if (!this.hostId) {
        console.error('hostId 未设置，无法加载数据')
        return
      }
      
      // 确保使用默认30分钟范围
      this.timeRangeType = 'preset'
      this.activeDuration = '30m'
      
      console.log('默认时间范围设置:', {
        type: this.timeRangeType,
        duration: this.activeDuration
      })
      
      // 停止之前的刷新
      this.stopRefresh()
      
      // 立即加载数据，不等待nextTick
      console.log('开始加载初始数据...')
      this.fetchMetrics().then(() => {
        console.log('初始数据加载完成')
      }).catch(error => {
        console.error('初始数据加载失败:', error)
      })
      
      // 设置定时器，10秒刷新一次
      this.refreshInterval = setInterval(() => {
        console.log('定时刷新数据...')
        this.fetchMetrics().catch(error => {
          console.error('定时刷新失败:', error)
        })
      }, 10000)
      
      console.log('定时器已设置，间隔10秒')
      console.log('=== startRefresh 结束 ===')
    },
    stopRefresh() {
      if (this.refreshInterval) {
        clearInterval(this.refreshInterval)
        this.refreshInterval = null
      }
    },
    mounted() {
      window.addEventListener('resize', this.handleResize)
    },
  beforeDestroy() {
    console.log('MonitorDialog beforeDestroy - 清理资源')
    window.removeEventListener('resize', this.handleResize)
    this.stopRefresh()
    this.chartInstances.forEach(chart => {
      if (chart) {
        chart.dispose()
      }
    })
    this.chartInstances = []
  },
    handleResize() {
      this.chartInstances.forEach(chart => {
        chart && chart.resize()
      })
    },
    selectPresetDuration(duration) {
      this.activeDuration = duration
      this.fetchMetrics()
    },
    handleTimeRangeChange() {
      this.fetchMetrics()
    },
    handleRawDateChange(dates) {
      if (!dates || dates.length !== 2) return
      
      // 手动格式化日期对象
      const formatDate = (date) => {
        const year = date.getFullYear()
        const month = String(date.getMonth() + 1).padStart(2, '0')
        const day = String(date.getDate()).padStart(2, '0')
        const hours = String(date.getHours()).padStart(2, '0')
        const minutes = String(date.getMinutes()).padStart(2, '0')
        const seconds = String(date.getSeconds()).padStart(2, '0')
        
        return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
      }
      
      this.customTimeRange = [
        formatDate(dates[0]),
        formatDate(dates[1])
      ]
      
      console.log('格式化后的时间范围:', this.customTimeRange)
    },
    
    handleCustomTimeChange() {
      // 保留此方法以兼容旧代码
    },
    handleQueryClick() {
      console.log('查询按钮点击事件触发')
      console.log('当前时间范围:', this.customTimeRange)
      
      if (this.customTimeRange && this.customTimeRange.length === 2) {
        console.log('准备发送请求...')
        this.fetchMetrics()
      } else {
        console.warn('时间范围无效，无法发送请求')
        this.$message.warning('请选择完整的时间范围')
      }
    },
    handleClose() {
      this.stopRefresh()
      this.chartInstances.forEach(chart => {
        chart && chart.dispose()
      })
      this.$emit('change', false)
    }
  },
  mounted() {
    console.log('MonitorDialog mounted, visible:', this.visible, 'hostId:', this.hostId)
    
    // 添加resize监听器
    window.addEventListener('resize', this.handleResize)
    
    // 由于v-if的原因，组件创建后就说明是可见的，立即加载数据
    if (this.hostId) {
      console.log('Component mounted, immediately starting refresh')
      this.timeRangeType = 'preset'
      this.activeDuration = '30m'
      
      // 使用简短延迟确保DOM准备好
      this.$nextTick(() => {
        setTimeout(() => {
          this.startRefresh()
        }, 200) // 稍微增加延迟确保图表容器完全准备
      })
    } else {
      console.warn('Component mounted but hostId is not provided')
    }
  },
  watch: {
    visible: {
      handler(newVal, oldVal) {
        console.log('visible prop changed:', newVal, 'oldVal:', oldVal)
        if (newVal && !oldVal) {
          // 只有在从 false 变为 true 时才执行（但由于v-if，这种情况很少发生）
          console.log('Dialog became visible, starting refresh immediately')
          this.timeRangeType = 'preset'
          this.activeDuration = '30m'
          this.startRefresh()
        } else if (!newVal && oldVal) {
          // 只有在从 true 变为 false 时才执行
          console.log('Dialog became hidden, stopping refresh')
          this.stopRefresh()
        }
      }
    },
    // 添加hostId的监听，确保在hostId变化时也能加载数据
    hostId: {
      handler(newVal, oldVal) {
        console.log('hostId changed:', newVal, 'oldVal:', oldVal)
        if (newVal && newVal !== oldVal) {
          console.log('HostId changed, restarting refresh')
          this.timeRangeType = 'preset'
          this.activeDuration = '30m'
          this.$nextTick(() => {
            this.startRefresh()
          })
        }
      }
    }
  }
}
</script>

<style scoped>
.time-selector {
  margin-bottom: 20px;
}

.preset-buttons {
  margin-top: 10px;
}

.preset-buttons .el-button {
  margin-right: 10px;
  margin-bottom: 10px;
}

.custom-time {
  margin-top: 10px;
}

.monitor-row {
  margin-bottom: 20px;
}

.metric-card {
  margin-bottom: 20px;
  height: 220px;
  background-color: #1a1a2e;
  border: 1px solid #00d8ff;
  color: #00d8ff;
}

.metric-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  color: #00d8ff;
}

.metric-title {
  font-weight: bold;
  font-size: 14px;
  color: #00d8ff;
}

.disk-io-tags {
  display: flex;
  align-items: center;
}

.network-io-tags {
  display: flex;
  align-items: center;
}

.metric-chart {
  height: 150px;
}

/deep/ .monitor-dialog {
  background-color: #121212;
  border: 1px solid #00d8ff;
}

/deep/ .monitor-dialog .el-dialog__title {
  color: #00d8ff;
}

/deep/ .monitor-dialog .el-dialog__header {
  border-bottom: 1px solid #00d8ff;
}

/deep/ .monitor-dialog .el-dialog__body {
  background-color: #121212;
  color: #00d8ff;
}

/deep/ .monitor-dialog .el-radio-button__inner {
  background-color: #1a1a2e;
  color: #00d8ff;
  border-color: #00d8ff;
}

/deep/ .monitor-dialog .el-radio-button__orig-radio:checked + .el-radio-button__inner {
  background-color: #00d8ff;
  color: #121212;
  border-color: #00d8ff;
}

/deep/ .monitor-dialog .el-button {
  background-color: #1a1a2e;
  color: #00d8ff;
  border-color: #00d8ff;
}

/deep/ .monitor-dialog .el-button:hover {
  background-color: #00d8ff;
  color: #121212;
}

/deep/ .monitor-dialog .el-button--primary {
  background-color: #00d8ff;
  color: #121212;
  border-color: #00d8ff;
}

/deep/ .monitor-dialog .el-button--primary:hover {
  background-color: #00a8d8;
  color: #121212;
}

/deep/ .monitor-dialog .el-date-editor .el-range-input {
  background-color: #1a1a2e;
  color: #00d8ff;
}

/deep/ .monitor-dialog .el-date-editor .el-range-separator {
  color: #00d8ff;
}

/deep/ .monitor-dialog .el-tag {
  background-color: #1a1a2e;
  color: #00d8ff;
  border-color: #00d8ff;
}
</style>
