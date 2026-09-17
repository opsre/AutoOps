<template>
  <div class="cmdb-snmp-management">
    <el-card shadow="hover" class="snmp-card">
      <template #header>
        <div class="card-header">
          <span class="title">SNMP设备管理</span>
          <el-tag type="success">自动采集模式</el-tag>
        </div>
      </template>

      <!-- 搜索表单 -->
      <div class="search-section">
        <el-form :inline="true" :model="queryParams" ref="queryForm" class="search-form">
          <el-form-item label="设备名称">
            <el-input
              v-model="queryParams.deviceName"
              placeholder="请输入设备名称"
              clearable
              size="small"
              style="width: 200px"
              @keyup.enter="handleQuery"
            />
          </el-form-item>
          <el-form-item label="设备类型">
            <el-select
              v-model="queryParams.deviceType"
              placeholder="请选择设备类型"
              clearable
              size="small"
              style="width: 150px">
              <el-option label="交换机" value="switch" />
            </el-select>
          </el-form-item>
          <el-form-item label="IP地址">
            <el-input
              v-model="queryParams.ipAddress"
              placeholder="请输入IP地址"
              clearable
              size="small"
              style="width: 180px"
              @keyup.enter="handleQuery"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" size="small" @click="handleQuery">
              <el-icon><Search /></el-icon>
              搜索
            </el-button>
            <el-button type="warning" size="small" @click="resetQuery">
              <el-icon><Refresh /></el-icon>
              重置
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 操作按钮 -->
      <div class="action-section">
        <el-button
          type="success"
          size="small"
          v-authority="['cmdb:snmp:add']"
          @click="addDeviceDialogVisible = true">
          <el-icon><Plus /></el-icon>
          新增资产
        </el-button>
        <el-button
          type="primary"
          size="small"
          @click="handleRefresh">
          <el-icon><Refresh /></el-icon>
          刷新列表
        </el-button>
      </div>

      <!-- 数据表格 -->
      <div class="table-section">
        <el-table
          :data="deviceList"
          v-loading="loading"
          stripe
          style="width: 100%"
          class="device-table">
          <el-table-column label="设备名称" prop="deviceName" min-width="180">
            <template #default="{ row }">
              <div class="device-name-cell">
                <el-icon class="device-icon" :class="getDeviceIconClass(row.deviceType)">
                  <component :is="getDeviceIcon(row.deviceType)" />
                </el-icon>
                <span class="device-name">{{ row.deviceName }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="设备类型" prop="deviceType" width="120">
            <template #default="{ row }">
              <div class="device-type-cell">
                <img v-if="row.deviceType === 'switch'" src="@/assets/image/交换机.svg" class="device-type-icon" />
                <img v-else-if="row.deviceType === 'windows'" src="@/assets/image/Windows.svg" class="device-type-icon" />
                <span class="device-type-text">{{ getDeviceTypeText(row.deviceType) }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="IP地址" prop="ipAddress" width="100" />
          <el-table-column label="监控状态" prop="monitorEnabled" width="80">
            <template #default="{ row }">
              <el-switch
                v-model="row.monitorEnabled"
                :active-value="1"
                :inactive-value="0"
                @change="handleMonitorToggle(row)"
                v-authority="['cmdb:snmp:update']"
              />
            </template>
          </el-table-column>
          <el-table-column label="设备状态" prop="status" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 'online' ? 'success' : 'danger'" size="small">
                {{ row.status === 'online' ? '在线' : '离线' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="CPU" prop="cpuUsage" width="80">
            <template #default="{ row }">
              <el-tag v-if="row.cpuUsage !== undefined" :type="getResourceUsageType(row.cpuUsage)" size="small">
                {{ row.cpuUsage }}%
              </el-tag>
              <span v-else style="color: #909399;">-</span>
            </template>
          </el-table-column>
          <el-table-column label="内存" prop="memoryUsage" width="80">
            <template #default="{ row }">
              <el-tag v-if="row.memoryUsage !== undefined" :type="getResourceUsageType(row.memoryUsage)" size="small">
                {{ row.memoryUsage }}%
              </el-tag>
              <span v-else style="color: #909399;">-</span>
            </template>
          </el-table-column>
          <el-table-column label="位置" prop="location" min-width="120">
            <template #default="{ row }">
              <span>{{ row.location || '-' }}</span>
            </template>
          </el-table-column>
          <el-table-column label="负责人" prop="owner" min-width="100">
            <template #default="{ row }">
              <span>{{ row.owner || '-' }}</span>
            </template>
          </el-table-column>
          <el-table-column label="备注" prop="remark" min-width="150" show-overflow-tooltip>
            <template #default="{ row }">
              <span>{{ row.remark || '-' }}</span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="400" fixed="right">
            <template #default="{ row }">
              <div class="operation-buttons">
                <el-tooltip content="监控" placement="top">
                  <el-button
                    type="primary"
                    size="small"
                    circle
                    @click="handleOpenMonitor(row)"
                    v-authority="['cmdb:snmp:collect']">
                    <el-icon><TrendCharts /></el-icon>
                  </el-button>
                </el-tooltip>

                <!-- Windows 设备显示 Exporter 管理按钮 -->
                <el-tooltip v-if="row.deviceType === 'windows'" content="Exporter管理" placement="top">
                  <el-button
                    type="warning"
                    size="small"
                    circle
                    @click="handleShowExporterDialog(row)"
                    v-authority="['cmdb:snmp:collect']">
                    <el-icon><Setting /></el-icon>
                  </el-button>
                </el-tooltip>

                <!-- Windows 设备显示 RDP 按钮 -->
                <el-tooltip v-if="row.deviceType === 'windows'" content="RDP远程桌面" placement="top">
                  <el-button
                    type="success"
                    size="small"
                    circle
                    @click="handleOpenRDP(row)"
                    v-authority="['cmdb:snmp:ssh']">
                    <el-icon><Monitor /></el-icon>
                  </el-button>
                </el-tooltip>

                <!-- Linux/交换机设备显示 SSH 按钮 -->
                <el-tooltip v-if="row.deviceType !== 'windows'" content="SSH终端" placement="top">
                  <el-button
                    type="success"
                    size="small"
                    circle
                    @click="handleOpenSSH(row)"
                    v-authority="['cmdb:snmp:ssh']">
                    <el-icon><Monitor /></el-icon>
                  </el-button>
                </el-tooltip>

                <el-tooltip content="编辑" placement="top">
                  <el-button
                    type="warning"
                    size="small"
                    circle
                    @click="showEditDeviceDialog(row.id)"
                    v-authority="['cmdb:snmp:update']">
                    <el-icon><Edit /></el-icon>
                  </el-button>
                </el-tooltip>
                <el-tooltip content="删除" placement="top">
                  <el-button
                    type="danger"
                    size="small"
                    circle
                    @click="handleDeviceDelete(row)"
                    v-authority="['cmdb:snmp:delete']">
                    <el-icon><Delete /></el-icon>
                  </el-button>
                </el-tooltip>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 分页 -->
      <div class="pagination-section">
        <el-pagination
          v-model:current-page="queryParams.page"
          v-model:page-size="queryParams.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!--新增资产-->
    <el-dialog title="新增资产" v-model="addDeviceDialogVisible" width="35%" @close="addDeviceDialogClosed">
      <el-form :model="addDeviceForm" :rules="addDeviceFormRules" ref="addDeviceFormRef" label-width="110px">
        <el-form-item label="资产名称" prop="deviceName">
          <el-input v-model="addDeviceForm.deviceName" placeholder="例如: 核心交换机01"></el-input>
        </el-form-item>
        <el-form-item label="资产类型" prop="deviceType">
          <el-select v-model="addDeviceForm.deviceType" placeholder="请选择设备类型" style="width: 100%">
            <el-option label="交换机 (Switch)" value="switch" />
          </el-select>
        </el-form-item>
        <el-form-item label="IP地址" prop="ipAddress">
          <el-input v-model="addDeviceForm.ipAddress" placeholder="例如: 192.168.1.1"></el-input>
        </el-form-item>

        <!-- SNMP配置(仅交换机显示) -->
        <template v-if="addDeviceForm.deviceType === 'switch'">
          <el-divider content-position="left">SNMP配置</el-divider>
          <el-row :gutter="20">
            <el-col :span="8">
              <el-form-item label="版本" prop="snmpVersion" label-width="50px">
                <el-input v-model="addDeviceForm.snmpVersion" placeholder="v2c" disabled size="small"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="Community" prop="snmpCommunity" label-width="80px">
                <el-input v-model="addDeviceForm.snmpCommunity" placeholder="public" size="small"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="端口" prop="snmpPort" label-width="50px">
                <el-input-number v-model="addDeviceForm.snmpPort" :min="1" :max="65535" placeholder="161" size="small" style="width: 100%"></el-input-number>
              </el-form-item>
            </el-col>
          </el-row>
        </template>

        <!-- 业务信息区域 -->
        <el-divider content-position="left">业务信息</el-divider>
        <el-form-item label="位置" prop="location">
          <el-input v-model="addDeviceForm.location" placeholder="例如: 机房A-01"></el-input>
        </el-form-item>
        <el-form-item label="负责人" prop="owner">
          <el-input v-model="addDeviceForm.owner" placeholder="例如: 张三"></el-input>
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="addDeviceForm.remark" type="textarea" :rows="2" placeholder="例如: 核心交换机"></el-input>
        </el-form-item>

        <!-- 远程访问配置(可选) -->
        <el-divider content-position="left">远程访问配置(可选)</el-divider>
        <el-row :gutter="20">

          <el-col :span="8">
            <el-form-item label="用户名" prop="remoteUsername" label-width="60px">
              <el-input v-model="addDeviceForm.remoteUsername" :placeholder="addDeviceForm.deviceType === 'windows' ? 'Administrator' : 'admin'" size="small"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="密码" prop="remotePassword" label-width="50px">
              <el-input v-model="addDeviceForm.remotePassword" type="password" show-password placeholder="密码" size="small"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="addDeviceForm.deviceType === 'windows' ? 'RDP端口' : 'SSH端口'" prop="remotePort" label-width="70px">
              <el-input-number v-model="addDeviceForm.remotePort" :min="1" :max="65535" :placeholder="addDeviceForm.deviceType === 'windows' ? '3389' : '22'" size="small" style="width: 100%"></el-input-number>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="域" prop="remoteDomain" v-if="addDeviceForm.deviceType === 'windows'">
          <el-input v-model="addDeviceForm.remoteDomain" placeholder="例如: WORKGROUP(可选)"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="addDeviceDialogVisible = false">取 消</el-button>
          <el-button type="primary" @click="addDevice">确 定</el-button>
        </span>
      </template>
    </el-dialog>

    <!--修改SNMP设备-->
    <el-dialog title="编辑资产" v-model="editDeviceDialogVisible" width="35%">
      <el-form :model="deviceInfo" :rules="editDeviceFormRules" ref="editDeviceFormRef" label-width="110px">
        <el-form-item label="资产名称" prop="deviceName">
          <el-input v-model="deviceInfo.deviceName"></el-input>
        </el-form-item>
        <el-form-item label="资产类型" prop="deviceType">
          <el-select v-model="deviceInfo.deviceType" placeholder="请选择资产类型" style="width: 100%">
            <el-option label="交换机 (Switch)" value="switch" />
          </el-select>
        </el-form-item>
        <el-form-item label="IP地址" prop="ipAddress">
          <el-input v-model="deviceInfo.ipAddress" placeholder="例如: 192.168.1.1"></el-input>
        </el-form-item>

        <!-- SNMP配置(仅交换机显示) -->
        <template v-if="deviceInfo.deviceType === 'switch'">
          <el-divider content-position="left">SNMP配置</el-divider>
          <el-row :gutter="20">
            <el-col :span="8">
              <el-form-item label="版本" prop="snmpVersion" label-width="50px">
                <el-input v-model="deviceInfo.snmpVersion" placeholder="v2c" disabled size="small"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="Community" prop="snmpCommunity" label-width="80px">
                <el-input v-model="deviceInfo.snmpCommunity" size="small"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="端口" prop="snmpPort" label-width="50px">
                <el-input-number v-model="deviceInfo.snmpPort" :min="1" :max="65535" placeholder="161" size="small" style="width: 100%"></el-input-number>
              </el-form-item>
            </el-col>
          </el-row>
        </template>

        <!-- 业务信息区域 -->
        <el-divider content-position="left">业务信息</el-divider>
        <el-form-item label="位置" prop="location">
          <el-input v-model="deviceInfo.location" placeholder="例如: 机房A-01"></el-input>
        </el-form-item>
        <el-form-item label="负责人" prop="owner">
          <el-input v-model="deviceInfo.owner" placeholder="例如: 张三"></el-input>
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="deviceInfo.remark" type="textarea" :rows="2" placeholder="例如: 核心交换机"></el-input>
        </el-form-item>

        <!-- 远程访问配置(可选) -->
        <el-divider content-position="left">远程访问配置(可选)</el-divider>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item :label="deviceInfo.deviceType === 'windows' ? 'RDP端口' : 'SSH端口'" prop="remotePort" label-width="70px">
              <el-input-number v-model="deviceInfo.remotePort" :min="1" :max="65535" :placeholder="deviceInfo.deviceType === 'windows' ? '3389' : '22'" size="small" style="width: 100%"></el-input-number>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="用户名" prop="remoteUsername" label-width="60px">
              <el-input v-model="deviceInfo.remoteUsername" :placeholder="deviceInfo.deviceType === 'windows' ? 'Administrator' : 'admin'" size="small"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="密码" prop="remotePassword" label-width="50px">
              <el-input v-model="deviceInfo.remotePassword" type="password" show-password placeholder="密码" size="small"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="域" prop="remoteDomain" v-if="deviceInfo.deviceType === 'windows'">
          <el-input v-model="deviceInfo.remoteDomain" placeholder="例如: WORKGROUP(可选)"></el-input>
        </el-form-item>
        <el-form-item label="启用监控" prop="monitorEnabled">
          <el-switch v-model="deviceInfo.monitorEnabled" :active-value="1" :inactive-value="0"></el-switch>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="editDeviceDialogVisible = false">取 消</el-button>
          <el-button type="primary" @click="editDevice">确 定</el-button>
        </span>
      </template>
    </el-dialog>

    <!--Windows Exporter管理对话框-->
    <el-dialog title="Windows Exporter 管理" v-model="exporterDialogVisible" width="600px">
      <div class="exporter-management">
        <el-card shadow="hover" class="status-card">
          <template #header>
            <div class="card-header">
              <span>资产信息</span>
            </div>
          </template>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="资产名称">{{ currentDevice?.deviceName }}</el-descriptions-item>
            <el-descriptions-item label="IP地址">{{ currentDevice?.ipAddress }}</el-descriptions-item>
            <el-descriptions-item label="状态" :span="2">
              <el-tag v-if="exporterStatus.loading" type="info">查询中...</el-tag>
              <el-tag v-else-if="exporterStatus.status === 0" type="warning">未安装</el-tag>
              <el-tag v-else-if="exporterStatus.status === 1" type="success">已安装运行中</el-tag>
              <el-tag v-else-if="exporterStatus.status === 2" type="primary">安装中</el-tag>
              <el-tag v-else-if="exporterStatus.status === 3" type="info">已安装未运行</el-tag>
              <el-tag v-else-if="exporterStatus.status === 4" type="danger">安装失败</el-tag>
              <el-tag v-else-if="exporterStatus.status === 5" type="primary">启动中</el-tag>
              <el-tag v-else type="info">未知状态</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="版本信息" :span="2" v-if="exporterStatus.version">
              {{ exporterStatus.version }}
            </el-descriptions-item>
          </el-descriptions>
        </el-card>

        <div class="action-buttons">
          <el-button
            v-if="exporterStatus.status === 0 || exporterStatus.status === 4"
            type="success"
            :loading="exporterInstalling"
            @click="handleInstallExporter">
            <el-icon><Download /></el-icon>
            安装 Windows Exporter
          </el-button>
          <el-button
            v-if="exporterStatus.status === 1 || exporterStatus.status === 3"
            type="danger"
            :loading="exporterUninstalling"
            @click="handleUninstallExporter">
            <el-icon><Delete /></el-icon>
            卸载 Windows Exporter
          </el-button>
          <el-button @click="handleRefreshExporterStatus">
            <el-icon><Refresh /></el-icon>
            刷新状态
          </el-button>
        </div>

        <el-alert
          v-if="exporterStatus.status === 1"
          title="提示"
          type="info"
          :closable="false"
          show-icon
          style="margin-top: 20px;">
          Windows Exporter 已安装并配置为主动推送模式，每60秒自动推送指标到 Pushgateway。
        </el-alert>
        <el-alert
          v-if="exporterStatus.status === 4"
          title="安装失败"
          type="error"
          :closable="false"
          show-icon
          style="margin-top: 20px;">
          Windows Exporter 安装失败，请检查资产连接或权限设置后重试。
        </el-alert>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="exporterDialogVisible = false">关 闭</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import API from '@/api/cmdb'
import {
  Search,
  Refresh,
  Plus,
  Edit,
  Delete,
  Clock,
  Monitor,
  Router,
  Connection,
  DataLine,
  Setting,
  Download,
  TrendCharts
} from '@element-plus/icons-vue'

export default {
  components: {
    Search,
    Refresh,
    Plus,
    Edit,
    Delete,
    Clock,
    Monitor,
    Router,
    Connection,
    DataLine,
    Setting,
    Download,
    TrendCharts
  },
  data() {
    return {
      queryParams: {
        page: 1,
        pageSize: 10,
        deviceName: '',
        deviceType: '',
        ipAddress: ''
      },
      loading: true,
      deviceList: [],
      total: 0,
      addDeviceDialogVisible: false,
      addDeviceForm: {
        deviceName: '',
        deviceType: 'switch',
        ipAddress: '',
        snmpVersion: 'v2c',
        snmpCommunity: 'public',
        snmpPort: 161,
        location: '',
        owner: '',
        remark: '',
        remotePort: null,
        remoteUsername: '',
        remotePassword: '',
        remoteDomain: ''
      },
      addDeviceFormRules: {
        deviceName: [{ required: true, message: '请输入资产名称', trigger: 'blur' }],
        deviceType: [{ required: true, message: '请选择资产类型', trigger: 'change' }],
        ipAddress: [
          { required: true, message: '请输入IP地址', trigger: 'blur' },
          {
            pattern: /^(\d{1,3}\.){3}\d{1,3}$/,
            message: 'IP地址格式不正确',
            trigger: 'blur'
          }
        ]
      },
      editDeviceDialogVisible: false,
      deviceInfo: {},
      editDeviceFormRules: {
        deviceName: [{ required: true, message: '请输入资产名称', trigger: 'blur' }],
        deviceType: [{ required: true, message: '请选择资产类型', trigger: 'change' }],
        ipAddress: [
          { required: true, message: '请输入IP地址', trigger: 'blur' },
          {
            pattern: /^(\d{1,3}\.){3}\d{1,3}$/,
            message: 'IP地址格式不正确',
            trigger: 'blur'
          }
        ]
      },
      // Windows Exporter 相关
      currentDevice: null,
      exporterDialogVisible: false,
      exporterStatus: {
        loading: false,
        status: 0, // 0=未安装 1=已安装运行中 2=安装中 3=已安装未运行 4=安装失败 5=启动中
        version: ''
      },
      exporterInstalling: false,
      exporterUninstalling: false
    }
  },
  methods: {
    // 获取资产列表
    async getList() {
      this.loading = true
      try {
        const { data: res } = await API.getSnmpDevices({
          params: this.queryParams
        })

        if (res.code !== 200) {
          throw new Error(res.message || '获取数据失败')
        }

        // 后端返回的是 items 而不是 list
        this.deviceList = res.data.items || res.data.list || []
        this.total = res.data.total || 0

        // 批量获取资源使用率数据
        if (this.deviceList.length > 0) {
          await this.loadBatchResourceUsage()
        }

      } catch (error) {
        console.error('获取资产列表失败:', error)
        this.$message.error(`获取资产列表失败: ${error?.message || '未知错误'}`)
      } finally {
        this.loading = false
      }
    },

    // 批量获取资源使用率
    async loadBatchResourceUsage() {
      try {
        // 按设备类型分组
        const switchDevices = this.deviceList.filter(d => d.deviceType === 'switch')
        const windowsDevices = this.deviceList.filter(d => d.deviceType === 'windows')

        // 批量获取交换机资源使用率
        if (switchDevices.length > 0) {
          const switchIds = switchDevices.map(d => d.id).join(',')
          const { data: res } = await API.getSnmpDevicesBatchResource(switchIds)

          if (res.code === 200 && res.data.items) {
            const resourceMap = {}
            res.data.items.forEach(item => {
              resourceMap[item.deviceId] = {
                cpuUsage: item.cpuUsage,
                memoryUsage: item.memoryUsage
              }
            })

            // 更新交换机设备的资源使用率
            switchDevices.forEach(device => {
              if (resourceMap[device.id]) {
                device.cpuUsage = resourceMap[device.id].cpuUsage
                device.memoryUsage = resourceMap[device.id].memoryUsage
              }
            })
          }
        }

        // 批量获取Windows主机实时CPU和内存数据
        if (windowsDevices.length > 0) {
          const windowsIds = windowsDevices.map(d => d.id).join(',')
          const { data: res } = await API.getWindowsDevicesBatchRealtime(windowsIds)

          if (res.code === 200 && res.data) {
            const resourceMap = {}
            // 后端返回的是数组,且字段名为 cpuPercent 和 memPercent
            const items = Array.isArray(res.data) ? res.data : []
            items.forEach(item => {
              resourceMap[item.deviceId] = {
                cpuUsage: item.cpuPercent !== undefined ? Number(item.cpuPercent).toFixed(2) : undefined,
                memoryUsage: item.memPercent !== undefined ? Number(item.memPercent).toFixed(2) : undefined
              }
            })

            // 更新Windows设备的资源使用率
            windowsDevices.forEach(device => {
              if (resourceMap[device.id]) {
                device.cpuUsage = resourceMap[device.id].cpuUsage
                device.memoryUsage = resourceMap[device.id].memoryUsage
              }
            })
          }
        }
      } catch (error) {
        console.error('获取资源使用率失败:', error)
        // 不显示错误消息，静默失败
      }
    },

    // 根据资源使用率返回标签类型
    getResourceUsageType(usage) {
      if (usage === undefined || usage === null) return 'info'
      if (usage > 80) return 'danger'
      if (usage > 50) return 'warning'
      return 'success'
    },

    // 搜索
    handleQuery() {
      this.queryParams.page = 1
      this.getList()
    },

    // 重置搜索
    resetQuery() {
      this.queryParams = {
        page: 1,
        pageSize: 10,
        deviceName: '',
        deviceType: '',
        ipAddress: ''
      }
      this.getList()
      this.$message.success("重置成功")
    },

    // 刷新列表
    handleRefresh() {
      this.getList()
      this.$message.success("刷新成功")
    },

    // 分页大小变化
    handleSizeChange(val) {
      this.queryParams.pageSize = val
      this.queryParams.page = 1
      this.getList()
    },

    // 当前页变化
    handleCurrentChange(val) {
      this.queryParams.page = val
      this.getList()
    },

    // 监听新增资产对话框关闭
    addDeviceDialogClosed() {
      this.$refs.addDeviceFormRef?.resetFields()
      this.addDeviceForm = {
        deviceName: '',
        deviceType: 'switch',
        ipAddress: '',
        snmpVersion: 'v2c',
        snmpCommunity: 'public',
        snmpPort: 161,
        location: '',
        owner: '',
        remark: '',
        remotePort: null,
        remoteUsername: '',
        remotePassword: '',
        remoteDomain: ''
      }
    },

    // 新增设备
    async addDevice() {
      try {
        // 表单验证
        await this.$refs.addDeviceFormRef.validate()

        // 准备提交数据(只包含必填和有值的字段)
        const deviceData = {
          deviceName: this.addDeviceForm.deviceName,
          deviceType: this.addDeviceForm.deviceType,
          ipAddress: this.addDeviceForm.ipAddress
        }

        // SNMP配置(仅交换机需要,可选字段,有值才传)
        if (this.addDeviceForm.deviceType === 'switch') {
          if (this.addDeviceForm.snmpVersion) {
            deviceData.snmpVersion = this.addDeviceForm.snmpVersion
          }
          if (this.addDeviceForm.snmpCommunity) {
            deviceData.snmpCommunity = this.addDeviceForm.snmpCommunity
          }
          if (this.addDeviceForm.snmpPort) {
            deviceData.snmpPort = this.addDeviceForm.snmpPort
          }
        }

        // 业务信息(可选字段,有值才传)
        if (this.addDeviceForm.location) {
          deviceData.location = this.addDeviceForm.location
        }
        if (this.addDeviceForm.owner) {
          deviceData.owner = this.addDeviceForm.owner
        }
        if (this.addDeviceForm.remark) {
          deviceData.remark = this.addDeviceForm.remark
        }

        // 远程访问配置(可选字段,有值才传)
        if (this.addDeviceForm.remotePort) {
          deviceData.remotePort = this.addDeviceForm.remotePort
        }
        if (this.addDeviceForm.remoteUsername) {
          deviceData.remoteUsername = this.addDeviceForm.remoteUsername
        }
        if (this.addDeviceForm.remotePassword) {
          deviceData.remotePassword = this.addDeviceForm.remotePassword
        }
        if (this.addDeviceForm.remoteDomain) {
          deviceData.remoteDomain = this.addDeviceForm.remoteDomain
        }

        // 发送创建请求
        const { data: res } = await API.createSnmpDevice(deviceData)

        // 判断是否成功
        if (res.code === 200) {
          this.$message.success('新增SNMP资产成功,系统将自动开始采集监控数据')
          this.addDeviceDialogVisible = false
          await this.getList()
        } else {
          this.$message.error(res.message || '新增资产失败')
        }
      } catch (error) {
        console.error('新增资产失败:', error)
        // 如果是表单验证失败，不显示错误消息
        if (error !== false && typeof error === 'object') {
          this.$message.error(error?.response?.data?.message || error?.message || '新增资产失败，请检查网络连接')
        }
      }
    },

    // 展示编辑对话框
    async showEditDeviceDialog(id) {
      try {
        const { data: res } = await API.getSnmpDeviceById(id)

        if (res.code !== 200) {
          throw new Error(res.message || '获取资产信息失败')
        }

        // 直接使用后端返回的数据，确保所有字段都存在
        this.deviceInfo = {
          ...res.data,
          snmpPort: res.data.snmpPort || 161,
          snmpVersion: res.data.snmpVersion || 'v2c',
          snmpCommunity: res.data.snmpCommunity || 'public',
          location: res.data.location || '',
          owner: res.data.owner || '',
          remark: res.data.remark || '',
          // 映射SSH/RDP字段到前端统一的remote字段
          remotePort: res.data.deviceType === 'windows' ? (res.data.rdpPort || null) : (res.data.sshPort || null),
          remoteUsername: res.data.deviceType === 'windows' ? (res.data.rdpUsername || '') : (res.data.sshUsername || ''),
          remotePassword: res.data.deviceType === 'windows' ? (res.data.rdpPassword || '') : (res.data.sshPassword || ''),
          remoteDomain: res.data.rdpDomain || '',
          monitorEnabled: res.data.monitorEnabled !== undefined ? res.data.monitorEnabled : 1
        }

        this.editDeviceDialogVisible = true

      } catch (error) {
        console.error('获取设备信息失败:', error)
        this.$message.error(`获取设备信息失败: ${error?.message || '未知错误'}`)
      }
    },

    // 修改设备信息
    async editDevice() {
      try {
        await this.$refs.editDeviceFormRef.validate()

        // 准备提交数据(所有字段都是可选的,按需更新)
        const deviceData = {}

        // 基础信息
        if (this.deviceInfo.deviceName) deviceData.deviceName = this.deviceInfo.deviceName
        if (this.deviceInfo.deviceType) deviceData.deviceType = this.deviceInfo.deviceType
        if (this.deviceInfo.ipAddress) deviceData.ipAddress = this.deviceInfo.ipAddress

        // SNMP配置(仅交换机需要)
        if (this.deviceInfo.deviceType === 'switch') {
          if (this.deviceInfo.snmpVersion) deviceData.snmpVersion = this.deviceInfo.snmpVersion
          if (this.deviceInfo.snmpCommunity) deviceData.snmpCommunity = this.deviceInfo.snmpCommunity
          if (this.deviceInfo.snmpPort) deviceData.snmpPort = this.deviceInfo.snmpPort
        }

        // 业务信息
        if (this.deviceInfo.location !== undefined) deviceData.location = this.deviceInfo.location
        if (this.deviceInfo.owner !== undefined) deviceData.owner = this.deviceInfo.owner
        if (this.deviceInfo.remark !== undefined) deviceData.remark = this.deviceInfo.remark

        // 远程访问配置
        if (this.deviceInfo.remotePort) deviceData.remotePort = this.deviceInfo.remotePort
        if (this.deviceInfo.remoteUsername !== undefined) deviceData.remoteUsername = this.deviceInfo.remoteUsername
        if (this.deviceInfo.remotePassword !== undefined) deviceData.remotePassword = this.deviceInfo.remotePassword
        if (this.deviceInfo.remoteDomain !== undefined && this.deviceInfo.deviceType === 'windows') {
          deviceData.remoteDomain = this.deviceInfo.remoteDomain
        }

        // 监控状态
        if (this.deviceInfo.monitorEnabled !== undefined) deviceData.monitorEnabled = this.deviceInfo.monitorEnabled

        const { data: res } = await API.updateSnmpDevice(this.deviceInfo.id, deviceData)

        if (res.code === 200) {
          this.$message.success("修改设备成功")
          this.editDeviceDialogVisible = false
          await this.getList()
        } else {
          this.$message.error(res.message || '更新设备失败')
        }

      } catch (error) {
        console.error('修改设备失败:', error)
        if (error !== false && typeof error === 'object') {
          this.$message.error(error?.response?.data?.message || error?.message || '修改设备失败，请检查网络连接')
        }
      }
    },

    // 删除设备
    async handleDeviceDelete(row) {
      try {
        const confirmResult = await this.$confirm(
          `是否确认删除设备"${row.deviceName}"? 删除后将停止采集该设备的监控数据。`,
          '提示',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        ).catch(err => err)

        if (confirmResult !== 'confirm') {
          return this.$message.info('已取消删除')
        }

        const { data: res } = await API.deleteSnmpDevice(row.id)

        if (res.code === 200) {
          this.$message.success('删除设备成功')
          await this.getList()
        } else {
          this.$message.error(res.message || '删除设备失败')
        }
      } catch (error) {
        console.error('删除设备失败:', error)
        this.$message.error(error?.response?.data?.message || error?.message || '删除设备失败，请检查网络连接')
      }
    },

    // 立即采集
    async handleCollect(id) {
      try {
        this.$message.info('正在触发采集...')

        const { data: res } = await API.collectSnmpDevice(id)

        if (res.code === 200) {
          this.$message.success('采集任务已触发,请稍后查看监控数据')

          // 刷新列表以更新状态
          setTimeout(() => {
            this.getList()
          }, 2000)
        } else {
          this.$message.error(res.message || '采集失败')
        }

      } catch (error) {
        console.error('采集失败:', error)
        this.$message.error(error?.response?.data?.message || error?.message || '采集失败')
      }
    },

    // 切换监控状态
    async handleMonitorToggle(row) {
      try {
        // 使用专门的启用/停用接口
        const apiCall = row.monitorEnabled === 1
          ? API.enableSnmpDeviceMonitor(row.id)
          : API.disableSnmpDeviceMonitor(row.id)

        const { data: res } = await apiCall

        if (res.code === 200) {
          const status = row.monitorEnabled === 1 ? '启用' : '停用'
          this.$message.success(`已${status}监控`)
        } else {
          // 如果更新失败，恢复开关状态
          row.monitorEnabled = row.monitorEnabled === 1 ? 0 : 1
          this.$message.error(res.message || '更新失败')
        }

      } catch (error) {
        console.error('切换监控状态失败:', error)
        // 恢复开关状态
        row.monitorEnabled = row.monitorEnabled === 1 ? 0 : 1
        this.$message.error(error?.response?.data?.message || error?.message || '操作失败')
      }
    },

    // 打开 RDP 远程桌面（仅Windows）
    handleOpenRDP(row) {
      this.$router.push({
        path: '/cmdb/network/rdp-terminal',
        query: {
          deviceId: row.id,
          deviceName: row.deviceName,
          ipAddress: row.ipAddress
        }
      })
    },

    // 打开 SSH 终端（Linux/交换机）
    handleOpenSSH(row) {
      this.$router.push({
        path: '/cmdb/network/ssh-terminal',
        query: {
          deviceId: row.id,
          deviceName: row.deviceName,
          ipAddress: row.ipAddress
        }
      })
    },

    // 打开监控页面（按设备类型跳转）
    handleOpenMonitor(row) {
      const path = row.deviceType === 'windows' ? `/cmdb/windows/${row.id}` : `/cmdb/switch/${row.id}`
      this.$router.push({ path })
    },

    // 获取设备图标
    getDeviceIcon(deviceType) {
      const iconMap = {
        'switch': 'Connection',
        'router': 'Router',
        'windows': 'Monitor'
      }
      return iconMap[deviceType] || 'Monitor'
    },

    // 获取设备图标样式类
    getDeviceIconClass(deviceType) {
      const classMap = {
        'switch': 'device-icon-switch',
        'router': 'device-icon-router',
        'windows': 'device-icon-windows'
      }
      return classMap[deviceType] || 'device-icon-default'
    },

    // 获取设备类型标签类型
    getDeviceTypeTag(deviceType) {
      const tagMap = {
        'switch': 'success',
        'router': 'warning',
        'windows': 'primary'
      }
      return tagMap[deviceType] || 'info'
    },

    // 获取设备类型文本
    getDeviceTypeText(deviceType) {
      const textMap = {
        'switch': '交换机',
        'router': '路由器',
        'windows': 'Windows'
      }
      return textMap[deviceType] || deviceType
    },

    // 显示 Windows Exporter 管理对话框
    async handleShowExporterDialog(row) {
      this.currentDevice = row
      this.exporterDialogVisible = true
      await this.handleRefreshExporterStatus()
    },

    // 刷新 Exporter 状态
    async handleRefreshExporterStatus() {
      if (!this.currentDevice) return

      try {
        this.exporterStatus.loading = true
        const { data: res } = await API.getWindowsExporterStatus(this.currentDevice.id)

        if (res.code === 200) {
          this.exporterStatus = {
            loading: false,
            status: res.data.status !== undefined ? res.data.status : 0,
            version: res.data.version || ''
          }
        } else {
          throw new Error(res.message || '查询状态失败')
        }
      } catch (error) {
        console.error('查询 Exporter 状态失败:', error)
        this.$message.error(`查询状态失败: ${error?.message || '未知错误'}`)
        this.exporterStatus = {
          loading: false,
          status: 0,
          version: ''
        }
      }
    },

    // 安装 Windows Exporter
    async handleInstallExporter() {
      if (!this.currentDevice) return

      try {
        const confirmResult = await this.$confirm(
          `确认在设备 "${this.currentDevice.deviceName}" 上安装 Windows Exporter？\n安装过程需要几分钟，请耐心等待。`,
          '确认安装',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        ).catch(err => err)

        if (confirmResult !== 'confirm') {
          return this.$message.info('已取消安装')
        }

        this.exporterInstalling = true
        this.$message.info('正在安装 Windows Exporter，请稍候...')

        const { data: res } = await API.installWindowsExporter({
          deviceId: this.currentDevice.id
        })

        if (res.code === 200) {
          this.$message.success('Windows Exporter 安装成功！')
          await this.handleRefreshExporterStatus()
        } else {
          throw new Error(res.message || '安装失败')
        }
      } catch (error) {
        console.error('安装 Exporter 失败:', error)
        this.$message.error(`安装失败: ${error?.response?.data?.message || error?.message || '未知错误'}`)
      } finally {
        this.exporterInstalling = false
      }
    },

    // 卸载 Windows Exporter
    async handleUninstallExporter() {
      if (!this.currentDevice) return

      try {
        const confirmResult = await this.$confirm(
          `确认从设备 "${this.currentDevice.deviceName}" 卸载 Windows Exporter？\n卸载后将停止指标采集。`,
          '确认卸载',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        ).catch(err => err)

        if (confirmResult !== 'confirm') {
          return this.$message.info('已取消卸载')
        }

        this.exporterUninstalling = true
        this.$message.info('正在卸载 Windows Exporter，请稍候...')

        const { data: res } = await API.uninstallWindowsExporter(this.currentDevice.id)

        if (res.code === 200) {
          this.$message.success('Windows Exporter 卸载成功！')
          await this.handleRefreshExporterStatus()
        } else {
          throw new Error(res.message || '卸载失败')
        }
      } catch (error) {
        console.error('卸载 Exporter 失败:', error)
        this.$message.error(`卸载失败: ${error?.response?.data?.message || error?.message || '未知错误'}`)
      } finally {
        this.exporterUninstalling = false
      }
    }
  },
  watch: {
    // 监听设备类型变化，自动设置默认端口(仅作为提示,不强制填写)
    'addDeviceForm.deviceType'() {
      // 不自动设置端口,让用户选择是否填写
      // 如果需要默认值，用户可以手动填写
    }
  },
  created() {
    this.getList()
  }
}
</script>

<style scoped>
/* 现代化科技感设计风格 - SNMP设备管理 */

.cmdb-snmp-management {
  padding: 5px;
  min-height: 100vh;
  background: #eef1f6;
}

.snmp-card {
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

/* 搜索区域样式 */
.search-section {
  margin-bottom: 20px;
  padding: 20px;
  background: rgba(103, 126, 234, 0.05);
  border-radius: 12px;
  border: 1px solid rgba(103, 126, 234, 0.1);
}

.search-form .el-form-item {
  margin-bottom: 0;
  margin-right: 16px;
}

.search-form .el-form-item__label {
  color: #606266;
  font-weight: 500;
}

/* 操作按钮区域 */
.action-section {
  margin-bottom: 20px;
  padding-left: 0;
}

/* 表格区域样式 */
.table-section {
  margin-top: 0;
}

.device-table {
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.device-table :deep(.el-table__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
}

.device-table :deep(.el-table__header th) {
  background: transparent !important;
  color: #2c3e50 !important;
  font-weight: 700 !important;
  border-bottom: none;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
}

.device-table :deep(.el-table__header th .cell) {
  color: #2c3e50 !important;
  font-weight: 700 !important;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
}

.device-table :deep(.el-table__row) {
  transition: all 0.3s ease;
}

.device-table :deep(.el-table__row:hover) {
  background-color: rgba(103, 126, 234, 0.05) !important;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

/* 设备名称样式 */
.device-name-cell {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.device-icon {
  font-size: 16px;
}

.device-icon-switch {
  color: #67c23a;
}

.device-icon-router {
  color: #e6a23c;
}

.device-icon-windows {
  color: #409eff;
}

.device-icon-default {
  color: #909399;
}

.device-name {
  font-weight: 600;
  color: #303133;
  font-size: 14px;
}

/* 设备类型样式 */
.device-type-cell {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.device-type-icon {
  width: 20px;
  height: 20px;
  object-fit: contain;
}

.device-type-text {
  font-weight: 500;
  color: #303133;
  font-size: 14px;
}

/* Community文本样式 */
.community-text {
  font-family: 'Monaco', 'Courier New', monospace;
  color: #606266;
  font-size: 13px;
}

.time-cell {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: #606266;
}

.time-icon {
  color: #909399;
  font-size: 14px;
}

.no-data {
  color: #c0c4cc;
  font-style: italic;
}

/* 操作按钮样式 */
.operation-buttons {
  display: flex;
  gap: 8px;
  justify-content: center;
}

.operation-buttons .el-button {
  transition: all 0.3s ease;
}

.operation-buttons .el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

/* 分页样式 */
.pagination-section {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

/* 表单提示样式 */
.form-tip {
  margin-left: 10px;
  color: #909399;
  font-size: 12px;
}

/* 标签样式优化 */
.el-tag {
  font-weight: 500;
  border-radius: 8px;
  border: none;
}

/* 按钮样式优化 */
.el-button {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* 输入框和选择器样式 */
.el-input :deep(.el-input__wrapper),
.el-select :deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(103, 126, 234, 0.2);
  border-radius: 8px;
  box-shadow: none;
  transition: all 0.3s ease;
}

.el-input :deep(.el-input__wrapper):hover,
.el-select :deep(.el-input__wrapper):hover {
  border-color: #c0c4cc;
}

.el-input :deep(.el-input__wrapper.is-focus),
.el-select :deep(.el-input__wrapper.is-focus) {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(103, 126, 234, 0.2);
  background: rgba(255, 255, 255, 1);
}

.el-input :deep(.el-input__inner),
.el-select :deep(.el-input__inner) {
  background: transparent;
  border: none;
  color: #2c3e50;
}

/* 加载遮罩样式 */
:deep(.el-loading-mask) {
  background-color: rgba(103, 126, 234, 0.1);
  backdrop-filter: blur(4px);
}

/* 对话框样式 */
.el-dialog {
  border-radius: 12px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* 命令结果样式 */
.command-result :deep(.el-textarea__inner) {
  font-family: 'Monaco', 'Courier New', monospace;
  background-color: #1e1e1e;
  color: #d4d4d4;
  border: none;
  font-size: 13px;
  line-height: 1.6;
}

/* Windows Exporter 管理样式 */
.exporter-management {
  padding: 10px 0;
}

.exporter-management .status-card {
  margin-bottom: 20px;
  border-radius: 12px;
}

.exporter-management .action-buttons {
  display: flex;
  gap: 12px;
  justify-content: center;
  margin-top: 20px;
}

/* Switch开关样式 */
:deep(.el-switch.is-checked .el-switch__core) {
  background-color: #67c23a;
}

/* 响应式布局 */
@media (max-width: 768px) {
  .cmdb-snmp-management {
    padding: 15px;
  }

  .title {
    font-size: 18px;
  }

  .search-section {
    padding: 16px;
  }

  .search-form .el-form-item {
    margin-bottom: 8px;
  }
}
</style>
