# Aniya Blog - 数据库初始化指南

## 数据库表结构

本项目支持 SQLite 和 MySQL 两种数据库。

### 表列表

| 表名 | 说明 |
|------|------|
| `users` | 用户表 |
| `posts` | 文章表 |
| `categories` | 分类表 |
| `tags` | 标签表 |
| `post_tags` | 文章标签关联表 |
| `comments` | 评论表 |
| `page_views` | 页面浏览记录表 |
| `links` | 友情链接表 |
| `configs` | 站点配置表 |
| `likes` | 点赞记录表 |

---

## MySQL 初始化

### 1. 创建数据库

```sql
CREATE DATABASE IF NOT EXISTS `aniya_blog` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 2. 执行初始化脚本

**方式一：使用命令行**

```bash
mysql -u root -p aniyablog < database/migrations/001_init.sql
mysql -u root -p aniyablog < database/migrations/002_seed.sql
```

**方式二：使用 MySQL 客户端**

```bash
mysql -u root -p
```

```sql
USE aniyablog;
SOURCE database/migrations/001_init.sql;
SOURCE database/migrations/002_seed.sql;
```

### 3. 配置后端连接

复制 `.env.example` 为 `.env` 并修改配置：

```env
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_DATABASE=aniya_blog
DB_USERNAME=root
DB_PASSWORD=your_password
```

---

## SQLite 初始化

SQLite 无需手动创建表，程序会自动创建数据库文件和表结构。

只需配置：

```env
DB_DRIVER=sqlite
DB_SOURCE=data/aniya.db
```

启动后端服务后会自动执行表迁移。

---

## 默认账户

执行 `002_seed.sql` 后，会创建默认管理员账户：

- **用户名**: admin
- **密码**: admin123
- **角色**: admin

> ⚠️ 首次登录后请立即修改密码

---

## 表结构说明

### users - 用户表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INT | 主键 |
| username | VARCHAR(50) | 用户名（唯一） |
| password | VARCHAR(255) | 密码（bcrypt 加密） |
| email | VARCHAR(100) | 邮箱（唯一） |
| avatar | VARCHAR(255) | 头像 URL |
| nickname | VARCHAR(50) | 昵称 |
| role | VARCHAR(20) | 角色：admin/editor/user |
| status | INT | 状态：1-活跃，0-禁用 |
| last_login_at | DATETIME | 最后登录时间 |
| last_login_ip | VARCHAR(50) | 最后登录 IP |

### posts - 文章表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INT | 主键 |
| title | VARCHAR(200) | 标题 |
| slug | VARCHAR(200) | URL 别名（唯一） |
| description | VARCHAR(500) | 描述 |
| content | TEXT | 内容（Markdown） |
| content_html | TEXT | 内容（HTML） |
| cover_image | VARCHAR(255) | 封面图 |
| author_id | INT | 作者 ID |
| status | INT | 状态：1-发布，0-草稿，2-归档 |
| published_at | DATETIME | 发布时间 |
| view_count | BIGINT | 浏览次数 |
| comment_count | BIGINT | 评论数 |
| like_count | BIGINT | 点赞数 |
| category_id | INT | 分类 ID |
| language | VARCHAR(10) | 语言 |
| is_top | TINYINT | 是否置顶 |
| custom_data | TEXT | 自定义数据（frontmatter） |

### categories - 分类表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INT | 主键 |
| name | VARCHAR(50) | 名称（唯一） |
| slug | VARCHAR(50) | URL 别名（唯一） |
| description | VARCHAR(255) | 描述 |
| parent_id | INT | 父分类 ID（支持多级分类） |
| sort_order | INT | 排序 |
| post_count | BIGINT | 文章数 |

### tags - 标签表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INT | 主键 |
| name | VARCHAR(50) | 名称（唯一） |
| slug | VARCHAR(50) | URL 别名（唯一） |
| description | VARCHAR(255) | 描述 |
| post_count | BIGINT | 文章数 |

### comments - 评论表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INT | 主键 |
| content | TEXT | 内容 |
| post_id | INT | 文章 ID |
| user_id | INT | 用户 ID |
| parent_id | INT | 父评论 ID（支持嵌套回复） |
| author_name | VARCHAR(50) | 作者名 |
| author_email | VARCHAR(100) | 作者邮箱 |
| author_url | VARCHAR(255) | 作者网站 |
| author_ip | VARCHAR(50) | 作者 IP |
| agent | VARCHAR(255) | User-Agent |
| status | INT | 状态：1-已通过，0-待审核，2-垃圾 |
| like_count | BIGINT | 点赞数 |
| is_admin | TINYINT | 是否管理员 |

### 其他表

- `page_views`: 页面浏览记录（用于统计）
- `links`: 友情链接
- `configs`: 站点配置（键值对）
- `likes`: 点赞记录
- `post_tags`: 文章与标签的多对多关联表

---

## 自动迁移

使用 Go 后端时，程序会自动执行表迁移：

```bash
cd backend
go run cmd/server/main.go
```

启动日志中看到 `Database initialized successfully` 即表示迁移成功。
