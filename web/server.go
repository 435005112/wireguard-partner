package web

import (
	"fmt"
	"net/http"
	"text/template"
	"wireguard-partner/internal/tunnel"
)

var htmlTemplates = map[string]string{
	"wizard": `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>WireGuard 伴侣 - 安装向导</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; background: #f5f5f5; }
        .container { max-width: 600px; margin: 50px auto; padding: 2rem; }
        .card { background: white; border-radius: 12px; padding: 2rem; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
        h1 { color: #2563eb; margin-bottom: 0.5rem; }
        h2 { color: #1f2937; margin: 1.5rem 0 1rem; }
        .step { display: flex; align-items: center; gap: 1rem; margin-bottom: 1rem; }
        .step-num { width: 32px; height: 32px; background: #e5e7eb; color: #6b7280; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-weight: bold; }
        .step.active .step-num { background: #2563eb; color: white; }
        .step.done .step-num { background: #10b981; color: white; }
        .form-group { margin-bottom: 1rem; }
        .form-group label { display: block; margin-bottom: 0.5rem; color: #374151; }
        .form-group input, .form-group select { width: 100%; padding: 0.75rem; border: 1px solid #d1d5db; border-radius: 8px; font-size: 1rem; }
        .btn { display: inline-block; padding: 0.75rem 1.5rem; background: #2563eb; color: white; border: none; border-radius: 8px; cursor: pointer; font-size: 1rem; }
        .btn:hover { background: #1d4ed8; }
        .tips { background: #fef3c7; padding: 1rem; border-radius: 8px; margin: 1rem 0; font-size: 0.875rem; color: #92400e; }
        .success { background: #d1fae5; padding: 1rem; border-radius: 8px; margin: 1rem 0; color: #065f46; }
    </style>
</head>
<body>
    <div class="container">
        <div class="card">
            <h1>WireGuard 伴侣</h1>
            <p style="color: #6b7280; margin-bottom: 2rem;">欢迎使用安装向导</p>
            
            <div class="step {{if eq .Step 1}}active{{else if lt .Step 1}}done{{else}}done{{end}}">
                <div class="step-num">1</div>
                <div>WireGuard 配置</div>
            </div>
            <div class="step {{if eq .Step 2}}active{{else if gt .Step 2}}done{{end}}">
                <div class="step-num">2</div>
                <div>DDNS 配置</div>
            </div>
            <div class="step {{if eq .Step 3}}active{{else if gt .Step 3}}done{{end}}">
                <div class="step-num">3</div>
                <div>协议封装配置</div>
            </div>
            
            {{if eq .Step 1}}
            <h2>第一步：WireGuard 配置</h2>
            <form method="POST" action="/wizard/wireguard">
                <div class="form-group">
                    <label>隧道名称</label>
                    <input type="text" name="name" placeholder="wg0" required>
                </div>
                <div class="form-group">
                    <label>服务器地址 (CIDR)</label>
                    <input type="text" name="address" placeholder="10.0.0.1/24" required>
                </div>
                <div class="form-group">
                    <label>监听端口</label>
                    <input type="number" name="port" value="51820" required>
                </div>
                <button type="submit" class="btn">下一步</button>
            </form>
            {{end}}
            
            {{if eq .Step 2}}
            <h2>第二步：DDNS 配置</h2>
            <form method="POST" action="/wizard/ddns">
                <div class="form-group">
                    <label>DNS服务商</label>
                    <select name="provider">
                        <option value="aliyun">阿里云</option>
                        <option value="cloudflare">Cloudflare</option>
                        <option value="dnspod">腾讯云DNSPod</option>
                    </select>
                </div>
                <div class="form-group">
                    <label>完整域名</label>
                    <input type="text" name="domain" placeholder="vpn.example.com">
                </div>
                <div class="form-group">
                    <label>API Key</label>
                    <input type="text" name="apiKey">
                </div>
                <div class="form-group">
                    <label>API Secret</label>
                    <input type="password" name="apiSecret">
                </div>
                <div class="tips">如不使用DDNS，可直接点击下一步跳过</div>
                <button type="submit" class="btn">下一步</button>
            </form>
            {{end}}
            
            {{if eq .Step 3}}
            <h2>第三步：协议封装配置</h2>
            <form method="POST" action="/wizard/tunnel">
                <div class="form-group">
                    <label>协议类型</label>
                    <select name="type">
                        <option value="none">不启用</option>
                        <option value="wstunnel">wstunnel</option>
                        <option value="udp2raw">udp2raw</option>
                    </select>
                </div>
                <div class="form-group">
                    <label>远程服务器</label>
                    <input type="text" name="server" placeholder="example.com:443">
                </div>
                <div class="form-group">
                    <label>端口</label>
                    <input type="number" name="port" value="443">
                </div>
                <button type="submit" class="btn">完成配置</button>
            </form>
            {{end}}
            
            {{if eq .Step 4}}
            <div class="success">
                <h2>配置完成！</h2>
                <p>配置文件已生成，点击下方按钮重启服务</p>
            </div>
            <a href="/wizard/restart" class="btn">重启服务</a>
            {{end}}
        </div>
    </div>
</body>
</html>`,

	"index": `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>WireGuard 伴侣</title>
    <meta http-equiv="refresh" content="30">
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: -apple-system, BlinkMacSystemFont, sans-serif; background: #f5f5f5; }
        .header { background: #2563eb; color: white; padding: 1rem 2rem; display: flex; justify-content: space-between; }
        .header h1 { font-size: 1.5rem; }
        .nav { background: white; padding: 0 2rem; border-bottom: 1px solid #e5e7eb; }
        .nav ul { display: flex; list-style: none; gap: 2rem; }
        .nav a { display: block; padding: 1rem 0; color: #6b7280; text-decoration: none; border-bottom: 2px solid transparent; }
        .nav a:hover, .nav a.active { color: #2563eb; border-color: #2563eb; }
        .container { padding: 2rem; max-width: 1200px; margin: 0 auto; }
        .card { background: white; border-radius: 8px; padding: 1.5rem; margin-bottom: 1.5rem; }
        .card h2 { margin-bottom: 1rem; color: #1f2937; }
        .status-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 1rem; }
        .status-item { padding: 1rem; background: #f9fafb; border-radius: 6px; text-align: center; }
        .status-item .label { color: #6b7280; font-size: 0.875rem; }
        .status-item .value { font-size: 1.5rem; font-weight: bold; margin-top: 0.5rem; }
        .status-item .value.ok { color: #10b981; }
        .status-item .value.stopped { color: #6b7280; }
    </style>
</head>
<body>
    <div class="header">
        <h1>WireGuard 伴侣</h1>
        <a href="/wizard" style="color: white; text-decoration: none;">重新配置</a>
    </div>
    <nav class="nav">
        <ul>
            <li><a href="/" class="active">概览</a></li>
            <li><a href="/wireguard">配置</a></li>
        </ul>
    </nav>
    <div class="container">
        <div class="card">
            <h2>模块状态 (30秒刷新)</h2>
            <div class="status-grid">
                <div class="status-item">
                    <div class="label">WireGuard 隧道</div>
                    <div class="value {{if .Status.Wireguard.Count}}ok{{else}}stopped{{end}}">{{.Status.Wireguard.Count}} 个</div>
                </div>
                <div class="status-item">
                    <div class="label">wstunnel</div>
                    <div class="value {{if .Status.WSTunnel.Running}}ok{{else}}stopped{{end}}">{{if .Status.WSTunnel.Running}}运行中{{else}}未运行{{end}}</div>
                </div>
                <div class="status-item">
                    <div class="label">udp2raw</div>
                    <div class="value {{if .Status.UDP2Raw.Running}}ok{{else}}stopped{{end}}">{{if .Status.UDP2Raw.Running}}运行中{{else}}未运行{{end}}</div>
                </div>
                <div class="status-item">
                    <div class="label">DDNS</div>
                    <div class="value {{if .DDNSEnabled}}ok{{else}}stopped{{end}}">{{if .DDNSEnabled}}已配置{{else}}未配置{{end}}</div>
                </div>
            </div>
        </div>
    </div>
</body>
</html>`,

	"wireguard": `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>配置 - WireGuard 伴侣</title>
    <style>
        body { font-family: -apple-system, sans-serif; background: #f5f5f5; margin: 0; }
        .header { background: #2563eb; color: white; padding: 1rem 2rem; }
        .nav { background: white; padding: 0 2rem; border-bottom: 1px solid #e5e7eb; }
        .nav ul { display: flex; list-style: none; gap: 2rem; margin: 0; }
        .nav a { display: block; padding: 1rem 0; color: #6b7280; text-decoration: none; border-bottom: 2px solid transparent; }
        .nav a:hover, .nav a.active { color: #2563eb; border-color: #2563eb; }
        .container { padding: 2rem; max-width: 800px; margin: 0 auto; }
        .card { background: white; border-radius: 8px; padding: 1.5rem; }
        .card h2 { margin-bottom: 1rem; color: #1f2937; }
        pre { background: #f3f4f6; padding: 1rem; border-radius: 6px; overflow-x: auto; }
    </style>
</head>
<body>
    <div class="header"><h1>WireGuard 伴侣</h1></div>
    <nav class="nav">
        <ul>
            <li><a href="/">概览</a></li>
            <li><a href="/wireguard" class="active">配置</a></li>
        </ul>
    </nav>
    <div class="container">
        <div class="card">
            <h2>当前配置</h2>
            <pre>{{.Config}}</pre>
        </div>
    </div>
</body>
</html>`,
}

// WizardState 向导状态
type WizardState struct {
	Step     int
	Wireguard struct {
		Name    string
		Address string
		Port    int
	}
	DDNS struct {
		Provider  string
		Domain    string
		APIKey    string
		APISecret string
	}
	Tunnel struct {
		Type     string
		Server   string
		Port     int
		Password string
	}
	Config string
}

// Server Web服务器
type Server struct {
	port      int
	mux       *http.ServeMux
	tmpls     map[string]*template.Template
	tunnelMgr *tunnel.TunnelMgr
	wizard    *WizardState
}

// NewServer 创建Web服务器
func NewServer(port int) *Server {
	s := &Server{
		port:      port,
		mux:       http.NewServeMux(),
		tmpls:     make(map[string]*template.Template),
		tunnelMgr: tunnel.NewTunnelMgr(),
		wizard:    &WizardState{Step: 1},
	}
	
	for name, html := range htmlTemplates {
		s.tmpls[name] = template.Must(template.New(name).Parse(html))
	}
	
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.mux.HandleFunc("/", s.handleIndex)
	s.mux.HandleFunc("/wizard", s.handleWizard)
	s.mux.HandleFunc("/wizard/wireguard", s.handleWizardWireguard)
	s.mux.HandleFunc("/wizard/ddns", s.handleWizardDDNS)
	s.mux.HandleFunc("/wizard/tunnel", s.handleWizardTunnel)
	s.mux.HandleFunc("/wizard/restart", s.handleWizardRestart)
	s.mux.HandleFunc("/wireguard", s.handleWireGuard)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	status := s.tunnelMgr.GetModuleStatus()
	ddnsEnabled := s.wizard.DDNS.Domain != ""
	s.tmpls["index"].Execute(w, map[string]interface{}{
		"Status":     status,
		"DDNSEnabled": ddnsEnabled,
	})
}

func (s *Server) handleWizard(w http.ResponseWriter, r *http.Request) {
	// 重置向导
	s.wizard = &WizardState{Step: 1}
	s.tmpls["wizard"].Execute(w, s.wizard)
}

func (s *Server) handleWizardWireguard(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	s.wizard.Step = 2
	s.wizard.Wireguard.Name = r.Form.Get("name")
	s.wizard.Wireguard.Address = r.Form.Get("address")
	fmt.Sscanf(r.Form.Get("port"), "%d", &s.wizard.Wireguard.Port)
	s.generateConfig()
	s.tmpls["wizard"].Execute(w, s.wizard)
}

func (s *Server) handleWizardDDNS(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	s.wizard.Step = 3
	s.wizard.DDNS.Provider = r.Form.Get("provider")
	s.wizard.DDNS.Domain = r.Form.Get("domain")
	s.wizard.DDNS.APIKey = r.Form.Get("apiKey")
	s.wizard.DDNS.APISecret = r.Form.Get("apiSecret")
	s.generateConfig()
	s.tmpls["wizard"].Execute(w, s.wizard)
}

func (s *Server) handleWizardTunnel(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	s.wizard.Step = 4
	s.wizard.Tunnel.Type = r.Form.Get("type")
	s.wizard.Tunnel.Server = r.Form.Get("server")
	fmt.Sscanf(r.Form.Get("port"), "%d", &s.wizard.Tunnel.Port)
	s.wizard.Tunnel.Password = r.Form.Get("password")
	s.generateConfig()
	s.tmpls["wizard"].Execute(w, s.wizard)
}

func (s *Server) handleWizardRestart(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "服务重启中...(需手动执行 sudo wg-quick up %s)", s.wizard.Wireguard.Name)
}

func (s *Server) generateConfig() {
	s.wizard.Config = fmt.Sprintf(`# WireGuard 配置
[Interface]
Address = %s
ListenPort = %d

# DDNS: %s %s

# 协议封装: %s %s
`, s.wizard.Wireguard.Address, s.wizard.Wireguard.Port,
		s.wizard.DDNS.Provider, s.wizard.DDNS.Domain,
		s.wizard.Tunnel.Type, s.wizard.Tunnel.Server)
}

func (s *Server) handleWireGuard(w http.ResponseWriter, r *http.Request) {
	s.tmpls["wireguard"].Execute(w, map[string]interface{}{
		"Config": s.wizard.Config,
	})
}

// Run 运行服务器
func (s *Server) Run() error {
	addr := fmt.Sprintf(":%d", s.port)
	fmt.Printf("Web界面: http://localhost%s\n", addr)
	fmt.Printf("安装向导: http://localhost%s/wizard\n", addr)
	return http.ListenAndServe(addr, s.mux)
}
