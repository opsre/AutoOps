<template>
  <div class="personal-management">
    <el-row :gutter="20">
      <el-col :xs="24" :sm="24" :md="8" :lg="6" :xl="5" style="margin-bottom: 20px">
        <el-card shadow="hover" class="profile-card">
          <template #header>
            <div class="card-header">
              <span class="title">个人信息</span>
            </div>
          </template>
          
          <div class="profile-content">
            <div class="avatar-section">
              <el-avatar :src="adminDetail.icon" :size="130" :key="avatarKey" class="profile-avatar"></el-avatar>
            </div>
            <el-form label-width="100px" size="small" class="profile-form">
              <el-row>
                <el-col :span="24">
                  <el-form-item label="用户账号：">
                    <span class="profile-value">{{ adminDetail.username }}</span>
                  </el-form-item>
                </el-col>
                <el-col :span="24">
                  <el-form-item label="用户昵称：">
                    <span class="profile-value">{{ adminDetail.nickname }}</span>
                  </el-form-item>
                </el-col>
                <el-col :span="24">
                  <el-form-item label="用户邮箱：">
                    <span class="profile-value">{{ adminDetail.email }}</span>
                  </el-form-item>
                </el-col>
                <el-col :span="24">
                  <el-form-item label="用户电话：">
                    <span class="profile-value">{{ adminDetail.phone }}</span>
                  </el-form-item>
                </el-col>
                <el-col :span="24">
                  <el-form-item label="用户备注：">
                    <span class="profile-value">{{ adminDetail.note }}</span>
                  </el-form-item>
                </el-col>
                <el-col :span="24">
                  <el-form-item label="创建时间：">
                    <span class="profile-value">{{ adminDetail.createTime }}</span>
                  </el-form-item>
                </el-col>
              </el-row>
            </el-form>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :sm="24" :md="16" :lg="18" :xl="19">
        <el-card shadow="hover" class="detail-card">
          <template #header>
            <div class="card-header">
              <span class="title">基本资料</span>
            </div>
          </template>
          
          <div class="tabs-section">
            <el-tabs v-model="activeName" class="personal-tabs">
              <el-tab-pane label="基本资料" name="first">
                <el-form :model="adminDetail" :rules="editFormRules" ref="editFormRefForm" label-width="100px" class="edit-form">
                  <el-form-item label="用户头像" prop="icon">
                    <el-upload :headers="headers" class="avatar-uploader" :action="uploadIconUrl"
                               :show-file-list="false" :on-success="handleAvatarSuccess">
                      <img v-if="icon" :src="icon" :key="avatarKey" class="avatar" title="点击更换头像">
                      <img v-else :src="adminDetail.icon" :key="avatarKey" class="avatar" title="点击更换头像">
                    </el-upload>
                  </el-form-item>
                  <el-form-item label="用户账号" prop="username">
                    <el-input v-model="adminDetail.username" size="small" />
                  </el-form-item>
                  <el-form-item label="用户昵称" prop="nickname">
                    <el-input v-model="adminDetail.nickname" size="small" />
                  </el-form-item>
                  <el-form-item label="手机号码" prop="phone">
                    <el-input v-model="adminDetail.phone" maxlength="11" size="small" />
                  </el-form-item>
                  <el-form-item label="用户邮箱" prop="email">
                    <el-input v-model="adminDetail.email" maxlength="50" size="small" />
                  </el-form-item>
                  <el-form-item label="用户备注" prop="note">
                    <el-input v-model="adminDetail.note" size="small" />
                  </el-form-item>
                  <el-form-item class="form-actions">
                    <el-button type="primary" size="small" @click="submitFirst">保存</el-button>
                    <el-button size="small" @click="closeFirst">关闭</el-button>
                  </el-form-item>
                </el-form>
              </el-tab-pane>
              <el-tab-pane label="修改密码" name="second">
                <el-form :model="updateForm" :rules="updateFormRules" ref="updateFormRefForm" label-width="100px" class="edit-form">
                  <el-form-item label="旧密码" prop="password">
                    <el-input v-model="updateForm.password" placeholder="请输入旧密码" type="password" size="small" />
                  </el-form-item>
                  <el-form-item label="新密码" prop="newPassword">
                    <el-input v-model="updateForm.newPassword" placeholder="请输入新密码" type="password" size="small" />
                  </el-form-item>
                  <el-form-item label="确认密码" prop="resetPassword">
                    <el-input v-model="updateForm.resetPassword" placeholder="请确认密码" type="password" size="small" />
                  </el-form-item>
                  <el-form-item class="form-actions">
                    <el-button type="primary" size="small" @click="submitSecond">保存</el-button>
                    <el-button size="small" @click="closeSecond">关闭</el-button>
                  </el-form-item>
                </el-form>
              </el-tab-pane>
            </el-tabs>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import storage from '@/utils/storage'

export default {
  data() {
    return {
      adminDetail: storage.getItem("sysAdmin"),
      activeName: 'first',
      uploadIconUrl: '/api/v1/upload', // 直接使用字符串而不是调用方法
      headers: {
        Authorization: "Bearer " + storage.getItem("token"),
      },
      icon: '',
      avatarKey: Date.now(),
      editFormRules: {
        username: [{ required: true, message: '请输入用户账号', trigger: 'blur' }],
        email: [{ required: true, message: '请输入用户邮箱', trigger: 'blur' }],
        nickname: [{ required: true, message: '请输入用户昵称', trigger: 'blur' }],
        phone: [{ required: true, message: '请输入用户手机', trigger: 'blur' }],
        note: [{ required: true, message: '请输入用户备注', trigger: 'blur' }],
      },
      updateForm: {
        password: '',
        newPassword: '',
        resetPassword: ''
      },
      updateFormRules: {
        password: [{ required: true, message: '请输入旧密码', trigger: 'blur' }],
        newPassword: [{ required: true, message: '请输入新密码', trigger: 'blur' }],
        resetPassword: [{ required: true, message: '请输入重复密码', trigger: 'blur' }],
      }
    }
  },
  methods: {
    // 成功调用图片信息
    handleAvatarSuccess(res) {
      this.icon = res.data;
      this.adminDetail.icon = res.data;
      this.avatarKey = Date.now();
    },
    // 修改个人信息
    submitFirst() {
      this.$refs.editFormRefForm.validate(async valid => {
        if (!valid) return
        const { data: res } = await this.$api.adminUpdatePersonal({
          icon: this.icon === '' ? this.adminDetail.icon : this.icon,
          username: this.adminDetail.username,
          email: this.adminDetail.email,
          nickname: this.adminDetail.nickname,
          phone: this.adminDetail.phone,
          note: this.adminDetail.note,
        })
        if (res.code !== 200) {
          this.$message.error(res.message);
        } else {
          this.$storage.clearItem("sysAdmin")
          this.$store.commit('saveSysAdmin', res.data)
          this.adminDetail = res.data
          this.icon = ''
          this.avatarKey = Date.now()
          this.$message.success('修改用户成功')
        }
      })
    },
    closeFirst() {
      this.$router.push('/home')
    },
    // 修改个人密码
    submitSecond() {
      this.$refs.updateFormRefForm.validate(async valid => {
        if (!valid) return
        const { data: res } = await this.$api.adminUpdatePersonalPassword({
          password: this.updateForm.password,
          newPassword: this.updateForm.newPassword,
          resetPassword: this.updateForm.resetPassword
        })
        if (res.code !== 200) {
          this.$message.error(res.message);
        } else {
          this.$storage.clearAll()
          this.$router.push("/login")
          this.$message.success('修改密码成功')
        }
      })
    },
    // 关闭页面跳转到首页
    closeSecond() {
      this.$router.push('/home')
    }
  }
}
</script>

<style scoped>
.personal-management {
  padding: 20px;
  min-height: 100vh;
  background: #eef1f6;
}

.profile-card,
.detail-card {
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
   background-clip: text;
  -webkit-text-fill-color: transparent;
}

.profile-content {
  padding: 20px 0;
}

.avatar-section {
  text-align: center;
  margin-bottom: 24px;
}

.profile-avatar {
  border: 4px solid rgba(103, 126, 234, 0.2);
  box-shadow: 0 4px 16px rgba(103, 126, 234, 0.3);
  transition: all 0.3s ease;
}

.profile-avatar:hover {
  transform: scale(1.05);
  box-shadow: 0 6px 20px rgba(103, 126, 234, 0.4);
}

.profile-form .el-form-item {
  margin-bottom: 16px;
}

.profile-form .el-form-item__label {
  color: #606266;
  font-weight: 500;
}

.profile-value {
  color: #2c3e50;
  font-weight: 500;
  padding: 4px 8px;
  background: rgba(103, 126, 234, 0.1);
  border-radius: 6px;
  display: inline-block;
  min-width: 120px;
}

.tabs-section {
  padding: 20px 0;
}

.personal-tabs {
  background: transparent;
}

.personal-tabs :deep(.el-tabs__header) {
  margin-bottom: 24px;
}

.personal-tabs :deep(.el-tabs__nav-wrap::after) {
  background: linear-gradient(90deg, #667eea, #764ba2);
  height: 2px;
}

.personal-tabs :deep(.el-tabs__item) {
  color: #606266;
  font-weight: 500;
  font-size: 16px;
  transition: all 0.3s ease;
}

.personal-tabs :deep(.el-tabs__item:hover) {
  color: #667eea;
}

.personal-tabs :deep(.el-tabs__item.is-active) {
  color: #667eea;
  font-weight: 600;
}

.personal-tabs :deep(.el-tabs__active-bar) {
  background: linear-gradient(90deg, #667eea, #764ba2);
}

.edit-form {
  padding: 20px;
  background: rgba(103, 126, 234, 0.02);
  border-radius: 12px;
}

.edit-form .el-form-item {
  margin-bottom: 24px;
}

.edit-form .el-form-item__label {
  color: #606266;
  font-weight: 500;
}

.form-actions {
  text-align: center;
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid rgba(103, 126, 234, 0.1);
}

.form-actions .el-button {
  margin: 0 12px;
}

.avatar-uploader {
  border: 2px dashed rgba(103, 126, 234, 0.3);
  border-radius: 12px;
  cursor: pointer;
  width: 80px;
  height: 80px;
  position: relative;
  overflow: hidden;
  transition: all 0.3s ease;
  background: rgba(103, 126, 234, 0.05);
}

.avatar-uploader:hover {
  border-color: #667eea;
  background: rgba(103, 126, 234, 0.1);
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(103, 126, 234, 0.2);
}

.avatar-uploader .avatar {
  width: 80px;
  height: 80px;
  display: block;
  border-radius: 8px;
}

.el-button {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
  padding: 8px 24px;
}

.el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.el-input :deep(.el-input__inner) {
  border-radius: 8px;
  border: 1px solid rgba(103, 126, 234, 0.2);
  transition: all 0.3s ease;
  background: rgba(255, 255, 255, 0.8);
  color: #2c3e50;
}

.el-input :deep(.el-input__inner):focus {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(103, 126, 234, 0.2);
  background: rgba(255, 255, 255, 1);
}

.el-input :deep(.el-input__inner::placeholder) {
  color: rgba(44, 62, 80, 0.6);
}

.el-card :deep(.el-card__body) {
  padding: 24px;
}

@media (max-width: 768px) {
  .personal-management {
    padding: 10px;
  }
  
  .profile-avatar {
    width: 100px !important;
    height: 100px !important;
  }
  
  .edit-form {
    padding: 16px;
  }
  
  .edit-form .el-form-item__label {
    width: 80px !important;
  }
}
</style>
