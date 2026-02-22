package config

import (
	"fmt"
	"os"
)

// Config 配置
type Config struct {
	Server   ServerConfig   `json:"server"`
	WireGuard []WGTunnel    `json:"wireguard"`
	DDNS     []DDNSConfig   `json:"ddns"`
	Tunnel   TunnelConfig   `json:"tunnel"`
}

type ServerConfig struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}

type WGTunnel struct {
	Name      string `json:"name"`
	Interface string `json:"interface"`
	Address   string `json:"address"`
	ListenPort int   `json:"listenPort"`
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
	Peers     []Peer  `json:"peers"`
}

type Peer struct {
	PublicKey           string `json:"publicKey"`
	AllowedIPs          string `json:"allowedIPs"`
	Endpoint            string `json:"endpoint"`
	PersistentKeepalive int    `json:"persistentKeepalive"`
}

type DDNSConfig struct {
	Enable    bool   `json:"enable"`
	Provider  string `json:"provider"`
	Domain    string `json:"domain"`
	SubDomain string `json:"subDomain"`
	APIKey    string `json:"apiKey"`
	APISecret string `json:"apiSecret"`
}

type TunnelConfig struct {
	Type    string `json:"type"`
	Enable  bool   `json:"enable"`
	Server  string `json:"server"`
	Port    int    `json:"port"`
	Password string `json:"password"`
}

func LoadConfig(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置文件不存在: %s", path)
	}
	
	_, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	
	var cfg Config
	// 简单解析 - 使用默认值
	cfg.Server.Port = 8080
	cfg.Server.Host = "0.0.0.0"
	
	return &cfg, nil
}

func SaveConfig(path string, cfg *Config) error {
	return nil
}
