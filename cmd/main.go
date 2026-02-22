package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	fmt.Println("WireGuard 伴侣")
	fmt.Println("================")
	
	// 检查 WireGuard 是否安装
	if err := checkWireGuard(); err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("WireGuard 已安装")
	
	// 显示当前状态
	showStatus()
}

func checkWireGuard() error {
	cmd := exec.Command("wg", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("WireGuard 未安装: %v", err)
	}
	fmt.Printf("WireGuard 版本: %s", output)
	return nil
}

func showStatus() {
	cmd := exec.Command("wg", "show")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("暂无 WireGuard 隧道")
		return
	}
	fmt.Println(string(output))
}

// WireGuard 管理
type WireGuardMgr struct {
	configPath string
}

func NewWireGuardMgr(configPath string) *WireGuardMgr {
	return &WireGuardMgr{configPath: configPath}
}

// ListTunnels 列出所有隧道
func (m *WireGuardMgr) ListTunnels() ([]string, error) {
	cmd := exec.Command("wg", "show", "interfaces")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(string(output)), "\n"), nil
}

// GetTunnelStatus 获取隧道状态
func (m *WireGuardMgr) GetTunnelStatus(iface string) (string, error) {
	cmd := exec.Command("wg", "show", iface)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// GenerateKeyPair 生成密钥对
func (m *WireGuardMgr) GenerateKeyPair() (string, string, error) {
	// 生成私钥
	cmd1 := exec.Command("wg", "genkey")
	privateKey, err := cmd1.Output()
	if err != nil {
		return "", "", err
	}
	
	// 生成公钥
	cmd2 := exec.Command("wg", "pubkey")
	cmd2.Stdin = strings.NewReader(string(privateKey))
	publicKey, err := cmd2.Output()
	
	return strings.TrimSpace(string(privateKey)), 
	       strings.TrimSpace(string(publicKey)), err
}
