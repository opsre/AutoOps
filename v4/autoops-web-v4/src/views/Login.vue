<template>
  <div class="login_container">
    <!-- 保留原有的粒子效果 canvas -->
    <canvas id="particles" class="particles-canvas"></canvas>
    <!-- 保留原有的黑客代码雨 canvas -->
    <canvas id="matrix" class="matrix-canvas"></canvas>
    
    <!-- 在原有效果之上增加科技感装饰 -->
    <div class="floating-orbs">
      <div class="orb orb-1"></div>
      <div class="orb orb-2"></div>
      <div class="orb orb-3"></div>
    </div>
    
    <!-- 动态网格叠加 -->
    <div class="grid-overlay"></div>

    <div class="login_box">
      <!-- 玻璃拟态装饰 -->
      <div class="glass-decoration glass-1"></div>
      <div class="glass-decoration glass-2"></div>
      
      <el-form class="login_form" ref="loginFormRef" :rules="rules" :model="loginForm">
        <!-- 升级后的科技感标题 -->
        <div class="title-container">
          <div class="title-with-logo">
            <div class="logo-icon">
              <svg viewBox="0 0 24 24" class="tech-logo">
                <path d="M12 2L2 7v10l10 5 10-5V7z" fill="none" stroke="currentColor" stroke-width="2"/>
                <path d="m16.5 9.4-5 3-5-3" stroke="currentColor" stroke-width="2" fill="none"/>
                <path d="M7.5 12.75l5 3 5-3" stroke="currentColor" stroke-width="2" fill="none"/>
              </svg>
            </div>
            <div class="title">AutoOps天枢运维平台</div>
          </div>
          <div class="subtitle">Orchestrate Chaos. Automate Wisdom</div>
        </div>
        
        <!-- 升级后的账号输入 -->
        <el-form-item prop="username" class="form-item">
          <div class="input-wrapper">
            <div class="input-icon">
              <svg viewBox="0 0 24 24">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" stroke="currentColor" fill="none" stroke-width="2"/>
                <circle cx="12" cy="7" r="4" stroke="currentColor" fill="none" stroke-width="2"/>
              </svg>
            </div>
            <el-input placeholder="账号" v-model="loginForm.username" class="tech-input" clearable />
          </div>
        </el-form-item>
        
        <!-- 升级后的密码输入 -->
        <el-form-item prop="password" class="form-item">
          <div class="input-wrapper">
            <div class="input-icon">
              <svg viewBox="0 0 24 24">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2" stroke="currentColor" fill="none" stroke-width="2"/>
                <circle cx="12" cy="16" r="1" stroke="currentColor" fill="none" stroke-width="2"/>
                <path d="m7 11V7a5 5 0 0 1 10 0v4" stroke="currentColor" fill="none" stroke-width="2"/>
              </svg>
            </div>
            <el-input placeholder="密码" v-model="loginForm.password" type="password" show-password class="tech-input" clearable />
          </div>
        </el-form-item>
        
        <!-- 升级后的验证码输入 -->
        <el-form-item prop="image" class="form-item">
          <div class="captcha-wrapper">
            <div class="input-wrapper captcha-input">
              <div class="input-icon">
                <svg viewBox="0 0 24 24">
                  <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" stroke="currentColor" fill="none" stroke-width="2"/>
                  <circle cx="12" cy="12" r="3" stroke="currentColor" fill="none" stroke-width="2"/>
                </svg>
              </div>
              <el-input placeholder="验证码" maxlength="6" v-model="loginForm.image" class="tech-input" clearable />
            </div>
            <div class="captcha-image" @click="getCaptcha">
              <el-image class="captchaImg" :src="image" />
              <div class="refresh-overlay">
                <svg viewBox="0 0 24 24" class="refresh-icon">
                  <path d="M1 4v6h6" stroke="currentColor" fill="none" stroke-width="2"/>
                  <path d="M3.51 15a9 9 0 1 0 2.13-9.36L1 10" stroke="currentColor" fill="none" stroke-width="2"/>
                </svg>
              </div>
            </div>
          </div>
        </el-form-item>
        
        <!-- 升级后的按钮组 -->
        <el-form-item class="button-group">
          <div class="tech-buttons">
            <button class="tech-btn primary-btn" @click="loginBtn" type="button">
              <span class="btn-text">登录系统</span>
              <div class="btn-glow"></div>
              <svg class="btn-icon" viewBox="0 0 24 24">
                <path d="m9 18 6-6-6-6" stroke="currentColor" fill="none" stroke-width="2"/>
              </svg>
            </button>
            <button class="tech-btn secondary-btn" @click="resetLoginForm" type="button">
              <span class="btn-text">重置表单</span>
              <div class="btn-glow"></div>
              <svg class="btn-icon" viewBox="0 0 24 24">
                <path d="M1 4v6h6" stroke="currentColor" fill="none" stroke-width="2"/>
                <path d="M3.51 15a9 9 0 1 0 2.13-9.36L1 10" stroke="currentColor" fill="none" stroke-width="2"/>
              </svg>
            </button>
          </div>
        </el-form-item>
      </el-form>
      
      <!-- 科技感底部装饰 -->
      <div class="bottom-decoration">
        <div class="tech-line"></div>
        <div class="version-info">执天机之枢，驭智能运维</div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "Login",
  data() {
    return {
      image: '',
      rules: {
        username: [
          { required: true, message: "请输入账号", trigger: "blur" }
        ],
        password: [
          { required: true, message: "请输入密码", trigger: "blur" }
        ],
        image: [
          { required: true, message: "请输入验证码", trigger: "blur" }
        ]
      },
      loginForm: {
        username: '',
        password: '',
        image: '',
        idKey: ''
      }
    }
  },
  methods: {
    // 保留原有的初始化黑客代码雨
    initMatrixRain() {
      const canvas = document.getElementById('matrix')
      const ctx = canvas.getContext('2d')

      canvas.height = window.innerHeight
      canvas.width = window.innerWidth

      const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=[]{}|;:,.<>?/~'.split('')
      const fontSize = 16
      const columns = canvas.width / fontSize
      ctx.font = `normal 400 ${fontSize}px monospace`
      const drops = Array(Math.floor(columns)).fill(0)

      function draw() {
        ctx.fillStyle = 'rgba(0, 0, 0, 0.05)'
        ctx.fillRect(0, 0, canvas.width, canvas.height)

        for (let i = 0; i < drops.length; i++) {
          const text = chars[Math.floor(Math.random() * chars.length)]
          ctx.fillStyle = '#0f0'
          ctx.fillText(text, i * fontSize, drops[i] * fontSize)
          
          if (drops[i] * fontSize > canvas.height && Math.random() > 0.975) {
            drops[i] = 0
          }
          drops[i]++
        }
      }

      setInterval(draw, 33)
    },

    // 获取验证码
    async getCaptcha() {
      const { data: res } = await this.$api.captcha()
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.image = res.data.image
        this.loginForm.idKey = res.data.idKey
      }
    },

    // 登录按钮
    loginBtn() {
      this.$refs.loginFormRef.validate(async valid => {
        if (valid) {
          const { data: res } = await this.$api.login(this.loginForm)

          console.log("【Login.vue】登录接口返回完整数据：", res)

          if (res.code !== 200) {
            this.$message.error(res.message)
          } else {
            this.$message.success("登录成功")
            console.log("【Login.vue】res.data 内容：", res.data)

            this.$store.commit('saveSysAdmin', res.data.sysAdmin)
            this.$store.commit('saveToken', res.data.token)
            this.$store.commit('saveLeftMenuList', res.data.leftMenuList)
            this.$store.commit('savePermissionList', res.data.permissionList)

            await this.$router.push("/home")
          }
        } else {
          return false
        }
      })
    },

    // 重置按钮
    resetLoginForm() {
      this.$refs.loginFormRef.resetFields()
    },

    // 保留原有的初始化粒子效果
    initParticles() {
      const canvas = document.getElementById('particles')
      const ctx = canvas.getContext('2d')
      
      // 设置canvas尺寸
      canvas.height = window.innerHeight
      canvas.width = window.innerWidth
      
      const particles = []
      const particleCount = 200  // 增加粒子数量
      const mouse = { x: null, y: null }
      
      // 监听鼠标移动
      window.addEventListener('mousemove', (e) => {
        mouse.x = e.x
        mouse.y = e.y
      })
      
      // 创建粒子
      for (let i = 0; i < particleCount; i++) {
        particles.push({
          x: Math.random() * canvas.width,
          y: Math.random() * canvas.height,
          size: Math.random() * 5 + 2,  // 增大粒子尺寸
          speedX: Math.random() * 3 - 1.5,
          speedY: Math.random() * 3 - 1.5,
          color: `rgba(200, 200, 200, ${Math.random() * 0.5 + 0.3})`  // 浅灰色粒子
        })
      }
      
      function animateParticles() {
        ctx.clearRect(0, 0, canvas.width, canvas.height)
        
        // 更新和绘制粒子
        for (let i = 0; i < particles.length; i++) {
          const p = particles[i]
          
          // 计算粒子到鼠标的距离
          let dx = mouse.x - p.x
          let dy = mouse.y - p.y
          let distance = Math.sqrt(dx * dx + dy * dy)
          
          // 鼠标吸引力效果
          if (distance < 100) {
            p.speedX += dx * 0.0005
            p.speedY += dy * 0.0005
          }
          
          // 更新位置
          p.x += p.speedX
          p.y += p.speedY
          
          // 边界检查
          if (p.x < 0 || p.x > canvas.width) p.speedX *= -1
          if (p.y < 0 || p.y > canvas.height) p.speedY *= -1
          
          // 绘制粒子
          ctx.fillStyle = p.color
          ctx.beginPath()
          ctx.arc(p.x, p.y, p.size, 0, Math.PI * 2)
          ctx.fill()
          
          // 绘制连线
          if (distance < 100) {
            ctx.strokeStyle = 'rgba(220, 220, 220, 0.5)'
            ctx.lineWidth = 0.3
            ctx.beginPath()
            ctx.moveTo(p.x, p.y)
            ctx.lineTo(mouse.x, mouse.y)
            ctx.stroke()
          }
        }
        
        requestAnimationFrame(animateParticles)
      }
      
      animateParticles()
    }
  },
  mounted() {
    this.getCaptcha()
    this.initMatrixRain() // 保留原有代码雨
    this.initParticles() // 保留原有粒子效果
  }
}
</script>

<style lang="less" scoped>
// =================== 保留你原有的所有样式 ===================
.login_container {
  position: relative;
  background-image: url("../assets/image/1.jpeg"); // 保留原有背景图
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  height: 100vh;
  width: 100vw;
  z-index: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

// 保留原有的代码雨样式
.matrix-canvas {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 1;
  opacity: 0.3;
  mix-blend-mode: lighten;
}

// 保留原有的粒子样式
.particles-canvas {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 0;
  pointer-events: none;
}

// =================== 新增科技感装饰（不影响原有效果）===================
// 新增浮动光球装饰
.floating-orbs {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 2;
  pointer-events: none;
}

.orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(40px);
  animation: float 8s ease-in-out infinite;
  
  &.orb-1 {
    width: 150px;
    height: 150px;
    background: radial-gradient(circle, rgba(0, 212, 255, 0.15) 0%, transparent 70%);
    top: 15%;
    left: 20%;
    animation-delay: 0s;
  }
  
  &.orb-2 {
    width: 200px;
    height: 200px;
    background: radial-gradient(circle, rgba(0, 153, 255, 0.1) 0%, transparent 70%);
    top: 60%;
    right: 15%;
    animation-delay: 3s;
  }
  
  &.orb-3 {
    width: 120px;
    height: 120px;
    background: radial-gradient(circle, rgba(64, 224, 255, 0.2) 0%, transparent 70%);
    bottom: 25%;
    left: 35%;
    animation-delay: 6s;
  }
}

@keyframes float {
  0%, 100% { transform: translateY(0px) rotate(0deg); }
  50% { transform: translateY(-20px) rotate(180deg); }
}

// 新增动态网格叠加
.grid-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image: 
    linear-gradient(rgba(0, 212, 255, 0.03) 1px, transparent 1px),
    linear-gradient(90deg, rgba(0, 212, 255, 0.03) 1px, transparent 1px);
  background-size: 50px 50px;
  animation: gridMove 15s linear infinite;
  z-index: 2;
  pointer-events: none;
}

@keyframes gridMove {
  0% { transform: translate(0, 0); }
  100% { transform: translate(50px, 50px); }
}

// =================== 升级登录框样式（保持原有尺寸位置）===================
.login_box {
  width: 400px; // 保持原有尺寸
  height: 380px; // 增加高度，避免按钮被遮挡
  background: rgba(0, 0, 0, 0.6); // 增加背景不透明度，使用深色背景
  backdrop-filter: blur(20px); // 增加背景模糊
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 20px; // 稍微增加圆角
  box-shadow: 
    0 8px 32px 0 rgba(0, 0, 0, 0.4),
    0 2px 16px 0 rgba(0, 212, 255, 0.1),
    inset 0 1px 0 0 rgba(255, 255, 255, 0.1);
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  z-index: 10;
  animation: boxGlow 4s ease-in-out infinite alternate;
  overflow: hidden;
}

@keyframes boxGlow {
  0% { 
    box-shadow: 
      0 8px 32px 0 rgba(0, 0, 0, 0.4),
      0 2px 16px 0 rgba(0, 212, 255, 0.1),
      inset 0 1px 0 0 rgba(255, 255, 255, 0.1);
  }
  100% { 
    box-shadow: 
      0 8px 32px 0 rgba(0, 0, 0, 0.6),
      0 4px 24px 0 rgba(0, 212, 255, 0.3),
      inset 0 1px 0 0 rgba(0, 212, 255, 0.2);
  }
}

// 玻璃装饰元素
.glass-decoration {
  position: absolute;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.1) 0%, transparent 100%);
  border-radius: 50%;
  filter: blur(1px);
  
  &.glass-1 {
    width: 80px;
    height: 80px;
    top: -15px;
    right: -15px;
  }
  
  &.glass-2 {
    width: 60px;
    height: 60px;
    bottom: -10px;
    left: -10px;
  }
}

.login_form {
  position: relative; // 改为相对定位，不再贴底部
  padding: 20px;
  box-sizing: border-box;
  width: 100%;
  
  .form-item {
    margin-bottom: 14px; // 减少间距，让布局更紧凑
    
    &:deep(.el-form-item__error) {
      color: #ff6b6b;
      font-size: 11px;
    }
  }
}

// =================== 科技感标题 ===================
.title-container {
  text-align: center;
  margin-bottom: 20px; // 减少底部间距
  
  .title-with-logo {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    margin-bottom: 4px;
    
    .logo-icon {
      width: 32px; // 稍微缩小图标
      height: 32px;
      
      .tech-logo {
        width: 100%;
        height: 100%;
        color: #00d4ff;
        filter: drop-shadow(0 0 6px rgba(0, 212, 255, 0.6));
        animation: logoGlow 3s ease-in-out infinite alternate;
      }
    }
    
    .title {
      font-size: 20px; // 稍微减小字体
      font-weight: 700;
      color: #ffffff;
      text-shadow: 0 0 15px rgba(0, 212, 255, 0.5);
      background: linear-gradient(135deg, #ffffff 0%, #00d4ff 100%);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
    }
  }
  
  .subtitle {
    font-size: 10px; // 稍微减小副标题
    color: rgba(255, 255, 255, 0.5);
    font-weight: 300;
    letter-spacing: 0.5px;
  }
}

@keyframes logoGlow {
  0% { filter: drop-shadow(0 0 6px rgba(0, 212, 255, 0.6)); }
  100% { filter: drop-shadow(0 0 12px rgba(0, 212, 255, 0.9)); }
}

// =================== 升级输入框样式 ===================
.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
  background: rgba(255, 255, 255, 0.12);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  padding: 0 12px;
  transition: all 0.3s ease;
  width: 300px; // 账号和密码输入框宽度设为300px
  margin: 0 auto; // 居中对齐
  transform: translateX(0px); // 账号密码输入框回到居中位置
  
  &:hover {
    background: rgba(255, 255, 255, 0.18);
    border-color: rgba(0, 212, 255, 0.4);
    box-shadow: 0 0 15px rgba(0, 212, 255, 0.15);
  }
  
  &:focus-within {
    background: rgba(255, 255, 255, 0.22);
    border-color: #00d4ff;
    box-shadow: 0 0 15px rgba(0, 212, 255, 0.25);
  }
}

.input-icon {
  width: 16px;
  height: 16px;
  margin-right: 10px;
  color: rgba(255, 255, 255, 0.6);
  
  svg {
    width: 100%;
    height: 100%;
  }
}

.tech-input {
  flex: 1;
  
  &:deep(.el-input__wrapper) {
    background: transparent !important;
    border: none !important;
    box-shadow: none !important;
    padding: 8px 0;
  }
  
  &:deep(.el-input__inner) {
    color: #ffffff !important;
    font-size: 14px;
    background: transparent !important;
    
    &::placeholder {
      color: rgba(255, 255, 255, 0.4) !important;
    }
  }
  
  &:deep(.el-input__suffix) {
    .el-input__suffix-inner {
      color: rgba(255, 255, 255, 0.5);
    }
  }
}

// =================== 验证码区域 ===================
.captcha-wrapper {
  display: flex;
  gap: 10px;
  align-items: center;
  justify-content: center; // 居中对齐验证码区域
  transform: translateX(15px); // 验证码区域整体往右移动15px与账号密码输入框对齐
  
  .captcha-input {
    width: auto; // 保持验证码输入框原有长度，不受300px限制
    flex: none; // 取消flex:1，使用固定宽度
    
    // 彻底移除验证码输入框的所有变换和缩放效果
    * {
      transform: none !important;
      scale: 1 !important;
      transition: background 0.3s ease, border-color 0.3s ease, box-shadow 0.3s ease !important;
    }
    
    .input-wrapper {
      transform: none !important;
      transition: background 0.3s ease, border-color 0.3s ease, box-shadow 0.3s ease !important;
      
      &:hover,
      &:focus,
      &:focus-within,
      &:active {
        transform: none !important;
        scale: 1 !important;
      }
    }
    
    // 针对Element Plus输入框组件的样式覆盖
    &:deep(.el-input) {
      transform: none !important;
      scale: 1 !important;
      
      &:hover,
      &:focus,
      &:active,
      &.is-focus {
        transform: none !important;
        scale: 1 !important;
      }
      
      .el-input__wrapper {
        transform: none !important;
        scale: 1 !important;
        transition: background 0.3s ease, border-color 0.3s ease, box-shadow 0.3s ease !important;
        
        &:hover,
        &:focus,
        &:focus-within,
        &:active,
        &.is-focus,
        &.is-focused {
          transform: none !important;
          scale: 1 !important;
        }
      }
      
      .el-input__inner {
        transform: none !important;
        scale: 1 !important;
        transition: background 0.3s ease, border-color 0.3s ease, box-shadow 0.3s ease, color 0.3s ease !important;
        
        &:hover,
        &:focus,
        &:active,
        &:focus-within,
        &[type="text"]:focus,
        &[type="text"]:active {
          transform: none !important;
          scale: 1 !important;
        }
      }
      
      // 覆盖输入框在输入文字时的所有状态
      &.el-input--focus,
      &:focus-within,
      &:has(input:focus) {
        transform: none !important;
        scale: 1 !important;
        
        .el-input__wrapper,
        .el-input__inner {
          transform: none !important;
          scale: 1 !important;
        }
      }
    }
  }
  
  .captcha-image {
    position: relative;
    width: 120px; // 验证码图片宽度改为120px
    height: 36px;
    border-radius: 6px;
    overflow: hidden;
    cursor: pointer;
    transition: all 0.3s ease;
    background: linear-gradient(135deg, #00d4ff 0%, #0099ff 100%); // 保持与登录按钮相同的天蓝色背景
    
    &:hover .refresh-overlay {
      opacity: 1;
    }
    
    .captchaImg {
      width: 100%;
      height: 100%;
      object-fit: cover;
      border-radius: 6px;
      // 移除所有滤镜效果，保持原始验证码图片不变
      // 通过CSS在验证码图片上方添加深蓝色文字覆盖层
      position: relative;
    }
    
    // 添加文字颜色覆盖层
    .captcha-image::after {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      background: transparent;
      mix-blend-mode: multiply;
      background-image: 
        linear-gradient(45deg, transparent 30%, rgba(0, 0, 128, 0.8) 70%), 
        linear-gradient(-45deg, transparent 30%, rgba(0, 0, 128, 0.8) 70%);
      pointer-events: none;
    }
    
    .refresh-overlay {
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      background: rgba(0, 212, 255, 0.8);
      display: flex;
      align-items: center;
      justify-content: center;
      opacity: 0;
      transition: opacity 0.3s ease;
      
      .refresh-icon {
        width: 16px;
        height: 16px;
        color: white;
        animation: spin 1s linear infinite;
      }
    }
  }
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

// =================== 科技按钮 ===================
.button-group {
  margin-top: 16px; // 减少顶部间距
  margin-bottom: 16px; // 添加底部间距，不再贴边
}

.tech-buttons {
  display: flex;
  gap: 12px;
  width: calc(100% - 30px); // 总宽度减少30px（每个按钮减少15px）
  margin: 0 auto; // 居中对齐
}

.tech-btn {
  position: relative;
  flex: 1;
  height: 40px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  font-size: 14px;
  font-weight: 600;
  transition: all 0.3s ease;
  
  .btn-text {
    position: relative;
    z-index: 2;
  }
  
  .btn-icon {
    width: 14px;
    height: 14px;
    position: relative;
    z-index: 2;
    transition: transform 0.3s ease;
  }
  
  .btn-glow {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    opacity: 0;
    transition: opacity 0.3s ease;
  }
  
  &:hover {
    transform: translateY(-1px);
    
    .btn-icon {
      transform: translateX(2px);
    }
    
    .btn-glow {
      opacity: 1;
    }
  }
  
  &.primary-btn {
    background: linear-gradient(135deg, #00d4ff 0%, #0099ff 100%);
    color: white;
    box-shadow: 0 3px 12px rgba(0, 212, 255, 0.3);
    
    .btn-glow {
      background: linear-gradient(135deg, #40e0ff 0%, #20b9ff 100%);
    }
    
    &:hover {
      box-shadow: 0 5px 16px rgba(0, 212, 255, 0.4);
    }
  }
  
  &.secondary-btn {
    background: rgba(255, 255, 255, 0.1);
    color: rgba(255, 255, 255, 0.9);
    border: 1px solid rgba(255, 255, 255, 0.2);
    
    .btn-glow {
      background: rgba(255, 255, 255, 0.2);
    }
    
    &:hover {
      box-shadow: 0 3px 12px rgba(255, 255, 255, 0.1);
    }
  }
}

// =================== 底部装饰 ===================
.bottom-decoration {
  position: absolute;
  bottom: 8px; // 稍微靠近底部
  left: 50%;
  transform: translateX(-50%);
  text-align: center;
  
  .tech-line {
    width: 60px;
    height: 1px;
    background: linear-gradient(90deg, transparent 0%, #00d4ff 50%, transparent 100%);
    margin: 0 auto 8px;
    animation: pulse 2s ease-in-out infinite;
  }
  
  .version-info {
    display: inline-block;
    padding: 4px 10px;
    border-radius: 12px;
    font-size: 12px;
    font-weight: 500;
    letter-spacing: 1px;
    color: #E6F7FF;
    background: rgba(0, 212, 255, 0.12);
    border: none;
    box-shadow: 0 0 12px rgba(0, 212, 255, 0.25), inset 0 0 8px rgba(0, 212, 255, 0.15);
    text-shadow: 0 0 8px rgba(0, 212, 255, 0.7), 0 0 14px rgba(0, 212, 255, 0.5);
    animation: softGlow 2.8s ease-in-out infinite alternate;
  }
}

@keyframes pulse {
  0%, 100% { opacity: 0.5; transform: scaleX(1); }
  50% { opacity: 1; transform: scaleX(1.2); }
}

@keyframes softGlow {
  0% { box-shadow: 0 0 10px rgba(0, 212, 255, 0.18), inset 0 0 6px rgba(0, 212, 255, 0.12); }
  100% { box-shadow: 0 0 18px rgba(0, 212, 255, 0.35), inset 0 0 10px rgba(0, 212, 255, 0.2); }
}

// =================== 响应式适配 ===================
@media (max-width: 768px) {
  .login_box {
    width: 90%;
    max-width: 350px;
    height: auto;
    min-height: 300px;
  }
  
  .title-container .title {
    font-size: 18px;
  }
  
  .tech-buttons {
    flex-direction: column;
    gap: 8px;
  }
}
</style>
