#!/bin/bash

# Waline 评论系统快速部署脚本
# 使用方法：./scripts/waline-deploy.sh [start|stop|restart|logs|status]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
COMPOSE_FILE="$PROJECT_ROOT/waline-docker-compose.yml"
ENV_FILE="$PROJECT_ROOT/.env.waline"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查 Docker 是否安装
check_docker() {
    if ! command -v docker &> /dev/null; then
        error "Docker 未安装，请先安装 Docker"
        exit 1
    fi
    
    if ! command -v docker compose &> /dev/null; then
        error "Docker Compose 未安装，请先安装 Docker Compose"
        exit 1
    fi
    
    info "Docker 和 Docker Compose 已安装"
}

# 检查环境变量文件
check_env() {
    if [ ! -f "$ENV_FILE" ]; then
        warn "环境变量文件不存在，将使用默认配置"
        info "如需自定义配置，请复制 .env.waline.example 为 .env.waline 并修改"
    else
        info "使用环境变量文件：$ENV_FILE"
    fi
}

# 启动服务
start() {
    info "启动 Waline 评论系统..."
    
    if [ -f "$ENV_FILE" ]; then
        docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" up -d
    else
        docker compose -f "$COMPOSE_FILE" up -d
    fi
    
    info "等待服务启动..."
    sleep 5
    
    status
}

# 停止服务
stop() {
    info "停止 Waline 评论系统..."
    docker compose -f "$COMPOSE_FILE" down
    info "服务已停止"
}

# 重启服务
restart() {
    stop
    sleep 2
    start
}

# 查看日志
logs() {
    docker compose -f "$COMPOSE_FILE" logs -f waline
}

# 查看状态
status() {
    info "Waline 服务状态:"
    docker compose -f "$COMPOSE_FILE" ps
    
    echo ""
    info "服务访问地址:"
    echo "  - Waline 服务端：http://localhost:8383"
    echo "  - MySQL 数据库：localhost:3307"
    echo ""
    info "数据库连接信息:"
    echo "  - 主机：127.0.0.1"
    echo "  - 端口：3307"
    echo "  - 数据库：waline"
    echo "  - 用户：waline"
    echo "  - 密码：WalineUser@2026 (请在 .env.waline 中修改)"
}

# 备份数据库
backup() {
    BACKUP_DIR="$PROJECT_ROOT/data/backups"
    mkdir -p "$BACKUP_DIR"
    
    TIMESTAMP=$(date +%Y%m%d_%H%M%S)
    BACKUP_FILE="$BACKUP_DIR/waline_backup_$TIMESTAMP.sql"
    
    info "备份数据库到：$BACKUP_FILE"
    docker exec waline-mysql mysqldump -u waline -p"${MYSQL_PASSWORD:-WalineUser@2026}" waline > "$BACKUP_FILE"
    
    info "备份完成"
}

# 显示帮助
help() {
    echo "Waline 评论系统管理脚本"
    echo ""
    echo "使用方法：$0 [命令]"
    echo ""
    echo "命令:"
    echo "  start     启动服务"
    echo "  stop      停止服务"
    echo "  restart   重启服务"
    echo "  logs      查看日志"
    echo "  status    查看状态"
    echo "  backup    备份数据库"
    echo "  help      显示帮助"
    echo ""
    echo "示例:"
    echo "  $0 start     # 启动 Waline 服务"
    echo "  $0 logs      # 查看实时日志"
    echo "  $0 backup    # 备份数据库"
}

# 主函数
main() {
    cd "$PROJECT_ROOT"
    
    check_docker
    check_env
    
    case "${1:-start}" in
        start)
            start
            ;;
        stop)
            stop
            ;;
        restart)
            restart
            ;;
        logs)
            logs
            ;;
        status)
            status
            ;;
        backup)
            backup
            ;;
        help|--help|-h)
            help
            ;;
        *)
            error "未知命令：$1"
            help
            exit 1
            ;;
    esac
}

main "$@"
