<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'

const router = useRouter()
const { user, isLoggedIn, logout } = useAuth()

function handleLogout() {
  logout()
  router.push('/')
}
</script>

<template>
  <nav class="navbar">
    <div class="navbar-inner">
      <router-link to="/" class="navbar-brand">Blog</router-link>
      <div class="navbar-links">
        <template v-if="isLoggedIn">
          <span class="navbar-user">{{ user?.username }}</span>
          <router-link to="/articles/create">
            <button class="btn btn-primary btn-sm">写文章</button>
          </router-link>
          <button class="btn btn-outline btn-sm" @click="handleLogout">退出</button>
        </template>
        <template v-else>
          <router-link to="/login">
            <button class="btn btn-outline btn-sm">登录</button>
          </router-link>
          <router-link to="/register">
            <button class="btn btn-primary btn-sm">注册</button>
          </router-link>
        </template>
      </div>
    </div>
  </nav>
  <router-view />
</template>
