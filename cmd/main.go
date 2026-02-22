package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"wireguard-partner/internal/installer"
	"wireguard-partner/web"
)

func main() {
	isInstall := flag.Bool("install", false, "运行安装向导")
	flag.Parse()

	fmt.Println("WireGuard 伴侣")
	fmt.Println("================")

	// 如果指定了-install参数，运行安装向导
	if *isInstall {
		fmt.Println()
		fmt.Println("正在检测环境并安装依赖...")
		fmt.Println()
		if err := installer.QuickInstall(); err != nil {
			fmt.Printf("安装失败: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// 检查WireGuard是否安装
	if err := checkWireGuard(); err != nil {
		fmt.Printf("提示: %v\n", err)
		fmt.Println()
		fmt.Println("运行安装向导:")
		fmt.Println("  ./wg-partner -install")
		return
	}

	fmt.Println("✅ WireGuard 已安装")
	fmt.Println()
	
	// 启动Web服务器
	fmt.Println("启动Web界面...")
	server := web.NewServer(8080)
	fmt.Println("========================================")
	fmt.Println("Web界面: http://localhost:8080")
	fmt.Println("========================================")
	if err := server.Run(); err != nil {
		fmt.Printf("Web服务器错误: %v\n", err)
	}
}

func checkWireGuard() error {
	cmd := exec.Command("wg", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("WireGuard 未安装: %v", err)
	}
	fmt.Printf("WireGuard 版本: %s\n", output)
	return nil
}

// WireGuard 管理
type WireGuardMgr struct {
	configPath string
}

func NewWireGuardMgr(configPath string) *WireGuardMgr {
	return &WireGuardMgr{configPath: configPath}
}

func (m *WireGuardMgr) ListTunnels() ([]string, error) {
	cmd := exec.Command("wg", "show", "interfaces")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(string(output)), "\n"), nil
}

func (m *WireGuardMgr) GetTunnelStatus(iface string) (string, error) {
	cmd := exec.Command("wg", "show", iface)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func (m *WireGuardMgr) GenerateKeyPair() (string, string, error) {
	cmd1 := exec.Command("wg", "genkey")
	privateKey, err := cmd1.Output()
	if err != nil {
		return "", "", err
	}
	cmd2 := exec.Command("wg", "pubkey")
	cmd2.Stdin = strings.NewReader(string(privateKey))
	publicKey, err := cmd2.Output()
	return strings.TrimSpace(string(privateKey)), strings.TrimSpace(string(publicKey)), err
}
