# 🎉 评论系统配置完成 - Go 后端 + MySQL

## ✅ 配置总结

已完全移除 Waline 服务，改用 **Go 后端 API + MySQL** 存储评论数据。

---

## 📊 系统架构

```
┌─────────────────┐
│   博客前端      │
│  (Astro Pure)   │
└────────┬────────┘
         │ HTTP REST API
         │ http://localhost:8081/api/v1
         ▼
┌─────────────────┐
│  Go 后端服务    │
│   (Gin + GORM)  │
└────────┬────────┘
         │ MySQL 协议
         │ localhost:3306
         ▼
┌─────────────────┐
│   MySQL 数据库  │
│  (aniya_blog)   │
└─────────────────┘
```

---

## 🚀 快速启动

### 1. 确保 MySQL 服务运行

```bash
# 检查 MySQL 状态
systemctl status mysql

# 如果未启动
sudo systemctl start mysql
```

### 2. 创建数据库

```bash
# 登录 MySQL
mysql -u root -p

# 创建数据库
CREATE DATABASE IF NOT EXISTS `aniya_blog` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
exit;
```

### 3. 启动 Go 后端

```bash
cd /root/.openclaw/workspace/aniya-blog/backend
go run cmd/server/main.go
```

或者后台运行：

```bash
cd /root/.openclaw/workspace/aniya-blog/backend
nohup go run cmd/server/main.go > server.log 2>&1 &
```

### 4. 重启前端开发服务器

```bash
cd /root/.openclaw/workspace/aniya-blog
bun run dev
```

---

## 📝 API 端点

### 评论相关 API

| 方法 | 端点 | 说明 |
|------|------|------|
| `GET` | `/api/v1/posts/:post_id/comments` | 获取文章评论列表 |
| `POST` | `/api/v1/comments` | 创建评论 |
| `GET` | `/api/v1/comments/:id` | 获取评论详情 |
| `PUT` | `/api/v1/comments/:id` | 更新评论（需认证） |
| `DELETE` | `/api/v1/comments/:id` | 删除评论（需认证） |
| `POST` | `/api/v1/comments/:id/like` | 点赞评论 |

### 测试 API

```bash
# 获取文章评论（假设文章 ID 为 1）
curl http://localhost:8081/api/v1/posts/1/comments

# 创建评论
curl -X POST http://localhost:8081/api/v1/comments \
  -H "Content-Type: application/json" \
  -d '{
    "post_id": 1,
    "content": "测试评论",
    "author_name": "测试用户",
    "author_email": "test@example.com"
  }'
```

---

## 📁 文件变更

### 新增文件

```
src/components/comment/
└── CommentSection.astro    # 新的评论组件
```

### 修改文件

```
src/layouts/BlogPost.astro      # 使用新评论组件
src/site.config.ts              # 移除 Waline，添加 comment 配置
backend/.env                    # 启用 MySQL 配置
```

### 可删除的文件（Waline 相关）

```
waline-server/                  # 可删除
waline-docker-compose.yml       # 可删除
WALINE-*.md                     # 可删除
src/components/waline/          # 可删除或保留备用
```

---

## 🔧 配置说明

### 前端配置 (`src/site.config.ts`)

```typescript
comment: {
  enable: true,
  apiURL: 'http://localhost:8081/api/v1',  // Go 后端 API 地址
  // 生产环境：
  // apiURL: 'https://api.wth2jhl.online/api/v1',
}
```

### 后端配置 (`backend/.env`)

```bash
# 服务器配置
SERVER_PORT=8081
SERVER_MODE=debug

# MySQL 配置
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_DATABASE=aniya_blog
DB_USERNAME=root
DB_PASSWORD=AniyaBlog@2026
```

---

## 🗄️ 数据库表结构

评论数据存储在 `comments` 表中：

```sql
CREATE TABLE `comments` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `content` text NOT NULL,
  `post_id` int(11) unsigned NOT NULL,
  `user_id` int(11) unsigned DEFAULT NULL,
  `parent_id` int(11) unsigned DEFAULT NULL,
  `author_name` varchar(50) DEFAULT NULL,
  `author_email` varchar(100) DEFAULT NULL,
  `author_url` varchar(255) DEFAULT NULL,
  `author_ip` varchar(50) DEFAULT NULL,
  `agent` varchar(255) DEFAULT NULL,
  `status` int(11) DEFAULT 1,
  `like_count` bigint DEFAULT 0,
  `is_admin` tinyint(1) DEFAULT 0,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `post_id` (`post_id`),
  KEY `parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

---

## 🛠️ 常用命令

### 查看评论内容

```bash
mysql -u root -pAniyaBlog@2026 aniya_blog -e "SELECT * FROM comments ORDER BY created_at DESC LIMIT 10;"
```

### 删除所有评论（清空测试数据）

```bash
mysql -u root -pAniyaBlog@2026 aniya_blog -e "DELETE FROM comments;"
```

### 查看后端日志

```bash
tail -f /root/.openclaw/workspace/aniya-blog/backend/server.log
```

### 重启后端服务

```bash
pkill -f "go run cmd/server"
cd /root/.openclaw/workspace/aniya-blog/backend
nohup go run cmd/server/main.go > server.log 2>&1 &
```

---

## 🌐 生产环境部署

### 1. 配置反向代理（Nginx）

```nginx
# API 反向代理
server {
    listen 443 ssl;
    server_name api.wth2jhl.online;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://localhost:8081;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 2. 更新前端配置

```typescript
comment: {
  enable: true,
  apiURL: 'https://api.wth2jhl.online/api/v1',
}
```

### 3. 使用 systemd 管理后端服务

创建服务文件 `/etc/systemd/system/aniya-blog-backend.service`:

```ini
[Unit]
Description=Aniya Blog Backend Service
After=network.target mysql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/root/.openclaw/workspace/aniya-blog/backend
ExecStart=/usr/local/go/bin/go run cmd/server/main.go
Restart=on-failure
Environment=GIN_MODE=release

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable aniya-blog-backend
sudo systemctl start aniya-blog-backend
sudo systemctl status aniya-blog-backend
```

---

## 🔐 安全建议

1. **修改默认密码** - 更改 MySQL root 密码
2. **启用 HTTPS** - 生产环境必须使用 HTTPS
3. **配置防火墙** - 只开放必要端口（80, 443）
4. **JWT 密钥** - 修改 `JWT_SECRET` 为强随机字符串
5. **评论审核** - 可修改代码实现评论审核功能

---

## 🐛 故障排查

### 后端无法启动

```bash
# 检查 Go 依赖
cd /root/.openclaw/workspace/aniya-blog/backend
go mod download
go mod tidy

# 检查端口占用
lsof -i :8081

# 查看日志
tail -f server.log
```

### 数据库连接失败

```bash
# 测试 MySQL 连接
mysql -u root -pAniyaBlog@2026 -e "SELECT 1"

# 检查数据库是否存在
mysql -u root -pAniyaBlog@2026 -e "SHOW DATABASES LIKE 'aniya_blog';"
```

### 评论不显示

1. 检查浏览器控制台是否有错误
2. 验证 API 是否可访问：`curl http://localhost:8081/api/v1/posts/1/comments`
3. 确认 `apiURL` 配置正确

---

## ✅ 完成清单

- [x] 移除 Waline 服务端配置
- [x] 创建新的评论组件（调用 Go API）
- [x] 更新前端配置
- [x] 配置后端 MySQL 连接
- [x] 更新博客布局使用新组件
- [ ] 启动 Go 后端服务
- [ ] 重启前端开发服务器
- [ ] 测试评论功能
- [ ] 清理 Waline 相关文件

---

**更新时间**: 2026-03-26  
**架构**: Go Backend + MySQL  
**前端**: Astro Pure + 自定义评论组件
