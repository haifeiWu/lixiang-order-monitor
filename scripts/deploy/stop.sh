#!/bin/bash

# 理想汽车订单监控停止脚本

PROGRAM_NAME="lixiang-monitor"

echo "🛑 正在停止理想汽车订单监控..."

# 查找程序进程
PID=$(pgrep -f "$PROGRAM_NAME")

if [ -z "$PID" ]; then
    echo "ℹ️  程序未运行"
    exit 0
fi

# 停止程序
echo "📝 找到程序进程 (PID: $PID)"
kill $PID

# 等待程序停止
for i in {1..5}; do
    if ! kill -0 $PID 2>/dev/null; then
        echo "✅ 程序已停止"
        exit 0
    fi
    echo "⏳ 等待程序停止... ($i/5)"
    sleep 1
done

# 如果程序仍在运行，强制停止
if kill -0 $PID 2>/dev/null; then
    echo "⚡ 强制停止程序..."
    kill -9 $PID
    sleep 1
    if ! kill -0 $PID 2>/dev/null; then
        echo "✅ 程序已强制停止"
    else
        echo "❌ 无法停止程序，请手动处理"
        exit 1
    fi
fi