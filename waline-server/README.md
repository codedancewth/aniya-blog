# Waline Server - 本地 MySQL 部署

## 🚀 快速启动

### 1. 安装依赖

```bash
cd /root/.openclaw/workspace/aniya-blog/waline-server
npm install
```

### 2. 启动服务

```bash
# 开发模式
npm run dev

# 生产模式
npm start

# 后台运行（推荐）
nohup npm start > waline.log 2>&1 &
```

### 3. 验证服务

```bash
curl http://localhost:8383/ui
```

应该返回 Waline 管理界面。

---

## 📊 数据库信息

- **主机**: localhost
- **端口**: 3306
- **数据库**: waline
- **用户**: waline
- **密码**: WalineUser@2026

查看数据：
```bash
mysql -u waline -pWalineUser@2026 waline -e "SELECT * FROM comment;"
```

---

## 🔧 配置说明

编辑 `index.js` 修改配置：

```javascript
const config = {
  mysql: {
    host: 'localhost',
    port: 3306,
    database: 'waline',
    user: 'waline',
    password: '你的密码',
  },
  
  token: '你的管理员令牌', // 重要！生产环境必须修改
  
  siteName: 'wth2jhl.online',
  siteURL: 'https://wth2jhl.online',
};
```

---

## 🛠️ 常用命令

```bash
# 查看服务状态
ps aux | grep "waline-server"

# 停止服务
pkill -f "waline-server"

# 查看日志
tail -f waline.log

# 重启服务
pkill -f "waline-server" && nohup npm start > waline.log 2>&1 &
```

---

## 📝 前端配置

确保 `src/site.config.ts` 中的配置：

```typescript
waline: {
  enable: true,
  server: 'http://localhost:8383/',
  ...
}
```

---

## 🔐 安全提醒

1. ⚠️ **修改默认 TOKEN** - 在 `index.js` 中更改
2. ⚠️ **修改数据库密码** - 使用强密码
3. 🔒 **生产环境启用 HTTPS**
4. 🛡️ **配置防火墙** - 只开放必要端口

---

**更新时间**: 2026-03-26
