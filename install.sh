#!/bin/bash
set -e
echo "WireGuard伴侣 安装脚本"
detect_os() { OS=$(cat /etc/os-release 2>/dev/null | grep "^ID=" | cut -d= -f2 | tr -d '"'); echo "系统: $OS"; }
detect_arch() { ARCH=$(uname -m); [[ "$ARCH" == "x86_64" ]] && ARCH="amd64"; [[ "$ARCH" == "aarch64" ]] && ARCH="arm64"; echo "架构: $ARCH"; }
check_github() { curl -s --connect-timeout 3 -o /dev/null https://github.com && GITHUB_OK=true || GITHUB_OK=false; }
install_wg() { echo "安装WireGuard..."; apt update && apt install -y wireguard-tools 2>/dev/null || yum install -y wireguard-tools 2>/dev/null || echo "请手动安装"; }
install_wstunnel() { echo "安装wstunnel..."; curl -L -o /usr/local/bin/wstunnel "https://github.com/erebe/wstunnel/releases/latest/download/wstunnel-$(curl -s https://api.github.com/repos/erebe/wstunnel/releases/latest | grep tag_name | cut -d'"' -f4)-unknown-linux-${ARCH}" && chmod +x /usr/local/bin/wstunnel; }
install_udp2raw() { echo "安装udp2raw..."; curl -L -o /usr/local/bin/udp2raw "https://github.com/wangyu-/udp2raw/releases/latest/download/udp2raw_amd64" && chmod +x /usr/local/bin/udp2raw; }
detect_os; detect_arch; check_github; echo "1.WireGuard 2.wstunnel 3.udp2raw 4.全部 5.退出"; read -p "选择: " c; case $c in 1) install_wg ;; 2) install_wstunnel ;; 3) install_udp2raw ;; 4) install_wg && install_wstunnel && install_udp2raw ;; esac; echo "完成! go build -o wg-partner cmd/main.go && ./wg-partner"; echo "Web: http://<ip>:8080 账号:admin 密码:admin123"
