@echo off
echo 正在生成 Swagger 文档...

REM 进入项目根目录
cd /d "%~dp0.."

REM 删除旧的文档文件
if exist "docs\docs.go" del "docs\docs.go"
if exist "docs\swagger.json" del "docs\swagger.json"
if exist "docs\swagger.yaml" del "docs\swagger.yaml"

REM 运行 swag init 命令生成文档，指定搜索路径
swag init -g cmd/main.go -o docs --parseDependency --parseInternal

echo Swagger 文档生成完成！
echo 请访问: http://localhost:8080/swagger/index.html
pause