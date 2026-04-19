#!/bin/bash
# 前端更新脚本 - 仅更新前端静态文件

echo "重新构建前端..."
cd /workspace/pro/jz_web/frontend
sudo rm -rf dist
npm run build

echo "复制到项目目录..."
sudo cp -r dist/* /workspace/pro/jz_web/frontend/dist/

echo "重载Nginx..."
sudo systemctl reload nginx

echo "前端更新完成！"