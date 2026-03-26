// Waline Server - 本地 MySQL 配置
import pkg from '@waline/vercel';

// Waline 配置
const config = {
  // 数据库配置 - 本地 MySQL
  type: 'mysql',
  host: 'localhost',
  port: 3306,
  database: 'waline',
  user: 'waline',
  password: 'WalineUser@2026',
  
  // 站点配置
  siteName: 'wth2jhl.online',
  siteURL: 'https://wth2jhl.online',
  
  // 安全配置
  token: 'WalineSecureToken2026ChangeMe',
  
  // 评论限制
  limitPerMinute: 10,
  limitAll: 100,
  
  // 其他配置
  region: 'CN',
  language: 'zh-CN',
};

// 初始化 Waline
const handler = pkg(config);

// 创建 HTTP 服务器
import http from 'http';

const PORT = process.env.PORT || 8383;

const server = http.createServer(async (req, res) => {
  // 处理 CORS
  res.setHeader('Access-Control-Allow-Origin', '*');
  res.setHeader('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS');
  res.setHeader('Access-Control-Allow-Headers', 'Content-Type, Authorization');
  
  if (req.method === 'OPTIONS') {
    res.writeHead(200);
    res.end();
    return;
  }
  
  // 调用 Waline handler
  try {
    await handler(req, res);
  } catch (error) {
    console.error('Waline error:', error);
    res.writeHead(500);
    res.end('Internal Server Error');
  }
});

server.listen(PORT, () => {
  console.log(`✅ Waline Server 已启动`);
  console.log(`📍 地址：http://localhost:${PORT}`);
  console.log(`💾 数据库：MySQL (localhost:3306, database: waline)`);
  console.log(`\n测试：curl http://localhost:${PORT}/`);
});
