<template>
  <el-dialog
    :model-value="modelValue"
    width="40%"
    title="导入云主机"
    @update:model-value="val => $emit('update:modelValue', val)"
    @close="dialogClose"
  >
    <el-steps :active="active" align-center style="margin-bottom: 10%">
      <el-step title="选择云服务商" />
      <el-step title="输入凭据并确认" />
    </el-steps>

    <el-form :model="form" ref="formRef" :rules="formRules" label-position="right" label-width="155px">
      <!--第一步-->
      <div v-show="active == 1">
        <el-form-item prop="cloud" label-width="0px">
          <el-radio-group v-model="form.cloud" class="cloud-radio-group">
            <el-radio label="aliyun"><img src="@/assets/image/aliyun.png" style="width: 60px; height: 60px; vertical-align: middle"></el-radio>
            <el-radio label="tencent"><img src="@/assets/image/tengxun.png" style="width: 60px; height: 60px; vertical-align: middle"></el-radio>
            <el-radio label="baidu"><img src="@/assets/image/baidu.svg" style="width: 50px; height: 50px; vertical-align: middle"></el-radio>
            <el-radio label="huawei"><img src="@/assets/image/华为.svg" style="width: 50px; height: 50px; vertical-align: middle"></el-radio>
            <el-radio label="aws"><img src="@/assets/image/aws.svg" style="width: 60px; height: 60px; vertical-align: middle"></el-radio>
          </el-radio-group>
        </el-form-item>
      </div>
      <!--第二步-->
      <div v-show="active == 2">
        <el-form-item label="AccessKey ID：" prop="secret_id">
          <el-input v-model="form.secret_id" clearable></el-input>
        </el-form-item>
        <el-form-item label="AccessKey Secret：" prop="secret_key">
          <el-input v-model="form.secret_key" clearable></el-input>
          <el-link v-if="form.cloud == 'aliyun'" href="https://ram.console.aliyun.com/manage/ak" target="_blank" type="primary">如何获取AccessKey？</el-link>
          <el-link v-if="form.cloud == 'tencent'" href="https://console.cloud.tencent.com/cam/capi" target="_blank" type="primary">如何获取AccessKey？</el-link>
          <el-link v-if="form.cloud == 'huawei'" href="https://console.huaweicloud.com/iam/#/mine/accessKey" target="_blank" type="primary">如何获取AccessKey？</el-link>
          <el-link v-if="form.cloud == 'aws'" href="https://console.aws.amazon.com/iam/home#/security_credentials" target="_blank" type="primary">如何获取AccessKey？</el-link>
        </el-form-item>
      </div>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="dialogClose" class="cancel-btn">取消</el-button>
        <el-button type="primary" @click="dialogToggle('pre')" v-if="active > 1">上一步</el-button>
        <el-button type="primary" @click="submit" v-if="active == 2">确认导入</el-button>
        <el-button type="primary" @click="dialogToggle('next')" v-if="active < 2">下一步</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script>
export default {
  name: "CreateCloud",
  props: {
    modelValue: Boolean,
  },
  emits: ['update:modelValue', 'success'],
  data() {
    const savedForm = JSON.parse(localStorage.getItem('cloudHostForm') || '{}')
    return {
      active: 1,
      form: {
        'cloud': savedForm.cloud || '',
        'secret_id': savedForm.secret_id || '',
        'secret_key': savedForm.secret_key || ''
      },
      formRules: {
        cloud: [
          {required: true, message: '请选择', trigger: 'blur'}
        ],
        secret_id: [
          {required: true, message: '请输入密钥ID', trigger: 'blur'},
          {min: 20, message: '密钥ID长度应不小于20个字符', trigger: 'blur'}
        ],
        secret_key: [
          {required: true, message: '请输入密钥Key', trigger: 'blur'},
          {min: 20, message: '密钥Key长度应不小于20个字符', trigger: 'blur'}
        ]
      }
    }
  },
  methods: {
    dialogClose() {
      console.log('[CreateCloud] 对话框关闭事件触发')
      localStorage.setItem('cloudHostForm', JSON.stringify(this.form))
      this.$emit('update:modelValue', false)
      this.resetForm()
    },
    resetForm() {
      console.log('[CreateCloud] 重置表单状态')
      this.form = {
        'cloud': '',
        'secret_id': '',
        'secret_key': ''
      }
      this.active = 1
    },
    dialogToggle(action) {
      console.log(`[CreateCloud] 对话框切换: ${action}`)
      if (action === 'pre') {
        if (this.active-- < 2) {
          this.active = 1
        }
      } else if (action === 'next') {
        if (this.active === 1 && !this.form.cloud) {
          this.$message.warning('请先选择云服务商')
          return
        }
        if (this.active++ > 2) {
          this.active = 1
        }
      }
    },
    async submit() {
      try {
        console.log('[CreateCloud] 开始验证表单...')
        
        // 验证AccessKey和SecretKey是否输入
        if (!this.form.secret_id || !this.form.secret_id.trim()) {
          this.$message.warning('请输入正确的AccessKey ID')
          return
        }
        
        if (!this.form.secret_key || !this.form.secret_key.trim()) {
          this.$message.warning('请输入正确的AccessKey Secret')
          return
        }
        
        // 验证AccessKey ID长度
        if (this.form.secret_id.trim().length < 16) {
          this.$message.warning('请输入正确的AccessKey ID（长度不能少于16个字符）')
          return
        }
        
        // 验证AccessKey Secret长度
        if (this.form.secret_key.trim().length < 20) {
          this.$message.warning('请输入正确的AccessKey Secret（长度不能少于20个字符）')
          return
        }
        
        await this.$refs.formRef.validate()
        
        console.log('[CreateCloud] 表单验证通过，开始导入云主机...')
        let res
        if(this.form.cloud === 'aliyun') {
          console.log('[CreateCloud] 导入阿里云主机...')
          res = await this.$api.hostcloudcreatealiyun(this.form)
        } else if(this.form.cloud === 'tencent') {
          console.log('[CreateCloud] 导入腾讯云主机...')
          // 腾讯云API需要特定参数格式
          const tencentParams = {
            secret_id: this.form.secret_id,
            secret_key: this.form.secret_key
          }
          res = await this.$api.hostcloudcreatetencent(tencentParams)
        } else if(this.form.cloud === 'huawei') {
          console.log('[CreateCloud] 导入华为云主机...')
          // 华为云API需要 accessKey/secretKey（区域自动发现）
          const huaweiParams = {
            accessKey: this.form.secret_id,
            secretKey: this.form.secret_key
            // 可选: groupId
          }
          res = await this.$api.hostcloudcreatehuawei(huaweiParams)
        } else if(this.form.cloud === 'aws') {
          console.log('[CreateCloud] 导入AWS主机...')
          // AWS API使用 secret_id 和 secret_key，会自动扫描所有区域
          const awsParams = {
            secret_id: this.form.secret_id,
            secret_key: this.form.secret_key
          }
          res = await this.$api.hostcloudcreateaws(awsParams)
        }

        this.$message.success('已受理，后台处理中')
        console.log('[CreateCloud] 导入成功，重置表单状态...')
        this.resetForm()
        console.log('[CreateCloud] 触发success事件...')
        this.$emit('success', { 
          success: true,
          cloud: this.form.cloud,
          shouldClose: true
        })
        this.$emit('update:modelValue', false) // 确保对话框关闭
      } catch (error) {
        console.error('[CreateCloud] 导入失败:', error)
        if (!error.response || error.response.status !== 200) {
          this.$message.error('导入失败: ' + (error.message || '未知错误'))
        }
        this.$emit('success', {
          success: false,
          error: error.message
        })
      }
    }
  }
}
</script>

<style scoped>
.cancel-btn {
  margin-right: 10px;
}

.cloud-radio-group {
  display: flex;
  flex-wrap: nowrap;
  justify-content: center;
  align-items: center;
  gap: 20px;
}

.cloud-radio-group .el-radio {
  margin-right: 0;
}
</style>
