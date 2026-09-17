<template>
  <el-dialog
    title="TCP端口监控"
    v-model:visible="dialogVisible"
    width="60%"
    top="8vh"
    @close="handleClose"
    custom-class="tcp-port-monitor-dialog"
    :modal="true"
    :append-to-body="true"
  >
    <div class="tcp-port-monitor-header">
      <div class="host-info">
        <span class="host-name">{{ portData.hostName || '未知主机' }}</span>
        <span class="update-time">最后更新: {{ formatTime(portData.updateTime) }}</span>
        <div class="port-summary">
          <el-tag type="info" size="small">总端口数: {{ portData.total || 0 }}</el-tag>
          <el-tag type="success" size="small">监听中: {{ listeningPortsCount }}</el-tag>
        </div>
      </div>
      <el-button
        type="primary"
        size="mini"
        icon="Refresh"
        :loading="loading"
        @click="fetchPortData"
      >
        手动刷新
      </el-button>
    </div>

    <div class="port-filters">
      <el-row :gutter="12">
        <el-col :span="5">
          <el-input
            v-model="searchPort"
            placeholder="端口号"
            clearable
            prefix-icon="Search"
            size="small"
          />
        </el-col>
        <el-col :span="5">
          <el-input
            v-model="searchService"
            placeholder="服务名"
            clearable
            prefix-icon="Search"
            size="small"
          />
        </el-col>
        <el-col :span="4">
          <el-select
            v-model="statusFilter"
            placeholder="状态"
            clearable
            size="small"
          >
            <el-option label="全部" value="" />
            <el-option label="监听中" :value="1" />
            <el-option label="未监听" :value="0" />
          </el-select>
        </el-col>
        <el-col :span="5">
          <el-select
            v-model="resourceFilter"
            placeholder="资源筛选"
            clearable
            size="small"
          >
            <el-option label="全部" value="" />
            <el-option label="高CPU (>5%)" value="high-cpu" />
            <el-option label="高内存 (>5%)" value="high-mem" />
            <el-option label="高负载" value="high-load" />
          </el-select>
        </el-col>
        <el-col :span="5">
          <el-button type="primary" size="small" @click="exportPorts">
            <el-icon><Download /></el-icon>
            导出
          </el-button>
        </el-col>
      </el-row>
    </div>

    <el-table
      :data="filteredPorts"
      stripe
      border
      height="450"
      v-loading="loading"
      class="port-table"
      size="small"
    >
      <el-table-column label="端口" prop="port" sortable width="100" align="center">
        <template v-slot="scope">
          <el-tag
            :type="scope.row.status === 1 ? 'success' : 'info'"
            size="mini"
          >
            {{ scope.row.port }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="PID" prop="pid" sortable width="110" align="center">
        <template v-slot="scope">
          <span class="pid-text">{{ scope.row.pid }}</span>
        </template>
      </el-table-column>
      <el-table-column label="服务名称" prop="service" show-overflow-tooltip width="180">
        <template v-slot="scope">
          <div class="service-cell">
            <el-icon size="12"><Setting /></el-icon>
            <span class="service-name">{{ scope.row.service || '未知' }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="状态" prop="status" width="110" align="center">
        <template v-slot="scope">
          <el-tag
            :type="scope.row.status === 1 ? 'success' : 'danger'"
            size="mini"
          >
            <el-icon v-if="scope.row.status === 1" size="10"><CircleCheck /></el-icon>
            <el-icon v-else size="10"><CircleClose /></el-icon>
            {{ scope.row.status === 1 ? '监听' : '停止' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="CPU%" prop="cpuUsage" sortable width="130" align="center">
        <template v-slot="scope">
          <div v-if="scope.row.cpuUsage !== undefined && scope.row.cpuUsage !== null">
            <el-progress
              :percentage="Math.min(parseFloat(scope.row.cpuUsage.toFixed(2)), 100)"
              :color="getUsageColor(scope.row.cpuUsage)"
              :stroke-width="4"
              :show-text="true"
              :format="() => `${scope.row.cpuUsage.toFixed(1)}%`"
            />
          </div>
          <span v-else class="text-muted">-</span>
        </template>
      </el-table-column>
      <el-table-column label="内存%" prop="memUsage" sortable width="130" align="center">
        <template v-slot="scope">
          <div v-if="scope.row.memUsage !== undefined && scope.row.memUsage !== null">
            <el-progress
              :percentage="Math.min(parseFloat(scope.row.memUsage.toFixed(2)), 100)"
              :color="getUsageColor(scope.row.memUsage)"
              :stroke-width="4"
              :show-text="true"
              :format="() => `${scope.row.memUsage.toFixed(1)}%`"
            />
          </div>
          <span v-else class="text-muted">-</span>
        </template>
      </el-table-column>
      <el-table-column label="协议" width="80" align="center">
        <template v-slot="scope">
          <el-tag type="info" size="mini">TCP</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" min-width="90" align="center">
        <template v-slot="scope">
          <el-button
            type="primary"
            size="mini"
            icon="InfoFilled"
            circle
            @click="showPortDetails(scope.row)"
            title="查看详情"
          />
        </template>
      </el-table-column>
    </el-table>
    <!-- 端口详情弹窗 -->
    <el-dialog
      title="端口详情"
      v-model="showDetailDialog"
      width="50%"
      append-to-body
    >
      <el-descriptions :column="2" border v-if="selectedPort">
        <el-descriptions-item label="端口号">
          <el-tag :type="selectedPort.status === 1 ? 'success' : 'danger'">
            {{ selectedPort.port }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="进程ID">{{ selectedPort.pid }}</el-descriptions-item>
        <el-descriptions-item label="服务名称">{{ selectedPort.service || '未知服务' }}</el-descriptions-item>
        <el-descriptions-item label="监听状态">
          <el-tag :type="selectedPort.status === 1 ? 'success' : 'danger'">
            {{ selectedPort.status === 1 ? '监听中' : '未监听' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="CPU使用率" v-if="selectedPort.cpuUsage !== undefined">
          <el-progress
            :percentage="Math.min(parseFloat(selectedPort.cpuUsage.toFixed(2)), 100)"
            :color="getUsageColor(selectedPort.cpuUsage)"
            :stroke-width="8"
            :show-text="true"
            :format="() => `${selectedPort.cpuUsage.toFixed(2)}%`"
          />
        </el-descriptions-item>
        <el-descriptions-item label="内存使用率" v-if="selectedPort.memUsage !== undefined">
          <el-progress
            :percentage="Math.min(parseFloat(selectedPort.memUsage.toFixed(2)), 100)"
            :color="getUsageColor(selectedPort.memUsage)"
            :stroke-width="8"
            :show-text="true"
            :format="() => `${selectedPort.memUsage.toFixed(2)}%`"
          />
        </el-descriptions-item>
        <el-descriptions-item label="协议类型">TCP</el-descriptions-item>
        <el-descriptions-item label="最后更新">{{ formatTime(portData.updateTime) }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </el-dialog>
</template>

<script>

export default {
  name: 'TcpPortMonitorDialog',
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
      portData: {},
      refreshInterval: null,
      searchPort: '',
      searchService: '',
      statusFilter: '',
      resourceFilter: '',
      showDetailDialog: false,
      selectedPort: null
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
    listeningPortsCount() {
      if (!this.portData.ports) return 0
      return this.portData.ports.filter(port => port.status === 1).length
    },
    filteredPorts() {
      if (!this.portData.ports) return []

      let filtered = this.portData.ports

      // 按端口号搜索
      if (this.searchPort) {
        filtered = filtered.filter(port =>
          port.port.toString().includes(this.searchPort)
        )
      }

      // 按服务名搜索
      if (this.searchService) {
        filtered = filtered.filter(port =>
          (port.service || '').toLowerCase().includes(this.searchService.toLowerCase())
        )
      }

      // 按状态筛选
      if (this.statusFilter !== '') {
        filtered = filtered.filter(port => port.status === this.statusFilter)
      }

      // 按资源使用率筛选
      if (this.resourceFilter) {
        filtered = filtered.filter(port => {
          const cpuPercent = port.cpuUsage || 0
          const memPercent = port.memUsage || 0

          switch (this.resourceFilter) {
            case 'high-cpu':
              return cpuPercent > 5
            case 'high-mem':
              return memPercent > 5
            case 'high-load':
              return cpuPercent > 5 || memPercent > 5
            default:
              return true
          }
        })
      }

      return filtered
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
      if (newVal && this.visible) {
        this.startMonitoring()
      }
    }
  },
  methods: {
    async fetchPortData() {
      if (!this.hostId) return

      this.loading = true
      try {
        const { data: res } = await this.$api.getHostTcpPorts(this.hostId)
        if (res.code === 200) {
          this.portData = res.data
        } else {
          this.$message.error(res.message || '获取端口数据失败')
        }
      } catch (error) {
        console.error('获取端口数据失败:', error)
        this.$message.error('获取端口数据失败')
      } finally {
        this.loading = false
      }
    },

    formatTime(timestamp) {
      if (!timestamp) return '未知时间'
      const date = new Date(timestamp * 1000)
      return date.toLocaleString()
    },

    startMonitoring() {
      if (!this.hostId) {
        console.error('hostId 未设置，无法加载数据')
        return
      }

      this.stopMonitoring()
      this.fetchPortData()

      // 设置定时器，60秒刷新一次（端口变化相对较慢）
      this.refreshInterval = setInterval(() => {
        this.fetchPortData()
      }, 60000)
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
      this.resetFilters()
    },

    resetFilters() {
      this.searchPort = ''
      this.searchService = ''
      this.statusFilter = ''
      this.resourceFilter = ''
    },

    showPortDetails(port) {
      this.selectedPort = port
      this.showDetailDialog = true
    },

    getUsageColor(usage) {
      if (!usage) return '#909399'
      if (usage > 80) return '#F56C6C'
      if (usage > 60) return '#E6A23C'
      return '#67C23A'
    },

    exportPorts() {
      if (!this.filteredPorts || this.filteredPorts.length === 0) {
        this.$message.warning('没有数据可导出')
        return
      }

      // CSV导出 - 包含CPU和内存使用率
      const headers = ['端口号', '进程ID', '服务名称', '监听状态', 'CPU使用率(%)', '内存使用率(%)', '协议类型']
      const csvContent = [
        headers.join(','),
        ...this.filteredPorts.map(port => [
          port.port,
          port.pid,
          `"${port.service || '未知服务'}"`,
          port.status === 1 ? '监听中' : '未监听',
          port.cpuUsage !== undefined ? port.cpuUsage.toFixed(2) : '-',
          port.memUsage !== undefined ? port.memUsage.toFixed(2) : '-',
          'TCP'
        ].join(','))
      ].join('\n')

      const blob = new Blob(['\ufeff' + csvContent], { type: 'text/csv;charset=utf-8;' })
      const fileName = `${this.portData.hostName || 'unknown'}_tcp_ports_${new Date().toISOString().slice(0, 10)}.csv`

      const link = document.createElement('a')
      const url = URL.createObjectURL(blob)
      link.setAttribute('href', url)
      link.setAttribute('download', fileName)
      link.style.visibility = 'hidden'
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)

      this.$message.success('数据导出成功')
    }
  },

  mounted() {
    if (this.hostId) {
      this.$nextTick(() => {
        setTimeout(() => {
          this.startMonitoring()
        }, 100)
      })
    }
  },

  beforeUnmount() {
    this.stopMonitoring()
  }
}
</script>

<style scoped>
.tcp-port-monitor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 15px 20px;
  background: #1e1e1e;
  border-radius: 8px;
  border: 1px solid #2D3A4B;
}

.host-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.host-name {
  font-size: 16px;
  font-weight: bold;
  color: #00d8ff;
}

.update-time {
  font-size: 12px;
  color: #6E7079;
}

.port-summary {
  display: flex;
  gap: 10px;
}

.port-filters {
  margin-bottom: 20px;
  padding: 15px;
  background: #1e1e1e;
  border-radius: 8px;
  border: 1px solid #2D3A4B;
}

.port-table {
  background: #1e1e1e;
  border-radius: 8px;
  overflow: hidden;
  width: 100%;
}

.port-table .el-table {
  width: 100% !important;
}

.service-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.tcp-port-monitor-dialog .el-dialog {
  background-color: #1e1e1e !important;
}

.tcp-port-monitor-dialog .el-dialog__header {
  background-color: #1e1e1e !important;
  border-bottom: 1px solid #2D3A4B !important;
}

.tcp-port-monitor-dialog .el-dialog__title {
  color: #E4E7ED !important;
}

.tcp-port-monitor-dialog .el-dialog__body {
  background-color: #1e1e1e !important;
  color: #E4E7ED !important;
}

.port-table .el-table {
  background-color: #1e1e1e !important;
}

.port-table .el-table__header {
  background-color: #2D3A4B !important;
}

.port-table .el-table__header th {
  background-color: transparent !important;
  color: #E4E7ED !important;
  border-bottom: 1px solid #4A5568;
}

.port-table .el-table__row {
  background-color: #1e1e1e !important;
  color: #E4E7ED;
}

.port-table .el-table__row:hover {
  background-color: rgba(103, 126, 234, 0.1) !important;
}

.port-table .el-table__row--striped {
  background-color: rgba(255, 255, 255, 0.05) !important;
}

.port-table .el-table td {
  border-bottom: 1px solid #2D3A4B;
}

/* 进度条样式优化 */
.port-table .el-progress {
  margin: 2px 0;
}

.port-table .el-progress__text {
  color: #E4E7ED !important;
  font-size: 11px;
}

.text-muted {
  color: #909399;
  font-style: italic;
  font-size: 12px;
}

/* 紧凑的服务名称样式 */
.service-cell {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
}

.service-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 140px;
}

.pid-text {
  font-size: 12px;
  color: #E4E7ED;
}

/* 表格行高紧凑 */
.port-table .el-table--small .el-table__body td {
  padding: 4px 0;
}

/* 筛选区域样式优化 */
.port-filters {
  margin-bottom: 15px;
  padding: 12px 15px;
  background: #1e1e1e;
  border-radius: 8px;
  border: 1px solid #2D3A4B;
}
</style>