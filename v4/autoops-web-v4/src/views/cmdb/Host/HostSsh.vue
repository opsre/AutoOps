<template>
  <div class="ssh-container">
    <div class="ssh-layout">
      <div class="tree-section">
        <el-card shadow="never" style="height: 100%; border: none;">
          <div slot="header">
            <h3>资产分组</h3>
          </div>
          <el-input
            v-model="groupSearchText"
            placeholder="搜索分组"
            clearable
            size="medium"
            style="margin-bottom: 20px; width: 100%"
            @input="handleGroupSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>

          <el-tree
            ref="groupTree"
            :data="groupList"
            :props="defaultProps"
            node-key="id"
            :expanded-keys="expandedKeys"
            :highlight-current="true"
            @node-click="handleGroupClick"
            @node-expand="handleNodeExpand"
            @node-collapse="handleNodeCollapse"
            style="border: none; padding: 5px;"
          >
            <template v-slot="{ node, data }">
              <span :class="{ 'font-weight-bold': !data.parentId }">
                <template v-if="data.isHost">
                  <el-icon style="margin-right: 5px"><Platform /></el-icon>
                </template>
                <template v-else-if="data.parentId">
                  <el-icon style="margin-right: 5px"><FolderRemove /></el-icon>
                </template>
                <template v-else>
                  <el-icon v-if="expandedKeys.includes(node.key)" style="margin-right: 5px"><FolderOpened /></el-icon>
                  <el-icon v-else style="margin-right: 5px"><Folder /></el-icon>
                </template>
                {{ node.label }}
              </span>
            </template>
          </el-tree>
        </el-card>
      </div>

      <div class="terminal-section">
        <el-card shadow="never" style="height: 100%;">
          <div slot="header">
            <h3>SSH终端连接</h3>
            <div class="terminal-actions" v-if="currentHost">
              <el-button size="mini" @click="connectTerminal" :disabled="isConnecting || !currentHost">
                {{ isConnected ? '重新连接' : '连接' }}
              </el-button>
              <el-button size="mini" @click="disconnectTerminal" :disabled="!isConnected || isConnecting">
                断开
              </el-button>
            </div>
          </div>
          <div class="terminal-content">
            <div class="host-info" v-if="currentHost">
              <el-descriptions :column="1" border>
                <el-descriptions-item label="主机名称">{{ currentHost.hostName }}</el-descriptions-item>
                <el-descriptions-item label="连接地址">
                  {{ currentHost.sshName }}@{{ currentHost.sshIp }}:{{ currentHost.sshPort }}
                </el-descriptions-item>
              </el-descriptions>
            </div>
            <div id="terminal" ref="terminalEl" class="terminal-wrapper"></div>
          </div>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script>
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import { AttachAddon } from '@xterm/addon-attach'
import 'xterm/css/xterm.css'
import storage from '@/utils/storage'

export default {
  data() {
    return {
      groupSearchText: '',
      expandedKeys: [],
      groupList: [],
      defaultProps: {
        children: 'children',
        label: 'name'
      },
      currentHost: null,
      term: null,
      fitAddon: null,
      isConnecting: false,
      isConnected: false,
      socket: null
    }
  },
  mounted() {
    this.$nextTick(() => {
      if (this.$refs.terminalEl) {
        this.initTerminal()
      }
    })
    window.addEventListener('resize', this.handleResize)
  },
  beforeUnmount() {
    this.disconnectTerminal()
    window.removeEventListener('resize', this.handleResize)
    if (this.term) {
      this.term.dispose()
    }
  },
  created() {
    this.loadGroupList()
  },
  methods: {
    initTerminal() {
      if (!this.$refs.terminalEl) {
        console.error('Terminal container element not found')
        return
      }
      
      this.term = new Terminal({
        fontSize: 15,
        cursorBlink: true,
        cursorStyle: "underline",
        cols: 188,
        rows: 42,
        theme: {
          background: '#060101',
          foreground: 'white',
          cursor: 'help'
        },
        allowTransparency: true,
        convertEol: true,
        scrollback: 10000,
        rendererType: 'canvas',
        disableStdin: false,
        fontFamily: 'operator mono,SFMono-Regular,Consolas,Liberation Mono,Menlo,monospace',
        tabStopWidth: 8,
        bellStyle: "sound",
        rightClickSelectsWord: true
      })
      this.fitAddon = new FitAddon()
      this.term.loadAddon(this.fitAddon)
      this.term.open(this.$refs.terminalEl)
      this.fitAddon.fit()
      
      this.term.writeln('欢迎使用SSH终端')
      this.term.writeln('请从左侧选择主机进行连接')
      this.term.writeln('')
      
      // 确保只绑定一次onData事件
      if (!this._dataHandlerBound) {
        this.term.onData(data => {
          if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            // 更严格的防抖处理
            if (this._lastInputTime && Date.now() - this._lastInputTime < 30) {
              console.log('忽略重复输入:', data)
              return
            }
            this._lastInputTime = Date.now()
            console.log('发送输入:', data)
            this.socket.send(data)
          }
        })
        this._dataHandlerBound = true
      }
    },
    
    async connectTerminal() {
      if (!this.currentHost || this.isConnecting) return



      this.isConnecting = true
      try {
          const token = storage.getItem('token')
  

        // 使用当前页面的host，支持Docker部署
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
        const wsUrl = `${protocol}//${window.location.host}/api/v1/cmdb/hostssh/connect/${this.currentHost.id}?token=${encodeURIComponent(token)}`
        
        console.log('WebSocket连接URL:', wsUrl)
        // 创建WebSocket连接并添加认证头
        this.socket = new WebSocket(wsUrl)
        
        // 连接超时检查
        const timeout = setTimeout(() => {
          if (this.socket.readyState === WebSocket.CONNECTING) {
            this.socket.close()
            this.term.writeln('\r\n连接超时，请检查：')
            this.term.writeln('1. 后端WebSocket服务是否运行')
            this.term.writeln('2. 防火墙/网络配置是否允许WebSocket连接')
            this.term.writeln('3. Nginx是否配置了WebSocket代理')
          }
        }, 5000)
        
        this.socket.onopen = () => {
          clearTimeout(timeout)
          // Send authentication information
          const token = storage.getItem('token')
          if (token) {
            this.socket.send(JSON.stringify({
              type: 'auth',
              token: token
            }))
          }
          
          this.isConnected = true
          this.isConnecting = false
          this.term.writeln('终端已连接')
        }
        
        // 手动处理WebSocket消息
        this.socket.onmessage = (event) => {
          this.term.write(event.data)
        }
        
        // 窗口大小调整处理
        this.term.onResize((size) => {
          if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            this.socket.send(JSON.stringify({
              type: 'resize',
              cols: size.cols,
              rows: size.rows
            }))
          }
        })
        
        // 移除AttachAddon，避免自动转发
        
        // 初始调整大小
        this.fitAddon.fit()
        
        this.socket.onclose = (event) => {
          this.isConnected = false
          this.isConnecting = false
          this.term.writeln(`\r\n终端连接已关闭，代码: ${event.code}, 原因: ${event.reason}`)
        }
        
        this.socket.onerror = (error) => {
          this.isConnected = false
          this.isConnecting = false
          this.term.writeln('\r\nWebSocket连接失败，请检查:')
          this.term.writeln('1. 后端服务是否运行在127.0.0.1:8000')
          this.term.writeln('2. WebSocket路由是否配置正确')
          this.term.writeln('3. 防火墙是否允许8000端口连接')
          this.term.writeln(`错误详情: ${error.message || '未知错误'}`)
          console.error('WebSocket连接错误:', error)
        }
      } catch (error) {
        console.error('连接异常详情:', error)
        this.isConnecting = false
        this.term.writeln(`\r\n连接失败: ${error.message}`)
      }
    },
    
    disconnectTerminal() {
      if (this.socket) {
        this.socket.close()
        this.socket = null
      }
      this.isConnected = false
      this.isConnecting = false
      if (this.term) {
        this.term.writeln('\r\nSSH连接已手动断开')
      }
    },
    
    handleResize() {
      if (this.fitAddon) {
        this.fitAddon.fit()
      }
    },
    
    async loadGroupList() {
      try {
        const response = await this.$api.getGroupListWithHosts()
        if (response.data.code === 200) {
          this.groupList = response.data.data.map(group => ({
            ...group,
            children: group.children ? group.children.map(child => ({
              ...child,
              children: child.hosts ? child.hosts.map(host => ({
                id: host.id,
                name: host.hostName,
                isHost: true,
                hostData: host
              })) : []
            })) : []
          }))
          
          if (this.groupList.length > 0) {
            this.expandedKeys = [this.groupList[0].id]
          }
        }
      } catch (error) {
        console.error('加载分组列表失败:', error)
        this.$message.error('加载分组列表失败')
      }
    },
    
    handleGroupClick(node, element) {
      if (element.data.isHost) {
        this.currentHost = {
          id: element.data.hostData.id,
          hostName: element.data.hostData.hostName,
          sshName: element.data.hostData.sshName,
          sshIp: element.data.hostData.sshIp,
          sshPort: element.data.hostData.sshPort
        }
        this.term.clear()
        this.term.writeln(`已选择主机: ${this.currentHost.hostName}`)
        this.term.writeln('点击"连接"按钮开始SSH连接')
      }
    },
    
    handleGroupSearch() {
      if (!this.groupSearchText) {
        this.expandedKeys = []
        return
      }
      
      try {
        const findMatchingGroup = (groups, searchText) => {
          for (const group of groups) {
            if (group.name.includes(searchText)) {
              return group
            }
            if (group.children && group.children.length > 0) {
              const found = findMatchingGroup(group.children, searchText)
              if (found) return found
            }
          }
          return null
        }
        
        const matchingGroup = findMatchingGroup(this.groupList, this.groupSearchText)
        if (matchingGroup) {
          const findPath = (groups, targetId, path = []) => {
            for (const group of groups) {
              if (group.id === targetId) {
                return [...path, group.id]
              }
              if (group.children && group.children.length > 0) {
                const foundPath = findPath(group.children, targetId, [...path, group.id])
                if (foundPath) return foundPath
              }
            }
            return null
          }
          
          const expandPath = findPath(this.groupList, matchingGroup.id)
          if (expandPath) {
            this.expandedKeys = expandPath.slice(0, -1)
            this.$nextTick(() => {
              this.$refs.groupTree.setCurrentKey(matchingGroup.id)
            })
          }
        } else {
          this.$message.warning('未找到匹配的分组')
        }
      } catch (error) {
        console.error('搜索分组失败:', error)
        this.$message.error('搜索分组失败')
      }
    },
    
    handleNodeExpand(data, node) {
      if (!this.expandedKeys.includes(node.key)) {
        this.expandedKeys.push(node.key)
      }
    },
    
    handleNodeCollapse(data, node) {
      this.expandedKeys = this.expandedKeys.filter(key => key !== node.key)
    }
  }
}
</script>

<style scoped>
.ssh-container {
  height: calc(100vh - 120px);
  padding: 20px;
}

.ssh-layout {
  display: flex;
  height: 100%;
}

.tree-section {
  width: 300px;
  margin-right: 20px;
}

.terminal-section {
  flex: 1;
}

.terminal-content {
  height: 550px;
  display: flex;
  flex-direction: column;
}

.host-info {
  margin-bottom: 20px;
}

.terminal-wrapper {
  height: calc(100vh - 120px);
  width: 100%;
  padding: 10px;
  background: #060101;
  overflow: hidden;
  position: relative;
}

.terminal-wrapper .xterm {
  padding: 10px;
  width: 100% !important;
  height: 100% !important;
}

.terminal-wrapper .xterm-viewport {
  width: 100% !important;
  overflow-y: auto;
}

.font-weight-bold {
  font-weight: bold;
}

.terminal-actions {
  display: inline-block;
  margin-left: 20px;
}
</style>
