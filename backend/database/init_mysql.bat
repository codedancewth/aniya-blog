@echo off
REM Aniya Blog MySQL 初始化脚本 (Windows 版本)
REM 使用方法：init_mysql.bat [用户名] [密码] [数据库名]

SETLOCAL

REM 默认值
SET MYSQL_USER=%~1
IF "%MYSQL_USER%"=="" SET MYSQL_USER=root

SET MYSQL_PASS=%~2
SET MYSQL_DB=%~3
IF "%MYSQL_DB%"=="" SET MYSQL_DB=aniya_blog

REM 获取脚本所在目录
SET SCRIPT_DIR=%~dp0
SET MIGRATIONS_DIR=%SCRIPT_DIR%migrations

echo ========================================
echo Aniya Blog MySQL 初始化
echo ========================================
echo MySQL 用户：%MYSQL_USER%
echo 数据库名：%MYSQL_DB%
echo ========================================

REM 创建数据库
echo [1/3] 创建数据库...
mysql -u%MYSQL_USER% %MYSQL_PASS:~0,1% %MYSQL_PASS:~1% -e "CREATE DATABASE IF NOT EXISTS `%MYSQL_DB%` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
IF ERRORLEVEL 1 (
    echo 错误：创建数据库失败
    EXIT /B 1
)
echo 数据库创建成功

REM 执行表结构迁移
echo [2/3] 执行表结构迁移...
mysql -u%MYSQL_USER% %MYSQL_PASS:~0,1% %MYSQL_PASS:~1% %MYSQL_DB% < "%MIGRATIONS_DIR%001_init.sql"
IF ERRORLEVEL 1 (
    echo 错误：执行表结构迁移失败
    EXIT /B 1
)
echo 表结构迁移成功

REM 执行初始数据插入
echo [3/3] 执行初始数据插入...
mysql -u%MYSQL_USER% %MYSQL_PASS:~0,1% %MYSQL_PASS:~1% %MYSQL_DB% < "%MIGRATIONS_DIR%002_seed.sql"
IF ERRORLEVEL 1 (
    echo 警告：执行初始数据插入失败（可忽略）
)

echo ========================================
echo 初始化完成!
echo ========================================
echo.
echo 默认管理员账户:
echo   用户名：admin
echo   密码：admin123
echo.
echo 请在.env 文件中配置:
echo   DB_DRIVER=mysql
echo   DB_HOST=localhost
echo   DB_PORT=3306
echo   DB_DATABASE=%MYSQL_DB%
echo   DB_USERNAME=%MYSQL_USER%
echo   DB_PASSWORD=^<your_password^>
echo.

ENDLOCAL
