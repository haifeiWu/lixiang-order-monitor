# 🍪 Cookie 失效快速处理指南

## 🚨 收到 Cookie 失效告警？

如果您收到以下通知：

```
🚨 理想汽车 Cookie 已失效
您的理想汽车订单监控 Cookie 已失效，请及时更新！
```

**不要慌张！** 按照以下步骤快速解决：

---

## ⚡ 5 分钟快速修复

### 步骤 1: 打开理想汽车官网

在浏览器中访问：https://www.lixiang.com/

### 步骤 2: 登录账号

使用您的理想汽车账号登录（手机号 + 验证码）

### 步骤 3: 获取 Cookie

**方法 A - Chrome/Edge** (推荐)
1. 按 `F12` 打开开发者工具
2. 点击顶部的 `Network` (网络) 标签
3. 按 `F5` 刷新页面
4. 在左侧请求列表中点击任意 API 请求
5. 在右侧找到 `Request Headers` 区域
6. 找到 `Cookie:` 那一行
7. 点击并全选 Cookie 值（很长的字符串）
8. 右键 → 复制

**方法 B - Safari**
1. `Cmd+Option+I` 打开开发者工具
2. 点击 `Network` (网络)
3. 刷新页面
4. 选择任意请求
5. 查看 `Headers` → `Request` → `Cookie`
6. 复制完整 Cookie 字符串

**方法 C - Firefox**
1. `F12` 打开开发者工具
2. 点击 `网络` 标签
3. 刷新页面
4. 点击任意请求
5. 右侧 `标头` → `请求标头` → `Cookie`
6. 复制完整值

### 步骤 4: 更新配置文件

打开 `config.yaml` 文件，找到这一行：

```yaml
lixiang_cookies: "旧的 Cookie"
```

替换为刚才复制的新 Cookie：

```yaml
lixiang_cookies: "X-LX-Deviceid=xxx; X-LX-Token=xxx; authli_device_id=xxx; ..."
```

**注意事项：**
- ✅ 保留双引号
- ✅ 确保是完整的一行（可能很长）
- ✅ 不要有换行符
- ✅ Cookie 之间用 `; ` 分隔（分号+空格）

### 步骤 5: 保存并验证

1. **保存文件** (`Cmd+S` / `Ctrl+S`)

2. **等待 1-2 秒**（程序会自动检测变化）

3. **查看日志确认**：
   ```
   2025/01/17 15:30:00 📝 配置已更新 (版本: 2)
   2025/01/17 15:30:00 ✅ Cookie 验证成功
   ```

4. **完成！** 🎉

---

## 🔍 验证 Cookie 是否有效

### 方法 1: 查看程序日志

监控程序运行日志，查找：

```bash
# 成功标志
✅ Cookie 验证成功
当前预计交付时间: 2025-12-01

# 失败标志
⚠️  Cookie 验证失败 (状态码: 401, 连续失败: 1 次)
```

### 方法 2: 手动测试

运行测试脚本：

```bash
./test-cookie-expiry.sh
```

选择选项 `1` 进行 Cookie 失效检测测试。

---

## 🤔 常见问题

### Q1: 复制的 Cookie 很长，是正常的吗？

**A:** 是的！完整的 Cookie 通常有 **500-1000 个字符**，包含多个字段，这是正常的。

示例长度：
```
X-LX-Deviceid=xxx; X-CHJ-ChannelCode=xxx; share-uid=xxx; X-CHJ-SourceUrl=xxx; authli_device_id=xxx; authli_device_reported=xxx; X-LX-HeaderData=xxx; X-LX-PreviewToken=xxx; X-LX-Token=xxx
```

### Q2: 需要复制哪些字段？

**A:** 最简单的方法是 **复制完整的 Cookie 字符串**。如果要手动筛选，至少需要：

- ✅ `X-LX-Token` (必需，最重要)
- ✅ `authli_device_id` (必需)
- ✅ `X-LX-Deviceid` (必需)
- ⚪ 其他字段（建议全部保留）

### Q3: 更新后还是提示失效怎么办？

**排查步骤：**

1. **确认格式正确**
   ```yaml
   # ✅ 正确
   lixiang_cookies: "X-LX-Token=xxx; authli_device_id=xxx"
   
   # ❌ 错误（缺少引号）
   lixiang_cookies: X-LX-Token=xxx
   
   # ❌ 错误（换行）
   lixiang_cookies: "X-LX-Token=xxx
   authli_device_id=xxx"
   ```

2. **确认是否包含 X-LX-Token**
   ```bash
   grep "X-LX-Token" config.yaml
   ```

3. **确认配置已重新加载**
   - 查看日志是否有 "配置已更新" 消息

4. **重新登录并获取新 Cookie**
   - 退出理想汽车网站
   - 清除浏览器 Cookie
   - 重新登录获取

### Q4: Cookie 多久会失效？

**A:** 理想汽车的 Cookie 有效期取决于：

- **登录时选择"记住我"**: 可能持续 **7-30 天**
- **临时会话**: 可能只有 **1-7 天**
- **安全策略**: 异常登录可能提前失效

**建议：**
- 登录时勾选"记住我"
- 每周检查一次 Cookie 有效性
- 设置监控告警（本系统已自动配置）

### Q5: 能否自动更新 Cookie？

**A:** 由于安全限制，无法完全自动化登录和 Cookie 获取。但本系统提供：

- ✅ **自动检测失效**（3 次失败后告警）
- ✅ **配置热加载**（更新后无需重启）
- ✅ **详细告警指引**（通知包含完整更新步骤）

这是最佳的平衡方案，既保证安全性又最小化手动操作。

---

## 📱 移动端获取 Cookie

如果您没有电脑，可以在手机上获取：

### iOS (Safari)

1. 安装 **Web Inspector** App
2. 在 Safari 中打开理想汽车网站
3. 使用 Web Inspector 查看请求头
4. 复制 Cookie 字段

### Android (Chrome)

1. 在 Chrome 中打开 `chrome://inspect`
2. 访问理想汽车网站
3. 查看网络请求
4. 复制 Cookie

**更简单的方法：** 使用电脑浏览器获取，然后通过微信/邮件发送到手机更新配置。

---

## 🛡️ 安全提醒

### ⚠️  重要警告

1. **不要分享您的 Cookie**
   - Cookie 等同于您的登录凭证
   - 拥有 Cookie 的人可以访问您的账号

2. **保护配置文件**
   ```bash
   # 确保 config.yaml 不被上传到 Git
   echo "config.yaml" >> .gitignore
   ```

3. **定期更新密码**
   - 建议每月更新理想汽车账号密码
   - 密码更新后需要重新获取 Cookie

---

## 📞 需要帮助？

如果按照上述步骤仍无法解决，请：

1. **查看完整文档**
   ```bash
   cat COOKIE_MANAGEMENT.md
   ```

2. **运行诊断测试**
   ```bash
   ./test-cookie-expiry.sh
   ```

3. **查看程序日志**
   ```bash
   tail -f lixiang-monitor.log
   ```

4. **提交 Issue**
   - 描述问题现象
   - 附上日志（⚠️ 删除敏感信息）

---

## ✅ 成功案例

```
🎉 Cookie 更新成功！

日志输出：
2025/01/17 15:30:00 📝 配置已更新 (版本: 2)
2025/01/17 15:30:00 开始检查订单交付时间...
2025/01/17 15:30:00 当前预计交付时间: 2025-12-01
2025/01/17 15:30:00 ✅ Cookie 验证成功
2025/01/17 15:30:00 预计交付时间未更新

监控服务已恢复正常运行！
```

---

**祝您的理想汽车早日交付！** 🚗💨
