<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import * as articleApi from '@/api/article'
import type { Article } from '@/api/article'

const router = useRouter()
const articles = ref<Article[]>([])
const total = ref(0)
const page = ref(1)
const size = ref(10)
const keyword = ref('')
const searchKeyword = ref('')
const loading = ref(false)

async function fetchArticles() {
  loading.value = true
  try {
    const result = await articleApi.list({
      page: page.value,
      size: size.value,
      keyword: searchKeyword.value || undefined
    })
    articles.value = result.list || []
    total.value = result.total
  } catch {
    // 错误已在拦截器中处理
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  page.value = 1
  searchKeyword.value = keyword.value.trim()
  fetchArticles()
}

function handlePageChange(newPage: number) {
  page.value = newPage
  fetchArticles()
}

function formatDate(dateStr: string): string {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

const totalPages = () => Math.ceil(total.value / size.value) || 1

onMounted(() => {
  fetchArticles()
})
</script>

<template>
  <div class="page">
    <div class="search-bar">
      <input
        v-model="keyword"
        type="text"
        class="form-input"
        placeholder="搜索文章..."
        @keyup.enter="handleSearch"
      />
      <button class="btn btn-primary" @click="handleSearch" :disabled="loading">搜索</button>
    </div>

    <div v-if="loading" class="empty-state">加载中...</div>

    <template v-else-if="articles.length > 0">
      <div v-for="article in articles" :key="article.ID" class="card">
        <div class="card-title">
          <router-link :to="`/articles/${article.ID}`">{{ article.title }}</router-link>
        </div>
        <div v-if="article.summary" class="card-summary">{{ article.summary }}</div>
        <div class="card-meta">
          <span>{{ article.user?.username || '匿名' }}</span>
          <span>{{ formatDate(article.CreatedAt) }}</span>
          <span>{{ article.views }} 次阅读</span>
        </div>
      </div>

      <div class="pagination">
        <button
          class="btn btn-outline btn-sm"
          :disabled="page <= 1"
          @click="handlePageChange(page - 1)"
        >
          上一页
        </button>
        <span class="pagination-info">{{ page }} / {{ totalPages() }}</span>
        <button
          class="btn btn-outline btn-sm"
          :disabled="page >= totalPages()"
          @click="handlePageChange(page + 1)"
        >
          下一页
        </button>
      </div>
    </template>

    <div v-else class="empty-state">暂无文章</div>
  </div>
</template>
