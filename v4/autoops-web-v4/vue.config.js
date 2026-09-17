/**
 * vue配置中心
 */
const { defineConfig } = require('@vue/cli-service')
const webpack = require('webpack')
const MonacoEditorWebpackPlugin = require('monaco-editor-webpack-plugin')

// Force polling so file changes are detected reliably in Docker/VM mounts
process.env.CHOKIDAR_USEPOLLING = 'true'
process.env.WATCHPACK_POLLING = 'true'

module.exports = defineConfig({
  lintOnSave: false,   // 关闭校验
  productionSourceMap: false, // 关闭生产环境sourceMap
  publicPath:"/",
  configureWebpack: {
    plugins: [
      new MonacoEditorWebpackPlugin({
        languages: ['sql'], // 只打包SQL语言，减小体积
      })
    ],
    watchOptions: {
      ignored: /node_modules/,
      poll: 1000, // 轮询文件变更，解决部分环境下文件监听失效
    }
  },
  outputDir: "dist",
  assetsDir: "static",
  devServer: {
    port: 8080, //前端服务启动的端口
    host: "0.0.0.0",
    hot: true,
    liveReload: true,
    allowedHosts: 'all',
    https: false,
    open:false,
    watchFiles: {
      paths: ['src/**/*', 'public/**/*'],
      options: {
        usePolling: true,
        interval: 1000,
        ignored: /node_modules/,
      }
    },
    proxy: {
      "/api/v1": {
        target: "http://127.0.0.1:8000",  // 后端服务地址
        changeOrigin: true,
        ws: true  // ✅ 关键：启用 WebSocket 代理支持
        // 完全移除 pathRewrite，让 /api/v1 前缀直接传递给后端
      }
    },
    client: {
      overlay: false,  // 关闭全屏错误提示
      webSocketURL: 'auto://0.0.0.0:0/ws', // 适配不同宿主/代理环境的 HMR 连接
      progress: true
    }
  }
})
