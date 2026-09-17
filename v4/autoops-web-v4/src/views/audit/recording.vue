<template>
  <div class="recording-management">
    <el-card shadow="hover" class="recording-card">
      <template #header>
        <div class="card-header">
          <span class="title">会话录制管理</span>
        </div>
      </template>

      <!-- 搜索表单 -->
      <div class="search-section" v-show="showSearch">
        <el-form :model="queryParams" :inline="true" class="search-form">
          <el-form-item prop="keyword" label="关键词">
            <el-input v-model="queryParams.keyword" size="small" placeholder="主机名/IP/用户名"
                      style="width: 200px" clearable @keyup.enter="handleQuery" />
          </el-form-item>
          <el-form-item prop="status" label="状态">
            <el-select v-model="queryParams.status" size="small" placeholder="请选择" clearable style="width: 120px">
              <el-option label="录制中" :value="1" />
              <el-option label="已完成" :value="2" />
              <el-option label="失败" :value="3" />
            </el-select>
          </el-form-item>
          <el-form-item prop="riskLevel" label="风险等级">
            <el-select v-model="queryParams.riskLevel" size="small" placeholder="请选择" clearable style="width: 120px">
              <el-option label="正常" :value="0" />
              <el-option label="可疑" :value="1" />
              <el-option label="高危" :value="2" />
            </el-select>
          </el-form-item>
          <el-form-item prop="dateRange" label="时间范围">
            <el-date-picker v-model="queryParams.dateRange" size="small" type="datetimerange"
                            range-separator="至" start-placeholder="开始时间" end-placeholder="结束时间"
                            style="width: 360px" clearable />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" icon="Search" size="small" @click="handleQuery">搜索</el-button>
            <el-button type="warning" icon="Refresh" size="small" @click="resetQuery">重置</el-button>
            <el-button type="danger" plain icon="Delete" size="small" :disabled="multiple"
                       @click="batchHandleDelete">批量删除
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 列表区域 -->
      <div class="table-section">
        <el-table v-loading="Loading" :data="recordingList" stripe style="width: 100%"
                  class="recording-table" @selection-change="handleSelectionChange">
          <el-table-column type="selection" />
          <el-table-column label="ID" prop="id" width="80" align="center" />

          <el-table-column label="主机信息" min-width="200">
            <template #default="{ row }">
              <div class="host-info">{{ row.hostName }} ({{ row.hostIp }})</div>
              <div class="info-sub">SSH用户: {{ row.sshUser }}</div>
            </template>
          </el-table-column>

          <el-table-column prop="username" label="操作用户" min-width="100" />

          <el-table-column label="会话时间" min-width="180">
            <template #default="{ row }">
              <div class="time-info">开始: {{ formatTime(row.startTime) }}</div>
              <div class="info-sub">时长: {{ formatDuration(row.duration) }}</div>
            </template>
          </el-table-column>

          <el-table-column label="统计" min-width="180">
            <template #default="{ row }">
              <div class="stats-info">输入: {{ row.inputCount }} | 输出: {{ row.outputCount }}</div>
              <div class="info-sub">命令: {{ row.commandCount }} | 调整: {{ row.resizeCount }}</div>
            </template>
          </el-table-column>

          <el-table-column label="状态" min-width="90" align="center">
            <template #default="{ row }">
              <el-tag v-if="row.status === 1" type="info" size="small" effect="dark">录制中</el-tag>
              <el-tag v-else-if="row.status === 2" type="success" size="small" effect="dark">已完成</el-tag>
              <el-tag v-else-if="row.status === 3" type="danger" size="small" effect="dark">失败</el-tag>
            </template>
          </el-table-column>

          <el-table-column label="风险等级" min-width="100" align="center">
            <template #default="{ row }">
              <el-tag v-if="row.riskLevel === 0" type="success" size="small" effect="dark">正常</el-tag>
              <el-tag v-else-if="row.riskLevel === 1" type="warning" size="small" effect="dark">可疑</el-tag>
              <el-tag v-else-if="row.riskLevel === 2" type="danger" size="small" effect="dark">高危</el-tag>
            </template>
          </el-table-column>

          <el-table-column label="操作" width="220" fixed="right" align="center">
            <template #default="{ row }">
              <div class="operation-buttons">
                <el-tooltip content="回放" placement="top">
                  <el-button type="primary" icon="VideoPlay" size="small"
                             @click="handlePlay(row)" :disabled="row.status !== 2">回放</el-button>
                </el-tooltip>
                <el-tooltip content="下载" placement="top">
                  <el-button type="success" icon="Download" size="small"
                             @click="handleDownload(row)" :disabled="row.status !== 2">下载</el-button>
                </el-tooltip>
                <el-tooltip content="删除" placement="top">
                  <el-button type="danger" icon="Delete" size="small"
                             @click="handleDelete(row)">删除</el-button>
                </el-tooltip>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 分页区域 -->
      <div class="pagination-section">
        <el-pagination @size-change="handleSizeChange" @current-change="handleCurrentChange"
                       :current-page="queryParams.pageNum" :page-sizes="[10, 20, 50, 100]" :page-size="queryParams.pageSize"
                       layout="total, sizes, prev, pager, next, jumper" :total="total" />
      </div>

      <!-- Placeholder for dialog -->
    </el-card>

    <!-- 回放对话框 -->
    <el-dialog v-model="playDialog.visible" :title="playDialog.title" width="88%" top="4vh"
               class="recording-play-dialog" @opened="handlePlayDialogOpened" @close="handlePlayClose">
      <div class="play-body" v-loading="playDialog.loading" element-loading-text="加载中...">
        <el-row :gutter="16">
          <el-col :xs="24" :lg="17">
            <el-tabs v-model="playDialog.activeTab" class="play-tabs">
              <el-tab-pane name="play" label="回放">
                <el-alert v-if="playDialog.playerError" :title="playDialog.playerError"
                          type="error" show-icon :closable="false" style="margin-bottom: 10px" />
                <div class="player-container">
                  <div id="asciinema-player" ref="playerRef"></div>

                  <!-- 自定义播放控制栏 -->
                  <div class="player-controls" v-if="showControls">
                    <!-- 时间轴 -->
                    <div class="timeline-container">
                      <span class="time-label">{{ formatPlayTime(currentTime) }}</span>
                      <el-slider v-model="currentTime" :max="totalDuration" :min="0" :step="1"
                                 :show-tooltip="true" :format-tooltip="formatPlayTime"
                                 @input="handleTimelineInput" class="timeline-slider" />
                      <span class="time-label">{{ formatPlayTime(totalDuration) }}</span>
                    </div>

                    <!-- 控制按钮 -->
                    <div class="control-buttons">
                      <div class="button-group-left">
                        <el-button-group>
                          <el-button @click="handleRewind" size="small" :icon="DArrowLeft" title="快退5秒">快退</el-button>
                          <el-button @click="handlePlayPause" size="small"
                                     :icon="isPlaying ? VideoPause : VideoPlay" type="primary" title="播放/暂停">
                            {{ isPlaying ? '暂停' : '播放' }}
                          </el-button>
                          <el-button @click="handleFastForward" size="small" :icon="DArrowRight" title="快进5秒">快进</el-button>
                        </el-button-group>
                      </div>

                      <div class="button-group-right">
                        <span class="speed-label">播放速度</span>
                        <el-select v-model="playbackSpeed" @change="handleSpeedChange" size="small" class="speed-select">
                          <el-option label="0.5x" :value="0.5" />
                          <el-option label="1.0x" :value="1.0" />
                          <el-option label="1.5x" :value="1.5" />
                          <el-option label="2.0x" :value="2.0" />
                        </el-select>
                      </div>
                    </div>
                  </div>
                </div>
              </el-tab-pane>
              <el-tab-pane name="audit" label="命令审计">
                <el-table :data="commandList" max-height="520" stripe>
                  <el-table-column prop="command" label="命令" min-width="320" />
                  <el-table-column prop="executeTime" label="执行时间" width="180">
                    <template #default="{ row }">
                      {{ formatTime(row.executeTime) }}
                    </template>
                  </el-table-column>
                  <el-table-column label="类型" width="100">
                    <template #default="{ row }">
                      <el-tag v-if="row.isSensitive" type="danger">敏感</el-tag>
                      <el-tag v-else>普通</el-tag>
                    </template>
                  </el-table-column>
                  <el-table-column label="风险等级" width="100">
                    <template #default="{ row }">
                      <el-tag v-if="row.riskLevel === 0" type="success">正常</el-tag>
                      <el-tag v-else-if="row.riskLevel === 1" type="warning">可疑</el-tag>
                      <el-tag v-else-if="row.riskLevel === 2" type="danger">高危</el-tag>
                    </template>
                  </el-table-column>
                </el-table>
              </el-tab-pane>
            </el-tabs>
          </el-col>

          <el-col :xs="24" :lg="7">
            <el-card shadow="never" class="meta-card">
              <template #header>
                <div class="meta-header">
                  <span>会话信息</span>
                  <el-button v-if="playDialog.row" type="primary" plain size="small" icon="Download"
                             @click="handleDownload(playDialog.row)">下载
                  </el-button>
                </div>
              </template>

              <el-descriptions :column="1" size="small" border>
                <el-descriptions-item label="主机">
                  {{ playDialog.row?.hostName }} ({{ playDialog.row?.hostIp }})
                </el-descriptions-item>
                <el-descriptions-item label="SSH用户">
                  {{ playDialog.row?.sshUser }}
                </el-descriptions-item>
                <el-descriptions-item label="操作用户">
                  {{ playDialog.row?.username }}
                </el-descriptions-item>
                <el-descriptions-item label="开始时间">
                  {{ formatTime(playDialog.row?.startTime) }}
                </el-descriptions-item>
                <el-descriptions-item label="时长">
                  {{ formatDuration(playDialog.row?.duration) }}
                </el-descriptions-item>
                <el-descriptions-item label="分辨率">
                  {{ playDialog.row?.terminalWidth || 80 }} x {{ playDialog.row?.terminalHeight || 24 }}
                </el-descriptions-item>
              </el-descriptions>
            </el-card>
          </el-col>
        </el-row>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { VideoPlay, VideoPause, DArrowLeft, DArrowRight } from '@element-plus/icons-vue'
import RecordingAPI from '@/api/recording'

// 数据
const Loading = ref(true)
const ids = ref([])
const single = ref(true)
const multiple = ref(true)
const showSearch = ref(true)
const total = ref(0)
const queryParams = reactive({
  pageNum: 1,
  pageSize: 10,
  keyword: '',
  status: null,
  riskLevel: null,
  dateRange: []
})
const recordingList = ref([])

// 回放对话框
const playDialog = reactive({
  visible: false,
  title: '',
  recordingId: null,
  activeTab: 'play',
  loading: false,
  row: null,
  playerError: '',
  castText: ''
})
const playerRef = ref(null)
let playerInstance = null

// 播放器状态
const showControls = ref(false)
const isPlaying = ref(false)
const currentTime = ref(0)
const totalDuration = ref(0)
const playbackSpeed = ref(1.0)
let timeUpdateInterval = null

// 命令审计列表
const commandList = ref([])

// 多选框选中数据
const handleSelectionChange = (selection) => {
  ids.value = selection.map(item => item.id)
  single.value = selection.length !== 1
  multiple.value = !selection.length
}

// 查询列表
const getRecordingList = async () => {
  Loading.value = true
  try {
    const params = {
      page: queryParams.pageNum,
      pageSize: queryParams.pageSize,
      keyword: queryParams.keyword || undefined,
      status: queryParams.status,
      riskLevel: queryParams.riskLevel
    }

    // 处理时间范围
    if (queryParams.dateRange && queryParams.dateRange.length === 2) {
      params.startTime = queryParams.dateRange[0].toISOString()
      params.endTime = queryParams.dateRange[1].toISOString()
    }

    const response = await RecordingAPI.getRecordingList(params)
    const data = response.data?.data || {}
    recordingList.value = data.list || []
    total.value = data.total || 0
  } catch (error) {
    ElMessage.error('加载数据失败: ' + error.message)
  } finally {
    Loading.value = false
  }
}

// 搜索按钮操作
const handleQuery = () => {
  queryParams.pageNum = 1
  getRecordingList()
}

// 重置按钮操作
const resetQuery = () => {
  queryParams.pageNum = 1
  queryParams.pageSize = 10
  queryParams.keyword = ''
  queryParams.status = null
  queryParams.riskLevel = null
  queryParams.dateRange = []
  getRecordingList()
  ElMessage.success('重置成功')
}

// pageSize
const handleSizeChange = (newSize) => {
  queryParams.pageSize = newSize
  getRecordingList()
}

// pageNum
const handleCurrentChange = (newPage) => {
  queryParams.pageNum = newPage
  getRecordingList()
}

// 删除
const handleDelete = async (row) => {
  const confirmResult = await ElMessageBox.confirm(
    `确定要删除 ${row.hostName} 的会话录制吗?`,
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).catch(err => err)

  if (confirmResult !== 'confirm') {
    return ElMessage.info('已取消删除')
  }

  const { data: res } = await RecordingAPI.deleteRecording(row.id)
  if (res.code !== 200) {
    ElMessage.error(res.message)
  } else {
    ElMessage.success('删除成功')
    await getRecordingList()
  }
}

// 批量删除
const batchHandleDelete = async () => {
  const recordingIds = ids.value
  const confirmResult = await ElMessageBox.confirm(
    `确定要删除选中的 ${recordingIds.length} 条录制记录吗?`,
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).catch(err => err)

  if (confirmResult !== 'confirm') {
    return ElMessage.info('已取消删除')
  }

  const { data: res } = await RecordingAPI.batchDeleteRecording(recordingIds)
  if (res.code !== 200) {
    ElMessage.error(res.message)
  } else {
    ElMessage.success('删除成功')
    ids.value = []
    await getRecordingList()
  }
}

// 等待函数
const waitFor = async (predicate, timeoutMs = 2000, intervalMs = 50) => {
  const deadline = Date.now() + timeoutMs
  while (Date.now() < deadline) {
    const value = predicate()
    if (value) return value
    await new Promise(resolve => setTimeout(resolve, intervalMs))
  }
  return null
}

// 初始化播放器
const initPlayer = async () => {
  if (!playDialog.visible || playDialog.activeTab !== 'play') return
  if (!playDialog.castText) return

  await nextTick()
  const container = (await waitFor(() => playerRef.value, 3000)) || document.getElementById('asciinema-player')
  if (!container) {
    throw new Error('播放器容器未渲染（请稍后重试）')
  }

  await loadAsciinemaPlayer()
  if (!window.AsciinemaPlayer && typeof AsciinemaPlayer !== 'undefined') {
    window.AsciinemaPlayer = AsciinemaPlayer
  }
  if (!window.AsciinemaPlayer) {
    throw new Error('播放器未加载：window.AsciinemaPlayer 为空（请查看 Console 是否有 CSP/脚本报错）')
  }

  let cols = playDialog.row?.terminalWidth || 80
  let rows = playDialog.row?.terminalHeight || 24
  try {
    const headerLine = playDialog.castText.split('\n')[0]
    const header = JSON.parse(headerLine)
    if (header && Number.isFinite(header.width)) cols = header.width
    if (header && Number.isFinite(header.height)) rows = header.height
  } catch (_) {
    // ignore
  }

  container.innerHTML = ''
  playerInstance = window.AsciinemaPlayer.create(
    { data: playDialog.castText },
    container,
    {
      cols,
      rows,
      autoPlay: false,
      speed: 1.0,
      theme: 'monokai',
      fit: 'width'
    }
  )

  // 获取总时长
  try {
    const lines = playDialog.castText.trim().split('\n')
    if (lines.length > 1) {
      for (let i = lines.length - 1; i >= 1; i--) {
        const line = lines[i].trim()
        if (line.startsWith('[')) {
          const match = line.match(/^\[([0-9.]+),/)
          if (match) {
            totalDuration.value = Math.ceil(parseFloat(match[1]))
            break
          }
        }
      }
    }

    if (totalDuration.value === 0) {
      if (playerInstance.duration) {
        totalDuration.value = Math.floor(playerInstance.duration)
      } else if (playerInstance.getDuration) {
        totalDuration.value = Math.floor(playerInstance.getDuration())
      }
    }

    if (totalDuration.value === 0) {
      totalDuration.value = 100
    }
  } catch (e) {
    totalDuration.value = 100
  }

  showControls.value = true
  startTimeUpdate()
}

// 播放/暂停
const handlePlayPause = () => {
  if (!playerInstance) return

  if (isPlaying.value) {
    playerInstance.pause()
    isPlaying.value = false
  } else {
    playerInstance.play()
    isPlaying.value = true
  }
}

// 快退5秒
const handleRewind = () => {
  if (!playerInstance) return
  const newTime = Math.max(0, currentTime.value - 5)
  playerInstance.seek(newTime)
  currentTime.value = newTime
}

// 快进5秒
const handleFastForward = () => {
  if (!playerInstance) return
  const newTime = Math.min(totalDuration.value, currentTime.value + 5)
  playerInstance.seek(newTime)
  currentTime.value = newTime
}

// 修改播放速度
const handleSpeedChange = (speed) => {
  if (!playerInstance) return
  const currentPos = currentTime.value
  const wasPlaying = isPlaying.value

  const container = document.getElementById('asciinema-player')
  if (!container) return

  container.innerHTML = ''
  playerInstance = window.AsciinemaPlayer.create(
    { data: playDialog.castText },
    container,
    {
      cols: playDialog.row?.terminalWidth || 80,
      rows: playDialog.row?.terminalHeight || 24,
      autoPlay: wasPlaying,
      speed: speed,
      theme: 'monokai',
      fit: 'width',
      startAt: currentPos
    }
  )

  isPlaying.value = wasPlaying
}

// 时间轴拖动
const handleTimelineInput = (value) => {
  if (!playerInstance) return

  try {
    if (typeof playerInstance.seek === 'function') {
      playerInstance.seek(value)
    } else {
      const container = document.getElementById('asciinema-player')
      if (!container) return

      const wasPlaying = isPlaying.value
      container.innerHTML = ''
      playerInstance = window.AsciinemaPlayer.create(
        { data: playDialog.castText },
        container,
        {
          cols: playDialog.row?.terminalWidth || 80,
          rows: playDialog.row?.terminalHeight || 24,
          autoPlay: wasPlaying,
          speed: playbackSpeed.value,
          theme: 'monokai',
          fit: 'width',
          startAt: value
        }
      )
      isPlaying.value = wasPlaying
    }
  } catch (e) {
    console.error('跳转失败:', e)
  }
}

// 格式化播放时间
const formatPlayTime = (seconds) => {
  if (!seconds || seconds < 0) return '00:00'
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
}

// 启动时间更新定时器
const startTimeUpdate = () => {
  stopTimeUpdate()

  timeUpdateInterval = setInterval(() => {
    if (!playerInstance) return

    try {
      let time = 0

      if (typeof playerInstance.getCurrentTime === 'function') {
        time = playerInstance.getCurrentTime()
      } else if (playerInstance.currentTime !== undefined) {
        time = playerInstance.currentTime
      }

      const newTime = Math.floor(time)
      if (currentTime.value !== newTime) {
        currentTime.value = newTime
      }

      if (typeof playerInstance.isPaused === 'function') {
        const paused = playerInstance.isPaused()
        if (paused && isPlaying.value) {
          isPlaying.value = false
        } else if (!paused && !isPlaying.value) {
          isPlaying.value = true
        }
      }

      if (time >= totalDuration.value - 0.5 && isPlaying.value) {
        isPlaying.value = false
      }
    } catch (e) {
      console.error('时间更新失败:', e)
    }
  }, 200)
}

// 停止时间更新定时器
const stopTimeUpdate = () => {
  if (timeUpdateInterval) {
    clearInterval(timeUpdateInterval)
    timeUpdateInterval = null
  }
}

// 回放
const handlePlay = async (row) => {
  playDialog.visible = true
  playDialog.title = `会话回放 - ${row.hostName} (${row.hostIp})`
  playDialog.recordingId = row.id
  playDialog.activeTab = 'play'
  playDialog.row = row
  playDialog.loading = true
  playDialog.playerError = ''
  playDialog.castText = ''

  try {
    const castResponse = await RecordingAPI.getRecordingCast(row.id)
    const castText = castResponse.data
    if (!castText || typeof castText !== 'string') {
      throw new Error('cast内容为空')
    }
    if (!castText.trimStart().startsWith('{') || !castText.includes('"version"')) {
      throw new Error('cast格式不正确')
    }
    if (castText.split('\n').length < 2) {
      throw new Error('cast没有事件数据（只有头部）')
    }

    playDialog.castText = castText
    await initPlayer()

    const cmdResponse = await RecordingAPI.getCommandAuditList({
      page: 1,
      pageSize: 100,
      recordingId: row.id
    })
    commandList.value = cmdResponse.data?.data?.list || []
  } catch (error) {
    playDialog.playerError = error?.message || String(error)
    ElMessage.error('加载回放失败: ' + playDialog.playerError)
  } finally {
    playDialog.loading = false
  }
}

const handlePlayDialogOpened = async () => {
  try {
    await initPlayer()
  } catch (e) {
    playDialog.playerError = e?.message || String(e)
  }
}

watch(
  () => playDialog.activeTab,
  async () => {
    if (playDialog.activeTab !== 'play') return
    try {
      await initPlayer()
    } catch (e) {
      playDialog.playerError = e?.message || String(e)
    }
  }
)

// 关闭回放
const handlePlayClose = () => {
  stopTimeUpdate()
  showControls.value = false
  isPlaying.value = false
  currentTime.value = 0
  totalDuration.value = 0
  playbackSpeed.value = 1.0

  if (playerInstance) {
    playerInstance = null
  }
  commandList.value = []
  playDialog.row = null
  playDialog.recordingId = null
  playDialog.activeTab = 'play'
  playDialog.playerError = ''
  playDialog.castText = ''
  const container = document.getElementById('asciinema-player')
  if (container) {
    container.innerHTML = ''
  }
}

// 下载
const handleDownload = async (row) => {
  try {
    const response = await RecordingAPI.downloadRecording(row.id)

    const blob = new Blob([response.data], { type: 'application/octet-stream' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `session_${row.sessionId || row.id}.cast`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)

    ElMessage.success('下载成功')
  } catch (error) {
    ElMessage.error('下载失败: ' + error.message)
  }
}

// 格式化时间
const formatTime = (time) => {
  if (!time) return '-'
  const date = new Date(time)
  return date.toLocaleString('zh-CN', { hour12: false })
}

// 格式化时长
const formatDuration = (seconds) => {
  if (!seconds) return '0秒'
  const h = Math.floor(seconds / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  const s = seconds % 60

  const parts = []
  if (h > 0) parts.push(`${h}小时`)
  if (m > 0) parts.push(`${m}分`)
  if (s > 0 || parts.length === 0) parts.push(`${s}秒`)

  return parts.join('')
}

// 动态加载asciinema-player
const loadAsciinemaPlayer = () => {
  return new Promise((resolve, reject) => {
    if (window.AsciinemaPlayer || typeof AsciinemaPlayer !== 'undefined') {
      if (!window.AsciinemaPlayer && typeof AsciinemaPlayer !== 'undefined') {
        window.AsciinemaPlayer = AsciinemaPlayer
      }
      resolve()
      return
    }

    const baseUrl = (process.env.BASE_URL || '/').replace(/\/?$/, '/')
    const cssUrl = `${baseUrl}vendor/asciinema-player/asciinema-player.min.css`
    const jsUrl = `${baseUrl}vendor/asciinema-player/asciinema-player.min.js`

    const link = document.createElement('link')
    link.rel = 'stylesheet'
    link.href = cssUrl
    document.head.appendChild(link)

    fetch(jsUrl, { cache: 'no-store' })
      .then(async (res) => {
        const text = await res.text()
        const head = text.trimStart().slice(0, 80).toLowerCase()
        if (!res.ok) {
          throw new Error(`播放器资源加载失败(${res.status}): ${jsUrl}`)
        }
        if (head.startsWith('<!doctype') || head.startsWith('<html')) {
          throw new Error(`播放器资源返回了HTML（请检查静态资源映射/History回退规则）：${jsUrl}`)
        }

        const script = document.createElement('script')
        script.src = jsUrl
        script.async = true
        script.onload = () => {
          if (!window.AsciinemaPlayer && typeof AsciinemaPlayer !== 'undefined') {
            window.AsciinemaPlayer = AsciinemaPlayer
          }
          if (!window.AsciinemaPlayer) {
            reject(new Error(`播放器脚本已加载但未暴露全局 AsciinemaPlayer：${jsUrl}`))
            return
          }
          resolve()
        }
        script.onerror = () => reject(new Error(`播放器脚本加载失败: ${jsUrl}`))
        document.head.appendChild(script)
      })
      .catch(reject)
  })
}

// 生命周期
onMounted(() => {
  getRecordingList()
})

onUnmounted(() => {
  handlePlayClose()
})
</script>

<style scoped>
.recording-management {
  padding: 5px;
  min-height: 100vh;
  background: #eef1f6;
}

.recording-card {
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
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.search-section {
  margin-bottom: 24px;
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

.action-section {
  margin-bottom: 24px;
  padding: 12px 0;
}

.action-section .el-button {
  margin-right: 12px;
}

.table-section {
  margin-bottom: 24px;
}

.recording-table {
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.recording-table :deep(.el-table__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
}

.recording-table :deep(.el-table__header th) {
  background: transparent !important;
  color: #2c3e50 !important;
  font-weight: 700 !important;
  border-bottom: none;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
}

.recording-table :deep(.el-table__header th .cell) {
  color: #2c3e50 !important;
  font-weight: 700 !important;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
}

.recording-table :deep(.el-table__row) {
  transition: all 0.3s ease;
}

.recording-table :deep(.el-table__row:hover) {
  background-color: rgba(103, 126, 234, 0.05) !important;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.host-info {
  color: #303133;
  font-size: 14px;
  line-height: 1.5;
}

.time-info {
  color: #303133;
  font-size: 13px;
  line-height: 1.5;
}

.stats-info {
  color: #303133;
  font-size: 13px;
  line-height: 1.5;
}

.info-sub {
  font-size: 12px;
  color: #909399;
  margin-top: 2px;
  line-height: 1.4;
}

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

.pagination-section {
  display: flex;
  justify-content: center;
  padding: 20px 0;
}

.el-tag {
  font-weight: 500;
  border-radius: 8px;
  border: none;
}

.el-button {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.el-input :deep(.el-input__inner),
.el-select :deep(.el-input__inner),
.el-date-picker :deep(.el-input__inner) {
  border-radius: 8px;
  border: 1px solid rgba(103, 126, 234, 0.2);
  transition: all 0.3s ease;
}

.el-input :deep(.el-input__inner):focus,
.el-select :deep(.el-input__inner):focus,
.el-date-picker :deep(.el-input__inner):focus {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(103, 126, 234, 0.2);
}

.el-loading-mask {
  background-color: rgba(103, 126, 234, 0.1);
  backdrop-filter: blur(4px);
}

.el-pagination {
  background: transparent;
}

.el-pagination :deep(.btn-next),
.el-pagination :deep(.btn-prev),
.el-pagination :deep(.el-pager li) {
  border-radius: 6px;
  margin: 0 2px;
  transition: all 0.3s ease;
}

.el-pagination :deep(.btn-next):hover,
.el-pagination :deep(.btn-prev):hover,
.el-pagination :deep(.el-pager li):hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

/* 播放器相关样式 */
.player-container {
  background-color: #000;
  border-radius: 4px;
  padding: 10px;
  min-height: 460px;
  height: 56vh;
  display: flex;
  flex-direction: column;
}

#asciinema-player {
  width: 100%;
  flex: 1;
  min-height: 0;
}

.player-controls {
  margin-top: 12px;
  padding: 16px;
  background: linear-gradient(145deg, #1e1e1e 0%, #2a2a2a 100%);
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.timeline-container {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
  width: 100%;
}

.timeline-slider {
  flex: 1;
}

.time-label {
  color: #e0e0e0;
  font-size: 13px;
  font-family: 'Monaco', 'Menlo', 'Consolas', 'Courier New', monospace;
  font-weight: 500;
  min-width: 45px;
  text-align: center;
  background-color: rgba(255, 255, 255, 0.05);
  padding: 4px 8px;
  border-radius: 4px;
}

.timeline-container :deep(.el-slider__runway) {
  background-color: #3a3a3a;
  height: 6px;
  border-radius: 3px;
}

.timeline-container :deep(.el-slider__bar) {
  background: linear-gradient(90deg, #409eff 0%, #66b1ff 100%);
  height: 6px;
  border-radius: 3px;
}

.timeline-container :deep(.el-slider__button) {
  border: 2px solid #409eff;
  width: 16px;
  height: 16px;
  background-color: #fff;
  box-shadow: 0 2px 6px rgba(64, 158, 255, 0.4);
  transition: all 0.2s ease;
}

.timeline-container :deep(.el-slider__button:hover) {
  transform: scale(1.2);
  box-shadow: 0 3px 10px rgba(64, 158, 255, 0.6);
}

.control-buttons {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.button-group-left {
  display: flex;
  align-items: center;
}

.button-group-right {
  display: flex;
  align-items: center;
  gap: 10px;
}

.speed-label {
  color: #b0b0b0;
  font-size: 13px;
  white-space: nowrap;
}

.speed-select {
  width: 90px;
}

.speed-select :deep(.el-input__wrapper) {
  background-color: rgba(255, 255, 255, 0.08);
  box-shadow: 0 0 0 1px rgba(255, 255, 255, 0.1) inset;
}

.speed-select :deep(.el-input__inner) {
  color: #fff;
}

.control-buttons :deep(.el-button-group) {
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.2);
}

.control-buttons :deep(.el-button) {
  transition: all 0.2s ease;
}

.control-buttons :deep(.el-button:hover) {
  transform: translateY(-1px);
}

.control-buttons :deep(.el-button.is-disabled) {
  opacity: 0.5;
}

.meta-card :deep(.el-card__header) {
  padding: 10px 12px;
}

.meta-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.play-tabs :deep(.el-tabs__content) {
  padding-top: 8px;
}
</style>
