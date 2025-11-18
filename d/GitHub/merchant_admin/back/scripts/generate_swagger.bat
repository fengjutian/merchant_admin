@echo off
echo 正在生成 Swagger 文档...

REM 进入项目根目录
cd /d "%~dp0.."

REM 运行 swag init 命令生成文档
swag init

echo Swagger 文档生成完成！
echo 请访问: http://localhost:8080/swagger/index.html
pause