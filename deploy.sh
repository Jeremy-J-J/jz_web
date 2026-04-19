#!/bin/bash

# 电子资源展示平台 - 一键部署脚本

set -e

PROJECT_DIR="/workspace/pro/jz_web"
BACKEND_DIR="$PROJECT_DIR/backend"
FRONTEND_DIR="$PROJECT_DIR/frontend"
DATA_DIR="$PROJECT_DIR/data"
RUNTIME_DIR="$PROJECT_DIR/runtime"

echo "=========================================="
echo "  JZ Web 电子资源展示平台 - 一键部署"
echo "=========================================="

# 创建必要目录
echo "[1/8] 创建目录结构..."
mkdir -p "$DATA_DIR/db"
mkdir -p "$DATA_DIR/uploads"
mkdir -p "$RUNTIME_DIR/logs"

# 安装Go（如果未安装）
echo "[2/8] 检查并安装Go环境..."
if ! command -v go &> /dev/null; then
    echo "未检测到Go，正在安装..."
    cd /tmp

    # 使用国内镜像源
    echo "正在下载Go（使用国内镜像）..."
    GO_URL="https://golang.google.cn/dl/go1.21.10.linux-amd64.tar.gz"

    for i in 1 2 3; do
        echo "尝试下载... ($i/3)"
        if curl -fsSL --connect-timeout 30 --max-time 300 "$GO_URL" -o go.tar.gz 2>/dev/null; then
            echo "下载完成，正在解压..."
            sudo rm -rf /usr/local/go
            sudo tar -C /usr/local -xzf go.tar.gz
            rm go.tar.gz
            echo "Go安装完成"
            break
        else
            echo "下载失败，重试中..."
            sleep 5
        fi
    done
else
    echo "Go已安装: $(go version)"
fi

# 确保Go在PATH中
export PATH=$PATH:/usr/local/go/bin

# 配置Go模块代理（国内）
echo "配置Go模块代理..."
export GOPROXY=https://goproxy.cn,direct
export GO111MODULE=on

# 判断是否需要编译后端
if [ ! -f "$BACKEND_DIR/app" ]; then
    echo "[3/8] 编译后端..."
    cd "$BACKEND_DIR"
    go mod tidy
    go mod download
    CGO_ENABLED=0 go build -o app .
    echo "后端编译完成"
else
    echo "[3/8] 后端已编译，跳过"
fi

# 构建前端
echo "[4/8] 构建前端..."
cd "$FRONTEND_DIR"
npm install --silent 2>/dev/null
npm run build

# 配置Nginx
echo "[5/8] 配置Nginx..."
NGINX_CONF="/etc/nginx/sites-available/jz_web.conf"
sudo cp "$PROJECT_DIR/conf/nginx/jz_web.conf" "$NGINX_CONF"
sudo ln -sf "$NGINX_CONF" /etc/nginx/sites-enabled/jz_web
sudo rm -f /etc/nginx/sites-enabled/default 2>/dev/null || true
sudo nginx -t && echo "Nginx配置测试通过"

# 配置systemd
echo "[6/8] 配置systemd服务..."
SERVICE_FILE="/etc/systemd/system/jz-web.service"
sudo cp "$PROJECT_DIR/conf/systemd/jz-web.service" "$SERVICE_FILE"
sudo systemctl daemon-reload

# 启动服务
echo "[7/8] 启动服务..."
sudo systemctl enable jz-web 2>/dev/null || true
sudo systemctl restart jz-web
sudo systemctl reload nginx

# 防火墙
echo "[8/8] 检查防火墙..."
if command -v ufw &> /dev/null; then
    sudo ufw allow 8888/tcp 2>/dev/null || true
    echo "防火墙规则已添加"
fi

# 获取服务器IP
SERVER_IP=$(hostname -I | awk '{print $1}')

echo ""
echo "=========================================="
echo "  部署完成！"
echo "=========================================="
echo ""
echo "访问地址:"
echo "  前台首页: http://$SERVER_IP:8888"
echo "  管理后台: http://$SERVER_IP:8888/admin/login"
echo ""
echo "默认管理员账号: admin / admin123"
echo ""
echo "服务管理命令:"
echo "  查看后端状态: systemctl status jz-web"
echo "  重启后端:     systemctl restart jz-web"
echo "  查看日志:     journalctl -u jz-web -f"
echo "=========================================="