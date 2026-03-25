# Aniya Blog Backend

Aniya Blog 的 Go 语言后端服务，提供完整的博客管理 API。

## 技术栈

- **语言**: Go 1.21+
- **框架**: Gin
- **数据库**: SQLite / MySQL (通过 GORM)
- **认证**: JWT
- **文档**: Swagger

## 快速开始

### 环境要求

- Go 1.21 或更高版本
- Make (可选，用于使用 Makefile 命令)

### 安装依赖

```bash
cd backend
go mod download
go mod tidy
```

### 配置

复制环境变量文件并修改配置：

```bash
cp .env.example .env
```

#### SQLite 配置（默认）

```env
SERVER_PORT=8080
SERVER_MODE=debug

DB_DRIVER=sqlite
DB_SOURCE=data/aniya.db

JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRE_TIME=24h

ALLOW_ORIGIN=http://localhost:4321
```

#### MySQL 配置

```env
SERVER_PORT=8080
SERVER_MODE=debug

DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_DATABASE=aniya_blog
DB_USERNAME=root
DB_PASSWORD=your_password

JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRE_TIME=24h

ALLOW_ORIGIN=http://localhost:4321
```

### 数据库初始化

#### MySQL 数据库

**方式一：使用初始化脚本（推荐）**

Linux/macOS:
```bash
cd backend/database
./init_mysql.sh -u root -p your_password -d aniyablog
```

Windows:
```cmd
cd backend\\database
init_mysql.bat root your_password aniyablog
```

**方式二：手动执行 SQL**

```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS \`aniya_blog\` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 导入表结构
mysql -u root -p aniyablog < database/migrations/001_init.sql

# 导入初始数据
mysql -u root -p aniyablog < database/migrations/002_seed.sql
```

**方式三：自动迁移**

程序启动时会自动执行表迁移，只需正确配置数据库连接即可。

#### SQLite 数据库

SQLite 无需手动初始化，程序启动时会自动创建数据库文件和表结构。

详细数据库说明请参考：[database/README.md](database/README.md)

### 运行

使用 Make：

```bash
make run
```

或直接使用 Go：

```bash
go run cmd/server/main.go
```

### 构建

```bash
make build
# 或
go build -o bin/aniya-blog-server cmd/server/main.go
```

## API 文档

启动服务后访问：

- Swagger UI: http://localhost:8080/swagger/index.html
- API 基础路径：`/api/v1`

## API 接口

### 认证

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | /auth/login | 用户登录 | 否 |
| POST | /auth/register | 用户注册 | 否 |
| GET | /auth/me | 获取当前用户 | 是 |
| POST | /auth/refresh | 刷新令牌 | 是 |
| POST | /auth/change-password | 修改密码 | 是 |

### 文章

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | /posts | 获取文章列表 | 否 |
| GET | /posts/:id | 获取文章详情 | 否 |
| GET | /posts/slug/:slug | 根据 slug 获取文章 | 否 |
| GET | /posts/search?q=keyword | 搜索文章 | 否 |
| GET | /tags/:tagSlug/posts | 根据标签获取文章 | 否 |
| GET | /categories/:categorySlug/posts | 根据分类获取文章 | 否 |
| POST | /posts | 创建文章 | 是 |
| PUT | /posts/:id | 更新文章 | 是 |
| DELETE | /posts/:id | 删除文章 | 是 |

### 标签

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | /tags | 获取标签列表 | 否 |
| GET | /tags/all | 获取所有标签 | 否 |
| GET | /tags/:slug | 获取标签详情 | 否 |
| POST | /tags | 创建标签 | 是 |
| PUT | /tags/:slug | 更新标签 | 是 |
| DELETE | /tags/:slug | 删除标签 | 是 |

### 分类

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | /categories | 获取分类列表 | 否 |
| GET | /categories/tree | 获取分类树 | 否 |
| GET | /categories/:slug | 获取分类详情 | 否 |
| POST | /categories | 创建分类 | 是 |
| PUT | /categories/:slug | 更新分类 | 是 |
| DELETE | /categories/:slug | 删除分类 | 是 |

### 评论

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | /posts/:post_id/comments | 获取文章评论 | 否 |
| POST | /comments | 创建评论 | 否 |
| PUT | /comments/:id | 更新评论 | 是 |
| DELETE | /comments/:id | 删除评论 | 是 |
| POST | /comments/:id/like | 点赞评论 | 否 |

### 页面浏览

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | /pageviews | 记录页面浏览 | 否 |
| GET | /pageviews/count?path=/xxx | 获取页面浏览次数 | 否 |
| GET | /pageviews/counts?paths=/a,/b | 批量获取浏览次数 | 否 |
| GET | /pageviews/stats | 获取站点统计 | 否 |

## 项目结构

```
backend/
├── cmd/
│   └── server/           # 主程序入口
├── internal/
│   ├── config/           # 配置管理
│   ├── database/         # 数据库连接
│   ├── models/           # 数据模型
│   ├── handlers/         # HTTP 处理器
│   ├── middleware/       # 中间件
│   ├── repository/       # 数据访问层
│   ├── services/         # 业务逻辑层
│   └── utils/            # 工具函数
├── pkg/
│   ├── auth/             # 认证模块
│   ├── response/         # 统一响应
│   └── validator/        # 参数验证
├── api/
│   └── docs/             # API 文档
├── go.mod
├── go.sum
└── Makefile
```

## Make 命令

| 命令 | 描述 |
|------|------|
| `make deps` | 下载依赖 |
| `make build` | 构建二进制文件 |
| `make run` | 运行服务器 |
| `make test` | 运行测试 |
| `make clean` | 清理构建文件 |
| `make swagger` | 生成 Swagger 文档 |
| `make docker` | 构建 Docker 镜像 |
| `make fmt` | 格式化代码 |
| `make lint` | 运行代码检查 |

## Docker 部署

### 构建镜像

```bash
make docker
# 或
docker build -t aniya-blog-server:latest .
```

### 运行容器

```bash
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/logs:/app/logs \
  -e JWT_SECRET=your-secret \
  --name aniya-blog \
  aniya-blog-server:latest
```

## 前端集成

前端项目位于根目录，使用 Astro 框架。

### 配置 API 地址

在 `src/lib/api.ts` 中配置：

```typescript
export const API_CONFIG = {
  baseURL: import.meta.env.PROD ? '/api' : 'http://localhost:8080/api/v1',
  timeout: 10000,
}
```

### 使用示例

```typescript
import { postApi, authApi } from '@/lib/api'

// 获取文章列表
const { data } = await postApi.list(1, 10)

// 登录
const { data } = await authApi.login('username', 'password')
localStorage.setItem('token', data.token)
```

## 许可证

Apache 2.0
