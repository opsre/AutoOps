<template>
  <div class="ssl-cert-apply-container">
    <!-- 搜索和操作区域 -->
    <div class="search-section">
      <el-form :inline="true" :model="sslCertSearchForm" class="search-form">
        <el-form-item label="域名">
          <el-input
            v-model="sslCertSearchForm.domain"
            placeholder="搜索域名"
            size="small"
            clearable
            @keyup.enter="handleSearch"
          />
        </el-form-item>

        <el-form-item label="证书状态">
          <el-select
            v-model="sslCertSearchForm.status"
            placeholder="证书状态"
            size="small"
            clearable
            style="width: 120px"
          >
            <el-option label="待验证" :value="0" />
            <el-option label="已颁发" :value="1" />
            <el-option label="已过期" :value="2" />
          </el-select>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" size="small" @click="handleSearch" class="modern-btn primary-btn">
            <el-icon><Search /></el-icon>
            查询
          </el-button>
          <el-button size="small" @click="handleReset" class="modern-btn reset-btn">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 操作按钮区域 -->
    <div class="action-section">
      <el-button
        type="success"
        size="small"
        @click="handleOpenSslCertDialog"
        class="modern-btn success-btn"
      >
        <el-icon><Plus /></el-icon>
        申请证书
      </el-button>

      <el-button
        type="primary"
        size="small"
        @click="handleOpenUploadCertDialog"
        class="modern-btn primary-btn"
      >
        <el-icon><Upload /></el-icon>
        上传证书
      </el-button>

      <el-button
        size="small"
        @click="fetchSslCertList"
        class="modern-btn info-btn"
      >
        <el-icon><Refresh /></el-icon>
        刷新
      </el-button>
    </div>

    <!-- 列表 -->
    <div class="table-section">
      <el-table
        :data="sslCertList"
        border
        stripe
        :loading="sslCertLoading"
        class="cert-table"
        :header-cell-style="{ background: '#eef1f6', color: '#606266' }"
      >
        <el-table-column prop="id" label="ID" min-width="60" align="center" />
        <el-table-column prop="domain" label="域名" min-width="140" show-overflow-tooltip />
        <el-table-column prop="configName" label="阿里云配置" min-width="100" align="center" show-overflow-tooltip>
          <template #default="{ row }">
            <span style="font-size: 12px">{{ row.configName || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="certSource" label="证书来源" min-width="110" align="center">
          <template #default="{ row }">
            <el-tag size="small" type="info">{{ row.certSourceText || (row.certType === 1 ? 'ACME' : '阿里云CAS') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="caProvider" label="CA提供商" min-width="110" align="center" show-overflow-tooltip>
          <template #default="{ row }">
            <span style="font-size: 12px">{{ row.caProvider || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="algorithm" label="加密算法" min-width="90" align="center">
          <template #default="{ row }">
            <el-tag size="small" type="success">{{ row.algorithm || '-' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" min-width="90" align="center">
          <template #default="{ row }">
            <el-tag :type="getCertStatusType(row.status)" size="small">{{ row.statusText || '未知' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="expireTime" label="过期时间" min-width="140" align="center">
          <template #default="{ row }">
            <div style="font-size: 12px; line-height: 1.5">
              <div>{{ formatExpireTime(row.expireTime) }}</div>
              <el-tag v-if="row.daysLeft !== undefined && row.daysLeft > 0" :type="getDaysLeftType(row.daysLeft)" size="small" style="margin-top: 4px;">
                剩余 {{ row.daysLeft }} 天
              </el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="320" fixed="right" align="center">
          <template #default="{ row }">
            <div class="operation-buttons">
              <el-tooltip content="下载证书" placement="top">
                <el-button
                  type="success"
                  :icon="Download"
                  size="small"
                  @click="handleDownloadCert(row)"
                >下载证书</el-button>
              </el-tooltip>
              <el-tooltip content="部署" placement="top">
                <el-button
                  type="primary"
                  :icon="Upload"
                  size="small"
                  @click="handleDeployCert(row)"
                >部署</el-button>
              </el-tooltip>
              <el-tooltip content="部署日志" placement="top">
                <el-button
                  type="info"
                  :icon="Tickets"
                  size="small"
                  @click="handleViewDeployLog(row)"
                >部署日志</el-button>
              </el-tooltip>
              <el-tooltip content="删除" placement="top">
                <el-button
                  type="danger"
                  :icon="Delete"
                  size="small"
                  @click="handleDeleteSslCert(row)"
                >删除</el-button>
              </el-tooltip>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 分页 -->
    <div class="pagination-section">
      <el-pagination
        v-model:current-page="sslCertSearchForm.pageNum"
        v-model:page-size="sslCertSearchForm.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="sslCertTotal"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="handlePageChange"
        @size-change="handleSizeChange"
      />
    </div>

    <!-- SSL证书申请对话框 -->
    <ssl-cert-dialog
      v-model:visible="sslCertDialogVisible"
      :title="sslCertDialogTitle"
      :form="sslCertForm"
      :loading="sslCertLoading"
      @submit="handleApplySslCert"
    />

    <!-- 证书部署对话框 -->
    <deploy-dialog
      v-model:visible="deployDialogVisible"
      :cert="currentDeployCert"
      :form="deployForm"
      :loading="deployLoading"
      :host-options="hostOptions"
      :logs="deployLogs"
      @submit="handleConfirmDeploy"
      @close="handleCloseDeployDialog"
    />

    <!-- 部署日志对话框 -->
    <deploy-log-dialog
      v-model:visible="deployLogDialogVisible"
      :cert="currentLogCert"
      :logs="deployLogList"
      :total="deployLogTotal"
      :page="deployLogPage"
      :page-size="deployLogPageSize"
      :loading="deployLogLoading"
      @page-change="handleDeployLogPageChange"
      @size-change="handleDeployLogPageSizeChange"
      @close="handleCloseDeployLogDialog"
    />

    <!-- 上传证书对话框 -->
    <upload-cert-dialog
      v-model:visible="uploadCertDialogVisible"
      :loading="uploadCertLoading"
      :file-list="uploadCertFileList"
      @file-change="handleUploadCertFileChange"
      @file-remove="handleUploadCertFileRemove"
      @submit="handleSubmitUploadCert"
    />
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Upload,
  Refresh,
  Search,
  Download,
  Delete,
  Tickets
} from '@element-plus/icons-vue'
import { sslCertApi } from '@/api/monitor'
import cmdbApi from '@/api/cmdb'
import SslCertDialog from './SslCertDialog.vue'
import DeployDialog from './DeployDialog.vue'
import DeployLogDialog from './DeployLogDialog.vue'
import UploadCertDialog from './UploadCertDialog.vue'

// SSL证书搜索表单
const sslCertSearchForm = reactive({
  domain: '',
  status: -1,
  pageNum: 1,
  pageSize: 10
})

// SSL证书表单
const sslCertForm = reactive({
  domain: '',
  certType: 1,
  aliyunConfigId: null,
  email: '',
  caProvider: 'LetsEncrypt',
  algorithm: 'RSA2048',
  certSource: '',
  validateType: 'DNS'
})

// SSL证书列表状态
const sslCertList = ref([])
const sslCertTotal = ref(0)
const sslCertLoading = ref(false)
const sslCertDialogVisible = ref(false)
const sslCertDialogTitle = ref('申请SSL证书')

// 证书部署相关
const deployDialogVisible = ref(false)
const deployLoading = ref(false)
const deployForm = reactive({
  certId: null,
  hostId: null,
  deployPath: '/etc/nginx/ssl'
})
const hostOptions = ref([])
const deployLogs = ref([])
const currentDeployCert = ref(null)

// 部署日志相关
const deployLogDialogVisible = ref(false)
const deployLogLoading = ref(false)
const deployLogList = ref([])
const deployLogTotal = ref(0)
const deployLogPage = ref(1)
const deployLogPageSize = ref(10)
const currentLogCert = ref(null)

// 上传证书相关
const uploadCertDialogVisible = ref(false)
const uploadCertLoading = ref(false)
const uploadCertFileList = ref([])

// 判断是否ACME方式
const isACME = computed(() => {
  return sslCertForm.certType === 1
})

// 获取SSL证书列表
const fetchSslCertList = async () => {
  sslCertLoading.value = true
  try {
    const { data: res } = await sslCertApi.getList(sslCertSearchForm)
    if (res.code === 200) {
      sslCertList.value = Array.isArray(res.data) ? res.data : res.data?.list || []
      sslCertTotal.value = res.data?.total || sslCertList.value.length
    } else {
      ElMessage.error(res.message || '获取证书列表失败')
    }
  } catch (error) {
    console.error('获取证书列表失败:', error)
    ElMessage.error('获取证书列表失败')
  } finally {
    sslCertLoading.value = false
  }
}

// 获取主机列表
const fetchHostList = async () => {
  try {
    const response = await cmdbApi.getCmdbHostList({
      page: 1,
      pageSize: 100
    })
    if (response.data?.code === 200) {
      const hostList = response.data.data?.list || []
      hostOptions.value = hostList.map(host => ({
        value: host.id,
        label: `${host.hostName || host.name} (${host.privateIp || host.publicIp || host.sshIp})`
      }))
    }
  } catch (error) {
    console.error('获取主机列表失败:', error)
    ElMessage.error('获取主机列表失败')
  }
}

// 获取部署日志
const fetchDeployLogs = async () => {
  deployLogLoading.value = true
  try {
    const { data: res } = await sslCertApi.getDeployLogs({
      certId: currentLogCert.value.id,
      page: deployLogPage.value,
      pageSize: deployLogPageSize.value
    })
    if (res.code === 200) {
      deployLogList.value = res.data?.list || []
      deployLogTotal.value = res.data?.total || 0
    }
  } catch (error) {
    console.error('获取部署日志失败:', error)
    ElMessage.error('获取部署日志失败')
  } finally {
    deployLogLoading.value = false
  }
}

// 搜索
const handleSearch = () => {
  sslCertSearchForm.pageNum = 1
  fetchSslCertList()
}

// 重置
const handleReset = () => {
  sslCertSearchForm.domain = ''
  sslCertSearchForm.status = -1
  sslCertSearchForm.pageNum = 1
  fetchSslCertList()
}

// 分页变化
const handlePageChange = (page) => {
  sslCertSearchForm.pageNum = page
  fetchSslCertList()
}

const handleSizeChange = (size) => {
  sslCertSearchForm.pageSize = size
  fetchSslCertList()
}

// 打开SSL证书申请对话框
const handleOpenSslCertDialog = () => {
  sslCertDialogTitle.value = '申请SSL证书'
  sslCertForm.domain = ''
  sslCertForm.certType = 1
  sslCertForm.aliyunConfigId = null
  sslCertForm.email = ''
  sslCertForm.caProvider = 'LetsEncrypt'  // 默认 Let's Encrypt
  sslCertForm.algorithm = 'RSA2048'       // 默认 RSA 2048
  sslCertDialogVisible.value = true
}

// 申请SSL证书
const handleApplySslCert = async () => {
  // 验证必填字段
  if (!sslCertForm.aliyunConfigId) {
    ElMessage.warning('请选择阿里云配置')
    return
  }

  if (isACME.value && !sslCertForm.email) {
    ElMessage.warning('ACME方式申请证书需要填写邮箱')
    return
  }

  sslCertLoading.value = true
  try {
    // 构建请求数据
    const data = {
      domain: sslCertForm.domain,
      certType: sslCertForm.certType,
      aliyunConfigId: sslCertForm.aliyunConfigId
    }

    // ACME方式需要的字段，使用默认值
    if (isACME.value) {
      data.email = sslCertForm.email
      data.caProvider = 'LetsEncrypt'  // 默认使用 Let's Encrypt
      data.algorithm = 'RSA2048'       // 默认使用 RSA2048
      data.validateType = 'DNS'        // 默认使用 DNS 验证
    }

    const { data: res } = await sslCertApi.apply(data)
    if (res.code === 200) {
      const certId = res.data.certificateId || res.data.id
      const domain = res.data.domain || data.domain

      ElMessage.success({
        message: `证书申请已提交！证书ID: ${certId}`,
        duration: 3000
      })

      ElMessage.info({
        message: `域名 ${domain} 的证书正在后台处理中，请稍后刷新列表查看状态`,
        duration: 5000
      })

      sslCertDialogVisible.value = false

      // 重置表单
      sslCertForm.domain = ''
      sslCertForm.certType = 1
      sslCertForm.aliyunConfigId = null
      sslCertForm.email = ''

      fetchSslCertList()
      setTimeout(() => {
        fetchSslCertList()
      }, 10000)
    } else {
      ElMessage.error(res.message || '证书申请失败')
    }
  } catch (error) {
    console.error('证书申请失败:', error)
    ElMessage.error(error.response?.data?.message || '证书申请失败')
  } finally {
    sslCertLoading.value = false
  }
}

// 下载证书
const handleDownloadCert = async (row) => {
  try {
    const { data } = await sslCertApi.download(row.id, {
      responseType: 'blob'
    })
    const url = window.URL.createObjectURL(new Blob([data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `${row.domain}.zip`)
    document.body.appendChild(link)
    link.click()
    link.parentNode.removeChild(link)
  } catch (error) {
    console.error('下载证书失败:', error)
    ElMessage.error('下载证书失败')
  }
}

// 删除证书
const handleDeleteSslCert = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该证书吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const { data: res } = await sslCertApi.delete(row.id)
    if (res.code === 200) {
      ElMessage.success('删除成功')
      fetchSslCertList()
    } else {
      ElMessage.error(res.message || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 打开部署对话框
const handleDeployCert = async (row) => {
  currentDeployCert.value = row
  deployForm.certId = row.id
  deployForm.hostId = null
  deployForm.deployPath = '/etc/nginx/ssl'
  deployLogs.value = []

  await fetchHostList()
  deployDialogVisible.value = true
}

// 确认部署
const handleConfirmDeploy = async () => {
  try {
    deployLoading.value = true
    const { data: res } = await sslCertApi.deploy({
      certId: deployForm.certId,
      hostId: deployForm.hostId,
      deployPath: deployForm.deployPath
    })

    if (res.code === 200) {
      deployLogs.value = res.data?.logs || []
      ElMessage.success(res.data?.message || '证书部署成功')
      setTimeout(() => {
        deployDialogVisible.value = false
      }, 2000)
    } else {
      deployLogs.value = res.data?.logs || []
      ElMessage.error(res.message || '证书部署失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('证书部署失败:', error)
      deployLogs.value = error.response?.data?.data?.logs || []
      ElMessage.error(error.response?.data?.message || '证书部署失败')
    }
  } finally {
    deployLoading.value = false
  }
}

// 关闭部署对话框
const handleCloseDeployDialog = () => {
  deployDialogVisible.value = false
  deployForm.certId = null
  deployForm.hostId = null
  deployForm.deployPath = '/etc/nginx/ssl'
  deployLogs.value = []
  currentDeployCert.value = null
}

// 查看部署日志
const handleViewDeployLog = async (cert) => {
  currentLogCert.value = cert
  deployLogPage.value = 1
  deployLogPageSize.value = 10
  deployLogDialogVisible.value = true
  await fetchDeployLogs()
}

// 部署日志分页
const handleDeployLogPageChange = (page) => {
  deployLogPage.value = page
  fetchDeployLogs()
}

const handleDeployLogPageSizeChange = (size) => {
  deployLogPageSize.value = size
  fetchDeployLogs()
}

// 关闭部署日志对话框
const handleCloseDeployLogDialog = () => {
  deployLogDialogVisible.value = false
  deployLogList.value = []
  currentLogCert.value = null
}

// 打开上传证书对话框
const handleOpenUploadCertDialog = () => {
  uploadCertFileList.value = []
  uploadCertDialogVisible.value = true
}

// 文件变化处理
const handleUploadCertFileChange = (file, fileList) => {
  uploadCertFileList.value = fileList
}

// 文件移除处理
const handleUploadCertFileRemove = (file, fileList) => {
  uploadCertFileList.value = fileList
}

// 提交上传证书
const handleSubmitUploadCert = async () => {
  if (uploadCertFileList.value.length === 0) {
    ElMessage.warning('请选择要上传的证书文件')
    return
  }

  uploadCertLoading.value = true
  try {
    const formData = new FormData()
    formData.append('file', uploadCertFileList.value[0].raw)

    const { data: res } = await sslCertApi.upload(formData)

    if (res.code === 200) {
      ElMessage.success(res.data?.message || '证书上传成功')

      if (res.data?.certId && res.data?.domain) {
        ElMessage.success(`证书上传成功！域名: ${res.data.domain}, 证书ID: ${res.data.certId}`)
      }

      uploadCertDialogVisible.value = false
      uploadCertFileList.value = []
      fetchSslCertList()
    } else {
      ElMessage.error(res.message || '证书上传失败')
    }
  } catch (error) {
    console.error('证书上传失败:', error)
    ElMessage.error(error.response?.data?.message || '证书上传失败')
  } finally {
    uploadCertLoading.value = false
  }
}

// 工具函数
const formatDate = (dateObj) => {
  if (!dateObj) return '-'

  if (typeof dateObj === 'object' && dateObj.hasOwnProperty('Valid')) {
    if (!dateObj.Valid || !dateObj.Time) return '-'
    const dateStr = dateObj.Time
    if (dateStr.startsWith('0001-01-01')) return '-'
    return new Date(dateStr).toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  if (typeof dateObj === 'string') {
    if (!dateObj || dateObj.startsWith('0001-01-01')) return '-'
    return new Date(dateObj).toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  return '-'
}

// 格式化过期时间
const formatExpireTime = (expireTime) => {
  if (!expireTime || expireTime.startsWith('0001-01-01')) return '-'
  try {
    return new Date(expireTime).toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch (error) {
    return expireTime
  }
}

// 根据剩余天数获取标签颜色
const getDaysLeftType = (daysLeft) => {
  if (daysLeft <= 0) return 'danger'     // 已过期 - 红色
  if (daysLeft <= 7) return 'danger'     // 7天内 - 红色
  if (daysLeft <= 30) return 'warning'   // 30天内 - 橙色
  return 'success'                        // 正常 - 绿色
}

const getCertStatusType = (status) => {
  // 根据后端返回的状态码映射颜色
  // 0=待验证, 1=申请中, 2=已颁发, 3=已过期, 4=申请失败
  switch (status) {
    case 0:
      return 'warning'  // 待验证 - 橙色
    case 1:
      return 'info'     // 申请中 - 蓝色
    case 2:
      return 'success'  // 已颁发 - 绿色
    case 3:
      return 'danger'   // 已过期 - 红色
    case 4:
      return 'danger'   // 申请失败 - 红色
    default:
      return 'info'
  }
}

const getCertStatusText = (status) => {
  // 这个函数已不再使用，保留以防万一
  switch (status) {
    case 0:
      return '待验证'
    case 1:
      return '申请中'
    case 2:
      return '已颁发'
    case 3:
      return '已过期'
    case 4:
      return '申请失败'
    default:
      return '未知'
  }
}

// 组件挂载
onMounted(() => {
  fetchSslCertList()
})
</script>

<style scoped>
.ssl-cert-apply-container {
  padding: 0;
}

.search-section {
  margin-bottom: 16px;
  padding: 16px;
  background: rgba(103, 126, 234, 0.05);
  border-radius: 12px;
  border: 1px solid rgba(103, 126, 234, 0.1);
}

.search-form {
  margin: 0;
}

.search-form :deep(.el-form-item) {
  margin-bottom: 0;
  margin-right: 16px;
}

.search-form :deep(.el-form-item__label) {
  color: #606266;
  font-weight: 500;
}

.search-form :deep(.el-input__wrapper) {
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.search-form :deep(.el-select .el-input__wrapper) {
  border-radius: 8px;
}

.action-section {
  margin-bottom: 16px;
}

:deep(.modern-btn) {
  border-radius: 8px;
  padding: 8px 20px;
  font-weight: 500;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: none;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

:deep(.modern-btn:hover) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.2);
}

:deep(.primary-btn) {
  background: linear-gradient(45deg, #409EFF, #66B3FF);
  color: white;
}

:deep(.success-btn) {
  background: linear-gradient(45deg, #67C23A, #85CE61);
  color: white;
}

:deep(.reset-btn) {
  background: linear-gradient(45deg, #E6A23C, #EEBE77);
  color: white;
}

:deep(.info-btn) {
  background: linear-gradient(45deg, #909399, #b1b3b8);
  color: white;
}

.table-section {
  margin-top: 0;
}

.cert-table {
  background: rgba(255, 255, 255, 0.9);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  font-size: 12px;
}

.cert-table :deep(.el-table__header) {
  background: rgba(102, 126, 234, 0.1) !important;
}

.cert-table :deep(.el-table__header th) {
  background: rgba(102, 126, 234, 0.1) !important;
  color: #2c3e50 !important;
  font-weight: 600;
  border: none;
  padding: 8px 0;
  font-size: 13px;
}

.cert-table :deep(.el-table__body tr:hover > td) {
  background-color: rgba(102, 126, 234, 0.08) !important;
}

.cert-table :deep(.el-table td) {
  border: none;
  padding: 6px 0;
}

.cert-table :deep(.el-table .cell) {
  padding: 0 8px;
  line-height: 1.4;
}

.cert-table :deep(.el-table::before) {
  display: none;
}

.cert-table :deep(.el-table--border::after) {
  display: none;
}

/* 确保操作列固定在最上层 */
.cert-table :deep(.el-table__fixed-right) {
  z-index: 3 !important;
}

.cert-table :deep(.el-table__fixed-right-patch) {
  z-index: 3 !important;
}

/* 操作按钮样式 */
.operation-buttons {
  display: flex;
  gap: 6px;
  justify-content: center;
  align-items: center;
}

.operation-buttons .el-button {
  transition: all 0.3s ease;
}

.operation-buttons .el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

.cert-table :deep(.el-button.is-link) {
  font-weight: 500;
  transition: all 0.3s ease;
}

.cert-table :deep(.el-button.is-link:hover) {
  transform: translateY(-1px);
}

.pagination-section {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

:deep(.el-pagination) {
  padding: 0;
}

:deep(.el-pagination button),
:deep(.el-pagination .el-pager li) {
  border-radius: 8px;
  transition: all 0.3s ease;
}

:deep(.el-pagination button:hover),
:deep(.el-pagination .el-pager li:hover) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

:deep(.el-pagination .el-pager li.is-active) {
  background: linear-gradient(45deg, #409EFF, #66B3FF);
  color: white;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
}
</style>
