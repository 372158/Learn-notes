import { ref, computed } from 'vue'
import { login as loginApi, register as registerApi } from '@/api/auth'
import type { UserInfo } from '@/api/auth'

const token = ref<string | null>(localStorage.getItem('token'))
const user = ref<UserInfo | null>(null)

// 初始化时从 localStorage 恢复用户信息
const savedUser = localStorage.getItem('user')
if (savedUser) {
  try {
    user.value = JSON.parse(savedUser)
  } catch {
    localStorage.removeItem('user')
  }
}

export function useAuth() {
  const isLoggedIn = computed(() => !!token.value)

  async function login(username: string, password: string) {
    const result = await loginApi({ username, password })
    token.value = result.token
    user.value = result.user
    localStorage.setItem('token', result.token)
    localStorage.setItem('user', JSON.stringify(result.user))
  }

  async function register(username: string, password: string, email: string) {
    await registerApi({ username, password, email })
    // 注册成功后自动登录
    await login(username, password)
  }

  function logout() {
    token.value = null
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  return {
    token,
    user,
    isLoggedIn,
    login,
    register,
    logout
  }
}
