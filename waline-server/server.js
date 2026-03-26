const path = require('node:path');

// 加载环境变量
require('dotenv').config({
  path: path.join(__dirname, '.env'),
});

const Application = require('thinkjs');

const instance = new Application({
  ROOT_PATH: __dirname,
  APP_PATH: path.join(__dirname, 'node_modules/@waline/vercel/src'),
  VIEW_PATH: path.join(__dirname, 'node_modules/@waline/vercel/view'),
  RUNTIME_PATH: path.join(__dirname, 'runtime'),
  proxy: true,
  env: 'production',
});

instance.run();

console.log(`✅ Waline Server 已启动`);
console.log(`📍 地址：http://localhost:${process.env.PORT || 8383}`);
console.log(`💾 数据库：MySQL (localhost:3306, database: waline)`);
