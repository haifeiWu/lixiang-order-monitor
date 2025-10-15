#!/bin/bash

# 理想汽车订单监控通知测试脚本

echo "🧪 测试理想汽车订单监控通知功能"
echo "=================================="

# 检查配置文件
if [ ! -f "config.yaml" ]; then
    echo "❌ 配置文件 config.yaml 不存在！"
    echo "请先根据 README.md 的说明创建配置文件"
    exit 1
fi

# 检查程序是否已编译
if [ ! -f "lixiang-monitor" ]; then
    echo "📦 程序未编译，正在编译..."
    go build -o lixiang-monitor main.go
    if [ $? -ne 0 ]; then
        echo "❌ 编译失败！"
        exit 1
    fi
    echo "✅ 编译成功！"
fi

echo ""
echo "🔍 检查配置..."

# 检查通知配置
WEBHOOK_URL=$(grep "wechat_webhook_url:" config.yaml | awk '{print $2}' | tr -d '"')
SERVERCHAN_KEY=$(grep "serverchan_sendkey:" config.yaml | awk '{print $2}' | tr -d '"')

NOTIFICATION_COUNT=0

if [ "$WEBHOOK_URL" != "" ] && [ "$WEBHOOK_URL" != '""' ]; then
    echo "📱 微信群机器人: ✅ 已配置"
    NOTIFICATION_COUNT=$((NOTIFICATION_COUNT + 1))
    
    echo "🧪 测试微信群机器人通知..."
    curl -s -X POST "$WEBHOOK_URL" \
         -H "Content-Type: application/json" \
         -d '{
             "msgtype": "text",
             "text": {
                 "content": "🧪 理想汽车监控测试消息\n\n这是一条测试消息，用于验证微信群机器人通知功能是否正常工作。\n\n发送时间: '$(date '+%Y-%m-%d %H:%M:%S')'"
             }
         }' > /dev/null
    
    if [ $? -eq 0 ]; then
        echo "✅ 微信群机器人测试消息发送成功"
    else
        echo "❌ 微信群机器人测试消息发送失败"
    fi
else
    echo "📱 微信群机器人: ❌ 未配置"
fi

if [ "$SERVERCHAN_KEY" != "" ] && [ "$SERVERCHAN_KEY" != '""' ]; then
    echo "📧 ServerChan: ✅ 已配置"
    NOTIFICATION_COUNT=$((NOTIFICATION_COUNT + 1))
    
    echo "🧪 测试 ServerChan 通知..."
    SERVERCHAN_URL=$(grep "serverchan_baseurl:" config.yaml | awk '{print $2}' | tr -d '"')
    if [ "$SERVERCHAN_URL" == "" ] || [ "$SERVERCHAN_URL" == '""' ]; then
        SERVERCHAN_URL="https://sctapi.ftqq.com/"
    fi
    
    curl -s -X POST "${SERVERCHAN_URL}${SERVERCHAN_KEY}.send" \
         -d "title=🧪 理想汽车监控测试消息" \
         -d "desp=这是一条测试消息，用于验证 ServerChan 通知功能是否正常工作。%0A%0A发送时间: $(date '+%Y-%m-%d %H:%M:%S')" > /dev/null
    
    if [ $? -eq 0 ]; then
        echo "✅ ServerChan 测试消息发送成功"
    else
        echo "❌ ServerChan 测试消息发送失败"
    fi
else
    echo "📧 ServerChan: ❌ 未配置"
fi

echo ""

if [ $NOTIFICATION_COUNT -eq 0 ]; then
    echo "⚠️  警告: 未配置任何通知方式，无法发送测试消息"
    echo "请编辑 config.yaml 文件配置至少一种通知方式："
    echo "- 微信群机器人: 参考 WECHAT_SETUP.md"
    echo "- ServerChan: 参考 SERVERCHAN_SETUP.md"
else
    echo "🎉 测试完成！已向 $NOTIFICATION_COUNT 种通知方式发送测试消息"
    echo "请检查你的微信是否收到测试消息"
    echo ""
    echo "如果收到测试消息，说明通知功能配置正确！"
    echo "现在可以启动监控服务: ./start.sh"
fi