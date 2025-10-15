#!/bin/bash

# 理想汽车订单监控状态查看脚本

PROGRAM_NAME="lixiang-monitor"
LOG_FILE="monitor.log"

echo "📊 理想汽车订单监控状态"
echo "=========================="

# 检查程序是否在运行
PID=$(pgrep -f "$PROGRAM_NAME")
if [ ! -z "$PID" ]; then
    echo "🟢 状态: 运行中"
    echo "📍 进程ID: $PID"
    
    # 显示进程信息
    PS_INFO=$(ps -p $PID -o pid,ppid,etime,pcpu,pmem,cmd --no-headers)
    echo "💾 进程信息: $PS_INFO"
    
    # 显示启动时间
    START_TIME=$(ps -p $PID -o lstart --no-headers)
    echo "⏰ 启动时间: $START_TIME"
else
    echo "🔴 状态: 未运行"
fi

echo ""

# 检查日志文件
if [ -f "$LOG_FILE" ]; then
    echo "📋 最近日志 (最后10行):"
    echo "------------------------"
    tail -10 "$LOG_FILE"
    echo ""
    echo "📝 查看完整日志: tail -f $LOG_FILE"
else
    echo "📋 日志文件不存在"
fi

echo ""

# 检查配置文件
if [ -f "config.yaml" ]; then
    echo "⚙️  配置文件: ✅ 存在"
    
    # 检查关键配置项
    ORDER_ID=$(grep "order_id:" config.yaml | awk '{print $2}' | tr -d '"')
    WEBHOOK_URL=$(grep "wechat_webhook_url:" config.yaml | awk '{print $2}' | tr -d '"')
    SERVERCHAN_KEY=$(grep "serverchan_sendkey:" config.yaml | awk '{print $2}' | tr -d '"')
    CHECK_INTERVAL=$(grep "check_interval:" config.yaml | awk '{print $2}' | tr -d '"')
    LOCK_ORDER_TIME=$(grep "lock_order_time:" config.yaml | awk '{print $2 " " $3}' | tr -d '"')
    ESTIMATE_MIN=$(grep "estimate_weeks_min:" config.yaml | awk '{print $2}' | tr -d '"')
    ESTIMATE_MAX=$(grep "estimate_weeks_max:" config.yaml | awk '{print $2}' | tr -d '"')
    
    echo "📦 订单ID: $ORDER_ID"
    echo "🔄 检查间隔: $CHECK_INTERVAL"
    echo "📅 锁单时间: $LOCK_ORDER_TIME"
    echo "⏱️  预计交付: $ESTIMATE_MIN-$ESTIMATE_MAX 周"
    
    # 检查通知配置
    NOTIFICATION_COUNT=0
    if [ "$WEBHOOK_URL" != "" ] && [ "$WEBHOOK_URL" != '""' ]; then
        echo "📱 微信群机器人: ✅ 已配置"
        NOTIFICATION_COUNT=$((NOTIFICATION_COUNT + 1))
    else
        echo "📱 微信群机器人: ❌ 未配置"
    fi
    
    if [ "$SERVERCHAN_KEY" != "" ] && [ "$SERVERCHAN_KEY" != '""' ]; then
        echo "📧 ServerChan: ✅ 已配置"
        NOTIFICATION_COUNT=$((NOTIFICATION_COUNT + 1))
    else
        echo "📧 ServerChan: ❌ 未配置"
    fi
    
    if [ $NOTIFICATION_COUNT -eq 0 ]; then
        echo "⚠️  警告: 未配置任何通知方式"
    else
        echo "✅ 已配置 $NOTIFICATION_COUNT 种通知方式"
    fi
else
    echo "⚙️  配置文件: ❌ 不存在"
fi

echo ""
echo "常用命令:"
echo "🚀 启动: ./start.sh"
echo "🛑 停止: ./stop.sh"
echo "📊 状态: ./status.sh"
echo "🧪 测试: ./test-notification.sh"
echo "📋 日志: tail -f $LOG_FILE"