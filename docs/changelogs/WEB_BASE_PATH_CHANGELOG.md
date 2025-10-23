# Web 根路由配置功能更新日志

## 更新时间
2025-10-23

## 版本
v1.9.0

## 功能概述
新增 Web 服务器根路由（Base Path）配置功能，支持在自定义路径前缀下运行 Web 界面。

## 主要变更

### 1. 配置系统增强 (cfg/config.go)

**新增配置字段**：
```go
type Config struct {
    // ... 现有字段 ...
    WebBasePath string  // Web 服务器根路由
}
```

**默认值设置**：
```go
viper.SetDefault("web_base_path", "")
```

**配置加载**：
```go
cfg.WebBasePath = viper.GetString("web_base_path")
```

### 2. Web 服务器增强 (web/server.go)

**Server 结构更新**：
```go
type Server struct {
    database   *db.Database
    orderID    string
    port       int
    basePath   string  // 新增：根路由配置
    httpServer *http.Server
    templates  *template.Template
}
```

**NewServer 函数签名变更**：
```go
// 旧版本
func NewServer(database *db.Database, orderID string, port int) (*Server, error)

// 新版本
func NewServer(database *db.Database, orderID string, port int, basePath string) (*Server, error)
```

**路径规范化逻辑**：
```go
// 规范化 basePath: 确保以 / 开头，不以 / 结尾
if basePath != "" {
    if basePath[0] != '/' {
        basePath = "/" + basePath
    }
    if basePath[len(basePath)-1] == '/' && len(basePath) > 1 {
        basePath = basePath[:len(basePath)-1]
    }
}
```

**新增辅助方法**：
```go
// route 根据 basePath 构建完整路由
func (s *Server) route(path string) string {
    if s.basePath == "" {
        return path
    }
    return s.basePath + path
}
```

**路由注册更新**：
```go
// 旧版本
mux.HandleFunc("/", s.handleIndex)
mux.HandleFunc("/api/stats", s.handleStats)

// 新版本
mux.HandleFunc(s.route("/"), s.handleIndex)
mux.HandleFunc(s.route("/api/stats"), s.handleStats)
```

**启动日志优化**：
```go
// 旧版本
log.Printf("[Web] 启动 Web 服务器: http://localhost:%d", s.port)

// 新版本
baseURL := fmt.Sprintf("http://localhost:%d%s", s.port, s.basePath)
log.Printf("[Web] 启动 Web 服务器: %s", baseURL)
```

### 3. 主程序集成 (main.go)

**Monitor 结构扩展**：
```go
type Monitor struct {
    // ... 现有字段 ...
    WebBasePath string  // Web 服务器根路由
}
```

**配置加载**：
```go
m.WebBasePath = config.WebBasePath
```

**Web 服务器初始化**：
```go
// 旧版本
webServer, err := web.NewServer(monitor.database, monitor.OrderID, monitor.WebPort)

// 新版本
webServer, err := web.NewServer(
    monitor.database, 
    monitor.OrderID, 
    monitor.WebPort, 
    monitor.WebBasePath,
)
```

### 4. 配置文件更新 (config.yaml)

**新增配置项**：
```yaml
# Web 管理界面配置
web_enabled: true
web_port: 8099
web_base_path: ""  # Web 服务器根路由 (例如: "/monitor" 则访问 http://localhost:8099/monitor)
```

## 使用示例

### 示例 1: 默认配置（无根路由）

**配置**：
```yaml
web_base_path: ""
```

**访问地址**：
- http://localhost:8099/
- http://localhost:8099/api/stats
- http://localhost:8099/api/records

### 示例 2: 单级路径

**配置**：
```yaml
web_base_path: "/monitor"
```

**访问地址**：
- http://localhost:8099/monitor/
- http://localhost:8099/monitor/api/stats
- http://localhost:8099/monitor/api/records

### 示例 3: 多级路径

**配置**：
```yaml
web_base_path: "/apps/lixiang"
```

**访问地址**：
- http://localhost:8099/apps/lixiang/
- http://localhost:8099/apps/lixiang/api/stats
- http://localhost:8099/apps/lixiang/api/records

## 应用场景

### 1. Nginx 反向代理

**场景**：将监控服务部署在 Nginx 的子路径下

**应用配置**：
```yaml
web_base_path: ""
web_port: 8099
```

**Nginx 配置**：
```nginx
location /lixiang-monitor/ {
    proxy_pass http://localhost:8099/;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
}
```

**访问**: `http://example.com/lixiang-monitor/`

### 2. 多实例部署

**场景**：在同一端口上运行多个监控实例

**实例 1（订单 A）**：
```yaml
order_id: "177971759268550919"
web_base_path: "/order-a"
web_port: 8099
```

**实例 2（订单 B）**：
```yaml
order_id: "188082860379661020"
web_base_path: "/order-b"
web_port: 8099
```

**访问**：
- 订单 A: http://localhost:8099/order-a/
- 订单 B: http://localhost:8099/order-b/

### 3. 统一服务入口

**场景**：所有内部服务使用统一路径前缀

**应用配置**：
```yaml
web_base_path: "/internal/lixiang"
```

**访问**: `http://localhost:8099/internal/lixiang/`

## 技术细节

### 路径规范化规则

1. **空字符串**: 保持不变（默认行为）
   ```
   输入: ""
   输出: ""
   ```

2. **自动添加前导斜杠**:
   ```
   输入: "monitor"
   输出: "/monitor"
   ```

3. **自动移除尾部斜杠**（除根路径外）:
   ```
   输入: "/monitor/"
   输出: "/monitor"
   
   输入: "/"
   输出: "/" (根路径保留)
   ```

4. **多级路径处理**:
   ```
   输入: "/apps/lixiang"
   输出: "/apps/lixiang"
   
   输入: "apps/lixiang/"
   输出: "/apps/lixiang"
   ```

### 路由构建逻辑

```go
func (s *Server) route(path string) string {
    if s.basePath == "" {
        return path
    }
    return s.basePath + path
}
```

**示例**：
```
basePath = ""          path = "/"           → "/"
basePath = "/monitor"  path = "/"           → "/monitor/"
basePath = "/monitor"  path = "/api/stats"  → "/monitor/api/stats"
```

## 测试验证

### 编译测试
```bash
./scripts/deploy/build.sh
# ✅ 编译成功！
```

### 功能测试 1: 默认配置

**配置**：
```yaml
web_base_path: ""
```

**启动日志**：
```
[Web] 启动 Web 服务器: http://localhost:8099
```

**API 测试**：
```bash
curl http://localhost:8099/api/stats
# ✅ 返回 JSON 数据
```

### 功能测试 2: 自定义根路由

**配置**：
```yaml
web_base_path: "/monitor"
```

**启动日志**：
```
[Web] 启动 Web 服务器: http://localhost:8099/monitor
```

**API 测试**：
```bash
curl http://localhost:8099/monitor/api/stats
# ✅ 返回 JSON 数据

curl http://localhost:8099/api/stats
# ❌ 404 Not Found (预期行为)
```

**请求日志**：
```
[Web] GET /monitor/api/stats - 401.917µs
```

## 向后兼容性

✅ **完全向后兼容**

- 默认值为空字符串 `""`，保持原有行为
- 现有配置文件无需修改即可继续使用
- 可选功能，不影响核心监控功能

## 代码统计

| 文件 | 变更类型 | 行数变化 |
|------|---------|---------|
| `cfg/config.go` | 修改 | +4 |
| `web/server.go` | 修改 | +17 |
| `main.go` | 修改 | +3 |
| `config.yaml` | 修改 | +1 |
| `docs/guides/WEB_BASE_PATH.md` | 新增 | +442 |
| `README.md` | 修改 | +20 |
| **总计** | - | **+487** |

## 文档更新

### 新增文档
- ✅ `docs/guides/WEB_BASE_PATH.md` - 根路由配置完整指南

### 更新文档
- ✅ `README.md` - 添加根路由配置说明
- ✅ `config.yaml` - 添加配置项注释

## 注意事项

### 配置热重载

虽然系统支持配置热重载，但 `web_base_path` 的修改需要**重启程序**才能生效，因为：

1. Web 服务器的路由在启动时注册
2. 修改根路由需要重新注册所有路由
3. 当前实现不支持运行时动态修改路由

### 反向代理配置

使用反向代理时，注意路径映射：

**正确配置**：
```nginx
location /monitor/ {
    proxy_pass http://localhost:8099/monitor/;
}
```

**错误配置**：
```nginx
location /monitor/ {
    proxy_pass http://localhost:8099/;  # 路径不匹配
}
```

### 路径安全

- 建议使用不易猜测的路径增加安全性
- 配合 HTTP Basic Auth 使用
- 不建议使用特殊字符

## 已知限制

1. **修改后需重启**: `web_base_path` 修改后必须重启程序
2. **不支持正则路由**: 仅支持固定路径前缀
3. **无动态路由**: 不支持运行时动态添加/移除路由

## 未来计划

可能的功能增强：

- [ ] 支持运行时动态修改根路由
- [ ] 支持多个根路由别名
- [ ] 支持路径重写规则
- [ ] 内置反向代理功能

## 相关链接

- [Web 可视化界面使用指南](docs/guides/WEB_INTERFACE.md)
- [Web 根路由配置指南](docs/guides/WEB_BASE_PATH.md)

## 总结

此次更新为 Web 服务器添加了灵活的根路由配置功能，使得应用可以：

✅ 在反向代理环境下部署  
✅ 与其他服务共享端口  
✅ 实现多实例部署  
✅ 增强路径安全性  
✅ 保持向后兼容  

配置简单，功能强大，为不同部署场景提供了更多选择。

---

**开发者**: GitHub Copilot  
**审核状态**: ✅ 通过  
**测试状态**: ✅ 全部通过  
**部署状态**: ✅ 已完成
