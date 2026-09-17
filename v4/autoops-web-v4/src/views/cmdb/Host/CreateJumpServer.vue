<template>
  <el-dialog
    v-model="dialogVisible"
    title="JumpServer主机导入"
    width="500px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
    >
      <el-form-item label="Jump地址" prop="url">
        <el-input
          v-model="form.url"
          placeholder="请输入Jump地址,如: https://jump.dding.net/"
          clearable
        />
      </el-form-item>

      <el-form-item label="Jump账号" prop="accessKey">
        <el-input
          v-model="form.accessKey"
          placeholder="请输入用户名"
          clearable
          show-password
        />
      </el-form-item>

      <el-form-item label="Jump密码" prop="secretKey">
        <el-input
          v-model="form.secretKey"
          placeholder="请输入密码"
          clearable
          show-password
        />
      </el-form-item>

      <el-alert
        title="提示"
        type="info"
        :closable="false"
        style="margin-bottom: 20px"
      >
        <template #default>
          <div>同步任务将在后台执行,请稍后刷新主机列表查看导入结果</div>
        </template>
      </el-alert>
    </el-form>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button
        type="primary"
        @click="handleSubmit"
        :loading="loading"
      >
        {{ loading ? '同步中...' : '开始同步' }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script>
export default {
  name: 'CreateJumpServer',
  props: {
    modelValue: {
      type: Boolean,
      default: false
    }
  },
  emits: ['update:modelValue', 'success'],
  data() {
    return {
      loading: false,
      form: {
        url: '',
        accessKey: '',
        secretKey: ''
      },
      rules: {
        url: [
          { required: true, message: '请输入JumpServer地址', trigger: 'blur' },
          {
            pattern: /^https?:\/\/.+/,
            message: '请输入有效的URL地址',
            trigger: 'blur'
          }
        ],
        accessKey: [
          { required: true, message: '请输入Access Key', trigger: 'blur' }
        ],
        secretKey: [
          { required: true, message: '请输入Secret Key', trigger: 'blur' }
        ]
      }
    }
  },
  computed: {
    dialogVisible: {
      get() {
        return this.modelValue
      },
      set(val) {
        this.$emit('update:modelValue', val)
      }
    }
  },
  methods: {
    async handleSubmit() {
      try {
        await this.$refs.formRef.validate()

        this.loading = true

        const { data: res } = await this.$api.syncHostFromJumpServer({
          url: this.form.url,
          accessKey: this.form.accessKey,
          secretKey: this.form.secretKey
        })

        if (res.code === 200) {
          this.$message.success(res.data?.message || 'JumpServer主机同步任务已启动')
          this.dialogVisible = false
          this.$emit('success')
          this.resetForm()
        } else {
          this.$message.error(res.message || '同步失败')
        }
      } catch (error) {
        if (error !== false) {
          console.error('JumpServer同步失败:', error)
          this.$message.error('同步失败: ' + (error.message || '未知错误'))
        }
      } finally {
        this.loading = false
      }
    },

    handleClose() {
      this.dialogVisible = false
      this.resetForm()
    },

    resetForm() {
      this.$refs.formRef?.resetFields()
      this.form = {
        url: '',
        accessKey: '',
        secretKey: ''
      }
    }
  }
}
</script>

<style scoped>
.el-alert {
  margin-top: 10px;
}

.el-form-item {
  margin-bottom: 20px;
}
</style>
