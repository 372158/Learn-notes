import client from './client'

export interface Article {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  DeletedAt: string | null
  title: string
  content: string
  summary: string
  views: number
  user_id: number
  user: {
    ID: number
    username: string
    email: string
    role: string
  } | null
}

export interface ArticleListParams {
  page?: number
  size?: number
  keyword?: string
}

export interface ArticleListResult {
  list: Article[]
  total: number
  page: number
  size: number
}

export async function list(params: ArticleListParams): Promise<ArticleListResult> {
  const res = await client.get('/api/v1/articles', { params })
  return res.data as ArticleListResult
}

export async function detail(id: number): Promise<Article> {
  const res: any = await client.get(`/api/v1/articles/${id}`)
  return res.msg as Article
}

export interface ArticleCreateParams {
  title: string
  content: string
  summary?: string
}

export async function create(data: ArticleCreateParams): Promise<Article> {
  const res = await client.post('/api/v1/articles', data)
  return res.data as Article
}

export async function update(id: number, data: ArticleCreateParams): Promise<Article> {
  const res = await client.put(`/api/v1/articles/${id}`, data)
  return res.data as Article
}

export async function remove(id: number): Promise<void> {
  await client.delete(`/api/v1/articles/${id}`)
}
