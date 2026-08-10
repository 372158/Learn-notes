import client from './client'

export interface RegisterParams {
  username: string
  password: string
  email: string
}

export interface LoginParams {
  username: string
  password: string
}

export interface UserInfo {
  id: number
  username: string
  email: string
  role: string
}

export interface LoginResult {
  token: string
  user: UserInfo
}

export async function register(data: RegisterParams): Promise<LoginResult> {
  const res = await client.post('/api/v1/users/register', data)
  return { token: '', user: null as any }
}

export async function login(data: LoginParams): Promise<LoginResult> {
  const res = await client.post('/api/v1/users/login', data)
  return res.data as LoginResult
}
