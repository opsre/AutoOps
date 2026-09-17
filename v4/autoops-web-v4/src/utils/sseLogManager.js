/**
 * SSE 实时日志管理工具类
 * 用于管理 Server-Sent Events 连接和日志流处理
 */
export class SSELogManager {
  constructor() {
    this.eventSource = null
    this.isConnected = false
    this.listeners = new Map()
    this.reconnectAttempts = 0
    this.maxReconnectAttempts = 5
    this.reconnectDelay = 1000
  }

  /**
   * 连接到SSE日志流
   * @param {string} url - SSE端点URL
   * @param {Object} options - 连接选项
   */
  connect(url, options = {}) {
    if (this.eventSource) {
      this.disconnect()
    }

    try {
      console.log('正在创建SSE连接:', url)
      this.eventSource = new EventSource(url, {
        withCredentials: true,
        ...options
      })

      this.setupEventListeners()
      
      return new Promise((resolve, reject) => {
        const openHandler = () => {
          console.log('SSE连接成功建立')
          this.isConnected = true
          this.reconnectAttempts = 0
          this.emit('connected')
          resolve(this.eventSource)
        }

        const errorHandler = (error) => {
          console.error('SSE连接错误:', {
            error,
            readyState: this.eventSource?.readyState,
            url: url,
            reconnectAttempts: this.reconnectAttempts
          })
          this.emit('error', error)
          if (this.reconnectAttempts < this.maxReconnectAttempts) {
            console.log(`SSE重试连接 ${this.reconnectAttempts + 1}/${this.maxReconnectAttempts}`)
            this.scheduleReconnect()
          } else {
            console.error('SSE连接失败，已达到最大重试次数')
            reject(new Error('SSE连接失败，已达到最大重试次数'))
          }
        }

        this.eventSource.addEventListener('open', openHandler, { once: true })
        this.eventSource.addEventListener('error', errorHandler)
        
        // 添加超时检查
        setTimeout(() => {
          if (!this.isConnected) {
            console.error('SSE连接超时')
            errorHandler(new Error('连接超时'))
          }
        }, 10000) // 10秒超时
      })
    } catch (error) {
      console.error('SSE连接创建失败:', error)
      throw error
    }
  }

  /**
   * 设置事件监听器
   */
  setupEventListeners() {
    if (!this.eventSource) return

    // 监听日志数据
    this.eventSource.addEventListener('log', (event) => {
      try {
        const data = JSON.parse(event.data)
        this.emit('log', data)
      } catch (error) {
        console.warn('解析日志数据失败:', error)
        this.emit('log', { content: event.data, timestamp: new Date().toISOString() })
      }
    })

    // 监听状态更新
    this.eventSource.addEventListener('status', (event) => {
      try {
        const data = JSON.parse(event.data)
        this.emit('status', data)
      } catch (error) {
        console.warn('解析状态数据失败:', error)
      }
    })

    // 监听完成事件
    this.eventSource.addEventListener('complete', (event) => {
      try {
        const data = JSON.parse(event.data)
        this.emit('complete', data)
        this.disconnect()
      } catch (error) {
        console.warn('解析完成数据失败:', error)
      }
    })

    // 监听连接错误
    this.eventSource.addEventListener('error', (event) => {
      this.isConnected = false
      console.error('SSE连接错误:', event)
      this.emit('error', event)
    })

    // 监听连接关闭
    this.eventSource.addEventListener('close', () => {
      this.isConnected = false
      this.emit('disconnected')
    })
  }

  /**
   * 断开SSE连接
   */
  disconnect() {
    if (this.eventSource) {
      this.eventSource.close()
      this.eventSource = null
      this.isConnected = false
      this.emit('disconnected')
    }
  }

  /**
   * 计划重连
   */
  scheduleReconnect() {
    this.reconnectAttempts++
    const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1)
    
    console.log(`SSE连接断开，${delay}ms后进行第${this.reconnectAttempts}次重连...`)
    
    setTimeout(() => {
      if (!this.isConnected && this.eventSource) {
        // 重新创建连接
        const url = this.eventSource.url
        this.connect(url)
      }
    }, delay)
  }

  /**
   * 添加事件监听器
   * @param {string} event - 事件类型
   * @param {Function} callback - 回调函数
   */
  on(event, callback) {
    if (!this.listeners.has(event)) {
      this.listeners.set(event, new Set())
    }
    this.listeners.get(event).add(callback)
  }

  /**
   * 移除事件监听器
   * @param {string} event - 事件类型
   * @param {Function} callback - 回调函数
   */
  off(event, callback) {
    if (this.listeners.has(event)) {
      this.listeners.get(event).delete(callback)
    }
  }

  /**
   * 触发事件
   * @param {string} event - 事件类型
   * @param {any} data - 事件数据
   */
  emit(event, data) {
    if (this.listeners.has(event)) {
      this.listeners.get(event).forEach(callback => {
        try {
          callback(data)
        } catch (error) {
          console.error(`SSE事件处理器错误 [${event}]:`, error)
        }
      })
    }
  }

  /**
   * 获取连接状态
   */
  getConnectionState() {
    return {
      isConnected: this.isConnected,
      readyState: this.eventSource?.readyState,
      reconnectAttempts: this.reconnectAttempts
    }
  }
}

// 默认导出单例
export default new SSELogManager()