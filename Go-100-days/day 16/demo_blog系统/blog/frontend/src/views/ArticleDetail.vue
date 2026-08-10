<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import * as articleApi from '@/api/article'
import type { Article } from '@/api/article'

const route = useRoute()
const router = useRouter()

const article = ref<Article | null>(null)
const loading = ref(true)
const currentUserId = ref<number | null>(null)

// 获取当前登录用户ID
const savedUser = localStorage.getItem('user')
if (savedUser) {
  try {
    currentUserId.value = JSON.parse(savedUser).id
  } catch { /* ignore */ }
}

const isAuthor = () => currentUserId.value !== null && currentUserId.value === article.value?.user_id

function formatDate(dateStr: string): string {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

async function handleDelete() {
  if (!article.value) return
  if (!confirm('确定要删除这篇文章吗？')) return
  try {
    await articleApi.remove(article.value.ID)
    router.push('/')
  } catch {
    // 错误已在拦截器中处理
  }
}

onMounted(async () => {
  try {
    const id = Number(route.params.id)
    article.value = await articleApi.detail(id)
  } catch {
    // 错误已在拦截器中处理
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="page">
    <router-link to="/" class="back-link">&larr; 返回首页</router-link>

    <div v-if="loading" class="empty-state">加载中...</div>

    <div v-else-if="article" class="article-detail">
      <h1>{{ article.title }}</h1>
      <div class="article-meta">
        <span>{{ article.user?.username || '匿名' }}</span>
        <span>{{ formatDate(article.CreatedAt) }}</span>
        <span>{{ article.views }} 次阅读</span>
      </div>
      <div class="article-content">{{ article.content }}</div>

      <div v-if="isAuthor()" class="article-actions">
        <router-link :to="`/articles/${article.ID}/edit`">
          <button class="btn btn-outline btn-sm">编辑</button>
        </router-link>
        <button class="btn btn-danger btn-sm" @click="handleDelete">删除</button>
      </div>
    </div>

    <div v-else class="empty-state">文章不存在</div>
  </div>
</template>
