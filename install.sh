#!/bin/bash
set -e

echo "========================================"
echo "  WireGuard伴侣 - 一键安装"
echo "========================================"

# 检测系统
detect_os() {
    if [[ -f /etc/os-release ]]; then
        . /etc/os-release
        echo "检测到系统: $ID"
    fi
}

# 检测架构
detect_arch() {
    ARCH=$(uname -m)
    echo "架构: $ARCH"
}

# 检测Docker
check_docker() {
    if command -v docker &> /dev/null; then
        echo "Docker: 已安装"
        return 0
    else
        echo "Docker: 未安装，开始安装..."
        install_docker
        return 1
    fi
}

# 安装Docker
install_docker() {
    if command -v apt-get &> /dev/null; then
        apt-get update
        apt-get install -y ca-certificates curl gnupg
        install -m 0755 -d /etc/apt/keyrings
        curl -fsSL https://download.docker.com/linux/$ID/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
        chmod a+r /etc/apt/keyrings/docker.gpg
        echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/$ID $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
        apt-get update
        apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
    elif command -v yum &> /dev/null; then
        yum install -y yum-utils
        yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
        yum install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
    fi
    systemctl start docker || service docker start
    systemctl enable docker || true
    echo "Docker 安装完成"
}

# 安装wg-easy
install_wg_easy() {
    echo "安装 wg-easy..."
    if docker ps -a | grep -q wg-easy; then
        echo "wg-easy 已安装"
        docker start wg-easy 2>/dev/null || true
        return
    fi
    docker network create wg-network 2>/dev/null || true
    docker run -d \
        --name wg-easy \
        -e WG_HOST=$(curl -s ifconfig.me) \
        -e PASSWORD=admin123 \
        -p 51820:51820/udp \
        -p 51821:51821/tcp \
        -v /etc/wireguard:/etc/wireguard \
        --restart=unless-stopped \
        --network wg-network \
        weejewel/wg-easy:latest
    echo "wg-easy 安装完成"
}

# 安装ddns-go
install_ddns_go() {
    echo "安装 ddns-go..."
    if pgrep ddns-go &> /dev/null; then
        echo "ddns-go 已安装"
        return
    fi
    
    ARCH=$(uname -m)
    case $ARCH in
        x86_64) ARCH="x86_64" ;;
        aarch64) ARCH="arm64" ;;
    esac
    
    VERSION=$(curl -s https://api.github.com/repos/jeessy2/ddns-go/releases/latest | grep -o '"tag_name": "v[^"]*' | cut -d'"' -f4)
    URL="https://github.com/jeessy2/ddns-go/releases/download${VERSION}/ddns-go_${VERSION}_linux_${ARCH}.tar.gz"
    
    curl -L -o /tmp/ddns-go.tar.gz "$URL"
    tar -xzf /tmp/ddns-go.tar.gz -C /tmp
    mv /tmp/ddns-go/ddns-go /usr/local/bin/
    chmod +x /usr/local/bin/ddns-go
    ddns-go -s install
    echo "ddns-go 安装完成"
}

# 安装wstunnel
install_wstunnel() {
    echo "安装 wstunnel..."
    if command -v wstunnel &> /dev/null; then
        echo "wstunnel 已安装"
        return
    fi
    
    ARCH=$(uname -m)
    VERSION=$(curl -s https://api.github.com/repos/erebe/wstunnel/releases/latest | grep -o '"tag_name": "v[^"]*' | cut -d'"' -f4)
    
    if [[ "$ARCH" == "x86_64" ]]; then
        ARCH="x86_64"
    elif [[ "$ARCH" == "aarch64" ]]; then
        ARCH="arm64"
    fi
    
    curl -L -o /usr/local/bin/wstunnel "https://github.com/erebe/wstunnel/releases/download${VERSION}/wstunnel-${VERSION}-unknown-linux-${ARCH}"
    chmod +x /usr/local/bin/wstunnel
    echo "wstunnel 安装完成"
}

# 主菜单
main() {
    detect_os
    detect_arch
    check_docker
    
    echo ""
    echo "========================================"
    echo "请选择安装:"
    echo "========================================"
    echo "1. 安装 WireGuard (wg-easy)"
    echo "2. 安装 DDNS (ddns-go)"
    echo "3. 安装 协议封装 (wstunnel)"
    echo "4. 全部安装"
    echo "5. 退出"
    echo ""
    
    read -p "请选择 (1-5): " choice
    
    case $choice in
        1) install_wg_easy ;;
        2) install_ddns_go ;;
        3) install_wstunnel ;;
        4) 
            install_wg_easy
            install_ddns_go
            install_wstunnel
            ;;
        5) exit 0 ;;
    esac
    
    echo ""
    echo "========================================"
    echo "安装完成!"
    echo "========================================"
    echo "wg-easy: http://<ip>:51821 (账号: admin, 密码: admin123)"
    echo "ddns-go: http://<ip>:9876"
    echo "========================================"
}

main "$@"
