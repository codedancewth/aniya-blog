# ✅ Waline 本地 MySQL 部署完成！

## 🎉 当前状态

**Waline 评论服务已启动并正常运行！**

- ✅ MySQL 服务：运行中 (localhost:3306)
- ✅ Waline 服务：运行中 (http://localhost:8360)
- ✅ 数据库表：已创建 (wl_Comment, wl_Counter, wl_UserInfo)
- ✅ 前端配置：已更新 (指向本地服务)

---

## 📊 服务信息

| 服务 | 地址 | 状态 |
|------|------|------|
| Waline Server | http://localhost:8360 | ✅ 运行中 |
| MySQL | localhost:3306 | ✅ 运行中 |
| 数据库 | waline | ✅ 已连接 |

---

## 🚀 测试评论功能

### 1. 测试 API

```bash
# 测试评论 API
curl "http://localhost:8360/comment?type=latest&path=/test"

# 应该返回：
# {"page":1,"totalPages":0,"pageSize":10,"count":0,"data":[]}
```

### 2. 重启博客开发服务器

```bash
# 停止当前运行的 dev 服务器（Ctrl+C）
cd /root/.openclaw/workspace/aniya-blog
bun run dev
```

### 3. 访问博客页面

打开浏览器访问博客，在任意文章页面应该可以看到评论框。

---

## 📁 文件位置

```
/root/.openclaw/workspace/aniya-blog/
├── waline-server/           # Waline 服务端
│   ├── server.js           # 启动脚本
│   ├── .env                # 环境变量配置
│   ├── package.json        # Node.js 依赖
│   └── data/               # 数据目录
├── src/
│   └── site.config.ts      # 前端配置（已更新）
└── WALINE-README.md        # 本文档
```

---

## 🛠️ 管理命令

### 查看 Waline 服务状态

```bash
ps aux | grep "node server.js" | grep waline
```

### 停止 Waline 服务

```bash
pkill -f "waline-server"
```

### 重启 Waline 服务

```bash
cd /root/.openclaw/workspace/aniya-blog/waline-server
pkill -f "waline-server"
node server.js &
```

### 查看评论内容

```bash
mysql -u waline -pWalineUser@2026 waline -e "SELECT * FROM wl_Comment;"
```

---

## 🔐 安全配置

### 修改管理员 TOKEN

编辑 `waline-server/.env`：

```bash
TOKEN=你的安全令牌_随机字符串
```

然后重启服务。

### 修改数据库密码

```bash
# 修改 MySQL 密码
mysql -u root -e "ALTER USER 'waline'@'localhost' IDENTIFIED WITH mysql_native_password BY '新密码'; FLUSH PRIVILEGES;"

# 更新 waline-server/.env 中的密码
# 重启 Waline 服务
```

---

## 🌐 生产环境部署

### 1. 使用 Nginx 反向代理

```nginx
server {
    listen 443 ssl;
    server_name waline.wth2jhl.online;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://localhost:8360;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 2. 更新前端配置

编辑 `src/site.config.ts`：

```typescript
waline: {
  enable: true,
  server: 'https://waline.wth2jhl.online/',
  ...
}
```

### 3. 使用 PM2 管理进程（推荐）

```bash
# 安装 PM2
npm install -g pm2

# 启动 Waline
cd /root/.openclaw/workspace/aniya-blog/waline-server
pm2 start server.js --name waline-server

# 开机自启
pm2 startup
pm2 save
```

---

## 📝 环境变量说明

`waline-server/.env` 文件配置：

```bash
# 数据库配置
TYPE=mysql
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_DB=waline
MYSQL_USER=waline
MYSQL_PASSWORD=WalineUser@2026

# 站点配置
SITE_NAME=wth2jhl.online
SITE_URL=https://wth2jhl.online

# 安全配置
TOKEN=WalineSecureToken2026ChangeMe

# 评论限制
LIMIT_PER_MINUTE=10
LIMIT_ALL=100

# 其他配置
REGION=CN
LANGUAGE=zh-CN

# 服务器端口
PORT=8360
```

---

## 🐛 故障排查

### Waline 服务未启动

```bash
# 检查进程
ps aux | grep "node server.js"

# 手动启动
cd /root/.openclaw/workspace/aniya-blog/waline-server
node server.js
```

### 数据库连接失败

```bash
# 测试 MySQL 连接
mysql -u waline -pWalineUser@2026 waline

# 检查用户权限
mysql -u root -e "SELECT User, Host FROM mysql.user WHERE User='waline';"
```

### 评论不显示

1. 检查浏览器控制台是否有错误
2. 确认 `server` 地址正确
3. 验证 Waline 服务：`curl http://localhost:8360/comment?type=latest&path=/test`

---

## ✅ 完成清单

- [x] MySQL 服务运行
- [x] Waline 数据库创建
- [x] 表结构创建 (wl_Comment, wl_Counter, wl_UserInfo)
- [x] Waline Server 启动
- [x] 前端配置更新
- [ ] 重启博客开发服务器
- [ ] 测试评论功能
- [ ] 修改默认 TOKEN（生产环境）
- [ ] 配置 HTTPS（生产环境）

---

**更新时间**: 2026-03-26 09:48  
**服务端口**: 8360  
**数据库**: MySQL (waline)
