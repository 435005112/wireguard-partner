package installer

import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
)

// Detection 环境检测结果
type Detection struct {
	OS             string `json:"os"`
	Arch           string `json:"arch"`
	WireGuard      bool   `json:"wireguard"`
	WSTunnel       bool   `json:"wstunnel"`
	UDP2Raw        bool   `json:"udp2raw"`
	DDNSGo         bool   `json:"ddnsgo"`
	GitHubReachable bool  `json:"githubReachable"`
	Location       string `json:"location"`
}

// Detector 环境检测器
type Detector struct{}

// NewDetector 创建检测器
func NewDetector() *Detector {
	return &Detector{}
}

// Detect 检测系统环境
func (d *Detector) Detect() (*Detection, error) {
	result := &Detection{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	result.WireGuard = d.checkCommand("wg")
	result.WSTunnel = d.checkCommand("wstunnel")
	result.UDP2Raw = d.checkCommand("udp2raw")
	result.DDNSGo = d.checkCommand("ddns-go")

	result.GitHubReachable = d.checkGitHub()
	if result.GitHubReachable {
		result.Location = "Global"
	} else {
		result.Location = "CN"
	}

	return result, nil
}

func (d *Detector) checkCommand(name string) bool {
	cmd := exec.Command("which", name)
	return cmd.Run() == nil
}

func (d *Detector) checkGitHub() bool {
	client := &http.Client{Timeout: 5}
	resp, err := client.Get("https://github.com")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

// Installer 自动安装器
type Installer struct {
	Location string
}

// NewInstaller 创建安装器
func NewInstaller(location string) *Installer {
	return &Installer{Location: location}
}

// QuickInstall 快速安装
func QuickInstall() error {
	detector := NewDetector()
	detection, err := detector.Detect()
	if err != nil {
		return err
	}

	installer := NewInstaller(detection.Location)

	fmt.Printf("检测结果:\n")
	fmt.Printf("  系统: %s %s\n", detection.OS, detection.Arch)
	fmt.Printf("  服务器: %s\n", detection.Location)
	fmt.Printf("  WireGuard: %v\n", detection.WireGuard)
	fmt.Printf("  wstunnel: %v\n", detection.WSTunnel)
	fmt.Printf("  udp2raw: %v\n", detection.UDP2Raw)
	fmt.Println()

	if !detection.WireGuard {
		fmt.Println("安装 WireGuard...")
		installer.InstallWireGuard()
	}

	if !detection.WSTunnel {
		fmt.Println("安装 wstunnel...")
		installer.InstallWSTunnel()
	}

	if !detection.UDP2Raw {
		fmt.Println("安装 udp2raw...")
		installer.InstallUDP2Raw()
	}

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("            安装完成!")
	fmt.Println("========================================")
	fmt.Println("启动: ./wg-partner")
	fmt.Println("Web: http://<ip>:8080")
	fmt.Println("账号: admin / admin123")

	return nil
}

// InstallWireGuard 安装WireGuard
func (i *Installer) InstallWireGuard() error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		if i.isDebian() {
			cmd = exec.Command("sudo", "apt", "update")
			cmd.Run()
			cmd = exec.Command("sudo", "apt", "install", "-y", "wireguard-tools")
		} else if i.isRHEL() {
			cmd = exec.Command("sudo", "yum", "install", "-y", "wireguard-tools")
		}
	}

	if cmd == nil {
		fmt.Println("不支持的系统")
		return nil
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("安装失败: %v\n%s\n", err, output)
		return err
	}
	fmt.Println("✅ WireGuard 安装完成")
	return nil
}

// InstallWSTunnel 安装wstunnel
func (i *Installer) InstallWSTunnel() error {
	arch := runtime.GOARCH
	if arch == "amd64" {
		arch = "x86_64"
	}

	baseURL := "https://github.com/erebe/wstunnel/releases/latest/download"
	if i.Location == "CN" {
		baseURL = "https://ghproxy.com/https://github.com/erebe/wstunnel/releases/latest/download"
	}

	url := fmt.Sprintf("%s/wstunnel-%s-unknown-linux-%s", baseURL, "v7.0.0", arch)

	cmd := exec.Command("sudo", "curl", "-L", "-o", "/usr/local/bin/wstunnel", url)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("下载失败: %v\n%s\n", err, output)
		return err
	}

	exec.Command("sudo", "chmod", "+x", "/usr/local/bin/wstunnel").Run()
	fmt.Println("✅ wstunnel 安装完成")
	return nil
}

// InstallUDP2Raw 安装udp2raw
func (i *Installer) InstallUDP2Raw() error {
	baseURL := "https://github.com/wangyu-/udp2raw/releases/latest/download"
	if i.Location == "CN" {
		baseURL = "https://ghproxy.com/https://github.com/wangyu-/udp2raw/releases/latest/download"
	}

	url := fmt.Sprintf("%s/udp2raw_amd64", baseURL)

	cmd := exec.Command("sudo", "curl", "-L", "-o", "/usr/local/bin/udp2raw", url)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("下载失败: %v\n%s\n", err, output)
		return err
	}

	exec.Command("sudo", "chmod", "+x", "/usr/local/bin/udp2raw").Run()
	fmt.Println("✅ udp2raw 安装完成")
	return nil
}

func (i *Installer) isDebian() bool {
	cmd := exec.Command("which", "apt")
	return cmd.Run() == nil
}

func (i *Installer) isRHEL() bool {
	cmd := exec.Command("which", "yum")
	return cmd.Run() == nil
}

var Detect = NewDetector().Detect
