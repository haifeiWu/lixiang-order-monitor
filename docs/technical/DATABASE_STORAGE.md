# 数据库功能说明

## 概述

理想汽车订单监控系统已集成 SQLite 数据库，用于持久化存储历史监控记录。使用纯 Go 实现的 SQLite 驱动（modernc.org/sqlite），无需 CGO 编译，保持项目的轻量化和跨平台特性。

## 功能特性

### ✅ 已实现功能

1. **自动记录保存**
   - 每次订单检查自动保存记录
   - 记录包含：订单 ID、预计交付时间、检查时间、通知状态等
   - 异常情况不影响主程序运行

2. **完整信息追踪**
   - 订单 ID
   - 预计交付时间
   - 锁单时间
   - 检查时间
   - 是否临近交付
   - 临近提示信息
   - 时间是否变化
   - 之前的预计时间
   - 是否发送通知

3. **数据查询能力**
   - 获取最新记录
   - 获取指定订单的所有记录
   - 获取时间发生变化的记录
   - 统计记录总数

## 技术实现

### 数据库设计

**表名**: `delivery_records`

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | INTEGER PRIMARY KEY | 自增主键 |
| order_id | TEXT | 订单 ID |
| estimate_time | TEXT | 预计交付时间 |
| lock_order_time | DATETIME | 锁单时间 |
| check_time | DATETIME | 检查时间 |
| is_approaching | BOOLEAN | 是否临近交付 |
| approach_message | TEXT | 临近提示信息 |
| time_changed | BOOLEAN | 时间是否变化 |
| previous_estimate | TEXT | 之前的预计时间 |
| notification_sent | BOOLEAN | 是否发送通知 |
| created_at | DATETIME | 记录创建时间 |

**索引**:
- `idx_order_id`: 订单 ID 索引
- `idx_check_time`: 检查时间索引
- `idx_created_at`: 创建时间索引

### 代码结构

```
db/
└── database.go      # 数据库管理器
    ├── Database     # 数据库连接管理
    ├── New()        # 初始化数据库
    ├── SaveDeliveryRecord()        # 保存记录
    ├── GetLatestRecord()           # 获取最新记录
    ├── GetRecordsByOrderID()       # 获取订单所有记录
    ├── GetRecordsCount()           # 获取记录总数
    ├── GetTimeChangedRecords()     # 获取时间变更记录
    └── Close()                     # 关闭连接
```

## 使用示例

### 1. 自动记录（无需手动操作）

程序运行时会自动保存每次检查的记录：

```
2025/10/23 16:34:54 [DB] 已保存交付记录: order_id=177971759268550919, estimate_time=预计6-8周内交付, is_approaching=false
```

### 2. 查询历史记录

使用提供的查询脚本：

```bash
./scripts/query-db.sh
```

输出示例：
```
📊 理想汽车订单监控 - 历史记录
================================

📈 记录统计
--------------------------------
总记录数: 25
订单数: 1
时间变更次数: 3
通知发送次数: 12

📋 最近 10 条记录
--------------------------------
id  order_id  estimate_time      check_time           approaching  changed  notified
--  --------  ----------------   ------------------   -----------  -------  --------
25  550919    预计6-8周内交付    2025-10-23 16:34     否           否       是
24  550919    预计6-8周内交付    2025-10-23 04:34     否           否       是
...

📊 时间变更历史
--------------------------------
check_time           旧时间              新时间
------------------   -----------------   -----------------
2025-10-22 08:15     预计7-9周内交付     预计6-8周内交付
2025-10-20 14:23     预计8-10周内交付    预计7-9周内交付
```

### 3. 使用 sqlite3 直接查询

```bash
# 进入数据库
sqlite3 lixiang-monitor.db

# 查询所有记录
SELECT * FROM delivery_records ORDER BY check_time DESC LIMIT 10;

# 查询时间变更记录
SELECT check_time, previous_estimate, estimate_time 
FROM delivery_records 
WHERE time_changed = 1;

# 统计通知发送情况
SELECT 
    DATE(check_time) as date,
    COUNT(*) as checks,
    SUM(notification_sent) as notifications
FROM delivery_records
GROUP BY DATE(check_time);
```

## 数据库文件

- **位置**: `./lixiang-monitor.db`
- **大小**: ~24KB（初始）
- **备份**: 建议定期备份此文件

### 备份方法

```bash
# 手动备份
cp lixiang-monitor.db lixiang-monitor.db.backup

# 自动备份（添加到 cron）
0 2 * * * cp /path/to/lixiang-monitor.db /path/to/backups/lixiang-monitor-$(date +\%Y\%m\%d).db
```

## 性能考虑

1. **轻量级**: 使用 SQLite 内嵌数据库，无需额外数据库服务
2. **高效**: 建立索引，查询性能优异
3. **可靠**: 自动提交事务，数据安全
4. **容错**: 数据库操作失败不影响主程序运行

## 依赖说明

### modernc.org/sqlite

- **版本**: v1.39.1
- **特点**: 
  - 纯 Go 实现
  - 无需 CGO
  - 跨平台兼容
  - 完整的 SQLite 3 支持

## 故障排查

### 数据库初始化失败

**症状**: 
```
⚠️  数据库初始化失败: ... (历史记录功能将不可用)
```

**解决方法**:
1. 检查文件权限
2. 确保磁盘空间充足
3. 检查是否已有其他进程占用数据库文件

### 记录保存失败

**症状**:
```
保存交付记录失败: ...
```

**解决方法**:
1. 不影响主程序运行，仅记录日志
2. 检查数据库文件是否损坏
3. 尝试删除数据库文件，程序会自动重建

## 未来扩展

可能的功能扩展方向：

1. **数据分析**
   - 交付时间变化趋势分析
   - 预测准确性统计
   - 通知效果分析

2. **Web 界面**
   - 历史记录可视化
   - 交付时间趋势图表
   - 统计报表

3. **数据导出**
   - CSV 导出
   - JSON 导出
   - Excel 报表

4. **数据清理**
   - 自动清理过期记录
   - 数据归档功能
   - 数据压缩

## 总结

SQLite 数据库集成为监控系统提供了可靠的历史数据存储能力，同时保持了系统的轻量化特性。所有操作对用户透明，无需额外配置即可使用。

---

**更新时间**: 2025-10-23  
**版本**: 1.0.0
