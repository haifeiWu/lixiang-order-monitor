package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite" // 纯 Go 实现的 SQLite 驱动，无需 CGO
)

// DeliveryRecord 交付记录
type DeliveryRecord struct {
	ID               int       `json:"id"`
	OrderID          string    `json:"order_id"`
	EstimateTime     string    `json:"estimate_time"`     // 预计交付时间
	LockOrderTime    time.Time `json:"lock_order_time"`   // 锁单时间
	CheckTime        time.Time `json:"check_time"`        // 检查时间
	IsApproaching    bool      `json:"is_approaching"`    // 是否临近交付
	ApproachMessage  string    `json:"approach_message"`  // 临近提示信息
	TimeChanged      bool      `json:"time_changed"`      // 时间是否变化
	PreviousEstimate string    `json:"previous_estimate"` // 之前的预计时间
	NotificationSent bool      `json:"notification_sent"` // 是否发送了通知
	CreatedAt        time.Time `json:"created_at"`
}

// Database 数据库管理器
type Database struct {
	db *sql.DB
}

// New 创建数据库实例
func New(dbPath string) (*Database, error) {
	// 打开数据库连接
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("打开数据库失败: %w", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	database := &Database{db: db}

	// 初始化表结构
	if err := database.initTables(); err != nil {
		return nil, fmt.Errorf("初始化表结构失败: %w", err)
	}

	log.Printf("[DB] 数据库初始化成功: %s", dbPath)
	return database, nil
}

// initTables 初始化数据库表
func (d *Database) initTables() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS delivery_records (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		order_id TEXT NOT NULL,
		estimate_time TEXT NOT NULL,
		lock_order_time DATETIME NOT NULL,
		check_time DATETIME NOT NULL,
		is_approaching BOOLEAN NOT NULL DEFAULT 0,
		approach_message TEXT,
		time_changed BOOLEAN NOT NULL DEFAULT 0,
		previous_estimate TEXT,
		notification_sent BOOLEAN NOT NULL DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_order_id ON delivery_records(order_id);
	CREATE INDEX IF NOT EXISTS idx_check_time ON delivery_records(check_time);
	CREATE INDEX IF NOT EXISTS idx_created_at ON delivery_records(created_at);
	`

	_, err := d.db.Exec(createTableSQL)
	return err
}

// SaveDeliveryRecord 保存交付记录
func (d *Database) SaveDeliveryRecord(record *DeliveryRecord) error {
	query := `
	INSERT INTO delivery_records (
		order_id, estimate_time, lock_order_time, check_time,
		is_approaching, approach_message, time_changed, 
		previous_estimate, notification_sent, created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := d.db.Exec(query,
		record.OrderID,
		record.EstimateTime,
		record.LockOrderTime,
		record.CheckTime,
		record.IsApproaching,
		record.ApproachMessage,
		record.TimeChanged,
		record.PreviousEstimate,
		record.NotificationSent,
		record.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("保存记录失败: %w", err)
	}

	log.Printf("[DB] 已保存交付记录: order_id=%s, estimate_time=%s, is_approaching=%v",
		record.OrderID, record.EstimateTime, record.IsApproaching)
	return nil
}

// GetLatestRecord 获取指定订单的最新记录
func (d *Database) GetLatestRecord(orderID string) (*DeliveryRecord, error) {
	query := `
	SELECT id, order_id, estimate_time, lock_order_time, check_time,
		   is_approaching, approach_message, time_changed, 
		   previous_estimate, notification_sent, created_at
	FROM delivery_records
	WHERE order_id = ?
	ORDER BY check_time DESC
	LIMIT 1
	`

	record := &DeliveryRecord{}
	err := d.db.QueryRow(query, orderID).Scan(
		&record.ID,
		&record.OrderID,
		&record.EstimateTime,
		&record.LockOrderTime,
		&record.CheckTime,
		&record.IsApproaching,
		&record.ApproachMessage,
		&record.TimeChanged,
		&record.PreviousEstimate,
		&record.NotificationSent,
		&record.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // 没有记录
	}

	if err != nil {
		return nil, fmt.Errorf("查询记录失败: %w", err)
	}

	return record, nil
}

// GetRecordsByOrderID 获取指定订单的所有记录
func (d *Database) GetRecordsByOrderID(orderID string, limit int) ([]*DeliveryRecord, error) {
	query := `
	SELECT id, order_id, estimate_time, lock_order_time, check_time,
		   is_approaching, approach_message, time_changed, 
		   previous_estimate, notification_sent, created_at
	FROM delivery_records
	WHERE order_id = ?
	ORDER BY check_time DESC
	LIMIT ?
	`

	rows, err := d.db.Query(query, orderID, limit)
	if err != nil {
		return nil, fmt.Errorf("查询记录失败: %w", err)
	}
	defer rows.Close()

	var records []*DeliveryRecord
	for rows.Next() {
		record := &DeliveryRecord{}
		err := rows.Scan(
			&record.ID,
			&record.OrderID,
			&record.EstimateTime,
			&record.LockOrderTime,
			&record.CheckTime,
			&record.IsApproaching,
			&record.ApproachMessage,
			&record.TimeChanged,
			&record.PreviousEstimate,
			&record.NotificationSent,
			&record.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描记录失败: %w", err)
		}
		records = append(records, record)
	}

	return records, nil
}

// GetRecordsCount 获取记录总数
func (d *Database) GetRecordsCount(orderID string) (int, error) {
	query := `SELECT COUNT(*) FROM delivery_records WHERE order_id = ?`

	var count int
	err := d.db.QueryRow(query, orderID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("查询记录数失败: %w", err)
	}

	return count, nil
}

// GetTimeChangedRecords 获取时间发生变化的记录
func (d *Database) GetTimeChangedRecords(orderID string, limit int) ([]*DeliveryRecord, error) {
	query := `
	SELECT id, order_id, estimate_time, lock_order_time, check_time,
		   is_approaching, approach_message, time_changed, 
		   previous_estimate, notification_sent, created_at
	FROM delivery_records
	WHERE order_id = ? AND time_changed = 1
	ORDER BY check_time DESC
	LIMIT ?
	`

	rows, err := d.db.Query(query, orderID, limit)
	if err != nil {
		return nil, fmt.Errorf("查询记录失败: %w", err)
	}
	defer rows.Close()

	var records []*DeliveryRecord
	for rows.Next() {
		record := &DeliveryRecord{}
		err := rows.Scan(
			&record.ID,
			&record.OrderID,
			&record.EstimateTime,
			&record.LockOrderTime,
			&record.CheckTime,
			&record.IsApproaching,
			&record.ApproachMessage,
			&record.TimeChanged,
			&record.PreviousEstimate,
			&record.NotificationSent,
			&record.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描记录失败: %w", err)
		}
		records = append(records, record)
	}

	return records, nil
}

// Close 关闭数据库连接
func (d *Database) Close() error {
	if d.db != nil {
		log.Println("[DB] 关闭数据库连接")
		return d.db.Close()
	}
	return nil
}
