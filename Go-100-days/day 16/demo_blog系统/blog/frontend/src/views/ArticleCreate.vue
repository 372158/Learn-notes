<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import * as articleApi from '@/api/article'

const router = useRouter()

const title = ref('')
const summary = ref('')
const content = ref('')
const error = ref('')
const loading = ref(false)

async function handleCreate() {
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
    const article = await articleApi.create({
      title: title.value.trim(),
      content: content.value.trim(),
      summary: summary.value.trim() || undefined
    })
    router.push(`/articles/${article.ID}`)
  } catch {
    // 错误已在拦截器中处理
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="page">
    <router-link to="/" class="back-link">&larr; 返回首页</router-link>
    <h1 class="page-title" style="margin-bottom: 24px">写文章</h1>

    <form @submit.prevent="handleCreate">
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
        {{ loading ? '发布中...' : '发布文章' }}
      </button>
    </form>
  </div>
</template>
