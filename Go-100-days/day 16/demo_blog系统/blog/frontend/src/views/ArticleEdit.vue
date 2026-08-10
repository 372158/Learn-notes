<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import * as articleApi from '@/api/article'

const route = useRoute()
const router = useRouter()

const title = ref('')
const summary = ref('')
const content = ref('')
const error = ref('')
const loading = ref(false)
const pageLoading = ref(true)

async function fetchArticle() {
  try {
    const id = Number(route.params.id)
    const article = await articleApi.detail(id)
    title.value = article.title
    summary.value = article.summary || ''
    content.value = article.content

    // 检查是否为作者
    const savedUser = localStorage.getItem('user')
    if (savedUser) {
      const user = JSON.parse(savedUser)
      if (user.id !== article.user_id) {
        alert('无权编辑此文章')
        router.push(`/articles/${id}`)
        return
      }
    }
  } catch {
    // 错误已在拦截器中处理
  } finally {
    pageLoading.value = false
  }
}

async function handleUpdate() {
  error.value = ''
  if (!title.value.trim()) {
    error.value = '请输入文章标题'
    return
  }
  if (!content.value.trim()) {
    error.value = '请输入文章内容'
    return
  }

  loading.value = true
  try {
    const id = Number(route.params.id)
    await articleApi.update(id, {
      title: title.value.trim(),
      content: content.value.trim(),
      summary: summary.value.trim() || undefined
    })
    router.push(`/articles/${id}`)
  } catch {
    // 错误已在拦截器中处理
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchArticle()
})
</script>

<template>
  <div class="page">
    <div v-if="pageLoading" class="empty-state">加载中...</div>

    <template v-else>
      <router-link to="/" class="back-link">&larr; 返回首页</router-link>
      <h1 class="page-title" style="margin-bottom: 24px">编辑文章</h1>

      <form @submit.prevent="handleUpdate">
        <div class="form-group">
          <label class="form-label">标题 *</label>
          <input v-model="title" type="text" class="form-input" placeholder="请输入文章标题" />
        </div>
        <div class="form-group">
          <label class="form-label">摘要</label>
          <input v-model="summary" type="text" class="form-input" placeholder="可选，简短的摘要" />
        </div>
        <div class="form-group">
          <label class="form-label">内容 *</label>
          <textarea v-model="content" class="form-input" placeholder="请输入文章内容"></textarea>
        </div>
        <p v-if="error" class="form-error">{{ error }}</p>
        <button type="submit" class="btn btn-primary" :disabled="loading">
          {{ loading ? '保存中...' : '保存修改' }}
        </button>
      </form>
    </template>
  </div>
</template>
