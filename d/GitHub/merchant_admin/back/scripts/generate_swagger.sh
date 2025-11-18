#!/bin/bash

# 生成 Swagger 文档脚本
echo "正在生成 Swagger 文档..."

# 进入项目根目录
cd "$(dirname "$0")/.."

# 运行 swag init 命令生成文档
swag init

echo "Swagger 文档生成完成！"
echo "请访问: http://localhost:8080/swagger/index.html"