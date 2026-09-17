<template>
  <el-dialog
    title="Excel导入主机"
    v-model="dialogVisible"
    width="25%"
    @close="handleClose"
  >
    <el-form :model="form" label-width="120px">
      <!-- 模板下载 -->
      <el-form-item label="模板下载">
        <el-button
          type="primary"
          @click="downloadTemplate"
          :loading="downloading"
        >
          <el-icon><Download /></el-icon>下载模板
        </el-button>
      </el-form-item>

      <!-- 选择分组 -->
      <el-form-item label="选择分组" prop="groupId">
        <el-cascader
          v-model="form.groupId"
          :options="groupList"
          :props="{
            checkStrictly: true,
            value: 'id',
            label: 'name',
            children: 'children',
            disabled: (node) => node.children && node.children.length > 0
          }"
          placeholder="请选择分组"
          style="width: 100%"
          clearable
          filterable
          @change="handleGroupChange"
        ></el-cascader>
      </el-form-item>

      <!-- 上传Excel -->
      <el-form-item label="上传Excel" prop="file">
        <el-upload
          class="upload-demo"
          :auto-upload="false"
          :on-change="handleFileChange"
          :on-remove="handleFileRemove"
          :show-file-list="false"
          accept=".xlsx,.xls"
        >
          <el-button type="primary">选择文件</el-button>
          <template #tip>
            <div class="el-upload__tip" v-if="form.file">
              已选择: {{ form.file.name }}
              <el-button
                type="danger"
                text
                icon="Close"
                circle
                size="small"
                @click.stop="handleFileRemove"
                style="margin-left: 8px"
              />
            </div>
            <div class="el-upload__tip" style="color: #999; margin-top: 15px">
              请上传Excel文件（.xlsx或.xls格式）
            </div>
          </template>
        </el-upload>
      </el-form-item>

      <el-progress 
        v-if="uploading"
        :percentage="uploadProgress" 
        :status="uploadProgress === 100 ? 'success' : ''"
      />
    </el-form>

    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button
        type="primary"
        @click="submitImport"
        :loading="uploading"
        :disabled="!form.file || !form.groupId"
      >
        导入主机
      </el-button>
    </template>
  </el-dialog>
</template>

<script>
import { Download } from '@element-plus/icons-vue'

export default {
  name: 'CreateExcel',
  components: { Download },
  props: {
    modelValue: {
      type: Boolean,
      default: false
    }
  },
  emits: ['update:modelValue'],
  data() {
    return {
      form: {
        groupId: '',
        file: null
      },
      groupList: [],
      downloading: false,
      uploading: false,
      uploadProgress: 0
    }
  },
  computed: {
    dialogVisible: {
      get() {
        return this.modelValue
      },
      set(value) {
        this.$emit('update:modelValue', value)
      }
    }
  },
  created() {
    this.getGroupList()
  },
  methods: {
    // 获取分组列表
    async getGroupList() {
      try {
        const { data: res } = await this.$api.getAllCmdbGroups()
        if (res.code === 200) {
          this.groupList = res.data
          // 设置默认分组为业务组
          const businessGroup = this.groupList.find(group => group.name === '业务组')
          if (businessGroup) {
            this.form.groupId = businessGroup.id
          }
        }
      } catch (error) {
        console.error('获取分组列表失败:', error)
        this.$message.error('获取分组列表失败')
      }
    },

    // 处理分组选择变化
    handleGroupChange(value) {
      if (value && value.length > 0) {
        // 取最后一级作为选中分组ID
        this.form.groupId = value[value.length - 1]
      } else {
        // 如果没有选择分组，设置默认分组
        const defaultGroup = this.groupList.find(item => item.isDefault)
        if (defaultGroup) {
          this.form.groupId = defaultGroup.id
        }
      }
    },

    // 下载模板
    async downloadTemplate() {
      this.downloading = true
      try {
        const response = await this.$api.DownloadHostTemplate()
        
        // 创建下载链接
        const url = window.URL.createObjectURL(new Blob([response.data]))
        const link = document.createElement('a')
        link.href = url
        link.setAttribute('download', '主机导入模板.xlsx')
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
        
        this.$message.success('模板下载成功')
      } catch (error) {
        console.error('模板下载失败:', error)
        this.$message.error('模板下载失败')
      } finally {
        this.downloading = false
      }
    },

    // 文件选择
    handleFileChange(file) {
      this.form.file = file.raw
    },

    // 文件移除
    handleFileRemove() {
      this.form.file = null
    },

    // 提交导入
    async submitImport() {
      if (!this.form.file) {
        return this.$message.warning('请选择要上传的文件')
      }
      if (!this.form.groupId) {
        return this.$message.warning('请选择分组')
      }

      this.uploading = true
      this.uploadProgress = 0

      try {
        const formData = new FormData()
        formData.append('file', this.form.file)
        formData.append('groupId', this.form.groupId)

        const { data: res } = await this.$api.ImportHostsFromExcel(formData, {
          onUploadProgress: progressEvent => {
            this.uploadProgress = Math.round(
              (progressEvent.loaded * 100) / progressEvent.total
            )
          }
        })

        if (res.code === 200) {
          this.$message.success('导入成功')
          this.dialogVisible = false
          this.$emit('success') // 通知父组件刷新数据
        } else {
          this.$message.error(res.message || '导入失败')
        }
      } catch (error) {
        console.error('导入失败:', error)
        this.$message.error('导入失败')
      } finally {
        this.uploading = false
      }
    },

    // 关闭对话框
    handleClose() {
      this.form = {
        groupId: '',
        file: null
      }
      this.uploadProgress = 0
    }
  }
}
</script>
