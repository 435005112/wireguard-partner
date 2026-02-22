package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// VersionInfo 版本信息
type VersionInfo struct {
	Current  string
	Latest   string
	Download string
	NeedUpdate bool
}

// Checker 版本检测
type Checker struct{}

// NewChecker 创建版本检测器
func NewChecker() *Checker {
	return &Checker{}
}

// CheckWireGuard 检测 WireGuard 版本
func (c *Checker) CheckWireGuard() (*VersionInfo, error) {
	// 本地版本
	cmd := exec.Command("wg", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("WireGuard 未安装: %v", err)
	}
	current := strings.TrimSpace(string(output))
	
	// WireGuard 内核版本无法通过简单方式检测
	// 这里只返回本地版本
	return &VersionInfo{
		Current:    current,
		Latest:     "内核模块",
		NeedUpdate: false,
	}, nil
}

// CheckWSTunnel 检测 wstunnel 版本
func (c *Checker) CheckWSTunnel() (*VersionInfo, error) {
	// 本地版本
	cmd := exec.Command("wstunnel", "--version")
	output, err := cmd.CombinedOutput()
	current := strings.TrimSpace(string(output))
	if err != nil {
		current = "未安装"
	}
	
	// 获取最新版本
	latest, download, err := c.getLatestFromGitHub("erebe", "wstunnel")
	if err != nil {
		return &VersionInfo{Current: current, Latest: "检测失败"}, err
	}
	
	return &VersionInfo{
		Current:    current,
		Latest:     latest,
		Download:   download,
		NeedUpdate: current != latest && latest != "检测失败",
	}, nil
}

// CheckUDP2Raw 检测 udp2raw 版本
func (c *Checker) CheckUDP2Raw() (*VersionInfo, error) {
	// 本地版本
	cmd := exec.Command("udp2raw", "--version")
	output, err := cmd.CombinedOutput()
	current := strings.TrimSpace(string(output))
	if err != nil {
		current = "未安装"
	}
	
	// 获取最新版本
	latest, download, err := c.getLatestFromGitHub("wangyu-", "udp2raw")
	if err != nil {
		return &VersionInfo{Current: current, Latest: "检测失败"}, err
	}
	
	return &VersionInfo{
		Current:    current,
		Latest:     latest,
		Download:   download,
		NeedUpdate: current != latest && latest != "检测失败",
	}, nil
}

// getLatestFromGitHub 从 GitHub 获取最新版本
func (c *Checker) getLatestFromGitHub(owner, repo string) (string, string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	
	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	
	var release map[string]interface{}
	if err := json.Unmarshal(body, &release); err != nil {
		return "", "", err
	}
	
	tagName := release["tag_name"].(string)
	
	// 查找对应平台的下裁链接
	assets := release["assets"].([]interface{})
	var download string
	runtimeOS := runtime.GOOS
	runtimeARCH := runtime.GOARCH
	
	arch := "amd64"
	if runtimeARCH == "arm64" {
		arch = "arm64"
	}
	
	for _, a := range assets {
		asset := a.(map[string]interface{})
		name := asset["name"].(string)
		if strings.Contains(name, runtimeOS) && strings.Contains(name, arch) {
			download = asset["browser_download_url"].(string)
			break
		}
	}
	
	return tagName, download, nil
}

// DownloadFile 下载文件
func (c *Checker) DownloadFile(url, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	
	_, err = io.Copy(out, resp.Body)
	return err
}
