package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
	"wireguard-partner/internal/tunnel"
)

var htmlTemplates = map[string]string{
	"index": `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>WireGuard 伴侣</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background: #f5f5f5; }
        .header { background: #2563eb; color: white; padding: 1rem 2rem; }
        .header h1 { font-size: 1.5rem; }
        .nav { background: white; padding: 0 2rem; border-bottom: 1px solid #e5e7eb; }
        .nav ul { display: flex; list-style: none; gap: 2rem; }
        .nav a { display: block; padding: 1rem 0; color: #6b7280; text-decoration: none; border-bottom: 2px solid transparent; }
        .nav a:hover, .nav a.active { color: #2563eb; border-color: #2563eb; }
        .container { padding: 2rem; max-width: 1200px; margin: 0 auto; }
        .card { background: white; border-radius: 8px; padding: 1.5rem; margin-bottom: 1.5rem; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
        .card h2 { margin-bottom: 1rem; color: #1f2937; }
        .status-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 1rem; }
        .status-item { padding: 1rem; background: #f9fafb; border-radius: 6px; text-align: center; }
        .status-item .label { color: #6b7280; font-size: 0.875rem; }
        .status-item .value { font-size: 1.5rem; font-weight: bold; margin-top: 0.5rem; }
        .status-item .value.ok { color: #10b981; }
        .status-item .value.error { color: #ef4444; }
        .status-item .value.running { color: #10b981; }
        .status-item .value.stopped { color: #6b7280; }
        .btn { display: inline-block; padding: 0.5rem 1rem; background: #2563eb; color: white; border: none; border-radius: 6px; cursor: pointer; text-decoration: none; }
        .btn:hover { background: #1d4ed8; }
        .btn-danger { background: #ef4444; }
        .btn-danger:hover { background: #dc2626; }
        .form-group { margin-bottom: 1rem; }
        .form-group label { display: block; margin-bottom: 0.5rem; color: #374151; }
        .form-group input, .form-group select { width: 100%; padding: 0.5rem; border: 1px solid #d1d5db; border-radius: 6px; }
        table { width: 100%; border-collapse: collapse; }
        th, td { padding: 0.75rem; text-align: left; border-bottom: 1px solid #e5e7eb; }
        th { background: #f9fafb; font-weight: 600; }
    </style>
</head>
<body>
    <div class="header">
        <h1>WireGuard 伴侣</h1>
    </div>
    <nav class="nav">
        <ul>
            <li><a href="/" class="active">概览</a></li>
            <li><a href="/wireguard">WireGuard</a></li>
            <li><a href="/ddns">DDNS</a></li>
            <li><a href="/tunnel">协议封装</a></li>
            <li><a href="/monitor">监控</a></li>
        </ul>
    </nav>
    <div class="container">
        <div class="card">
            <h2>模块状态</h2>
            <div class="status-grid">
                <div class="status-item">
                    <div class="label">WireGuard 隧道</div>
                    <div class="value {{if .Status.Wireguard.Count}}ok{{else}}error{{end}}">{{.Status.Wireguard.Count}} 个</div>
                </div>
                <div class="status-item">
                    <div class="label">wstunnel</div>
                    <div class="value {{if .Status.WSTunnel.Running}}running{{else}}stopped{{end}}">{{if .Status.WSTunnel.Running}}运行中{{else}}未运行{{end}}</div>
                </div>
                <div class="status-item">
                    <div class="label">udp2raw</div>
                    <div class="value {{if .Status.UDP2Raw.Running}}running{{else}}stopped{{end}}">{{if .Status.UDP2Raw.Running}}运行中{{else}}未运行{{end}}</div>
                </div>
            </div>
        </div>
        <div class="card">
            <h2>快速操作</h2>
            <p><a href="/wireguard" class="btn">管理 WireGuard 隧道</a></p>
        </div>
    </div>
</body>
</html>`,

	"wireguard": `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>WireGuard - WireGuard 伴侣</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background: #f5f5f5; margin: 0; }
        .header { background: #2563eb; color: white; padding: 1rem 2rem; }
        .nav { background: white; padding: 0 2rem; border-bottom: 1px solid #e5e7eb; }
        .nav ul { display: flex; list-style: none; gap: 2rem; margin: 0; }
        .nav a { display: block; padding: 1rem 0; color: #6b7280; text-decoration: none; border-bottom: 2px solid transparent; }
        .nav a:hover, .nav a.active { color: #2563eb; border-color: #2563eb; }
        .container { padding: 2rem; max-width: 1200px; margin: 0 auto; }
        .card { background: white; border-radius: 8px; padding: 1.5rem; margin-bottom: 1.5rem; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
        .card h2 { margin-bottom: 1rem; color: #1f2937; }
        .btn { display: inline-block; padding: 0.5rem 1rem; background: #2563eb; color: white; border: none; border-radius: 6px; cursor: pointer; text-decoration: none; }
        .btn:hover { background: #1d4ed8; }
        .btn-danger { background: #ef4444; }
        .btn-danger:hover { background: #dc2626; }
        .form-group { margin-bottom: 1rem; }
        .form-group label { display: block; margin-bottom: 0.5rem; color: #374151; }
        .form-group input { width: 100%; padding: 0.5rem; border: 1px solid #d1d5db; border-radius: 6px; }
        table { width: 100%; border-collapse: collapse; }
        th, td { padding: 0.75rem; text-align: left; border-bottom: 1px solid #e5e7eb; }
        th { background: #f9fafb; font-weight: 600; }
    </style>
</head>
<body>
    <div class="header"><h1>WireGuard 伴侣</h1></div>
    <nav class="nav">
        <ul>
            <li><a href="/">概览</a></li>
            <li><a href="/wireguard" class="active">WireGuard</a></li>
            <li><a href="/ddns">DDNS</a></li>
            <li><a href="/tunnel">协议封装</a></li>
            <li><a href="/monitor">监控</a></li>
        </ul>
    </nav>
    <div class="container">
        <div class="card">
            <h2>WireGuard 隧道列表</h2>
            <table>
                <thead>
                    <tr>
                        <th>接口名称</th>
                        <th>状态</th>
                        <th>操作</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .Tunnels}}
                    <tr>
                        <td>{{.}}</td>
                        <td>活跃</td>
                        <td><a href="/wireguard/delete/{{.}}" class="btn btn-danger" onclick="return confirm('确定删除?')">删除</a></td>
                    </tr>
                    {{else}}
                    <tr><td colspan="3">暂无隧道</td></tr>
                    {{end}}
                </tbody>
            </table>
        </div>
        <div class="card">
            <h2>创建新隧道</h2>
            <form method="POST" action="/api/wireguard/create">
                <div class="form-group">
                    <label>接口名称</label>
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
                <button type="submit" class="btn">创建隧道</button>
            </form>
        </div>
    </div>
</body>
</html>`,

	"ddns": `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>DDNS - WireGuard 伴侣</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background: #f5f5f5; margin: 0; }
        .header { background: #2563eb; color: white; padding: 1rem 2rem; }
        .nav { background: white; padding: 0 2rem; border-bottom: 1px solid #e5e7eb; }
        .nav ul { display: flex; list-style: none; gap: 2rem; margin: 0; }
        .nav a { display: block; padding: 1rem 0; color: #6b7280; text-decoration: none; border-bottom: 2px solid transparent; }
        .nav a:hover, .nav a.active { color: #2563eb; border-color: #2563eb; }
        .container { padding: 2rem; max-width: 1200px; margin: 0 auto; }
        .card { background: white; border-radius: 8px; padding: 1.5rem; margin-bottom: 1.5rem; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
        .card h2 { margin-bottom: 1rem; color: #1f2937; }
        .btn { display: inline-block; padding: 0.5rem 1rem; background: #2563eb; color: white; border: none; border-radius: 6px; cursor: pointer; }
        .btn:hover { background: #1d4ed8; }
        .form-group { margin-bottom: 1rem; }
        .form-group label { display: block; margin-bottom: 0.5rem; color: #374151; }
        .form-group input, .form-group select { width: 100%; padding: 0.5rem; border: 1px solid #d1d5db; border-radius: 6px; }
    </style>
</head>
<body>
    <div class="header"><h1>WireGuard 伴侣</h1></div>
    <nav class="nav">
        <ul>
            <li><a href="/">概览</a></li>
            <li><a href="/wireguard">WireGuard</a></li>
            <li><a href="/ddns" class="active">DDNS</a></li>
            <li><a href="/tunnel">协议封装</a></li>
            <li><a href="/monitor">监控</a></li>
        </ul>
    </nav>
    <div class="container">
        <div class="card">
            <h2>DDNS 配置</h2>
            <form method="POST" action="/api/ddns/config">
                <div class="form-group">
                    <label>DNS服务商</label>
                    <select name="provider">
                        <option value="aliyun">阿里云</option>
                        <option value="cloudflare">Cloudflare</option>
                        <option value="dnspod">腾讯云DNSPod</option>
                        <option value="huawei">华为云</option>
                    </select>
                </div>
                <div class="form-group">
                    <label>域名</label>
                    <input type="text" name="domain" placeholder="example.com">
                </div>
                <div class="form-group">
                    <label>子域名</label>
                    <input type="text" name="subdomain" placeholder="vpn">
                </div>
                <div class="form-group">
                    <label>API Key</label>
                    <input type="text" name="apiKey" placeholder="API Key">
                </div>
                <div class="form-group">
                    <label>API Secret</label>
                    <input type="password" name="apiSecret" placeholder="API Secret">
                </div>
                <button type="submit" class="btn">保存配置</button>
            </form>
        </div>
        <div class="card">
            <h2>说明</h2>
            <p>保存配置后，DDNS服务将自动启动并监听端口 9876</p>
            <p>访问 http://{{.Host}}:9876 管理DDNS</p>
        </div>
    </div>
</body>
</html>`,

	"tunnel": `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>协议封装 - WireGuard 伴侣</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background: #f5f5f5; margin: 0; }
        .header { background: #2563eb; color: white; padding: 1rem 2rem; }
        .nav { background: white; padding: 0 2rem; border-bottom: 1px solid #e5e7eb; }
        .nav ul { display: flex; list-style: none; gap: 2rem; margin: 0; }
        .nav a { display: block; padding: 1rem 0; color: #6b7280; text-decoration: none; border-bottom: 2px solid transparent; }
        .nav a:hover, .nav a.active { color: #2563eb; border-color: #2563eb; }
        .container { padding: 2rem; max-width: 1200px; margin: 0 auto; }
        .card { background: white; border-radius: 8px; padding: 1.5rem; margin-bottom: 1.5rem; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
        .card h2 { margin-bottom: 1rem; color: #1f2937; }
        .btn { display: inline-block; padding: 0.5rem 1rem; background: #2563eb; color: white; border: none; border-radius: 6px; cursor: pointer; }
        .btn:hover { background: #1d4ed8; }
        .form-group { margin-bottom: 1rem; }
        .form-group label { display: block; margin-bottom: 0.5rem; color: #374151; }
        .form-group input, .form-group select { width: 100%; padding: 0.5rem; border: 1px solid #d1d5db; border-radius: 6px; }
    </style>
</head>
<body>
    <div class="header"><h1>WireGuard 伴侣</h1></div>
    <nav class="nav">
        <ul>
            <li><a href="/">概览</a></li>
            <li><a href="/wireguard">WireGuard</a></li>
            <li><a href="/ddns">DDNS</a></li>
            <li><a href="/tunnel" class="active">协议封装</a></li>
            <li><a href="/monitor">监控</a></li>
        </ul>
    </nav>
    <div class="container">
        <div class="card">
            <h2>wstunnel 配置</h2>
            <form method="POST" action="/api/tunnel/ws/config">
                <div class="form-group">
                    <label>启用</label>
                    <select name="enable">
                        <option value="false">关闭</option>
                        <option value="true">启用</option>
                    </select>
                </div>
                <div class="form-group">
                    <label>本地监听端口</label>
                    <input type="number" name="localPort" value="51820">
                </div>
                <div class="form-group">
                    <label>远程服务器</label>
                    <input type="text" name="remote" placeholder="example.com:443">
                </div>
                <button type="submit" class="btn">保存配置</button>
            </form>
        </div>
        <div class="card">
            <h2>udp2raw 配置</h2>
            <form method="POST" action="/api/tunnel/udp2raw/config">
                <div class="form-group">
                    <label>启用</label>
                    <select name="enable">
                        <option value="false">关闭</option>
                        <option value="true">启用</option>
                    </select>
                </div>
                <div class="form-group">
                    <label>本地端口</label>
                    <input type="number" name="localPort" value="51820">
                </div>
                <div class="form-group">
                    <label>远程服务器</label>
                    <input type="text" name="remote" placeholder="example.com:443">
                </div>
                <div class="form-group">
                    <label>密码</label>
                    <input type="password" name="password" placeholder="密码">
                </div>
                <button type="submit" class="btn">保存配置</button>
            </form>
        </div>
    </div>
</body>
</html>`,

	"monitor": `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>监控 - WireGuard 伴侣</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background: #f5f5f5; margin: 0; }
        .header { background: #2563eb; color: white; padding: 1rem 2rem; }
        .nav { background: white; padding: 0 2rem; border-bottom: 1px solid #e5e7eb; }
        .nav ul { display: flex; list-style: none; gap: 2rem; margin: 0; }
        .nav a { display: block; padding: 1rem 0; color: #6b7280; text-decoration: none; border-bottom: 2px solid transparent; }
        .nav a:hover, .nav a.active { color: #2563eb; border-color: #2563eb; }
        .container { padding: 2rem; max-width: 1200px; margin: 0 auto; }
        .card { background: white; border-radius: 8px; padding: 1.5rem; margin-bottom: 1.5rem; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
        .card h2 { margin-bottom: 1rem; color: #1f2937; }
        .status-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 1rem; }
        .status-item { padding: 1.5rem; background: #f9fafb; border-radius: 8px; text-align: center; }
        .status-item .name { font-size: 1.25rem; font-weight: bold; margin-bottom: 0.5rem; }
        .status-item .info { color: #6b7280; margin-bottom: 0.5rem; }
        .status-item .state { font-size: 1.5rem; font-weight: bold; }
        .status-item .state.running { color: #10b981; }
        .status-item .state.stopped { color: #6b7280; }
        .status-item .state.error { color: #ef4444; }
        .btn { display: inline-block; padding: 0.5rem 1rem; background: #2563eb; color: white; border: none; border-radius: 6px; cursor: pointer; text-decoration: none; }
        .btn-danger { background: #ef4444; }
    </style>
</head>
<body>
    <div class="header"><h1>WireGuard 伴侣</h1></div>
    <nav class="nav">
        <ul>
            <li><a href="/">概览</a></li>
            <li><a href="/wireguard">WireGuard</a></li>
            <li><a href="/ddns">DDNS</a></li>
            <li><a href="/tunnel">协议封装</a></li>
            <li><a href="/monitor" class="active">监控</a></li>
        </ul>
    </nav>
    <div class="container">
        <div class="card">
            <h2>模块运行状态</h2>
            <div class="status-grid">
                <div class="status-item">
                    <div class="name">WireGuard</div>
                    <div class="info">隧道数: {{.Status.Wireguard.Count}}</div>
                    {{if .Status.Wireguard.Count}}
                    <div class="state running">运行中</div>
                    {{else}}
                    <div class="state stopped">无隧道</div>
                    {{end}}
                </div>
                <div class="status-item">
                    <div class="name">wstunnel</div>
                    <div class="info">WebSocket隧道</div>
                    {{if .Status.WSTunnel.Running}}
                    <div class="state running">运行中</div>
                    {{else}}
                    <div class="state stopped">未运行</div>
                    {{end}}
                </div>
                <div class="status-item">
                    <div class="name">udp2raw</div>
                    <div class="info">UDP伪TCP</div>
                    {{if .Status.UDP2Raw.Running}}
                    <div class="state running">运行中</div>
                    {{else}}
                    <div class="state stopped">未运行</div>
                    {{end}}
                </div>
            </div>
        </div>
    </div>
</body>
</html>`,
}

// Server Web服务器
type Server struct {
	port     int
	mux      *http.ServeMux
	tmpls    map[string]*template.Template
	tunnelMgr *tunnel.TunnelMgr
}

// NewServer 创建Web服务器
func NewServer(port int) *Server {
	s := &Server{
		port:      port,
		mux:       http.NewServeMux(),
		tmpls:     make(map[string]*template.Template),
		tunnelMgr: tunnel.NewTunnelMgr(),
	}
	
	for name, html := range htmlTemplates {
		s.tmpls[name] = template.Must(template.New(name).Parse(html))
	}
	
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.mux.HandleFunc("/", s.handleIndex)
	s.mux.HandleFunc("/wireguard", s.handleWireGuard)
	s.mux.HandleFunc("/ddns", s.handleDDNS)
	s.mux.HandleFunc("/tunnel", s.handleTunnel)
	s.mux.HandleFunc("/monitor", s.handleMonitor)
	s.mux.HandleFunc("/api/status", s.handleAPIStatus)
	s.mux.HandleFunc("/api/wireguard/create", s.handleWireGuardCreate)
	s.mux.HandleFunc("/api/wireguard/delete/", s.handleWireGuardDelete)
	s.mux.HandleFunc("/api/ddns/config", s.handleDDNSConfig)
	s.mux.HandleFunc("/api/tunnel/ws/config", s.handleTunnelWSConfig)
	s.mux.HandleFunc("/api/tunnel/udp2raw/config", s.handleTunnelUDP2RawConfig)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	status := s.tunnelMgr.GetModuleStatus()
	s.tmpls["index"].Execute(w, map[string]interface{}{
		"Status": status,
	})
}

func (s *Server) handleWireGuard(w http.ResponseWriter, r *http.Request) {
	tunnels, _ := s.tunnelMgr.ListWireGuardTunnels()
	s.tmpls["wireguard"].Execute(w, map[string]interface{}{
		"Tunnels": tunnels,
	})
}

func (s *Server) handleDDNS(w http.ResponseWriter, r *http.Request) {
	s.tmpls["ddns"].Execute(w, map[string]interface{}{
		"Host": r.Host,
	})
}

func (s *Server) handleTunnel(w http.ResponseWriter, r *http.Request) {
	s.tmpls["tunnel"].Execute(w, nil)
}

func (s *Server) handleMonitor(w http.ResponseWriter, r *http.Request) {
	status := s.tunnelMgr.GetModuleStatus()
	s.tmpls["monitor"].Execute(w, map[string]interface{}{
		"Status": status,
	})
}

func (s *Server) handleAPIStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.tunnelMgr.GetModuleStatus())
}

func (s *Server) handleWireGuardCreate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.Form.Get("name")
	address := r.Form.Get("address")
	port := 51820
	fmt.Sscanf(r.Form.Get("port"), "%d", &port)
	
	err := s.tunnelMgr.CreateWireGuardTunnel(name, address, port)
	if err != nil {
		fmt.Fprintf(w, "创建失败: %v", err)
		return
	}
	
	http.Redirect(w, r, "/wireguard", http.StatusFound)
}

func (s *Server) handleWireGuardDelete(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/api/wireguard/delete/"):]
	err := s.tunnelMgr.DeleteWireGuardTunnel(name)
	if err != nil {
		fmt.Fprintf(w, "删除失败: %v", err)
		return
	}
	http.Redirect(w, r, "/wireguard", http.StatusFound)
}

func (s *Server) handleDDNSConfig(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Fprintf(w, "DDNS配置已保存 (暂未实现)")
}

func (s *Server) handleTunnelWSConfig(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Fprintf(w, "wstunnel配置已保存 (暂未实现)")
}

func (s *Server) handleTunnelUDP2RawConfig(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Fprintf(w, "udp2raw配置已保存 (暂未实现)")
}

// Run 运行服务器
func (s *Server) Run() error {
	addr := fmt.Sprintf(":%d", s.port)
	fmt.Printf("Web界面: http://localhost%s\n", addr)
	return http.ListenAndServe(addr, s.mux)
}
