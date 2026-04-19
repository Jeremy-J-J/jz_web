<template>
  <div class="login-page">
    <div class="login-box">
      <h1>管理员登录</h1>
      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <input v-model="username" type="text" placeholder="用户名" required>
        </div>
        <div class="form-group">
          <input v-model="password" type="password" placeholder="密码" required>
        </div>
        <div v-if="errorMsg" class="error-msg">{{ errorMsg }}</div>
        <button type="submit" :disabled="loading">{{ loading ? '登录中...' : '登录' }}</button>
      </form>
      <router-link to="/" class="back-home">&larr; 返回首页</router-link>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { authApi } from '../api'

const router = useRouter()
const username = ref('')
const password = ref('')
const loading = ref(false)
const errorMsg = ref('')

const handleLogin = async () => {
  loading.value = true
  errorMsg.value = ''
  try {
    const res = await authApi.login(username.value, password.value)
    localStorage.setItem('token', res.token)
    localStorage.setItem('admin', JSON.stringify(res.admin))
    router.push('/admin')
  } catch (e) {
    errorMsg.value = e.response?.data?.error || '登录失败'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page { min-height: 100vh; display: flex; align-items: center; justify-content: center; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
.login-box { background: #fff; padding: 48px; border-radius: 8px; box-shadow: 0 8px 32px rgba(0,0,0,0.2); width: 100%; max-width: 400px; }
.login-box h1 { text-align: center; font-size: 24px; margin-bottom: 32px; color: #333; }
.form-group { margin-bottom: 20px; }
.form-group input { width: 100%; padding: 14px 16px; border: 1px solid #ddd; border-radius: 4px; font-size: 16px; }
.form-group input:focus { outline: none; border-color: #667eea; }
.error-msg { color: #ff4d4f; font-size: 14px; margin-bottom: 16px; text-align: center; }
button { width: 100%; padding: 14px; background: #667eea; color: #fff; border: none; border-radius: 4px; font-size: 16px; cursor: pointer; }
button:hover:not(:disabled) { background: #5a6fd6; }
button:disabled { opacity: 0.6; cursor: not-allowed; }
.back-home { display: block; text-align: center; margin-top: 20px; color: #999; text-decoration: none; font-size: 14px; }
.back-home:hover { color: #667eea; }
</style>