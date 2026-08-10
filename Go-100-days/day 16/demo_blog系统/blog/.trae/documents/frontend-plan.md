# Blog 前端界面实现计划

## 摘要

基于已有的 Go + Gin 后端系统，使用 **Vue 3 + Vite + TypeScript** 构建一个风格简洁的博客前端界面，包含文章浏览、用户认证、文章管理三大模块。

---

## 一、后端 API 总览（已存在）

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/v1/users/register` | 注册 | 否 |
| POST | `/api/v1/users/login` | 登录，返回 JWT | 否 |
| GET | `/api/v1/articles` | 文章列表（分页+搜索） | 否 |
| GET | `/api/v1/articles/:id` | 文章详情（自动+阅读量） | 否 |
| POST | `/api/v1/articles` | 创建文章 | 是 |
| PUT | `/api/v1/articles/:id` | 更新文章（仅作者） | 是 |
| DELETE | `/api/v1/articles/:id` | 删除文章（仅作者） | 是 |

统一响应格式：
```json
{ "code": 0, "msg": "...", "data": {...} }
```
`code === 0` 表示成功，非 0 表示错误。

---

## 二、前端项目结构

```
blog/
├── frontend/                    # 前端项目根目录（与后端并列）
│   ├── index.html
│   ├── package.json
│   ├── tsconfig.json
│   ├── vite.config.ts
│   └── src/
│       ├── main.ts              # 入口，挂载 App + Router
│       ├── App.vue              # 根组件（Navbar + RouterView）
│       ├── router/
│       │   └── index.ts         # 路由定义 + 导航守卫
│       ├── api/
│       │   ├── client.ts        # Axios 实例（baseURL、拦截器）
│       │   ├── auth.ts          # 注册 / 登录 API
│       │   └── article.ts       # 文章 CRUD API
│       ├── composables/
│       │   └── useAuth.ts       # 认证状态管理（token、user、login/logout）
│       ├── views/
│       │   ├── HomePage.vue     # 首页：文章列表 + 搜索 + 分页
│       │   ├── ArticleDetail.vue # 文章详情页
│       │   ├── LoginPage.vue    # 登录页
│       │   ├── RegisterPage.vue # 注册页
│       │   ├── ArticleCreate.vue # 创建文章页
│       │   └── ArticleEdit.vue  # 编辑文章页
│       └── styles/
│           └── global.css       # 全局样式（简洁风格）
```

---

## 三、技术选型

| 类别 | 选型 | 说明 |
|------|------|------|
| 构建工具 | Vite 5 | 快速开发体验 |
| 框架 | Vue 3 (Composition API) | `<script setup lang="ts">` 写法 |
| 语言 | TypeScript | 类型安全 |
| 路由 | Vue Router 4 | SPA 路由 |
| HTTP | Axios | 请求拦截、Token 自动携带 |
| 状态管理 | Composable (useAuth) | 项目规模小，无需 Pinia |
| UI 库 | 无，纯手写 CSS | 保持简洁风格，无需重型组件库 |
| 样式 | CSS Variables | 统一配色、间距 |

---

## 四、路由设计

| 路径 | 页面 | 组件 | 认证 |
|------|------|------|------|
| `/` | 文章列表（首页） | `HomePage.vue` | 否 |
| `/articles/:id` | 文章详情 | `ArticleDetail.vue` | 否 |
| `/login` | 登录 | `LoginPage.vue` | 否 |
| `/register` | 注册 | `RegisterPage.vue` | 否 |
| `/articles/create` | 创建文章 | `ArticleCreate.vue` | 是 |
| `/articles/:id/edit` | 编辑文章 | `ArticleEdit.vue` | 是 |

---

## 五、详细实现

### 5.1 项目初始化

**文件：`frontend/`（整个目录新建）**

使用 Vite 创建 Vue 3 + TypeScript 项目：
```bash
npm create vite@latest frontend -- --template vue-ts
cd frontend
npm install
npm install axios vue-router@4
```

### 5.2 API 层

#### `src/api/client.ts` — Axios 实例

- 创建 Axios 实例，`baseURL` 默认 `http://localhost:8080`
- **请求拦截器**：从 `localStorage` 读取 token，自动添加到 `Authorization: Bearer <token>` 请求头
- **响应拦截器**：统一处理 `code !== 0` 的错误，弹窗提示 `msg`；401 时清除 token 并跳转登录页

#### `src/api/auth.ts` — 认证 API

- `register(data: {username, password, email})` → `POST /api/v1/users/register`
- `login(data: {username, password})` → `POST /api/v1/users/login`

#### `src/api/article.ts` — 文章 API

- `list(params: {page, size, keyword?})` → `GET /api/v1/articles`
- `detail(id: number)` → `GET /api/v1/articles/:id`
- `create(data: {title, content, summary?})` → `POST /api/v1/articles`
- `update(id: number, data)` → `PUT /api/v1/articles/:id`
- `remove(id: number)` → `DELETE /api/v1/articles/:id`

### 5.3 认证状态管理

#### `src/composables/useAuth.ts`

使用 Vue 3 `ref` + `computed` 管理全局认证状态：

- **状态**：`token`（`ref<string | null>`，初始化从 `localStorage` 读取）、`user`（`ref`，登录后存储）
- **计算属性**：`isLoggedIn`（`computed`，基于 `token` 是否有值）
- **方法**：
  - `login(username, password)` — 调 API → 存 token 到 `localStorage` + `ref`
  - `register(username, password, email)` — 调 API → 自动登录
  - `logout()` — 清除 `localStorage` 中的 token + user，重置状态
  - `fetchUser()` — 从 token 解析用户信息（可选）

通过 `provide/inject` 或直接导出响应式单例供全局使用。

### 5.4 路由

#### `src/router/index.ts`

- 定义 6 条路由（见上文路由设计表）
- **导航守卫 `beforeEach`**：
  - 访问 `/articles/create` 和 `/articles/:id/edit` 时检查是否登录
  - 未登录 → 跳转 `/login` 并携带 `redirect` 参数
  - 已登录访问 `/login` 或 `/register` → 重定向到 `/`

### 5.5 页面组件

#### `App.vue` — 根布局

- **Navbar**：左侧 Blog 名称（点击跳首页），右侧根据登录状态显示：
  - 未登录：`登录` / `注册` 按钮
  - 已登录：用户名 + `写文章` / `退出` 按钮
- **RouterView**：页面内容区域
- Container 最大宽度 960px，居中

#### `HomePage.vue` — 首页（文章列表）

- 顶部搜索栏：输入框 + 搜索按钮（支持关键词模糊搜索）
- 文章卡片列表：每张卡片显示标题、摘要、作者、发布时间、阅读量
- 点击卡片跳转到 `/articles/:id`
- 底部分页组件（上一页/下一页/页码）
- 空状态：暂无文章时显示提示

#### `ArticleDetail.vue` — 文章详情

- 显示标题、作者、发布时间、阅读量
- 正文内容（HTML/markdown 渲染，简单起见按纯文本或简单 HTML 展示）
- 返回按钮 → 回到首页
- 如果是作者本人，显示 `编辑` 和 `删除` 按钮（需确认）

#### `LoginPage.vue` — 登录

- 表单：用户名 + 密码
- 登录按钮 + 注册链接
- 登录成功 → 跳转回 `redirect` 参数指定的页面或首页
- 错误提示（用户名或密码错误）

#### `RegisterPage.vue` — 注册

- 表单：用户名 + 邮箱 + 密码 + 确认密码
- 注册按钮 + 登录链接
- 前端校验：用户名 ≥3 位、邮箱格式、密码 ≥6 位、两次密码一致
- 注册成功 → 自动登录 → 跳转首页

#### `ArticleCreate.vue` — 创建文章

- 表单：标题（必填）、摘要（选填）、正文（必填，textarea）
- 提交按钮
- 创建成功 → 跳转到文章详情页

#### `ArticleEdit.vue` — 编辑文章

- 进入时通过 API 加载已有文章数据填充表单
- 表单同创建页
- 提交更新 → 跳转到文章详情页
- 检查是否为文章作者，非作者提示无权编辑

### 5.6 样式方案

#### `src/styles/global.css`

- 使用 CSS Variables 定义主题色：
  - `--primary: #333`（主文字色）
  - `--bg: #fff`（背景色）
  - `--accent: #409eff`（链接/按钮色）
  - `--border: #e5e5e5`（边框色）
- 简洁风格：白底 + 适量留白 + 细边框 + 无多余装饰
- 响应式：最大宽度 960px，移动端适配

### 5.7 Vite 配置

#### `vite.config.ts`

- 配置代理（开发环境）：将 `/api` 开头的请求代理到 `http://localhost:8080`，避免跨域问题

```ts
export default defineConfig({
  server: {
    proxy: {
      '/api': 'http://localhost:8080'
    }
  }
})
```

---

## 六、数据流示意

```
用户操作 → Vue 组件 → api 层（Axios）→ 后端 Go API
                               ↓
                         响应拦截器处理
                               ↓
                     code===0 → 返回 data
                     code!==0 → alert(msg)
                     401      → 跳转登录
```

---

## 七、实施步骤

1. **初始化项目**：Vite 创建 Vue 3 + TS 项目，安装依赖（axios、vue-router）
2. **创建 API 层**：`client.ts`（Axios 实例 + 拦截器）→ `auth.ts` → `article.ts`
3. **创建认证模块**：`useAuth.ts` composable
4. **创建路由**：`router/index.ts`（路由表 + 导航守卫）
5. **创建全局样式**：`styles/global.css`
6. **创建 App.vue**：根布局（Navbar + RouterView）
7. **创建页面组件**（按依赖顺序）：
   - `LoginPage.vue` + `RegisterPage.vue`（认证页）
   - `HomePage.vue`（文章列表）
   - `ArticleDetail.vue`（文章详情）
   - `ArticleCreate.vue` + `ArticleEdit.vue`（文章管理）
8. **配置 Vite 代理**到后端
9. **整体测试**：启动后端 + 前端，验证完整流程

---

## 八、不做的内容

- 不添加评论功能（后端暂无此 API）
- 不添加 Markdown 编辑器（使用 textarea 即可，保持简洁）
- 不添加用户个人信息编辑页（后端暂无此 API）
- 不添加管理后台（无角色管理，仅作者可管理自己的文章）
- 不添加图片上传（后端暂无此 API）
- 不使用 UI 组件库（保持轻量简洁）

---

## 九、验证方式

1. 启动后端：`go run main.go`（确保 MySQL 已启动）
2. 启动前端：`cd frontend && npm run dev`
3. 验证流程：
   - 访问首页 → 看到空列表或已有文章
   - 注册新用户 → 自动登录
   - 创建文章 → 展示在首页
   - 查看文章详情 → 阅读量增加
   - 编辑自己的文章
   - 删除自己的文章
   - 退出登录 → 无法访问创建/编辑页
   - 搜索文章 → 关键词过滤
   - 分页浏览
