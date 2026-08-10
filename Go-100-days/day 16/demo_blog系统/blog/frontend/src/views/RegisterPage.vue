<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'

const router = useRouter()
const { register } = useAuth()

const username = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const error = ref('')
const loading = ref(false)

async function handleRegister() {
  error.value = ''

  if (!username.value.trim()) {
    error.value = '请输入用户名'
    return
  }
  if (username.value.trim().length < 3) {
    error.value = '用户名至少需要3个字符'
    return
  }
  if (!email.value.trim()) {
    error.value = '请输入邮箱'
    return
  }
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(email.value.trim())) {
    error.value = '请输入有效的邮箱地址'
    return
  }
  if (!password.value) {
    error.value = '请输入密码'
    return
  }
  if (password.value.length < 6) {
    error.value = '密码至少需要6个字符'
    return
  }
  if (password.value !== confirmPassword.value) {
    error.value = '两次输入的密码不一致'
    return
  }

  loading.value = true
  try {
    await register(username.value.trim(), password.value, email.value.trim())
    router.push('/')
  } catch {
    // 错误已在拦截器中处理
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="auth-form">
    <h2>注册</h2>
    <form @submit.prevent="handleRegister">
      <div class="form-group">
        <label class="form-label">用户名</label>
        <input
          v-model="username"
          type="text"
          class="form-input"
          placeholder="至少3个字符"
          autocomplete="username"
        />
      </div>
      <div class="form-group">
        <label class="form-label">邮箱</label>
        <input
          v-model="email"
          type="email"
          class="form-input"
          placeholder="请输入邮箱"
          autocomplete="email"
        />
      </div>
      <div class="form-group">
        <label class="form-label">密码</label>
        <input
          v-model="password"
          type="password"
          class="form-input"
          placeholder="至少6个字符"
          autocomplete="new-password"
        />
      </div>
      <div class="form-group">
        <label class="form-label">确认密码</label>
        <input
          v-model="confirmPassword"
          type="password"
          class="form-input"
          placeholder="再次输入密码"
          autocomplete="new-password"
        />
      </div>
      <p v-if="error" class="form-error">{{ error }}</p>
      <button type="submit" class="btn btn-primary" :disabled="loading">
        {{ loading ? '注册中...' : '注册' }}
      </button>
    </form>
    <p class="auth-switch">
      已有账号？<router-link to="/login">立即登录</router-link>
    </p>
  </div>
</template>
