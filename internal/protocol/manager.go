package protocol

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ProtocolType 协议类型
type ProtocolType string

const (
	ProtocolWSTunnel  ProtocolType = "wstunnel"
	ProtocolUDP2Raw   ProtocolType = "udp2raw"
)

// RoleType 角色类型
type RoleType string

const (
	RoleClient RoleType = "client"
	RoleServer RoleType = "server"
)

// WSTunnelConfig wstunnel配置
type WSTunnelConfig struct {
	// 必填
	Name     string `json:"name"`     // 配置名称
	Server   string `json:"server"`  // 服务器地址
	Port     int    `json:"port"`    // 服务器端口
	LocalTun string `json:"localTun"` // 本地隧道 -L 参数

	// 基本
	UseTLS   bool   `json:"useTLS"`   // 使用TLS
	TLSConf  TLSConfig `json:"tlsConf"` // TLS配置

	// 高级
	SNI            string `json:"sni"`             // TLS SNI
	PathPrefix     string `json:"pathPrefix"`      // HTTP升级路径前缀
	LogLevel       string `json:"logLevel"`       // 日志级别
	HTTPProxy      string `json:"httpProxy"`       // HTTP代理
	HTTPProxyLogin    string `json:"httpProxyLogin"` // 代理用户名
	HTTPProxyPass  string `json:"httpProxyPass"`  // 代理密码
	WebsocketPing int    `json:"websocketPing"`  // Websocket ping频率(秒)
}

// TLSConfig TLS配置
type TLSConfig struct {
	CertFile string `json:"certFile"` // 证书文件
	KeyFile  string `json:"keyFile"`  // 私钥文件
}

// UDP2RawConfig udp2raw配置
type UDP2RawConfig struct {
	// 必填
	Name      string `json:"name"`       // 配置名称
	Server    string `json:"server"`    // 服务器地址
	ServerPort int   `json:"serverPort"` // 服务器端口
	LocalAddr string `json:"localAddr"` // 本地监听地址
	RemoteAddr string `json:"remoteAddr"` // 远程目标地址
	Key       string `json:"key"`        // 通信密钥

	// 基本
	Mode    string `json:"mode"`     // 传输模式 faketcp/udp/icmp/easy-faketcp
	Cipher  string `json:"cipher"`   // 加密方式
	Auth    string `json:"auth"`     // 认证方式
	AutoIPT bool   `json:"autoIPT"`  // 自动iptables

	// 高级
	LogLevel      int    `json:"logLevel"`       // 日志级别 0-6
	SourceIP      string `json:"sourceIP"`       // 强制源IP
	DisableReplay bool   `json:"disableReplay"`  // 禁用防重放
	Device        string `json:"device"`         // 绑定设备
	SockBuf       int    `json:"sockBuf"`        // socket缓冲区
}

// ProtocolConfig 统一配置结构
type ProtocolConfig struct {
	Type        ProtocolType    `json:"type"`         // 协议类型
	Role        RoleType        `json:"role"`         // 角色
	WSTunnel    *WSTunnelConfig `json:"wstunnel,omitempty"`    // wstunnel配置
	UDP2Raw     *UDP2RawConfig  `json:"udp2raw,omitempty"`     // udp2raw配置
	Status      string          `json:"status"`      // running/stopped
	PID         int             `json:"pid"`         // 进程PID
	LastStart   time.Time       `json:"lastStart"`   // 最后启动时间
}

// Manager 协议管理器
type Manager struct {
	configDir string
	configs   map[string]*ProtocolConfig
}

// NewManager 创建管理器
func NewManager() *Manager {
	home := os.Getenv("HOME")
	if home == "" {
		home = "/root"
	}
	configDir := filepath.Join(home, ".wireguard-partner", "protocols")
	os.MkdirAll(configDir, 0755)

	return &Manager{
		configDir: configDir,
		configs:   make(map[string]*ProtocolConfig),
	}
}

// Load 加载配置
func (m *Manager) Load() error {
	files, err := os.ReadDir(m.configDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if filepath.Ext(f.Name()) == ".json" {
			path := filepath.Join(m.configDir, f.Name())
			data, err := os.ReadFile(path)
			if err != nil {
				continue
			}
			var cfg ProtocolConfig
			if err := json.Unmarshal(data, &cfg); err == nil {
				// 检查进程状态
				cfg.Status = m.checkStatus(cfg.PID)
				m.configs[cfg.GetName()] = &cfg
			}
		}
	}
	return nil
}

// Save 保存配置
func (m *Manager) Save(cfg *ProtocolConfig) error {
	name := cfg.GetName()
	path := filepath.Join(m.configDir, name+".json")
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Delete 删除配置
func (m *Manager) Delete(name string) error {
	// 先停止服务
	cfg := m.configs[name]
	if cfg != nil && cfg.Status == "running" {
		m.Stop(name)
	}

	path := filepath.Join(m.configDir, name+".json")
	return os.Remove(path)
}

// List 列出所有配置
func (m *Manager) List() []*ProtocolConfig {
	var list []*ProtocolConfig
	for _, cfg := range m.configs {
		// 刷新状态
		cfg.Status = m.checkStatus(cfg.PID)
		list = append(list, cfg)
	}
	return list
}

// Get 获取单个配置
func (m *Manager) Get(name string) *ProtocolConfig {
	cfg := m.configs[name]
	if cfg != nil {
		cfg.Status = m.checkStatus(cfg.PID)
	}
	return cfg
}

// GetName 获取配置名称
func (c *ProtocolConfig) GetName() string {
	if c.WSTunnel != nil {
		return c.WSTunnel.Name
	}
	if c.UDP2Raw != nil {
		return c.UDP2Raw.Name
	}
	return ""
}

// BuildCommand 构建启动命令
func (c *ProtocolConfig) BuildCommand() string {
	if c.WSTunnel != nil {
		return c.buildWSTunnelCommand()
	}
	if c.UDP2Raw != nil {
		return c.buildUDP2RawCommand()
	}
	return ""
}

func (c *ProtocolConfig) buildWSTunnelCommand() string {
	cfg := c.WSTunnel
	parts := []string{"wstunnel"}

	if c.Role == RoleClient {
		parts = append(parts, "client")
	} else {
		parts = append(parts, "server")
	}

	// 服务器地址
	scheme := "ws"
	if cfg.UseTLS {
		scheme = "wss"
	}
	serverAddr := fmt.Sprintf("%s://%s:%d", scheme, cfg.Server, cfg.Port)
	parts = append(parts, serverAddr)

	// -L 参数
	if cfg.LocalTun != "" {
		parts = append(parts, "-L", cfg.LocalTun)
	}

	// -r 参数 (服务端用)
	if c.Role == RoleServer && cfg.WSTunnel != nil {
		// 服务端需要 -r 指定转发目标
	}

	// TLS
	if cfg.TLSConf.CertFile != "" {
		parts = append(parts, "--tls-certificate", cfg.TLSConf.CertFile)
	}
	if cfg.TLSConf.KeyFile != "" {
		parts = append(parts, "--tls-private-key", cfg.TLSConf.KeyFile)
	}

	// SNI
	if cfg.SNI != "" {
		parts = append(parts, "--tls-sni-override", cfg.SNI)
	} else if cfg.UseTLS {
		parts = append(parts, "--tls-sni-disable")
	}

	// 路径前缀
	if cfg.PathPrefix != "" && cfg.PathPrefix != "v1" {
		parts = append(parts, "--http-upgrade-path-prefix", cfg.PathPrefix)
	}

	// 日志级别
	if cfg.LogLevel != "" {
		parts = append(parts, "--log-lvl", cfg.LogLevel)
	}

	// HTTP代理
	if cfg.HTTPProxy != "" {
		parts = append(parts, "--http-proxy", cfg.HTTPProxy)
	}
	if cfg.HTTPProxyLogin != "" {
		parts = append(parts, "--http-proxy-login", cfg.HTTPProxyLogin)
	}
	if cfg.HTTPProxyPass != "" {
		parts = append(parts, "--http-proxy-password", cfg.HTTPProxyPass)
	}

	return strings.Join(parts, " ")
}

func (c *ProtocolConfig) buildUDP2RawCommand() string {
	cfg := c.UDP2Raw
	parts := []string{"udp2raw"}

	// 角色
	if c.Role == RoleClient {
		parts = append(parts, "-c")
	} else {
		parts = append(parts, "-s")
	}

	// 监听地址
	parts = append(parts, "-l", cfg.LocalAddr)

	// 远程地址
	parts = append(parts, "-r", cfg.Server+":"+fmt.Sprintf("%d", cfg.ServerPort))

	// 密钥
	if cfg.Key != "" {
		parts = append(parts, "-k", cfg.Key)
	}

	// 模式
	if cfg.Mode != "" {
		parts = append(parts, "--raw-mode", cfg.Mode)
	}

	// 加密
	if cfg.Cipher != "" && cfg.Cipher != "none" {
		parts = append(parts, "--cipher-mode", cfg.Cipher)
	}

	// 认证
	if cfg.Auth != "" && cfg.Auth != "none" {
		parts = append(parts, "--auth-mode", cfg.Auth)
	}

	// 自动iptables
	if cfg.AutoIPT {
		parts = append(parts, "-a")
	}

	// 高级选项
	if cfg.LogLevel > 0 {
		parts = append(parts, "--log-level", fmt.Sprintf("%d", cfg.LogLevel))
	}

	if cfg.SourceIP != "" {
		parts = append(parts, "--source-ip", cfg.SourceIP)
	}

	if cfg.DisableReplay {
		parts = append(parts, "--disable-anti-replay")
	}

	if cfg.Device != "" {
		parts = append(parts, "--dev", cfg.Device)
	}

	if cfg.SockBuf > 0 {
		parts = append(parts, "--sock-buf", fmt.Sprintf("%d", cfg.SockBuf))
	}

	return strings.Join(parts, " ")
}

// BuildService 生成systemd service
func (c *ProtocolConfig) BuildService() string {
	name := c.GetName()
	cmd := c.BuildCommand()

	service := fmt.Sprintf(`[Unit]
Description=Protocol tunnel: %s
After=network.target

[Service]
Type=simple
User=root
ExecStart=/bin/sh -c '%s'
Restart=on-failure
RestartSec=5s
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
`, name, cmd)

	return service
}

// Start 启动服务
func (m *Manager) Start(name string) error {
	cfg := m.configs[name]
	if cfg == nil {
		return fmt.Errorf("配置不存在: %s", name)
	}

	// 生成service文件
	serviceContent := cfg.BuildService()
	serviceName := fmt.Sprintf("wg-tunnel@%s.service", name)
	servicePath := fmt.Sprintf("/etc/systemd/system/%s", serviceName)

	// 写入service文件
	if err := os.WriteFile(servicePath, []byte(serviceContent), 0644); err != nil {
		return fmt.Errorf("写入service失败: %v", err)
	}

	// reload systemd
	exec.Command("systemctl", "daemon-reload").Run()

	// 启动服务
	cmd := exec.Command("systemctl", "start", serviceName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("启动失败: %v, %s", err, string(output))
	}

	// 获取PID
	cmd = exec.Command("systemctl", "show", serviceName, "-p", "MainPID", "--value")
	if pidStr, err := cmd.Output(); err == nil {
		fmt.Sscanf(strings.TrimSpace(string(pidStr)), "%d", &cfg.PID)
	}

	cfg.Status = "running"
	cfg.LastStart = time.Now()
	return m.Save(cfg)
}

// Stop 停止服务
func (m *Manager) Stop(name string) error {
	cfg := m.configs[name]
	if cfg == nil {
		return fmt.Errorf("配置不存在: %s", name)
	}

	serviceName := fmt.Sprintf("wg-tunnel@%s.service", name)
	cmd := exec.Command("systemctl", "stop", serviceName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("停止失败: %v, %s", err, string(output))
	}

	cfg.Status = "stopped"
	cfg.PID = 0
	return m.Save(cfg)
}

// Restart 重启服务
func (m *Manager) Restart(name string) error {
	m.Stop(name)
	time.Sleep(500 * time.Millisecond)
	return m.Start(name)
}

// GetLog 获取日志
func (m *Manager) GetLog(name string, lines int) (string, error) {
	if lines <= 0 {
		lines = 100
	}

	serviceName := fmt.Sprintf("wg-tunnel@%s.service", name)
	cmd := exec.Command("journalctl", "-u", serviceName, "-n", fmt.Sprintf("%d", lines), "--no-pager", "-o", "short-iso")
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// checkStatus 检查进程状态
func (m *Manager) checkStatus(pid int) string {
	if pid <= 0 {
		return "stopped"
	}
	cmd := exec.Command("ps", "-p", fmt.Sprintf("%d", pid), "-o", "pid=")
	if err := cmd.Run(); err != nil {
		return "stopped"
	}
	return "running"
}

// InstallTools 检查并提示安装工具
func (m *Manager) InstallTools() map[string]bool {
	status := make(map[string]bool)

	// 检查 wstunnel
	if _, err := exec.LookPath("wstunnel"); err == nil {
		status["wstunnel"] = true
	} else {
		status["wstunnel"] = false
	}

	// 检查 udp2raw
	if _, err := exec.LookPath("udp2raw"); err == nil {
		status["udp2raw"] = true
	} else {
		status["udp2raw"] = false
	}

	return status
}
