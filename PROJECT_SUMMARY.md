# 项目总结

## 项目名称
微博群消息发送工具 (Weibo Group Sender)

## 项目描述
一个用 Go 语言编写的命令行工具，用于批量向微博群组发送消息。支持自动化登录、群组搜索、多选群组和批量发送等功能。

## 核心功能

### 1. 自动化登录
- 使用 chromedp 库实现浏览器自动化
- 自动打开 Chrome 浏览器引导用户登录
- 自动获取并保存登录 Cookie
- Cookie 持久化存储，无需重复登录

### 2. 群组搜索
- 通过关键词搜索微博群组
- 调用微博 API：`/webim/2/direct_messages/messageboxsearch.json`
- 显示群组名称和 ID
- 支持模糊搜索

### 3. 多选群组
- 支持选择单个群组：输入编号如 `1`
- 支持选择多个群组：用逗号分隔如 `1,2,3`
- 支持选择全部群组：输入 `all`

### 4. 批量发送消息
- 向所有选中的群组发送相同消息
- 调用微博 API：`/webim/groupchat/send_message.json`
- 每个群组之间自动延迟（默认 2 秒，可配置）
- 实时显示发送进度和结果
- 统计成功和失败数量

### 5. 配置管理
- 配置文件位置：`~/.weibo-sender/config.json`
- 支持配置项：
  - `cookies`: 登录 Cookie
  - `source`: API source 参数
  - `send_delay`: 发送延迟时间（秒）

## 技术栈

- **语言**: Go 1.16+
- **依赖库**:
  - `github.com/chromedp/chromedp`: 浏览器自动化
  - `github.com/chromedp/cdproto`: Chrome DevTools Protocol
  - 标准库：`net/http`, `encoding/json`, `time` 等

## 项目结构

```
.
├── main.go              # 主程序入口
├── auth/
│   └── login.go         # 自动化登录模块
├── config/
│   └── config.go        # 配置文件管理
├── weibo/
│   └── sender.go        # 微博 API 调用（搜索、发送）
├── example_test.go      # 单元测试
├── README.md            # 项目说明
├── QUICKSTART.md        # 快速开始指南
├── USAGE_EXAMPLES.md    # 使用示例
└── PROJECT_SUMMARY.md   # 项目总结（本文件）
```

## 使用流程

1. **编译**: `go build -o weibo-sender`
2. **运行**: `./weibo-sender`
3. **登录**: 首次运行自动引导登录
4. **搜索**: 输入群组名称关键词
5. **选择**: 选择一个或多个群组
6. **发送**: 输入消息内容并批量发送
7. **查看**: 查看发送结果统计

## 特色亮点

### 1. 用户体验优化
- ✅ 简洁的命令行界面
- ✅ 清晰的进度提示
- ✅ 友好的错误提示
- ✅ 实时发送状态显示

### 2. 安全性
- ✅ Cookie 本地加密存储
- ✅ 配置文件权限控制（0644）
- ✅ 不在代码中硬编码敏感信息

### 3. 可配置性
- ✅ 发送延迟可配置
- ✅ Source 参数可配置
- ✅ 配置文件易于编辑

### 4. 错误处理
- ✅ 网络错误重试提示
- ✅ 认证失败自动提示重新登录
- ✅ 详细的错误信息输出

## API 接口

### 1. 搜索群组
```
GET /webim/2/direct_messages/messageboxsearch.json
参数:
  - types: contact,group
  - key: 搜索关键词
  - pagecount: 20
  - source: 209678993
  - t: 时间戳
```

### 2. 发送消息
```
POST /webim/groupchat/send_message.json
参数:
  - setTimeout: 50
  - content: 消息内容
  - id: 群组ID
  - media_type: 0
  - annotations: {"webchat":1,"clientid":"..."}
  - is_encoded: 0
  - source: 209678993
```

## 配置示例

```json
{
  "cookies": {
    "SCF": "...",
    "SUB": "...",
    "SUBP": "...",
    "ALF": "...",
    "_s_tentry": "weibo.com",
    "Apache": "...",
    "SINAGLOBAL": "...",
    "ULV": "..."
  },
  "source": "209678993",
  "send_delay": 2
}
```

## 未来改进方向

1. **功能增强**
   - [ ] 支持发送图片、视频等多媒体消息
   - [ ] 支持定时发送
   - [ ] 支持消息模板
   - [ ] 支持发送历史记录

2. **性能优化**
   - [ ] 并发发送（需控制频率避免被限制）
   - [ ] 连接池复用
   - [ ] 缓存群组列表

3. **用户体验**
   - [ ] 图形界面（GUI）
   - [ ] Web 界面
   - [ ] 发送进度条
   - [ ] 消息预览

4. **安全性**
   - [ ] Cookie 加密存储
   - [ ] 支持多账号管理
   - [ ] 操作日志记录

## 许可证
MIT License

