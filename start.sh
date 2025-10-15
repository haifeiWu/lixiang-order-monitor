#!/bin/bash

# 理想汽车订单监控启动脚本

PROGRAM_NAME="lixiang-monitor"
CONFIG_FILE="config.yaml"
LOG_FILE="monitor.log"

# 检查配置文件是否存在
if [ ! -f "$CONFIG_FILE" ]; then
    echo "❌ 配置文件 $CONFIG_FILE 不存在！"
    echo "请先根据 README.md 的说明创建配置文件"
    exit 1
fi

# 检查程序是否已编译
if [ ! -f "$PROGRAM_NAME" ]; then
    echo "📦 程序未编译，正在编译..."
    go build -o "$PROGRAM_NAME" main.go
    if [ $? -ne 0 ]; then
        echo "❌ 编译失败！"
        exit 1
    fi
    echo "✅ 编译成功！"
fi

# 检查是否已经在运行
PID=$(pgrep -f "$PROGRAM_NAME")
if [ ! -z "$PID" ]; then
    echo "⚠️  程序已经在运行 (PID: $PID)"
    echo "如需重启，请先运行: ./stop.sh"
    exit 1
fi

echo "🚀 启动理想汽车订单监控..."
echo "📊 日志文件: $LOG_FILE"
echo "⏹️  停止程序: ./stop.sh"
echo "📋 查看日志: tail -f $LOG_FILE"
echo ""

# 后台启动程序
nohup ./"$PROGRAM_NAME" > "$LOG_FILE" 2>&1 &
NEW_PID=$!

# 等待一下检查程序是否启动成功
sleep 2
if kill -0 $NEW_PID 2>/dev/null; then
    echo "✅ 程序启动成功！(PID: $NEW_PID)"
    echo "📱 请确保已正确配置微信 Webhook URL"
    echo ""
    echo "查看实时日志："
    echo "tail -f $LOG_FILE"
else
    echo "❌ 程序启动失败！请检查日志: $LOG_FILE"
    exit 1
fi