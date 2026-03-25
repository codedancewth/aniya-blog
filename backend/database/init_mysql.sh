#!/bin/bash

# Aniya Blog MySQL 初始化脚本
# 使用方法：./init_mysql.sh -u root -p -d aniyablog

# 默认值
MYSQL_USER="root"
MYSQL_PASS=""
MYSQL_DB="aniya_blog"

# 解析参数
while getopts "u:p:d:h" opt; do
  case $opt in
    u) MYSQL_USER="$OPTARG" ;;
    p) MYSQL_PASS="$OPTARG" ;;
    d) MYSQL_DB="$OPTARG" ;;
    h)
      echo "用法：$0 [-u 用户名] [-p 密码] [-d 数据库名]"
      echo "  -u MySQL 用户名 (默认：root)"
      echo "  -p MySQL 密码 (默认：空)"
      echo "  -d 数据库名 (默认：aniya_blog)"
      exit 0
      ;;
    \?) echo "无效选项：-$OPTARG" >&2; exit 1 ;;
  esac
done

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MIGRATIONS_DIR="$SCRIPT_DIR/migrations"

echo "========================================"
echo "Aniya Blog MySQL 初始化"
echo "========================================"
echo "MySQL 用户：$MYSQL_USER"
echo "数据库名：$MYSQL_DB"
echo "迁移目录：$MIGRATIONS_DIR"
echo "========================================"

# 创建数据库
echo "[1/3] 创建数据库..."
mysql -u"$MYSQL_USER" ${MYSQL_PASS:+-p"$MYSQL_PASS"} -e "CREATE DATABASE IF NOT EXISTS \`$MYSQL_DB\` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

if [ $? -ne 0 ]; then
  echo "错误：创建数据库失败"
  exit 1
fi
echo "数据库创建成功"

# 执行表结构迁移
echo "[2/3] 执行表结构迁移..."
mysql -u"$MYSQL_USER" ${MYSQL_PASS:+-p"$MYSQL_PASS"} "$MYSQL_DB" < "$MIGRATIONS_DIR/001_init.sql"

if [ $? -ne 0 ]; then
  echo "错误：执行表结构迁移失败"
  exit 1
fi
echo "表结构迁移成功"

# 执行初始数据插入
echo "[3/3] 执行初始数据插入..."
mysql -u"$MYSQL_USER" ${MYSQL_PASS:+-p"$MYSQL_PASS"} "$MYSQL_DB" < "$MIGRATIONS_DIR/002_seed.sql"

if [ $? -ne 0 ]; then
  echo "警告：执行初始数据插入失败（可忽略）"
fi

echo "========================================"
echo "初始化完成!"
echo "========================================"
echo ""
echo "默认管理员账户:"
echo "  用户名：admin"
echo "  密码：admin123"
echo ""
echo "请在 .env 文件中配置:"
echo "  DB_DRIVER=mysql"
echo "  DB_HOST=localhost"
echo "  DB_PORT=3306"
echo "  DB_DATABASE=$MYSQL_DB"
echo "  DB_USERNAME=$MYSQL_USER"
echo "  DB_PASSWORD=<your_password>"
echo ""
