package tunnel

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ProtocolType 协议类型
type ProtocolType string

const (
	ProtocolWSTunnel ProtocolType = "wstunnel"
	ProtocolUDP2Raw   ProtocolType = "udp2raw"
)

// TunnelMgr 隧道管理
type TunnelMgr struct{}

// NewTunnelMgr 创建隧道管理器
func NewTunnelMgr() *TunnelMgr {
	return &TunnelMgr{}
}

// CheckInstalled 检查工具是否已安装
func (m *TunnelMgr) CheckInstalled(p ProtocolType) (bool, error) {
	var name string
	
	switch p {
	case ProtocolWSTunnel:
		name = "wstunnel"
	case ProtocolUDP2Raw:
		name = "udp2raw"
	}
	
	_, err := exec.LookPath(name)
	if err != nil {
		return false, nil
	}
	return true, nil
}

// GetVersion 获取版本
func (m *TunnelMgr) GetVersion(p ProtocolType) (string, error) {
	var cmd *exec.Cmd
	
	switch p {
	case ProtocolWSTunnel:
		cmd = exec.Command("wstunnel", "--version")
	case ProtocolUDP2Raw:
		cmd = exec.Command("udp2raw", "--version")
	}
	
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

// CreateWireGuardTunnel 创建WireGuard隧道
func (m *TunnelMgr) CreateWireGuardTunnel(name, address string, port int) error {
	// 生成密钥
	privateKey, err := exec.Command("wg", "genkey").Output()
	if err != nil {
		return fmt.Errorf("生成私钥失败: %v", err)
	}
	
	pubKeyCmd := exec.Command("wg", "pubkey")
	pubKeyCmd.Stdin = strings.NewReader(string(privateKey))
	pubKeyOutput, err := pubKeyCmd.Output()
	if err != nil {
		return fmt.Errorf("生成公钥失败: %v", err)
	}
	
	_ = pubKeyOutput // 公钥待用
	
	config := fmt.Sprintf(`[Interface]
Address = %s
ListenPort = %d
PrivateKey = %s
SaveConfig = true

[Peer]
# PublicKey = 
# AllowedIPs = 0.0.0.0/0
`, address, port, strings.TrimSpace(string(privateKey)))
	
	// 创建配置文件
	configPath := fmt.Sprintf("/etc/wireguard/%s.conf", name)
	os.WriteFile(configPath, []byte(config), 0600)
	
	// 启动隧道
	cmd := exec.Command("wg-quick", "up", name)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("启动隧道失败: %v", err)
	}
	
	return nil
}

// DeleteWireGuardTunnel 删除WireGuard隧道
func (m *TunnelMgr) DeleteWireGuardTunnel(name string) error {
	cmd := exec.Command("wg-quick", "down", name)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("停止隧道失败: %v", err)
	}
	
	// 删除配置文件
	configPath := fmt.Sprintf("/etc/wireguard/%s.conf", name)
	os.Remove(configPath)
	
	return nil
}

// ListWireGuardTunnels 列出所有WireGuard隧道
func (m *TunnelMgr) ListWireGuardTunnels() ([]string, error) {
	cmd := exec.Command("wg", "show", "interfaces")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	
	interfaces := strings.Split(strings.TrimSpace(string(output)), "\n")
	var tunnels []string
	for _, i := range interfaces {
		if i != "" {
			tunnels = append(tunnels, i)
		}
	}
	return tunnels, nil
}

// GetTunnelStatus 获取隧道状态
func (m *TunnelMgr) GetTunnelStatus(iface string) (string, error) {
	cmd := exec.Command("wg", "show", iface)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// GetModuleStatus 获取模块运行状态
func (m *TunnelMgr) GetModuleStatus() map[string]interface{} {
	status := make(map[string]interface{})
	
	// WireGuard状态
	wgTunnels, _ := m.ListWireGuardTunnels()
	status["wireguard"] = map[string]interface{}{
		"installed": true,
		"tunnels":   wgTunnels,
		"count":     len(wgTunnels),
	}
	
	// wstunnel状态
	wsPID, _ := exec.Command("pgrep", "-f", "wstunnel").Output()
	status["wstunnel"] = map[string]interface{}{
		"installed": true,
		"running":   len(wsPID) > 0,
	}
	
	// udp2raw状态
	udpPID, _ := exec.Command("pgrep", "-f", "udp2raw").Output()
	status["udp2raw"] = map[string]interface{}{
		"installed": true,
		"running":   len(udpPID) > 0,
	}
	
	return status
}
