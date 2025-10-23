# Web 服务器根路由配置指南

## 功能概述

从 v1.9.0 开始，Web 服务器支持配置自定义根路由（Base Path），允许您在特定路径前缀下运行 Web 界面。这在以下场景非常有用：

- 🔄 **反向代理部署**：与其他服务共享端口
- 🌐 **多服务集成**：统一入口管理多个应用
- 🛡️ **路径隔离**：避免路由冲突
- 📦 **子目录部署**：在网站的子目录下运行

## 配置方法

### 1. 在 config.yaml 中配置

```yaml
# Web 管理界面配置
web_enabled: true
web_port: 8099
web_base_path: ""  # 根路由配置
```

### 2. 配置示例

#### 示例 1: 默认配置（无根路由）

```yaml
web_base_path: ""
```

**访问地址**：
- 主页：`http://localhost:8099/`
- API 统计：`http://localhost:8099/api/stats`
- 历史记录：`http://localhost:8099/api/records`
- 时间变更：`http://localhost:8099/api/time-changes`

#### 示例 2: 使用 /monitor 作为根路由

```yaml
web_base_path: "/monitor"
```

**访问地址**：
- 主页：`http://localhost:8099/monitor/`
- API 统计：`http://localhost:8099/monitor/api/stats`
- 历史记录：`http://localhost:8099/monitor/api/records`
- 时间变更：`http://localhost:8099/monitor/api/time-changes`

#### 示例 3: 多级路径

```yaml
web_base_path: "/apps/lixiang"
```

**访问地址**：
- 主页：`http://localhost:8099/apps/lixiang/`
- API 统计：`http://localhost:8099/apps/lixiang/api/stats`
- 历史记录：`http://localhost:8099/apps/lixiang/api/records`
- 时间变更：`http://localhost:8099/apps/lixiang/api/time-changes`

## 路径规范化

系统会自动规范化您配置的 `web_base_path`：

1. **自动添加前导斜杠**：
   - 配置：`monitor` → 实际：`/monitor`
   - 配置：`/monitor` → 实际：`/monitor`

2. **自动移除尾部斜杠**：
   - 配置：`/monitor/` → 实际：`/monitor`
   - 配置：`/monitor` → 实际：`/monitor`

3. **保留空字符串**：
   - 配置：`""` → 实际：`""`（默认行为，无根路由）

## 使用场景

### 场景 1: Nginx 反向代理

如果您使用 Nginx 作为反向代理，可以将监控服务部署在子路径下：

**Nginx 配置**：
```nginx
server {
    listen 80;
    server_name example.com;

    # 理想汽车监控服务
    location /lixiang-monitor/ {
        proxy_pass http://localhost:8099/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # 其他服务
    location /another-app/ {
        proxy_pass http://localhost:8100/;
    }
}
```

**应用配置**：
```yaml
web_base_path: ""
web_port: 8099
```

访问：`http://example.com/lixiang-monitor/`

### 场景 2: 统一服务入口

如果您希望所有内部服务都使用统一的路径前缀：

**应用配置**：
```yaml
web_base_path: "/internal/lixiang"
web_port: 8099
```

**Nginx 配置**：
```nginx
location /internal/ {
    proxy_pass http://localhost:8099/internal/;
}
```

访问：`http://example.com/internal/lixiang/`

### 场景 3: 多实例部署

如果您需要监控多个订单，可以运行多个实例：

**实例 1（订单 A）**：
```yaml
order_id: "177971759268550919"
web_port: 8099
web_base_path: "/order-a"
```

**实例 2（订单 B）**：
```yaml
order_id: "188082860379661020"
web_port: 8099
web_base_path: "/order-b"
```

访问：
- 订单 A：`http://localhost:8099/order-a/`
- 订单 B：`http://localhost:8099/order-b/`

## 启动日志

配置根路由后，启动日志会显示完整的访问地址：

### 无根路由
```
[Web] 启动 Web 服务器: http://localhost:8099
```

### 有根路由
```
[Web] 启动 Web 服务器: http://localhost:8099/monitor
```

## 注意事项

### ✅ 支持的配置

- ✅ 空字符串（默认，无根路由）
- ✅ 单级路径：`/monitor`
- ✅ 多级路径：`/apps/lixiang`
- ✅ 任意合法路径字符

### ⚠️ 不支持的配置

- ❌ 不包含 `/` 的路径会自动添加前导 `/`
- ❌ 末尾的 `/` 会被自动移除（除非是根路径 `/`）
- ❌ 特殊字符（如空格、`?`、`#` 等）不建议使用

### 🔒 安全建议

1. **路径保密性**：使用不易猜测的路径增加安全性
   ```yaml
   web_base_path: "/lx-mon-a8f3e2d1"
   ```

2. **反向代理认证**：配合 Nginx 的 HTTP Basic Auth
   ```nginx
   location /monitor/ {
       auth_basic "Restricted Access";
       auth_basic_user_file /etc/nginx/.htpasswd;
       proxy_pass http://localhost:8099/monitor/;
   }
   ```

3. **防火墙规则**：限制直接访问内部端口
   ```bash
   # 仅允许本地访问 8099 端口
   sudo ufw deny 8099
   sudo ufw allow from 127.0.0.1 to any port 8099
   ```

## 热重载支持

`web_base_path` 配置支持热重载：

1. 修改 `config.yaml` 中的 `web_base_path`
2. 保存文件
3. 系统自动检测配置变化
4. **需要重启程序**才能生效（Web 服务器需要重新初始化路由）

**注意**：虽然配置热重载功能会检测到变化，但 Web 服务器的路由在启动时已经注册，修改 `web_base_path` 后必须重启程序才能生效。

## 测试验证

### 1. 启动程序
```bash
./lixiang-monitor
```

### 2. 检查日志
查看启动日志中的 Web 服务器地址：
```
[Web] 启动 Web 服务器: http://localhost:8099/monitor
```

### 3. 测试 API
```bash
# 默认配置
curl http://localhost:8099/api/stats

# 使用根路由
curl http://localhost:8099/monitor/api/stats
```

### 4. 浏览器访问
```
# 默认配置
http://localhost:8099/

# 使用根路由
http://localhost:8099/monitor/
```

## 故障排除

### 问题 1: 404 Not Found

**症状**：访问页面返回 404
```bash
curl http://localhost:8099/monitor/
# 404 page not found
```

**原因**：
1. 配置的 `web_base_path` 与访问路径不匹配
2. 修改配置后未重启程序

**解决方法**：
1. 检查 `config.yaml` 中的 `web_base_path` 配置
2. 查看启动日志中显示的实际访问地址
3. 重启程序使配置生效

### 问题 2: 路径重复

**症状**：URL 中路径重复
```
http://localhost:8099/monitor/monitor/api/stats
```

**原因**：配置和访问路径都包含了相同的前缀

**解决方法**：
- 如果使用反向代理，确保 `proxy_pass` 配置正确
- 如果直接访问，使用正确的完整路径

### 问题 3: 反向代理 404

**症状**：通过 Nginx 访问返回 404

**Nginx 配置问题**：
```nginx
# ❌ 错误配置
location /monitor {
    proxy_pass http://localhost:8099/;
}
```

**正确配置**：
```nginx
# ✅ 正确配置（保留尾部斜杠）
location /monitor/ {
    proxy_pass http://localhost:8099/monitor/;
}
```

或者：
```nginx
# ✅ 使用 rewrite
location /monitor {
    rewrite ^/monitor/(.*)$ /$1 break;
    proxy_pass http://localhost:8099;
}
```

## 最佳实践

1. **开发环境**：使用默认配置（空 basePath）
   ```yaml
   web_base_path: ""
   ```

2. **生产环境**：使用有意义的路径
   ```yaml
   web_base_path: "/lixiang-monitor"
   ```

3. **多实例**：使用订单 ID 作为路径
   ```yaml
   web_base_path: "/order-177971759268550919"
   ```

4. **反向代理**：保持应用和代理配置一致
   ```yaml
   # 应用配置
   web_base_path: "/monitor"
   ```
   ```nginx
   # Nginx 配置
   location /monitor/ {
       proxy_pass http://localhost:8099/monitor/;
   }
   ```

## 代码示例

### Go 代码中的实现

**路径规范化**：
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

**路由构建**：
```go
func (s *Server) route(path string) string {
    if s.basePath == "" {
        return path
    }
    return s.basePath + path
}
```

**路由注册**：
```go
mux.HandleFunc(s.route("/"), s.handleIndex)
mux.HandleFunc(s.route("/api/stats"), s.handleStats)
mux.HandleFunc(s.route("/api/records"), s.handleRecords)
mux.HandleFunc(s.route("/api/time-changes"), s.handleTimeChanges)
```

## 更新日志

### v1.9.0 (2025-10-23)
- ✨ 新增 `web_base_path` 配置项
- ✨ 支持自定义根路由
- ✨ 自动路径规范化
- ✨ 启动日志显示完整访问地址
- 📝 新增配置文档

## 相关文档

- [Web 可视化界面使用指南](./WEB_INTERFACE.md)
- [反向代理部署指南](../technical/NGINX_PROXY.md)（待创建）
- [多实例部署指南](../technical/MULTI_INSTANCE.md)（待创建）

## 技术支持

如果您在使用过程中遇到问题，请：

1. 检查启动日志确认实际访问地址
2. 使用 curl 测试 API 是否正常
3. 检查反向代理配置（如果使用）
4. 查看本文档的"故障排除"部分

---

**功能版本**: v1.9.0  
**更新时间**: 2025-10-23  
**维护状态**: ✅ 活跃维护
