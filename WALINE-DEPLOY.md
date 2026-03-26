# Waline 评论系统部署指南

## 📦 概述

本文档说明如何在本地部署 Waline 评论系统，使用 MySQL 存储评论数据。

## 🚀 快速开始

### 1. 启动 Waline 服务

在项目根目录执行：

```bash
# 启动 Waline 和 MySQL 容器
docker compose -f waline-docker-compose.yml up -d

# 查看日志
docker compose -f waline-docker-compose.yml logs -f waline

# 停止服务
docker compose -f waline-docker-compose.yml down
```

### 2. 验证服务

访问 `http://localhost:8383` 应该能看到 Waline 服务正常运行。

### 3. 更新前端配置

前端配置已在 `src/site.config.ts` 中更新：

```typescript
waline: {
  enable: true,
  server: 'http://localhost:8383/',  // 本地开发
  // server: 'https://waline.wth2jhl.online/',  // 生产环境
  ...
}
```

### 4. 重启博客开发服务器

```bash
# 如果正在运行，先停止
# 然后重新启动
bun run dev
```

## 🔧 配置说明

### 环境变量

在 `waline-docker-compose.yml` 中可以修改以下配置：

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `SITE_NAME` | 站点名称 | wth2jhl.online |
| `SITE_URL` | 站点 URL | https://wth2jhl.online |
| `MYSQL_HOST` | MySQL 主机 | waline-db |
| `MYSQL_PORT` | MySQL 端口 | 3306 |
| `MYSQL_DB` | 数据库名 | waline |
| `MYSQL_USER` | 数据库用户 | waline |
| `MYSQL_PASSWORD` | 数据库密码 | WalineUser@2026 |
| `TOKEN` | 管理员令牌 | WalineSecureToken2026ChangeMe |
| `LIMIT_PER_MINUTE` | 每分钟评论限制 | 10 |
| `LIMIT_ALL` | 总评论限制 | 100 |

### 管理员配置

如果需要邮件通知功能，取消注释并配置：

```yaml
environment:
  MAIL_USER: your-email@example.com
  SMTP_SERVICE: Gmail  # 或 QQ, 163, 等
  SMTP_USER: your-email@example.com
  SMTP_PASS: your-smtp-password
```

## 📊 数据库管理

### 连接本地 MySQL

```bash
# 使用 MySQL 客户端连接
mysql -h 127.0.0.1 -P 3307 -u waline -p

# 密码：WalineUser@2026
```

### 查看数据

```sql
-- 查看所有评论
USE waline;
SELECT * FROM comment;

-- 查看统计数据
SELECT COUNT(*) as total_comments FROM comment;
```

### 备份数据

```bash
# 导出数据库
docker exec waline-mysql mysqldump -u waline -pWalineUser@2026 waline > waline-backup.sql

# 恢复数据库
docker exec -i waline-mysql mysql -u waline -pWalineUser@2026 waline < waline-backup.sql
```

## 🌐 生产环境部署

### 1. 修改服务器地址

在 `src/site.config.ts` 中：

```typescript
server: 'https://waline.wth2jhl.online/',  // 你的域名
```

### 2. 使用反向代理（Nginx 示例）

```nginx
server {
    listen 443 ssl;
    server_name waline.wth2jhl.online;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://localhost:8383;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 3. 安全建议

- ⚠️ **修改默认密码**：更改 MySQL 和管理员 TOKEN
- 🔒 **启用 HTTPS**：生产环境必须使用 HTTPS
- 🛡️ **配置防火墙**：只开放必要端口
- 📧 **配置邮件通知**：便于接收评论提醒

## 🐛 故障排查

### Waline 无法启动

```bash
# 查看容器日志
docker logs waline-server

# 检查 MySQL 是否就绪
docker logs waline-mysql
```

### 评论不显示

1. 检查浏览器控制台是否有错误
2. 确认 `server` 地址正确
3. 验证 Waline 服务是否可访问：`curl http://localhost:8383`

### 数据库连接失败

```bash
# 重启容器
docker compose -f waline-docker-compose.yml restart

# 检查网络
docker network inspect waline-network
```

## 📝 数据迁移

### 从远程 Waline 迁移

如果之前使用远程 Waline 服务，可以导出后导入：

1. 联系原服务管理员导出数据
2. 使用 Waline 官方迁移工具
3. 导入到本地 MySQL

## 🔗 相关文档

- [Waline 官方文档](https://waline.js.org/)
- [Waline Docker 部署](https://waline.js.org/guide/get-started/docker.html)
- [Waline 服务端配置](https://waline.js.org/reference/server/env.html)

---

**最后更新**: 2026-03-26
