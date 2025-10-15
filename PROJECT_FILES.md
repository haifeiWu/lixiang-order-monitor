# 项目文件说明

## 核心文件

| 文件 | 说明 |
|------|------|
| `main.go` | 主程序文件，包含监控逻辑和通知功能 |
| `go.mod` | Go 模块依赖文件 |
| `config.yaml` | 配置文件（需要用户自行配置） |
| `config.example.yaml` | 配置文件示例 |

## 可执行脚本

| 脚本 | 说明 |
|------|------|
| `build.sh` | 编译程序 |
| `start.sh` | 启动监控服务（后台运行） |
| `stop.sh` | 停止监控服务 |
| `status.sh` | 查看运行状态和配置信息 |
| `test-notification.sh` | 测试通知功能是否正常 |

## 文档文件

| 文件 | 说明 |
|------|------|
| `README.md` | 主要说明文档 |
| `WECHAT_SETUP.md` | 微信群机器人配置指南 |
| `SERVERCHAN_SETUP.md` | ServerChan 配置指南 |
| `PROJECT_FILES.md` | 本文件，项目文件说明 |

## 自动生成文件

| 文件 | 说明 |
|------|------|
| `lixiang-monitor` | 编译后的可执行文件 |
| `monitor.log` | 运行日志文件 |
| `go.sum` | Go 模块校验文件 |

## 配置优先级

1. **必填配置**：
   - `order_id`: 理想汽车订单 ID
   - `lixiang_cookies`: 理想汽车网站的 Cookie

2. **通知配置**（至少配置一种）：
   - `wechat_webhook_url`: 微信群机器人 Webhook URL
   - `serverchan_sendkey`: ServerChan 的 SendKey

3. **可选配置**：
   - `check_interval`: 检查间隔（默认 30 分钟）
   - `serverchan_baseurl`: ServerChan API 地址（有默认值）

## 使用流程

1. **初始设置**：
   ```bash
   # 1. 复制配置文件
   cp config.example.yaml config.yaml
   
   # 2. 编辑配置文件
   vim config.yaml
   
   # 3. 编译程序
   ./build.sh
   ```

2. **测试配置**：
   ```bash
   # 4. 测试通知功能
   ./test-notification.sh
   ```

3. **启动监控**：
   ```bash
   # 5. 启动服务
   ./start.sh
   
   # 6. 查看状态
   ./status.sh
   
   # 7. 查看日志
   tail -f monitor.log
   ```

4. **管理服务**：
   ```bash
   # 停止服务
   ./stop.sh
   
   # 重启服务
   ./stop.sh && ./start.sh
   ```

## 故障排除

1. **编译失败**：检查 Go 环境和依赖
2. **通知不工作**：使用 `./test-notification.sh` 测试
3. **获取订单失败**：检查 Cookie 是否过期
4. **程序无法启动**：查看 `monitor.log` 日志

## 安全注意事项

- `config.yaml` 包含敏感信息，已在 `.gitignore` 中排除
- 定期更新 Cookie（通常几天到几周会过期）
- 妥善保管 ServerChan SendKey 和微信 Webhook URL