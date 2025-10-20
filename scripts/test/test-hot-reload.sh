#!/bin/bash

# 配置热加载功能测试脚本

echo "=========================================="
echo "    理想汽车订单监控 - 配置热加载测试"
echo "=========================================="
echo ""

# 检查配置文件是否存在
if [ ! -f "config.yaml" ]; then
    echo "❌ 错误: config.yaml 文件不存在"
    exit 1
fi

# 备份配置文件
echo "📦 备份当前配置文件..."
cp config.yaml config.yaml.test-backup
echo "✅ 配置文件已备份到 config.yaml.test-backup"
echo ""

echo "📝 测试说明:"
echo "1. 此脚本将帮助您测试配置热加载功能"
echo "2. 建议在程序运行时执行此脚本"
echo "3. 修改配置后，观察程序日志确认是否自动重新加载"
echo ""

# 显示当前配置
echo "📋 当前配置:"
echo "----------------------------------------"
grep -E "^(order_id|enable_periodic_notify|notification_interval_hours):" config.yaml
echo "----------------------------------------"
echo ""

# 提示用户选择测试场景
echo "请选择要测试的场景:"
echo "1. 修改通知间隔时间"
echo "2. 切换定期通知开关"
echo "3. 恢复备份配置"
echo "0. 退出"
echo ""

read -p "请输入选项 (0-3): " choice

case $choice in
    1)
        echo ""
        read -p "请输入新的通知间隔时间（小时）[当前值请查看上方配置]: " hours
        if [[ "$hours" =~ ^[0-9]+$ ]]; then
            # 使用 sed 修改配置文件
            if [[ "$OSTYPE" == "darwin"* ]]; then
                # macOS
                sed -i '' "s/^notification_interval_hours:.*/notification_interval_hours: $hours/" config.yaml
            else
                # Linux
                sed -i "s/^notification_interval_hours:.*/notification_interval_hours: $hours/" config.yaml
            fi
            echo "✅ 通知间隔已修改为 $hours 小时"
            echo "👀 请查看程序日志，确认配置是否自动重新加载"
        else
            echo "❌ 输入无效，必须是数字"
        fi
        ;;
    2)
        echo ""
        current_value=$(grep "^enable_periodic_notify:" config.yaml | awk '{print $2}')
        if [ "$current_value" = "true" ]; then
            new_value="false"
        else
            new_value="true"
        fi
        
        if [[ "$OSTYPE" == "darwin"* ]]; then
            # macOS
            sed -i '' "s/^enable_periodic_notify:.*/enable_periodic_notify: $new_value/" config.yaml
        else
            # Linux
            sed -i "s/^enable_periodic_notify:.*/enable_periodic_notify: $new_value/" config.yaml
        fi
        echo "✅ 定期通知已切换为: $new_value"
        echo "👀 请查看程序日志，确认配置是否自动重新加载"
        ;;
    3)
        if [ -f "config.yaml.test-backup" ]; then
            cp config.yaml.test-backup config.yaml
            echo "✅ 配置已恢复到测试前的状态"
            echo "👀 请查看程序日志，确认配置是否自动重新加载"
        else
            echo "❌ 备份文件不存在"
        fi
        ;;
    0)
        echo "退出测试"
        ;;
    *)
        echo "❌ 无效选项"
        ;;
esac

echo ""
echo "=========================================="
echo "💡 提示:"
echo "  - 查看日志: tail -f lixiang-monitor.log"
echo "  - 查看进程: ps aux | grep lixiang-monitor"
echo "  - 备份文件: config.yaml.test-backup"
echo "=========================================="
