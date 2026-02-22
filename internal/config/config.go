package config

import (
	"fmt"
	"os"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	WireGuard []WGTunnel    `yaml:"wireguard"`
	DDNS     []DDNSConfig   `yaml:"ddns"`
	Tunnel   TunnelConfig   `yaml:"tunnel"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type WGTunnel struct {
	Name      string `yaml:"name"`
	Interface string `yaml:"interface"`
	Address   string `yaml:"address"`
	ListenPort int   `yaml:"listenPort"`
	PrivateKey string `yaml:"privateKey"`
	PublicKey  string `yaml:"publicKey"`
	Peers     []Peer  `yaml:"peers"`
}

type Peer struct {
	PublicKey           string `yaml:"publicKey"`
	AllowedIPs          string `yaml:"allowedIPs"`
	Endpoint            string `yaml:"endpoint"`
	PersistentKeepalive int    `yaml:"persistentKeepalive"`
}

type DDNSConfig struct {
	Enable    bool   `yaml:"enable"`
	Provider  string `yaml:"provider"`  // aliyun, cloudflare, etc
	Domain    string `yaml:"domain"`
	SubDomain string `yaml:"subDomain"`
	APIKey    string `yaml:"apiKey"`
	APISecret string `yaml:"apiSecret"`
}

type TunnelConfig struct {
	Type    string `yaml:"type"`    // wstunnel, udp2raw
	Enable  bool   `yaml:"enable"`
	Server  string `yaml:"server"`
	Port    int    `yaml:"port"`
	Password string `yaml:"password"`
}

func LoadConfig(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置文件不存在: %s", path)
	}
	
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	
	return &cfg, nil
}

func SaveConfig(path string, cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
