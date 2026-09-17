<template>
  <!-- 头像组件 -->
  <div class="user-info">
    <!-- 头像 -->
    <img
        :src="sysAdmin.icon || require('./../assets/image/touxiang.jpg')"
        alt="头像"
        class="user-avatar"
        @error="useDefaultAvatar"
    />

    <!-- 用户名下拉菜单 -->
    <el-dropdown trigger="click" @command="handleCommand">
      <span class="user-username-trigger">
        <span class="user-username">{{ sysAdmin.username }}</span>
        <el-icon class="dropdown-arrow"><ArrowDown /></el-icon>
      </span>

      <!-- 下拉菜单 -->
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item command="license">系统激活</el-dropdown-item>
          <el-dropdown-item command="profile">个人信息</el-dropdown-item>
          <el-dropdown-item command="logout">退出登录</el-dropdown-item>
          <el-dropdown-item disabled class="version-item">AuthOps:v4.0</el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>

    <el-dialog
      v-model="licenseDialogVisible"
      title="系统激活"
      width="640px"
      :close-on-click-modal="false"
    >
      <div class="license-dialog-content">
        <div class="license-status-card">
          <div class="license-status-row">
            <span class="license-label">当前状态</span>
            <el-tag :type="licenseStatusTagType" size="small">
              {{ licenseStatusText }}
            </el-tag>
          </div>
          <div class="license-status-row" v-if="licenseInfo.customer || licenseInfo.customerName">
            <span class="license-label">客户名称</span>
            <span class="license-value">{{ licenseInfo.customer || licenseInfo.customerName }}</span>
          </div>
          <div class="license-status-row" v-if="licenseInfo.licenseType">
            <span class="license-label">许可证类型</span>
            <span class="license-value">{{ licenseInfo.licenseType }}</span>
          </div>
          <div class="license-status-row" v-if="licenseInfo.validityText">
            <span class="license-label">授权状态</span>
            <span class="license-value">{{ licenseInfo.validityText }}</span>
          </div>
          <div class="license-status-row" v-if="licenseInfo.expireAt">
            <span class="license-label">到期时间</span>
            <span class="license-value">{{ licenseInfo.expireAt }}</span>
          </div>
        </div>

        <el-form label-width="90px">
          <el-form-item label="机器码">
            <div class="machine-code-box">
              {{ machineCode || '正在获取机器码' }}
            </div>
          </el-form-item>

          <el-form-item label="激活码">
            <el-input
              v-model="activationCode"
              type="textarea"
              :rows="8"
              placeholder="请输入离线激活码"
            />
          </el-form-item>
        </el-form>
      </div>

      <template #footer>
        <el-button @click="licenseDialogVisible = false">取消</el-button>
        <el-button @click="copyMachineCode">复制机器码</el-button>
        <el-button type="primary" :loading="licenseSubmitting" @click="submitActivation">
          激活
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import storage from '@/utils/storage'
import { ArrowDown } from '@element-plus/icons-vue'
import { licenseApi } from '@/api/license'

export default {
  name: "HeadImage",
  components: {
    ArrowDown
  },
  data() {
    return {
      sysAdmin: storage.getItem("sysAdmin") || {},
      licenseDialogVisible: false,
      licenseLoading: false,
      licenseSubmitting: false,
      machineCode: '',
      activationCode: '',
      licenseInfo: {}
    };
  },
  computed: {
    licenseStatusText() {
      const status = String(this.licenseInfo.status || '').toLowerCase()
      if (status === 'activated' || status === 'active' || this.licenseInfo.activated === true) return '已激活'
      if (status === 'expired' || this.licenseInfo.expired === true) return '已过期'
      return '未激活'
    },
    licenseStatusTagType() {
      const text = this.licenseStatusText
      if (text === '已激活') return 'success'
      if (text === '已过期') return 'danger'
      return 'warning'
    }
  },
  methods: {
    // 图片加载失败时使用默认头像
    useDefaultAvatar(e) {
      e.target.src = require('./../assets/image/touxiang.jpg');
    },
    // 菜单点击事件
    handleCommand(command) {
      if (command === 'logout') {
        // 调用新的 logout 方法
        this.logout();
      } else if (command === 'profile') {
        this.$router.push('/system/personal'); // 跳转到个人页面路由
      } else if (command === 'license') {
        this.openLicenseDialog()
      }
    },
    async openLicenseDialog() {
      this.licenseDialogVisible = true
      await Promise.allSettled([
        this.fetchLicenseStatus(),
        this.fetchMachineCode()
      ])
    },
    async fetchLicenseStatus() {
      this.licenseLoading = true
      try {
        const { data: res } = await licenseApi.getStatus()
        if (res.code === 200) {
          this.licenseInfo = res.data || {}
        } else {
          this.$message.error(res.message || '获取激活状态失败')
        }
      } catch (error) {
        console.error('获取激活状态失败:', error)
      } finally {
        this.licenseLoading = false
      }
    },
    async fetchMachineCode() {
      try {
        const { data: res } = await licenseApi.getMachineCode()
        if (res.code === 200) {
          this.machineCode = res.data?.machineCode || res.data?.code || res.data || ''
        } else {
          this.$message.error(res.message || '获取机器码失败')
        }
      } catch (error) {
        console.error('获取机器码失败:', error)
      }
    },
    async copyMachineCode() {
      if (!this.machineCode) {
        this.$message.warning('机器码为空')
        return
      }

      try {
        if (navigator.clipboard && window.isSecureContext) {
          await navigator.clipboard.writeText(this.machineCode)
        } else {
          const textArea = document.createElement('textarea')
          textArea.value = this.machineCode
          textArea.setAttribute('readonly', 'readonly')
          textArea.style.position = 'fixed'
          textArea.style.left = '-9999px'
          textArea.style.top = '0'
          document.body.appendChild(textArea)
          textArea.focus()
          textArea.select()

          const copied = document.execCommand('copy')
          document.body.removeChild(textArea)

          if (!copied) {
            throw new Error('execCommand copy failed')
          }
        }
        this.$message.success('机器码已复制')
      } catch (error) {
        console.error('复制机器码失败:', error)
        this.$message.error('复制机器码失败')
      }
    },
    async submitActivation() {
      const licenseKey = this.activationCode.trim()
      if (!licenseKey) {
        this.$message.warning('请输入激活码')
        return
      }

      this.licenseSubmitting = true
      try {
        const { data: res } = await licenseApi.activate({ licenseKey })
        if (res.code === 200) {
          this.$message.success(res.message || '激活成功')
          this.activationCode = ''
          await this.fetchLicenseStatus()
          this.licenseDialogVisible = false
        } else {
          this.$message.error(res.message || '激活失败')
        }
      } catch (error) {
        console.error('激活失败:', error)
      } finally {
        this.licenseSubmitting = false
      }
    },
    // 退出登录
    async logout() {
      const confirmResult = await this.$confirm('确定要退出登录吗, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).catch(err => err)

      if (confirmResult !== 'confirm') {
        return this.$message.info('已取消退出')
      }
      // 清除本地存储并跳转到登录页
      this.$storage.clearAll()
      this.$router.push('/login')
      this.$message.success('退出成功')
    }
  }
};
</script>

<style lang="less" scoped>
.user-info {
  position: fixed;
  right: 30px;
  top: 10px;
  display: flex;
  align-items: center;
  gap: 10px;
  z-index: 9999 !important;
}

.user-username {
  font-size: small;
  z-index: 9999;
}

.user-username-trigger {
  display: flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  padding: 6px 10px;
  border-radius: 6px;
  transition: all 0.3s ease;

  &:hover {
    background: rgba(102, 126, 234, 0.1);
  }

  .user-username {
    font-size: small;
    color: #2c3e50;
    font-weight: 500;
  }

  .dropdown-arrow {
    font-size: 12px;
    color: #909399;
    transition: transform 0.3s ease;
  }
}

.user-username-trigger:hover .dropdown-arrow {
  color: #667eea;
}

.user-avatar {
  cursor: pointer;
  width: 40px;
  height: 40px;
  border-radius: 50%; /* 圆形头像 */
  object-fit: cover; /* 防止变形 */
}

.license-dialog-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.license-status-card {
  padding: 14px 16px;
  border-radius: 10px;
  background: #f8fafc;
  border: 1px solid #ebeef5;
}

.machine-code-box {
  width: 100%;
  min-height: 96px;
  padding: 12px 14px;
  border-radius: 8px;
  border: 1px solid #b7eb8f;
  background: #f6ffed;
  color: #389e0d;
  font-family: 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-all;
  box-sizing: border-box;
}

.license-status-row {
  display: flex;
  align-items: center;
  gap: 12px;
  min-height: 28px;
}

.license-status-row + .license-status-row {
  margin-top: 8px;
}

.license-label {
  width: 72px;
  color: #606266;
  flex-shrink: 0;
}

.license-value {
  color: #303133;
  word-break: break-all;
}

:deep(.version-item) {
  color: #c0c4cc !important;
  font-size: 12px;
  cursor: not-allowed;
  text-align: center;
  padding-top: 8px;
  border-top: 1px solid #ebeef5;
}
</style>
