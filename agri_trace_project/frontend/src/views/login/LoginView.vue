<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-header">
        <h2 class="login-title">🍎 阿克苏苹果区块链溯源系统</h2>
        <p class="login-subtitle">Aksu Apple Blockchain Traceability System</p>
      </div>

      <el-tabs v-model="activeTab">
        <!-- 登录 -->
        <el-tab-pane label="用户登录" name="login">
          <el-form :model="loginForm" :rules="loginRules" ref="loginFormRef" label-position="top">
            <el-form-item label="用户名" prop="username">
              <el-input v-model="loginForm.username" placeholder="请输入用户名" prefix-icon="User" size="large" />
            </el-form-item>
            <el-form-item label="密码" prop="password">
              <el-input v-model="loginForm.password" type="password" placeholder="请输入密码（默认123456）" prefix-icon="Lock" size="large" show-password @keyup.enter="handleLogin" />
            </el-form-item>
            <el-form-item label="用户角色" prop="role">
              <el-select v-model="loginForm.role" placeholder="请选择您的角色" size="large" style="width:100%">
                <el-option v-for="r in roleOptions" :key="r.value" :label="r.label" :value="r.value" />
              </el-select>
            </el-form-item>
            <el-form-item label="验证码" prop="captcha">
              <div class="captcha-row">
                <el-input v-model="loginForm.captcha" placeholder="请输入验证码" size="large" style="flex:1" />
                <div class="captcha-box" @click="refreshCaptcha">
                  <span class="captcha-text">{{ captchaCode }}</span>
                  <span class="captcha-hint">点击刷新</span>
                </div>
              </div>
            </el-form-item>
            <el-button type="primary" size="large" style="width:100%;margin-top:8px" :loading="loading" @click="handleLogin">登 录</el-button>
          </el-form>
          <div class="login-footer">
            <el-link @click="activeTab='register'">没有账号？立即注册</el-link>
            <el-link type="info" @click="router.push('/trace')">消费者溯源查询</el-link>
          </div>
        </el-tab-pane>

        <!-- 注册 -->
        <el-tab-pane label="用户注册" name="register">
          <el-form :model="regForm" :rules="regRules" ref="regFormRef" label-position="top">
            <el-row :gutter="12">
              <el-col :span="12">
                <el-form-item label="用户名" prop="username">
                  <el-input v-model="regForm.username" placeholder="3-20位" size="large" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="真实姓名" prop="real_name">
                  <el-input v-model="regForm.real_name" placeholder="真实姓名" size="large" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="12">
              <el-col :span="12">
                <el-form-item label="密码" prop="password">
                  <el-input v-model="regForm.password" type="password" placeholder="至少6位" size="large" show-password />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="确认密码" prop="confirm_password">
                  <el-input v-model="regForm.confirm_password" type="password" placeholder="再次输入" size="large" show-password />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="12">
              <el-col :span="12">
                <el-form-item label="用户角色" prop="role">
                  <el-select v-model="regForm.role" placeholder="选择角色" size="large" style="width:100%">
                    <el-option v-for="r in roleOptions.filter(r=>r.value!=='admin')" :key="r.value" :label="r.label" :value="r.value" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="手机号" prop="phone">
                  <el-input v-model="regForm.phone" placeholder="手机号" size="large" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-form-item label="验证码" prop="captcha">
              <div class="captcha-row">
                <el-input v-model="regForm.captcha" placeholder="请输入验证码" size="large" style="flex:1" />
                <div class="captcha-box" @click="refreshCaptcha">
                  <span class="captcha-text">{{ captchaCode }}</span>
                  <span class="captcha-hint">点击刷新</span>
                </div>
              </div>
            </el-form-item>
            <el-button type="success" size="large" style="width:100%;margin-top:8px" :loading="loading" @click="handleRegister">注 册</el-button>
          </el-form>
          <div class="login-footer">
            <el-link @click="activeTab='login'">已有账号？返回登录</el-link>
          </div>
        </el-tab-pane>
      </el-tabs>

      <el-alert type="info" :closable="false" style="margin-top:16px">
        <template #title>
          <span style="font-size:12px">
            测试账号（密码均为 123456）：admin / farmer_wang / inspector_zhang / transporter_sun / retailer_zhao
          </span>
        </template>
      </el-alert>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { login, register, getCaptcha } from '../../api/index'
import { useUserStore } from '../../store/user'

const router = useRouter()
const userStore = useUserStore()
const activeTab = ref('login')
const loading = ref(false)
const captchaKey = ref('')
const captchaCode = ref('----')

const roleOptions = [
  { value: 'admin', label: '🔑 系统管理员' },
  { value: 'farmer', label: '🌾 种植户' },
  { value: 'inspector', label: '🔬 质检员' },
  { value: 'transporter', label: '🚛 物流商' },
  { value: 'retailer', label: '🏪 销售商' },
]

const loginForm = reactive({ username: '', password: '', role: '', captcha: '' })
const regForm = reactive({ username: '', password: '', confirm_password: '', real_name: '', role: '', phone: '', captcha: '' })
const loginFormRef = ref()
const regFormRef = ref()

const loginRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  role: [{ required: true, message: '请选择角色', trigger: 'change' }],
  captcha: [{ required: true, message: '请输入验证码', trigger: 'blur' }],
}

const regRules = {
  username: [{ required: true, min: 3, max: 20, message: '用户名3-20位', trigger: 'blur' }],
  real_name: [{ required: true, message: '请输入真实姓名', trigger: 'blur' }],
  password: [{ required: true, min: 6, message: '密码至少6位', trigger: 'blur' }],
  confirm_password: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { validator: (_: any, v: string, cb: Function) => { v !== regForm.password ? cb(new Error('两次密码不一致')) : cb() }, trigger: 'blur' }
  ],
  role: [{ required: true, message: '请选择角色', trigger: 'change' }],
  captcha: [{ required: true, message: '请输入验证码', trigger: 'blur' }],
}

async function refreshCaptcha() {
  try {
    const res: any = await getCaptcha()
    captchaKey.value = res.data?.key || ''
    captchaCode.value = res.data?.captcha_code || '----'
  } catch {
    captchaKey.value = 'mock_' + Date.now()
    captchaCode.value = String(Math.floor(1000 + Math.random() * 9000))
  }
}

async function handleLogin() {
  await loginFormRef.value?.validate(async (valid: boolean) => {
    if (!valid) return
    loading.value = true
    try {
      const res: any = await login({ username: loginForm.username, password: loginForm.password, captcha_key: captchaKey.value, captcha: loginForm.captcha })
      const data = res.data || res
      userStore.setToken(data.token)
      userStore.setUserInfo({ user_id: data.user_id, username: data.username, real_name: data.real_name, role: data.role })
      ElMessage.success('登录成功，欢迎回来！')
      router.push('/dashboard')
    } catch { refreshCaptcha() }
    finally { loading.value = false }
  })
}

async function handleRegister() {
  await regFormRef.value?.validate(async (valid: boolean) => {
    if (!valid) return
    loading.value = true
    try {
      await register({ username: regForm.username, password: regForm.password, real_name: regForm.real_name, role: regForm.role, phone: regForm.phone, captcha_key: captchaKey.value, captcha: regForm.captcha })
      ElMessage.success('注册成功，请登录')
      activeTab.value = 'login'
      loginForm.username = regForm.username
    } catch { refreshCaptcha() }
    finally { loading.value = false }
  })
}

onMounted(() => { refreshCaptcha() })
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #0d1b2a 0%, #1a3a5c 40%, #2d6a4f 100%);
  display: flex; align-items: center; justify-content: center;
}
.login-card {
  width: 500px; background: rgba(255,255,255,0.97);
  border-radius: 16px; padding: 36px 40px 28px;
  box-shadow: 0 20px 60px rgba(0,0,0,.4);
}
.login-header { text-align: center; margin-bottom: 20px; }
.login-title { font-size: 20px; font-weight: 700; color: #1a1a1a; margin: 0 0 4px; }
.login-subtitle { font-size: 12px; color: #999; margin: 0; }
.captcha-row { display: flex; gap: 12px; align-items: center; }
.captcha-box {
  width: 120px; height: 40px;
  background: linear-gradient(135deg, #1677ff, #52c41a);
  border-radius: 6px; display: flex; flex-direction: column;
  align-items: center; justify-content: center; cursor: pointer;
  flex-shrink: 0; user-select: none;
}
.captcha-text { color: #fff; font-size: 20px; font-weight: 700; letter-spacing: 4px; }
.captcha-hint { color: rgba(255,255,255,0.7); font-size: 10px; }
.login-footer { display: flex; justify-content: space-between; margin-top: 16px; font-size: 13px; }
</style>
