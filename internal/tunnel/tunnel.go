package tunnel

import (
	"fmt"
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
	var cmd *exec.Cmd
	var args []string
	
	switch p {
	case ProtocolWSTunnel:
		cmd = exec.Command("wstunnel", "--version")
	case ProtocolUDP2Raw:
		cmd = exec.Command("udp2raw", "--version")
	}
	
	if err := cmd.Run(); err != nil {
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

// StartWSTunnel 启动 wstunnel
func (m *TunnelMgr) StartWSTunnel(listenAddr, remoteAddr, tunnelType string) error {
	args := []string{
		"-l", listenAddr,
		"-r", remoteAddr,
	}
	
	if tunnelType == "stunnel" {
		args = append(args, "-t", "2")
	}
	
	cmd := exec.Command("wstunnel", args...)
	return cmd.Start()
}

// StartUDP2Raw 启动 udp2raw
func (m *TunnelMgr) StartUDP2Raw(mode, localAddr, remoteAddr, password string) error {
	args := []string{
		fmt.Sprintf("-%s", mode), // -s or -c
		fmt.Sprintf("-l%s", localAddr),
		fmt.Sprintf("-r%s", remoteAddr),
		"-k", password,
		"--raw-mode", "faketcp",
		"-a",
	}
	
	cmd := exec.Command("udp2raw_amd64", args...)
	return cmd.Start()
}

// StopProcess 停止进程
func (m *TunnelMgr) StopProcess(name string) error {
	cmd := exec.Command("pkill", name)
	return cmd.Run()
}

// GetProcessStatus 获取进程状态
func (m *TunnelMgr) GetProcessStatus(name string) (bool, error) {
	cmd := exec.Command("pgrep", "-f", name)
	err := cmd.Run()
	if err != nil {
		return false, nil // not running
	}
	return true, nil // running
}
