#!/bin/bash

# 查询数据库历史记录的脚本

DB_FILE="./lixiang-monitor.db"

if [ ! -f "$DB_FILE" ]; then
    echo "❌ 数据库文件不存在: $DB_FILE"
    exit 1
fi

echo "📊 理想汽车订单监控 - 历史记录"
echo "================================"
echo ""

# 检查是否安装了 sqlite3
if ! command -v sqlite3 &> /dev/null; then
    echo "❌ 未安装 sqlite3 命令"
    echo "请安装: brew install sqlite3"
    exit 1
fi

# 查询记录总数
echo "📈 记录统计"
echo "--------------------------------"
sqlite3 "$DB_FILE" <<EOF
SELECT 
    COUNT(*) as '总记录数',
    COUNT(DISTINCT order_id) as '订单数',
    COUNT(CASE WHEN time_changed = 1 THEN 1 END) as '时间变更次数',
    COUNT(CASE WHEN notification_sent = 1 THEN 1 END) as '通知发送次数'
FROM delivery_records;
EOF

echo ""
echo "📋 最近 10 条记录"
echo "--------------------------------"
sqlite3 -column -header "$DB_FILE" <<EOF
SELECT 
    id,
    substr(order_id, -6) as order_id,
    estimate_time,
    strftime('%Y-%m-%d %H:%M', check_time) as check_time,
    CASE WHEN is_approaching THEN '是' ELSE '否' END as approaching,
    CASE WHEN time_changed THEN '是' ELSE '否' END as changed,
    CASE WHEN notification_sent THEN '是' ELSE '否' END as notified
FROM delivery_records
ORDER BY check_time DESC
LIMIT 10;
EOF

echo ""
echo "📊 时间变更历史"
echo "--------------------------------"
sqlite3 -column -header "$DB_FILE" <<EOF
SELECT 
    strftime('%Y-%m-%d %H:%M', check_time) as check_time,
    previous_estimate as '旧时间',
    estimate_time as '新时间'
FROM delivery_records
WHERE time_changed = 1
ORDER BY check_time DESC
LIMIT 10;
EOF

echo ""
echo "✅ 查询完成"
