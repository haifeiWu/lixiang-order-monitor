package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"lixiang-monitor/db"
)

//go:embed templates/*
var templatesFS embed.FS

// Server Web 服务器
type Server struct {
	database   *db.Database
	orderID    string
	port       int
	httpServer *http.Server
	templates  *template.Template
}

// NewServer 创建 Web 服务器实例
func NewServer(database *db.Database, orderID string, port int) (*Server, error) {
	// 解析模板
	tmpl, err := template.ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("解析模板失败: %w", err)
	}

	server := &Server{
		database:  database,
		orderID:   orderID,
		port:      port,
		templates: tmpl,
	}

	return server, nil
}

// Start 启动 Web 服务器
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// 注册路由
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/api/stats", s.handleStats)
	mux.HandleFunc("/api/records", s.handleRecords)
	mux.HandleFunc("/api/time-changes", s.handleTimeChanges)

	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      s.logMiddleware(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("[Web] 启动 Web 服务器: http://localhost:%d", s.port)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("[Web] 服务器启动失败: %v", err)
		}
	}()

	return nil
}

// Stop 停止 Web 服务器
func (s *Server) Stop() error {
	if s.httpServer != nil {
		log.Println("[Web] 正在关闭 Web 服务器...")
		return s.httpServer.Close()
	}
	return nil
}

// logMiddleware 日志中间件
func (s *Server) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[Web] %s %s - %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// handleIndex 处理首页
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := map[string]interface{}{
		"OrderID": s.orderID,
		"Title":   "理想汽车订单监控",
	}

	if err := s.templates.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("[Web] 渲染模板失败: %v", err)
	}
}

// StatsResponse 统计响应
type StatsResponse struct {
	TotalRecords      int                `json:"total_records"`
	TimeChangedCount  int                `json:"time_changed_count"`
	NotificationCount int                `json:"notification_count"`
	LatestCheckTime   string             `json:"latest_check_time"`
	LatestEstimate    string             `json:"latest_estimate"`
	FirstCheckTime    string             `json:"first_check_time"`
	MonitoringDays    int                `json:"monitoring_days"`
	LatestRecord      *db.DeliveryRecord `json:"latest_record"`
}

// handleStats 处理统计数据
func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 获取记录总数
	totalRecords, err := s.database.GetRecordsCount(s.orderID)
	if err != nil {
		s.sendJSONError(w, "查询记录总数失败", http.StatusInternalServerError)
		return
	}

	// 获取最新记录
	latestRecord, err := s.database.GetLatestRecord(s.orderID)
	if err != nil {
		s.sendJSONError(w, "查询最新记录失败", http.StatusInternalServerError)
		return
	}

	// 获取所有记录用于统计
	allRecords, err := s.database.GetRecordsByOrderID(s.orderID, totalRecords)
	if err != nil {
		s.sendJSONError(w, "查询记录失败", http.StatusInternalServerError)
		return
	}

	// 统计时间变更和通知次数
	timeChangedCount := 0
	notificationCount := 0
	var firstCheckTime time.Time

	for i, record := range allRecords {
		if record.TimeChanged {
			timeChangedCount++
		}
		if record.NotificationSent {
			notificationCount++
		}
		if i == len(allRecords)-1 { // 最后一条是最早的
			firstCheckTime = record.CheckTime
		}
	}

	// 计算监控天数
	monitoringDays := 0
	if !firstCheckTime.IsZero() && latestRecord != nil {
		monitoringDays = int(latestRecord.CheckTime.Sub(firstCheckTime).Hours() / 24)
	}

	stats := StatsResponse{
		TotalRecords:      totalRecords,
		TimeChangedCount:  timeChangedCount,
		NotificationCount: notificationCount,
		MonitoringDays:    monitoringDays,
		LatestRecord:      latestRecord,
	}

	if latestRecord != nil {
		stats.LatestCheckTime = latestRecord.CheckTime.Format("2006-01-02 15:04:05")
		stats.LatestEstimate = latestRecord.EstimateTime
	}
	if !firstCheckTime.IsZero() {
		stats.FirstCheckTime = firstCheckTime.Format("2006-01-02 15:04:05")
	}

	json.NewEncoder(w).Encode(stats)
}

// handleRecords 处理记录查询
func (s *Server) handleRecords(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 获取分页参数
	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	records, err := s.database.GetRecordsByOrderID(s.orderID, limit)
	if err != nil {
		s.sendJSONError(w, "查询记录失败", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(records)
}

// handleTimeChanges 处理时间变更记录查询
func (s *Server) handleTimeChanges(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 获取分页参数
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	records, err := s.database.GetTimeChangedRecords(s.orderID, limit)
	if err != nil {
		s.sendJSONError(w, "查询时间变更记录失败", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(records)
}

// sendJSONError 发送 JSON 错误响应
func (s *Server) sendJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}
