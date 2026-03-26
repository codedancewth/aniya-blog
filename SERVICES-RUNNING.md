# 🎉 网站已启动 - 访问指南

## ✅ 服务状态

| 服务 | 状态 | 地址 | 说明 |
|------|------|------|------|
| **前端** | ✅ 运行中 | http://localhost:4321 | Astro 博客网站 |
| **后端 API** | ✅ 运行中 | http://localhost:8081 | Go 后端服务 |
| **MySQL** | ✅ 运行中 | localhost:3306 | 数据库 |

---

## 🌐 访问网站

### 本地访问

打开浏览器访问：

```
http://localhost:4321
```

### 局域网访问（从其他设备）

```
http://10.1.0.15:4321
```

---

## 📝 测试评论功能

### 步骤 1：访问博客文章

1. 打开 http://localhost:4321/blog
2. 点击任意文章

### 步骤 2：发表评论

1. 滚动到文章底部
2. 填写评论表单：
   - 昵称（必填）
   - 邮箱（必填）
   - 网站（可选）
   - 评论内容
3. 点击"提交评论"

### 步骤 3：查看评论

提交成功后，评论会立即显示在页面上。

---

## 🔧 测试 API

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

### 点赞评论

```bash
curl -X POST http://localhost:8081/api/v1/comments/1/like
```

---

## 📊 服务管理

### 查看服务状态

```bash
# 前端
ps aux | grep "astro dev" | grep -v grep

# 后端
ps aux | grep "go run" | grep -v grep

# MySQL
systemctl status mysql
```

### 停止服务

```bash
# 停止前端
pkill -f "astro dev"

# 停止后端
pkill -f "go run cmd/server"
```

### 重启服务

```bash
# 重启前端
cd /root/.openclaw/workspace/aniya-blog
npm run dev

# 重启后端
cd /root/.openclaw/workspace/aniya-blog/backend
go run cmd/server/main.go
```

---

## 🗄️ 查看数据库

### 登录 MySQL

```bash
mysql -u root -pAniyaBlog@2026
```

### 查看评论数据

```sql
USE aniya_blog;
SELECT id, author_name, author_email, content, created_at 
FROM comments 
ORDER BY created_at DESC 
LIMIT 10;
```

### 查看文章列表

```sql
SELECT id, title, slug, created_at FROM posts;
```

---

## 🌍 生产环境部署

### 1. 构建前端

```bash
cd /root/.openclaw/workspace/aniya-blog
npm run build
```

### 2. 配置 Nginx

```nginx
# 前端
server {
    listen 80;
    server_name wth2jhl.online;
    
    root /root/.openclaw/workspace/aniya-blog/dist;
    index index.html;
    
    location / {
        try_files $uri $uri/ /404.html;
    }
    
    # 静态资源缓存
    location ~* \.(jpg|jpeg|png|gif|ico|css|js)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}

# 后端 API
server {
    listen 80;
    server_name api.wth2jhl.online;
    
    location / {
        proxy_pass http://localhost:8081;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 3. 更新前端配置

编辑 `src/site.config.ts`：

```typescript
comment: {
  enable: true,
  apiURL: 'https://api.wth2jhl.online/api/v1',
}
```

### 4. 使用 systemd 管理后端

创建 `/etc/systemd/system/aniya-blog-backend.service`:

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
```

---

## 🐛 故障排查

### 前端无法访问

```bash
# 检查端口
lsof -i :4321

# 查看日志
ps aux | grep "astro dev"
```

### 后端 API 错误

```bash
# 检查端口
lsof -i :8081

# 测试健康检查
curl http://localhost:8081/health

# 查看日志
ps aux | grep "go run"
```

### 评论提交失败

1. 打开浏览器开发者工具（F12）
2. 查看 Console 是否有错误
3. 检查 Network 标签中的 API 请求
4. 验证后端 API 是否可访问

---

## 📋 快速检查清单

- [x] MySQL 服务运行
- [x] 数据库 aniya_blog 已创建
- [x] Go 后端运行（端口 8081）
- [x] 前端运行（端口 4321）
- [x] 评论组件已配置
- [x] API 可访问

---

## 🎯 现在可以做什么

1. **访问网站**: http://localhost:4321
2. **浏览博客**: http://localhost:4321/blog
3. **发表评论**: 打开任意文章，滚动到底部
4. **管理评论**: 通过 MySQL 数据库或后端 API

---

**更新时间**: 2026-03-26 10:10  
**前端**: Astro v5.18.1  
**后端**: Go + Gin  
**数据库**: MySQL 8.0
