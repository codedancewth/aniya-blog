# Waline 评论系统 - 快速开始

## 🎯 一句话说明

Waline 是一个轻量、安全且支持多种数据库的评论系统，本配置使用 **本地 MySQL** 存储评论数据。

## ⚡ 3 分钟快速部署

### 步骤 1: 启动服务

```bash
# 在项目根目录执行
./scripts/waline-deploy.sh start
```

或者使用 Docker Compose：

```bash
docker compose -f waline-docker-compose.yml up -d
```

### 步骤 2: 验证服务

```bash
# 查看服务状态
./scripts/waline-deploy.sh status

# 访问测试
curl http://localhost:8383
```

### 步骤 3: 重启博客

```bash
# 停止当前开发服务器（Ctrl+C）
# 重新启动
bun run dev
```

打开博客页面，评论功能应该已经可以使用了！

## 📁 文件说明

```
aniya-blog/
├── waline-docker-compose.yml    # Docker Compose 配置
├── .env.waline.example          # 环境变量模板
├── .env.waline                  # 自定义环境变量（需创建）
├── WALINE-DEPLOY.md             # 详细部署文档
├── scripts/
│   └── waline-deploy.sh         # 快速部署脚本
└── data/
    └── waline-mysql/            # MySQL 数据持久化目录
```

## 🔐 默认配置

| 项目 | 值 | 说明 |
|------|-----|------|
| Waline 服务地址 | http://localhost:8383 | 前端评论组件连接地址 |
| MySQL 端口 | 3307 | 外部访问端口 |
| 数据库名 | waline | 评论数据存储库 |
| 数据库用户 | waline | 数据库用户名 |
| 数据库密码 | WalineUser@2026 | 数据库密码（请修改！） |

## 🛠️ 常用命令

```bash
# 启动服务
./scripts/waline-deploy.sh start

# 停止服务
./scripts/waline-deploy.sh stop

# 重启服务
./scripts/waline-deploy.sh restart

# 查看日志
./scripts/waline-deploy.sh logs

# 查看状态
./scripts/waline-deploy.sh status

# 备份数据库
./scripts/waline-deploy.sh backup
```

## 📝 自定义配置

1. 复制环境变量模板：

```bash
cp .env.waline.example .env.waline
```

2. 编辑 `.env.waline`，修改你需要的配置：

```bash
# 修改管理员令牌（重要！）
TOKEN=YourSecureToken2026

# 修改数据库密码
MYSQL_PASSWORD=YourStrongPassword

# 配置邮件通知（可选）
MAIL_USER=your-email@example.com
SMTP_SERVICE=Gmail
SMTP_USER=your-email@example.com
SMTP_PASS=your-smtp-password
```

3. 重启服务：

```bash
./scripts/waline-deploy.sh restart
```

## 🌐 生产环境部署

### 1. 修改前端配置

编辑 `src/site.config.ts`：

```typescript
waline: {
  enable: true,
  server: 'https://waline.wth2jhl.online/',  // 你的域名
  ...
}
```

### 2. 配置反向代理

使用 Nginx 或 Caddy 将域名指向 Waline 服务（端口 8383）

### 3. 启用 HTTPS

生产环境必须使用 HTTPS，可以通过 Let's Encrypt 获取免费证书

### 4. 修改默认密码

⚠️ **重要**：生产环境必须修改所有默认密码！

## 🐛 遇到问题？

查看 [WALINE-DEPLOY.md](./WALINE-DEPLOY.md) 获取详细的故障排查指南

## 📚 更多资源

- [Waline 官方文档](https://waline.js.org/)
- [Waline 服务端配置参考](https://waline.js.org/reference/server/env.html)
- [Docker Compose 文档](https://docs.docker.com/compose/)

---

**提示**: 评论数据存储在 `data/waline-mysql/` 目录，请定期备份！
