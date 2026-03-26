# ✅ 评论系统配置完成！

## 🎉 配置总结

已完全移除 Waline 服务，改用 **Go 后端 API + MySQL** 存储评论数据。

---

## ✅ 当前状态

| 服务 | 状态 | 地址 |
|------|------|------|
| MySQL | ✅ 运行中 | localhost:3306 |
| Go 后端 | ✅ 运行中 | http://localhost:8081 |
| 数据库 | ✅ 已创建 | aniya_blog |
| 评论 API | ✅ 可用 | /api/v1/comments |

---

## 📊 系统架构

```
┌─────────────────┐
│   博客前端      │
│  (Astro Pure)   │
│                 │
│  评论组件       │
│  (CommentSection)│
└────────┬────────┘
         │ HTTP REST API
         │ POST /api/v1/comments
         │ GET /api/v1/posts/:id/comments
         ▼
┌─────────────────┐
│  Go 后端服务    │
│   (Gin + GORM)  │
│                 │
│  CommentHandler │
│  Comment Model  │
└────────┬────────┘
         │ MySQL 协议
         ▼
┌─────────────────┐
│   MySQL 数据库  │
│  (aniya_blog)   │
│                 │
│  comments 表    │
└─────────────────┘
```

---

## 🚀 下一步操作

### 1. 重启前端开发服务器

```bash
cd /root/.openclaw/workspace/aniya-blog
# 停止当前运行的 dev 服务器（Ctrl+C）
bun run dev
```

### 2. 测试评论功能

1. 打开博客文章页面
2. 滚动到页面底部
3. 填写评论表单（昵称、邮箱、内容）
4. 提交评论
5. 查看评论是否显示

---

## 📝 API 测试

### 创建测试评论

```bash
curl -X POST http://localhost:8081/api/v1/comments \
  -H "Content-Type: application/json" \
  -d '{
    "post_id": 1,
    "content": "这是一条测试评论",
    "author_name": "测试用户",
    "author_email": "test@example.com"
  }'
```

### 获取评论列表

```bash
curl http://localhost:8081/api/v1/posts/1/comments
```

---

## 📁 完成的修改

### 新增文件
- ✅ `src/components/comment/CommentSection.astro` - 新评论组件

### 修改文件
- ✅ `src/layouts/BlogPost.astro` - 使用新评论组件
- ✅ `src/site.config.ts` - 移除 Waline，添加 comment 配置
- ✅ `backend/.env` - 启用 MySQL 配置

### 可清理的文件
```bash
# Waline 相关文件（可选删除）
rm -rf waline-server/
rm -f waline-docker-compose.yml
rm -f WALINE-*.md
rm -rf src/components/waline/  # 或者保留备用
```

---

## 🔧 配置说明

### 前端配置

```typescript
// src/site.config.ts
comment: {
  enable: true,
  apiURL: 'http://localhost:8081/api/v1',
}
```

### 后端配置

```bash
# backend/.env
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_DATABASE=aniya_blog
DB_USERNAME=root
DB_PASSWORD=AniyaBlog@2026
```

---

## 🗄️ 数据库表

评论数据存储在 `comments` 表：

```sql
mysql -u root -pAniyaBlog@2026 aniya_blog -e "DESCRIBE comments;"
```

字段说明：
- `id` - 评论 ID
- `content` - 评论内容
- `post_id` - 文章 ID
- `author_name` - 作者昵称
- `author_email` - 作者邮箱
- `author_ip` - 作者 IP
- `status` - 状态（1=通过，0=待审核，2=垃圾）
- `like_count` - 点赞数
- `created_at` - 创建时间

---

## 🛠️ 常用命令

### 查看后端进程

```bash
ps aux | grep "go run" | grep -v grep
```

### 重启后端

```bash
pkill -f "go run cmd/server"
cd /root/.openclaw/workspace/aniya-blog/backend
nohup go run cmd/server/main.go > server.log 2>&1 &
```

### 查看评论内容

```bash
mysql -u root -pAniyaBlog@2026 aniya_blog -e "SELECT id, author_name, content, created_at FROM comments ORDER BY created_at DESC LIMIT 10;"
```

### 清空测试数据

```bash
mysql -u root -pAniyaBlog@2026 aniya_blog -e "DELETE FROM comments;"
```

---

## 🌐 生产环境部署

### 1. 修改 API 地址

```typescript
// src/site.config.ts
comment: {
  enable: true,
  apiURL: 'https://api.wth2jhl.online/api/v1',
}
```

### 2. 配置 Nginx 反向代理

```nginx
server {
    listen 443 ssl;
    server_name api.wth2jhl.online;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location /api/ {
        proxy_pass http://localhost:8081;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 3. 使用 systemd 管理后端

```bash
sudo systemctl enable aniya-blog-backend
sudo systemctl start aniya-blog-backend
```

---

## 🐛 故障排查

### 评论提交失败

1. 检查浏览器控制台错误
2. 验证后端 API：`curl http://localhost:8081/api/v1/posts/1/comments`
3. 确认 CORS 配置正确

### 后端无法启动

```bash
# 检查端口占用
lsof -i :8081

# 查看日志
tail -f backend/server.log
```

### 数据库连接失败

```bash
# 测试 MySQL 连接
mysql -u root -pAniyaBlog@2026 -e "SELECT 1"
```

---

## ✅ 完成清单

- [x] 创建 Go 后端评论 API
- [x] 配置 MySQL 数据库连接
- [x] 创建前端评论组件
- [x] 更新博客布局
- [x] 移除 Waline 配置
- [x] 启动 Go 后端服务
- [ ] 重启前端开发服务器
- [ ] 测试评论功能
- [ ] 清理 Waline 文件

---

**更新时间**: 2026-03-26 10:00  
**架构**: Go Backend (Gin) + MySQL  
**前端**: Astro Pure + 自定义评论组件  
**API**: RESTful (/api/v1/comments)
