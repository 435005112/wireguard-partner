# WireGuard 伴侣

统一管理 WireGuard + 协议转换工具

## 功能

### 协议转换模块（Web UI）
- **wstunnel** - WebSocket 隧道，支持客户端/服务端模式
- **udp2raw** - UDP 转加密 TCP/ICMP 流量

**特性：**
- 完整的参数配置界面（必填 + 高级选项）
- 每个参数配有 ? 帮助说明
- 协议/角色切换时表单自动联动
- 一键生成 systemd service 并管理
- 实时日志查看

### 旧版功能（CLI）
- WireGuard VPN 隧道管理
- DDNS 动态域名解析 (DDNS-GO)
- 版本检测与自动更新

## 快速开始

### 协议转换模块（推荐）

```bash
# 编译
go build -o wg-partner cmd/protocol.go

# 运行
./wg-partner -port 8080
```

访问 http://localhost:8080 打开 Web 界面

### 旧版 CLI

```bash
# 编译
go build -o wg-partner cmd/main.go

# 运行
./wg-partner -c config.yaml
```

## 协议转换模块使用

1. 打开 Web 界面
2. 点击"新建配置"
3. 选择协议类型（wstunnel / udp2raw）
4. 选择角色（客户端 / 服务端）
5. 填写必填参数，高级选项可展开配置
6. 保存后可启动/停止服务
7. 查看日志监控运行状态

## 配置示例

### wstunnel 客户端
```
服务器: tunnel.example.com
端口: 443
本地隧道: tcp://51820:localhost:51820
使用TLS: ✓
```

### udp2raw 客户端
```
服务器: 1.2.3.4
端口: 443
密钥: your-secret-key
模式: faketcp
自动iptables: ✓
```

## 目录结构

```
.
├── cmd/
│   ├── main.go       # 主程序入口
│   └── protocol.go   # 协议转换模块
├── internal/
│   ├── protocol/     # 协议管理核心逻辑
│   ├── tunnel/       # WireGuard 隧道管理
│   ├── monitor/      # 状态监控
│   └── ...
└── web/
    ├── protocol.go   # Web 服务器
    └── templates/    # HTML 模板
```

## License

MIT
