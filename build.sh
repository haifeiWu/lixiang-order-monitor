#!/bin/bash

echo "正在编译理想汽车订单监控程序..."

# 编译程序
go build -o lixiang-monitor main.go

if [ $? -eq 0 ]; then
    echo "✅ 编译成功！"
    echo ""
    echo "使用说明："
    echo "1. 编辑 config.yaml 文件，配置你的订单ID、微信Webhook URL和Cookie"
    echo "2. 运行程序: ./lixiang-monitor"
    echo "3. 或者后台运行: nohup ./lixiang-monitor > monitor.log 2>&1 &"
    echo ""
    echo "配置文件模板已生成为 config.yaml"
    echo "请根据README.md中的说明进行配置"
else
    echo "❌ 编译失败！请检查代码"
    exit 1
fi