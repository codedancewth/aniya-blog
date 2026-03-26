# Waline 本地 MySQL 配置 - 更新完成

## ✅ 已完成的配置

### 1. Docker Compose 配置文件
- 📄 `waline-docker-compose.yml` - Waline 服务 + MySQL 数据库编排配置
- 支持环境变量自定义
- 包含健康检查和数据持久化

### 2. 前端配置更新
- 📄 `src/site.config.ts` - Waline 客户端配置
- 已切换到本地地址：`http://localhost:8383/`
- 本地化文案：中文提示

### 3. 部署脚本
- 📄 `scripts/waline-deploy.sh` - 一键部署脚本
- 支持：start | stop | restart | logs | status | backup
- 已添加执行权限

### 4. 文档
- 📄 `WALINE-DEPLOY.md` - 详细部署指南
- 📄 `waline/README.md` - 快速开始指南
- 📄 `.env.waline.example` - 环境变量模板

---

## 🚀 下一步操作

### 立即启动（开发环境）

```bash
# 1. 进入项目目录
cd /root/.openclaw/workspace/aniya-blog

# 2. 启动 Waline 服务
./scripts/waline-deploy.sh start

# 或手动启动
docker compose -f waline-docker-compose.yml up -d

# 3. 验证服务
curl http://localhost:8383

# 4. 重启博客开发服务器
# 停止当前运行的 dev 服务器（Ctrl+C）
bun run dev
```

### 自定义配置（可选）

```bash
# 1. 复制环境变量模板
cp .env.waline.example .env.waline

# 2. 编辑配置（修改密码、TOKEN 等）
vim .env.waline

# 3. 重启服务
./scripts/waline-deploy.sh restart
```

---

## 📊 架构说明

```
┌─────────────────┐
│   博客前端      │
│  (Astro Pure)   │
└────────┬────────┘
         │ HTTP
         │ http://localhost:8383
         ▼
┌─────────────────┐
│  Waline 服务端  │
│   (Docker)      │
└────────┬────────┘
         │ MySQL 协议
         │ 内部网络
         ▼
┌─────────────────┐
│   MySQL 数据库  │
│  (Docker)       │
│  端口：3307     │
└─────────────────┘
```

---

## 🔐 安全提醒

### 开发环境
- ✅ 当前配置适合本地开发
- ⚠️ 默认密码仅用于测试

### 生产环境部署前必须：
1. ❗ 修改 `TOKEN`（管理员令牌）
2. ❗ 修改 MySQL 密码
3. ❗ 配置 HTTPS
4. ❗ 配置防火墙
5. ❗ 配置邮件通知（可选）

---

## 📁 新增文件清单

```
aniya-blog/
├── waline-docker-compose.yml        # [新增] Docker Compose 配置
├── .env.waline.example              # [新增] 环境变量模板
├── WALINE-DEPLOY.md                 # [新增] 详细部署文档
├── src/site.config.ts               # [已修改] Waline 配置
├── scripts/
│   └── waline-deploy.sh             # [新增] 部署脚本
└── waline/
    └── README.md                    # [新增] 快速开始指南
```

---

## 💡 使用提示

### 查看实时日志
```bash
./scripts/waline-deploy.sh logs
```

### 备份评论数据
```bash
./scripts/waline-deploy.sh backup
```

### 连接 MySQL 数据库
```bash
mysql -h 127.0.0.1 -P 3307 -u waline -p
# 密码：WalineUser@2026
```

### 查看评论内容
```sql
USE waline;
SELECT * FROM comment ORDER BY created_at DESC LIMIT 10;
```

---

## 🎉 配置完成！

现在你的 Waline 评论系统已经配置为使用本地 MySQL 存储。

启动服务后，博客的评论数据将保存在本地数据库中，完全由你控制！

---

**更新时间**: 2026-03-26  
**配置版本**: v1.0
