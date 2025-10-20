#!/bin/bash

# 理想汽车订单监控 - 定期通知功能测试脚本

echo "🚗 理想汽车订单监控 - 定期通知功能测试"
echo "================================================"
echo ""

# 检查配置文件
if [ ! -f "config.yaml" ]; then
    echo "⚠️  配置文件 config.yaml 不存在，创建示例配置..."
    cat > config.yaml << EOF
# 理想汽车订单监控配置
order_id: "177971759268550919"
check_interval: "@every 5m"  # 测试时使用较短间隔

# 定期通知配置
enable_periodic_notify: true           # 启用定期通知
notification_interval_hours: 1         # 测试时设置为1小时
always_notify_when_approaching: true   # 临近交付时总是通知

# 理想汽车API配置
lixiang_cookies: ""  # 需要从浏览器获取

# 通知配置 (至少配置一个)
serverchan_sendkey: ""     # Server酱的SendKey
wechat_webhook_url: ""     # 微信群机器人Webhook URL

# 预计交付配置
lock_order_time: "2024-09-27 13:08:00"  # 锁单时间
estimate_weeks_min: 7                   # 最少预计周数
estimate_weeks_max: 9                   # 最多预计周数
EOF
    echo "✅ 已创建示例配置文件 config.yaml"
    echo "📝 请根据实际情况修改配置文件中的参数"
    echo ""
fi

echo "📋 当前配置验证:"
echo "- 检查配置文件是否存在: $([ -f config.yaml ] && echo '✅' || echo '❌')"
echo "- 编译可执行文件..."

# 编译程序
if go build -o lixiang-monitor main.go; then
    echo "✅ 编译成功"
else
    echo "❌ 编译失败"
    exit 1
fi

echo ""
echo "🔧 新增功能说明:"
echo "1. ⏰ 定期通知功能"
echo "   - 可配置通知间隔 (notification_interval_hours)"
echo "   - 即使交付时间未更新也会定期发送状态报告"
echo ""
echo "2. 🎯 智能通知策略"
echo "   - 临近交付时间自动提醒"
echo "   - 支持配置是否总是在临近时通知"
echo ""
echo "3. 📊 增强的通知内容"
echo "   - 详细的交付状态分析"
echo "   - 进度百分比显示"
echo "   - 智能的时间状态判断"
echo ""
echo "📁 相关配置文件:"
echo "   config.yaml              - 主配置文件"
echo "   DELIVERY_OPTIMIZATION.md - 优化功能详细说明"
echo ""
echo "🚀 启动监控服务:"
echo "   ./lixiang-monitor"
echo ""
echo "注意：首次运行时请确保已正确配置："
echo "1. 理想汽车网站的 cookies"
echo "2. 至少一个通知渠道 (ServerChan 或微信机器人)"
echo "3. 正确的锁单时间和预计交付周数"