#!/bin/bash

# 设置日志文件路径
LOG_FILE="./deploy.log"

# 清空之前的日志内容
> $LOG_FILE

echo "开始构建程序喵..."
GOOS=linux GOARCH=amd64 go build -o rps_server 2>>$LOG_FILE
if [ $? -ne 0 ]; then
    echo "构建失败，请检查日志文件 $LOG_FILE."
    exit 1
else
    echo "构建成功了喵."
fi

echo "确保服务器上的目录存在喵..."
sshpass -p "1850560Dwc" ssh root@47.99.133.66 'mkdir -p /home/server/' 2>>$LOG_FILE
if [ $? -ne 0 ]; then
    echo "创建目录失败，请检查日志文件 $LOG_FILE."
    exit 1
else
    echo "目录创建成功或已存在喵."
fi
# 删除服务器上现有的同名文件，以避免权限或覆盖问题
echo "确保服务器上的程序已停止运行喵..."
sshpass -p "1850560Dwc" ssh root@47.99.133.66 'killall -9 rps_server'
echo "正在删除旧文件喵..."
sshpass -p "1850560Dwc" ssh root@47.99.133.66 'sudo rm -f /home/server/rps_server'

echo "正在上传到服务器喵..."
sshpass -p "1850560Dwc" scp ./rps_server root@47.99.133.66:/home/server/ 2>>$LOG_FILE
if [ $? -ne 0 ]; then
    echo "上传失败，请检查日志文件 $LOG_FILE."
    exit 1
else
    echo "上传成功了喵."
fi

echo "远程执行程序中，下面是前线报道喵..."
sshpass -p "1850560Dwc" ssh root@47.99.133.66 'chmod +x /home/server//rps_server && /home/server/rps_server' 2>>$LOG_FILE
if [ $? -ne 0 ]; then
    echo "执行失败，请检查日志文件 $LOG_FILE."
    exit 1
else
    echo "程序执行成功喵."
fi
