# WireGuard 伴侣

统一管理 WireGuard + DDNS + 协议封装工具

## 功能

- WireGuard VPN 隧道管理
- DDNS 动态域名解析 (DDNS-GO)
- 协议封装 (wstunnel, udp2raw)
- 版本检测与自动更新

## 快速开始

```bash
# 编译
go build -o wg-partner cmd/main.go

# 运行
./wg-partner -c config.yaml
```

## 配置

见 `configs/config.example.yaml`
