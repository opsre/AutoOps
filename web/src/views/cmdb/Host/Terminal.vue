<template>
  <div>
    <!-- 使用自定义弹窗替代el-dialog -->
    <div v-if="visible" class="custom-dialog-wrapper">
      <div class="custom-dialog-overlay" @click="handleOverlayClick"></div>
      <div class="custom-dialog" ref="customDialog" @click.stop>
        <div class="custom-dialog-header">
          <span class="custom-dialog-title">{{ drawerTitle }}</span>
          <button class="custom-dialog-close" @click="onClose">×</button>
        </div>
        <div class="custom-dialog-body">
          <div class="terminal-header">
            <el-button size="mini" v-authority="['cmdb:ecs:connecthost']" @click="connect" :disabled="isConnecting || !currentHost">
              {{ isConnected ? '重新连接' : '连接' }}
            </el-button>
            <el-button size="mini" @click="disconnect" :disabled="!isConnected || isConnecting">
              断开
            </el-button>
            <el-button size="mini" @click="clear" :disabled="!term">
              清屏
            </el-button>
            <el-button size="mini" class="close-button" @click="onClose">
              关闭
            </el-button>
          </div>
          <div id="xterm" class="xterm" style="height: calc(100% - 50px);" />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import "xterm/css/xterm.css";
import { Terminal } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import storage from '@/utils/storage'

export default {
  mounted() {
    window.addEventListener('resize', this.handleResize);
  },
  beforeUnmount() {
    window.removeEventListener('resize', this.handleResize);
    if (this.socket) {
      this.socket.close();
      this.socket = null;
    }
    // 不再销毁终端实例，以便可以重用
    this.isConnected = false;
    this.isConnecting = false;
    this._inputHandlerBound = false;
    this._termInitialized = false;
  },
  name: 'Terminal',
  props: {
    currentHost: {
      type: Object,
      default: null
    }
  },
  data() {
    return {
      visible: false,
      term: null,
      fitAddon: null,
      socket: null,
      isConnecting: false,
      isConnected: false,
      drawerTitle: 'SSH终端',
      _inputHandlerBound: false,
      _termInitialized: false
    };
  },
  methods: {
    async initTerm() {
      try {
        if (this._termInitialized) {
          return;
        }

        if (this.term) {
          try {
            if (this.fitAddon && this.fitAddon._disposed === false) {
              this.fitAddon.dispose();
            }
            this.term.dispose();
          } catch (e) {
            console.warn('Terminal dispose error:', e);
          }
        }

        this._termInitialized = true;

        this.term = new Terminal({
          fontSize: 12,
          cursorBlink: true,
          theme: {
            background: '#060101',
            foreground: 'white',
            cursor: 'help',
            selection: '#ffffff80'
          },
          fontFamily: 'Menlo, Monaco, monospace',
          cols: 250,
          rows: 30,
          convertEol: false,
          allowTransparency: false,
          windowsMode: false,
          lineHeight: 1.0,
          letterSpacing: 0,
          scrollback: 1000,
          tabStopWidth: 8,
          bellStyle: 'none',
          rendererType: 'canvas',
          disableStdin: false,
          rightClickSelectsWord: true,
          drawBoldTextInBrightColors: true,
          experimentalCharAtlas: 'dynamic',
          allowProposedApi: true,
          screenReaderMode: false,
          wordSeparator: ' ()[]{}\'"`',
          cursorStyle: 'underline',
          cursorBlink: true,
          scrollSensitivity: 1,
          fastScrollModifier: 'alt',
          macOptionClickForcesSelection: true,
          // 强制禁用全局设置覆盖
          overrideGlobalSettings: true,
          // 确保正确处理ANSI颜色代码
          allowProposedApi: true,
          experimentalCharAtlas: 'none'
        });

        this.fitAddon = new FitAddon();
        const terminalElement = document.getElementById("xterm");
        terminalElement.style.width = '100%';
        terminalElement.style.overflow = 'hidden';
        
        this.term.loadAddon(this.fitAddon);
        this.term.open(terminalElement);

        // 绑定输入事件 - 确保只绑定一次
        if (!this._inputHandlerBound) {
          this._inputHandlerBound = true;
          let lastInputTime = 0;
          this.term.onData(data => {
            const now = Date.now();
            // 详细输入日志
            console.group('Terminal Input');
            const hasNewline = data.includes('\n') || data.includes('\r');

            console.groupEnd();
            
            // 原有逻辑保持不变
            if (now - lastInputTime < 50) {
              return;
            }
            lastInputTime = now;

            if (this.socket && this.socket.readyState === WebSocket.OPEN) {
              // 添加输入内容处理逻辑
              if (hasNewline) {
              }
              this.socket.send(data);
            }
          });
        }

        setTimeout(() => {
          this.fitAddon.fit();
          this.term.focus();
        }, 100);

        this.term.clear();
        this.term.writeln('\x1B[1;32m欢迎使用SSH终端\x1B[0m');
        this.term.writeln('\x1B[1;34m请点击"连接"按钮开始SSH连接\x1B[0m');
        this.term.writeln('\x1B[1;33m终端宽度: ' + this.term.cols + ' 列\x1B[0m');

      } catch (error) {
        console.error('Terminal initialization failed:', error);
        throw error;
      }
    },

    // 其他方法保持不变...
    show() {
      this.visible = true;
      
      this.$nextTick(() => {
        this.applyCustomDialogStyles();
        if (!this.term || this.term._core._isDisposed) {
          this.initTerm();
        } else {
          this.term.clear();
          this.term.writeln('\x1B[1;32m终端已重新打开\x1B[0m');
          this.term.writeln('\x1B[1;34m可以点击"连接"按钮开始SSH连接\x1B[0m');
        }
      });
    },
    
    handleOverlayClick() {
      // 点击遮罩层不关闭弹框
    },
    
    applyCustomDialogStyles() {
      if (this.$refs.customDialog) {
        const dialog = this.$refs.customDialog
        dialog.style.position = 'fixed'
        dialog.style.left = '550px'  
        dialog.style.right = '30px'  
        dialog.style.top = '60px'
        dialog.style.bottom = '10px'
        dialog.style.backgroundColor = '#2a3f54'
        dialog.style.border = '3px solid #2a3f54'
        dialog.style.borderRadius = '15px'
        dialog.style.boxShadow = '0 10px 30px rgba(0, 0, 0, 0.3), inset 0 0 20px rgba(0, 0, 0, 0.2)'
        dialog.style.display = 'flex'
        dialog.style.flexDirection = 'column'
        dialog.style.zIndex = '2000'
      }
    },
    
    async connectTerminal() {
      if (!this.currentHost || this.isConnecting) return;
      
      if (this.socket) {
        this.disconnectTerminal();
        await new Promise(resolve => setTimeout(resolve, 500));
      }
      
      this.isConnecting = true;
      try {
        const token = storage.getItem('token');
        const getWsBaseUrl = () => {
          const baseUrl = (process.env.VUE_APP_API_BASE_URL || '').replace(/\/$/, '')
          if (baseUrl.startsWith('http')) {
            return baseUrl.replace(/^http/, 'ws')
          }
          const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
          return `${protocol}//${window.location.host}${baseUrl}`
        }

        const wsUrl = `${getWsBaseUrl()}/api/v1/cmdb/hostssh/connect/${this.currentHost.id}?token=${encodeURIComponent(token)}`;

        this.socket = new WebSocket(wsUrl);
        
        const timeout = setTimeout(() => {
          if (this.socket.readyState === WebSocket.CONNECTING) {
            this.socket.close();
            this.term.write('\r\n\x1B[31m连接超时，请检查:\x1B[0m\r\n');
            this.term.write('1. 后端WebSocket服务是否运行\r\n');
            this.term.write('2. 防火墙/网络配置是否允许连接\r\n');
          }
        }, 5000);
        
        this.socket.onopen = () => {
          clearTimeout(timeout);
          this.isConnected = true;
          this.isConnecting = false;
          this.drawerTitle = `SSH终端 - ${this.currentHost.hostName}`;
          this.term.write('\x1B[32m终端正在连接中\x1B[0m\r\n');
          this.fitTerminal();
        };
        
        this.socket.onmessage = e => {
          try {
            let data = e.data;
            // 详细输出日志
            console.group('Terminal Output');
            console.groupEnd();
            
            try {
              const jsonData = JSON.parse(data);
              // 过滤掉resize消息
              if (jsonData && jsonData.resize) {
                return;
              }
              if (jsonData && jsonData.data) {
                data = jsonData.data;
              }
            } catch (jsonError) {
              // 不是JSON格式，保持原样处理
            }
            
            // 添加输出内容处理逻辑
            const hasNewline = data.includes('\n') || data.includes('\r');
            
            // 处理自动换行问题
            if (hasNewline) {
            }
            this.term.write(data);
          } catch (error) {
            console.error('Terminal write error:', error);
          }
        };
        
        this.socket.onclose = e => {
          clearTimeout(timeout);
          this.isConnected = false;
          this.isConnecting = false;
          let message = `\r\n\r\n\x1B[31m连接超时已关闭，请检查配置 (code: ${e.code}`;
          if (e.reason) message += `, reason: ${e.reason}`;
          message += ')\x1B[0m\r\n';
          setTimeout(() => this.term.write(message), 200);
        };
        
        this.socket.onerror = error => {
          clearTimeout(timeout);
          this.isConnected = false;
          this.isConnecting = false;
          this.term.write('\r\n\x1B[31mWebSocket连接错误:\x1B[0m\r\n');
          this.term.write(`- ${error.message || '未知错误'}\r\n`);
          this.term.write('请检查:\r\n');
          this.term.write('1. 后端服务是否运行\r\n');
          this.term.write('2. WebSocket路由配置是否正确\r\n');
          this.term.write('3. 防火墙是否允许连接\r\n');
        };
        
        this.term.onResize(({cols, rows}) => {
          if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            this.socket.send(JSON.stringify({resize: [cols, rows]}));
          }
        });
      } catch (error) {
        this.isConnecting = false;
        this.term.write(`\r\n\x1B[31m连接失败: ${error.message}\x1B[0m\r\n`);
      }
    },
    
    fitTerminal() {
      if (this.fitAddon) {
        this.fitAddon.fit();
        this.term.focus();
      }
    },
    
    disconnectTerminal() {
      if (this.socket) {
        this.socket.close();
        this.socket = null;
      }
      this.isConnected = false;
      this.isConnecting = false;
      this.drawerTitle = 'SSH终端';
      
      if (this.term) {
        this.term.writeln('\r\nSSH连接已手动断开');
        this.term.writeln('可以点击"连接"按钮重新连接');
      }
    },
    
    clearTerminal() {
      if (this.term) {
        this.term.clear();
        this.term.writeln('终端已清屏');
      }
    },
    
    connect() {
      if (!this._termInitialized) {
        this.initTerm().then(() => {
          this.connectTerminal();
        });
      } else {
        this.connectTerminal();
      }
    },
    disconnect() {
      this.disconnectTerminal();
    },
    clear() {
      this.clearTerminal();
    },
    onClose() {
      this.visible = false;
      try {
        // 重置所有连接相关状态
        if (this.socket) {
          this.socket.close();
          this.socket = null;
        }
        this.isConnected = false;
        this.isConnecting = false;
        this._termInitialized = false; // 允许重新初始化
        
        if (this.term) {
          this.term.clear();
          this.term.writeln('\x1B[1;32m终端已关闭\x1B[0m');
          this.term.writeln('\x1B[1;34m可以重新打开终端连接\x1B[0m');
          
          // 确保终端实例不被销毁
          if (this.fitAddon && this.fitAddon._disposed === false) {
            this.fitAddon.dispose();
          }
        }
      } catch (error) {
        console.error('Error while closing terminal:', error);
      }
    },
    
    handleResize() {
      this.fitTerminal();
    },
    
    handleKeyEvent(arg) {
      if (arg.code === 'PageUp' && arg.type === 'keydown') {
        this.term.scrollPages(-1);
        return false;
      } else if (arg.code === 'PageDown' && arg.type === 'keydown') {
        this.term.scrollPages(1);
        return false;
      }
      return true;
    }
  }
}
</script>

<style scoped>
/* 样式保持不变 */
.custom-dialog-wrapper {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 1999;
}

.custom-dialog-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
}

.custom-dialog {
  position: absolute;
  left: 1000px;
  right: 10px;
  top: 30px;
  bottom: 30px;
  background-color: #2a3f54;
  border: 3px solid #2a3f54;
  border-radius: 15px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3), inset 0 0 20px rgba(0, 0, 0, 0.2);
  display: flex;
  flex-direction: column;
  z-index: 2000;
  overflow: hidden;
}

.custom-dialog-header {
  padding: 15px 20px;
  border-bottom: 1px solid #00ff88;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #2a3f54;
}

.custom-dialog-title {
  font-size: 18px;
  font-weight: bold;
  color: #00ff88;
}

.custom-dialog-close {
  background: transparent;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #ff4d4f;
}

.custom-dialog-close:hover {
  color: #ff7875;
}

.custom-dialog-body {
  flex: 1;
  overflow: auto;
  background-color: #2a3f54;
  display: flex;
  flex-direction: column;
}

.terminal-header {
  padding: 10px;
  border-bottom: 1px solid #00ff88;
  margin-bottom: 10px;
  background-color: #2a3f54;
}

/* 统一按钮样式 */
.terminal-header .el-button {
  background-color: #3a4f64 !important;
  border-color: #00ff88 !important;
  color: #00ff88 !important;
  padding: 4px 8px;
  font-size: 12px;
  min-height: 24px;
}

.terminal-header .el-button:hover {
  background-color: #4a5f74 !important;
  border-color: #33ff99 !important;
  color: #33ff99 !important;
}

.terminal-header .el-button:active {
  background-color: #2a3f54 !important;
  border-color: #00cc66 !important;
  color: #00cc66 !important;
}

/* 关闭按钮特殊样式 */
.terminal-header .close-button {
  background-color: #3a4f64 !important;
  border-color: #ff4d4f !important;
  color: #ff4d4f !important;
}

.terminal-header .close-button:hover {
  background-color: #4a5f74 !important;
  border-color: #ff7875 !important;
  color: #ff7875 !important;
}

.terminal-header .close-button:active {
  background-color: #2a3f54 !important;
  border-color: #d9363e !important;
  color: #d9363e !important;
}

.xterm {
  padding: 10px;
  background-color: #060101;
  white-space: nowrap; /* 防止自动换行 */
  overflow: hidden;
}
</style>
