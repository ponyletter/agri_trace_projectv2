<template>
  <div class="login-page">
    <div class="login-card">
      <!-- 顶部Logo区 -->
      <div class="login-header">
        <div class="logo-icon">🍎</div>
        <h2 class="login-title">阿克苏苹果区块链溯源系统</h2>
        <p class="login-subtitle">农产品质量安全可信溯源平台</p>
      </div>

      <el-tabs v-model="activeTab" class="login-tabs">
        <!-- ========== 登录 ========== -->
        <el-tab-pane label="账号登录" name="login">
          <el-form :model="loginForm" :rules="loginRules" ref="loginFormRef" label-position="top" class="login-form">
            <el-form-item label="用户名" prop="username">
              <el-input
                v-model="loginForm.username"
                placeholder="请输入用户名"
                prefix-icon="User"
                size="large"
                autocomplete="username"
              />
            </el-form-item>
            <el-form-item label="登录密码" prop="password">
              <el-input
                v-model="loginForm.password"
                type="password"
                placeholder="请输入登录密码"
                prefix-icon="Lock"
                size="large"
                show-password
                autocomplete="current-password"
                @keyup.enter="handleLogin"
              />
            </el-form-item>
            <el-form-item label="用户角色" prop="role">
              <el-select v-model="loginForm.role" placeholder="请选择您的角色" size="large" style="width:100%">
                <el-option v-for="r in roleOptions" :key="r.value" :label="r.label" :value="r.value" />
              </el-select>
            </el-form-item>
            <el-form-item label="验证码" prop="captcha">
              <div class="captcha-row">
                <el-input
                  v-model="loginForm.captcha"
                  placeholder="请输入右侧验证码"
                  size="large"
                  style="flex:1"
                  maxlength="6"
                  @keyup.enter="handleLogin"
                />
                <div class="captcha-img-box" @click="refreshCaptcha" title="点击刷新验证码">
                  <el-icon v-if="captchaLoading" class="captcha-loading"><Loading /></el-icon>
                  <img v-else-if="captchaImgSrc" :src="captchaImgSrc" alt="验证码" class="captcha-img" />
                  <span v-else class="captcha-text-fallback">{{ captchaCode }}</span>
                  <span class="captcha-refresh-tip">点击刷新</span>
                </div>
              </div>
            </el-form-item>
            <el-button
              type="primary"
              size="large"
              style="width:100%;margin-top:4px;font-size:16px;letter-spacing:4px"
              :loading="loading"
              @click="handleLogin"
            >登 录</el-button>
          </el-form>
          <div class="login-footer">
            <el-link type="primary" @click="activeTab='register'">没有账号？立即注册</el-link>
            <el-link type="info" @click="goTrace">消费者溯源查询</el-link>
          </div>
        </el-tab-pane>

        <!-- ========== 注册 ========== -->
        <el-tab-pane label="新用户注册" name="register">
          <el-form :model="regForm" :rules="regRules" ref="regFormRef" label-position="top" class="login-form">
            <el-row :gutter="12">
              <el-col :span="12">
                <el-form-item label="用户名" prop="username">
                  <el-input v-model="regForm.username" placeholder="3-20位字母或数字" size="large" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="真实姓名" prop="real_name">
                  <el-input v-model="regForm.real_name" placeholder="请输入真实姓名" size="large" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="12">
              <el-col :span="12">
                <el-form-item label="登录密码" prop="password">
                  <el-input v-model="regForm.password" type="password" placeholder="至少6位" size="large" show-password />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="确认密码" prop="confirm_password">
                  <el-input v-model="regForm.confirm_password" type="password" placeholder="再次输入密码" size="large" show-password />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="12">
              <el-col :span="12">
                <el-form-item label="用户角色" prop="role">
                  <el-select v-model="regForm.role" placeholder="请选择角色" size="large" style="width:100%">
                    <el-option v-for="r in roleOptions.filter(r => r.value !== 'admin')" :key="r.value" :label="r.label" :value="r.value" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="手机号码" prop="phone">
                  <el-input v-model="regForm.phone" placeholder="请输入手机号码" size="large" maxlength="11" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-form-item label="所属机构" prop="org_name">
              <el-input v-model="regForm.org_name" placeholder="请输入所属单位或机构名称" size="large" />
            </el-form-item>
            <el-form-item label="验证码" prop="captcha">
              <div class="captcha-row">
                <el-input
                  v-model="regForm.captcha"
                  placeholder="请输入右侧验证码"
                  size="large"
                  style="flex:1"
                  maxlength="6"
                />
                <div class="captcha-img-box" @click="refreshCaptcha" title="点击刷新验证码">
                  <el-icon v-if="captchaLoading" class="captcha-loading"><Loading /></el-icon>
                  <img v-else-if="captchaImgSrc" :src="captchaImgSrc" alt="验证码" class="captcha-img" />
                  <span v-else class="captcha-text-fallback">{{ captchaCode }}</span>
                  <span class="captcha-refresh-tip">点击刷新</span>
                </div>
              </div>
            </el-form-item>
            <el-button
              type="success"
              size="large"
              style="width:100%;margin-top:4px;font-size:16px;letter-spacing:4px"
              :loading="loading"
              @click="handleRegister"
            >注 册</el-button>
          </el-form>
          <div class="login-footer">
            <el-link type="primary" @click="activeTab='login'">已有账号？返回登录</el-link>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Loading } from '@element-plus/icons-vue'
import { login, register, getCaptcha } from '../../api/index'
import { useUserStore } from '../../store/user'

const router = useRouter()
const userStore = useUserStore()
const activeTab = ref('login')
const loading = ref(false)

// ========== 验证码状态 ==========
const captchaKey = ref('')
const captchaCode = ref('')       // 文字备用（后端直接返回文本时使用）
const captchaImgSrc = ref('')     // base64图片（后端返回图片时使用）
const captchaLoading = ref(false)

// ========== 角色选项（纯中文，无英文、无emoji） ==========
const roleOptions = [
  { value: 'admin',       label: '系统管理员' },
  { value: 'farmer',      label: '种植户' },
  { value: 'inspector',   label: '质检员' },
  { value: 'transporter', label: '物流商' },
  { value: 'retailer',    label: '销售商' },
]

// ========== 表单数据 ==========
const loginForm = reactive({ username: '', password: '', role: '', captcha: '' })
const regForm = reactive({
  username: '', password: '', confirm_password: '',
  real_name: '', role: '', phone: '', org_name: '', captcha: ''
})
const loginFormRef = ref()
const regFormRef = ref()

// ========== 表单校验规则 ==========
const loginRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入登录密码', trigger: 'blur' }],
  role:     [{ required: true, message: '请选择用户角色', trigger: 'change' }],
  captcha:  [{ required: true, message: '请输入验证码', trigger: 'blur' }],
}

const regRules = {
  username:  [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度为3-20位', trigger: 'blur' }
  ],
  real_name: [{ required: true, message: '请输入真实姓名', trigger: 'blur' }],
  password:  [{ required: true, min: 6, message: '密码至少6位', trigger: 'blur' }],
  confirm_password: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    {
      validator: (_: any, v: string, cb: Function) => {
        v !== regForm.password ? cb(new Error('两次输入的密码不一致')) : cb()
      },
      trigger: 'blur'
    }
  ],
  role:    [{ required: true, message: '请选择用户角色', trigger: 'change' }],
  phone:   [
    { required: true, message: '请输入手机号码', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ],
  captcha: [{ required: true, message: '请输入验证码', trigger: 'blur' }],
}

// ========== 刷新验证码（核心修复） ==========
// 问题根因：request.ts 的响应拦截器已将 response.data.data 解包后返回，
// 所以 getCaptcha() 直接返回 { key, captcha_code } 而非 { code, data: { key, captcha_code } }
// 原代码写的是 res.data?.key，但 res 本身已经是 data，导致取到 undefined，验证码永远不更新
async function refreshCaptcha() {
  captchaLoading.value = true
  captchaImgSrc.value = ''
  captchaCode.value = ''
  try {
    // getCaptcha() 返回值已经是 response.data.data，即 { key: "xxx", captcha_code: "1234" }
    const data: any = await getCaptcha()
    captchaKey.value = data?.key || ''
    if (data?.img) {
      // 后端返回 base64 图片时
      captchaImgSrc.value = data.img
    } else {
      // 后端返回文本验证码时（当前实现）
      captchaCode.value = data?.captcha_code || data?.code || '----'
    }
  } catch {
    // 网络异常降级为本地随机数字
    captchaKey.value = 'local_' + Date.now()
    captchaCode.value = String(Math.floor(1000 + Math.random() * 9000))
  } finally {
    captchaLoading.value = false
  }
}

// ========== 登录 ==========
async function handleLogin() {
  await loginFormRef.value?.validate(async (valid: boolean) => {
    if (!valid) return
    loading.value = true
    try {
      const data: any = await login({
        username:    loginForm.username,
        password:    loginForm.password,
        captcha_key: captchaKey.value,
        captcha:     loginForm.captcha,
      })
      userStore.setToken(data.token)
      userStore.setUserInfo({
        user_id:   data.user_id,
        username:  data.username,
        real_name: data.real_name,
        role:      data.role,
      })
      ElMessage.success('登录成功，欢迎回来！')
      router.push('/dashboard')
    } catch {
      await refreshCaptcha()
      loginForm.captcha = ''
    } finally {
      loading.value = false
    }
  })
}

// ========== 注册 ==========
async function handleRegister() {
  await regFormRef.value?.validate(async (valid: boolean) => {
    if (!valid) return
    loading.value = true
    try {
      await register({
        username:    regForm.username,
        password:    regForm.password,
        real_name:   regForm.real_name,
        role:        regForm.role,
        phone:       regForm.phone,
        org_name:    regForm.org_name,
        captcha_key: captchaKey.value,
        captcha:     regForm.captcha,
      })
      ElMessage.success('注册成功！请使用新账号登录')
      activeTab.value = 'login'
      loginForm.username = regForm.username
      await refreshCaptcha()
    } catch {
      await refreshCaptcha()
      regForm.captcha = ''
    } finally {
      loading.value = false
    }
  })
}

// ========== 跳转溯源查询 ==========
function goTrace() {
  router.push('/trace')
}

// ========== 初始化 ==========
onMounted(() => {
  refreshCaptcha()
})
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #0a1628 0%, #1a3a5c 45%, #1e5c3a 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}
.login-card {
  width: 520px;
  background: rgba(255, 255, 255, 0.98);
  border-radius: 20px;
  padding: 40px 44px 32px;
  box-shadow: 0 24px 80px rgba(0, 0, 0, 0.45);
}
.login-header {
  text-align: center;
  margin-bottom: 24px;
}
.logo-icon {
  font-size: 48px;
  line-height: 1;
  margin-bottom: 10px;
}
.login-title {
  font-size: 20px;
  font-weight: 700;
  color: #1a1a1a;
  margin: 0 0 6px;
  letter-spacing: 1px;
}
.login-subtitle {
  font-size: 13px;
  color: #888;
  margin: 0;
}
.login-tabs :deep(.el-tabs__header) {
  margin-bottom: 20px;
}
.login-form :deep(.el-form-item__label) {
  font-weight: 500;
  color: #333;
}
/* ===== 验证码区域 ===== */
.captcha-row {
  display: flex;
  gap: 12px;
  align-items: stretch;
  width: 100%;
}
.captcha-img-box {
  width: 130px;
  min-height: 40px;
  flex-shrink: 0;
  border-radius: 8px;
  border: 1.5px solid #dcdfe6;
  background: linear-gradient(135deg, #1677ff 0%, #36cfc9 100%);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  overflow: hidden;
  position: relative;
  transition: opacity 0.2s;
  user-select: none;
}
.captcha-img-box:hover {
  opacity: 0.85;
}
.captcha-img {
  width: 100%;
  height: 40px;
  object-fit: cover;
  display: block;
}
.captcha-text-fallback {
  color: #fff;
  font-size: 22px;
  font-weight: 800;
  letter-spacing: 6px;
  font-family: 'Courier New', monospace;
  text-shadow: 1px 1px 3px rgba(0,0,0,0.3);
}
.captcha-refresh-tip {
  position: absolute;
  bottom: 2px;
  font-size: 10px;
  color: rgba(255, 255, 255, 0.75);
  line-height: 1;
}
.captcha-loading {
  color: #fff;
  font-size: 20px;
  animation: spin 1s linear infinite;
}
@keyframes spin {
  from { transform: rotate(0deg); }
  to   { transform: rotate(360deg); }
}
/* ===== 底部链接 ===== */
.login-footer {
  display: flex;
  justify-content: space-between;
  margin-top: 16px;
  font-size: 13px;
}
</style>
