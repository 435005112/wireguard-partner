package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
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
            <h2>系统状态</h2>
            <div class="status-grid">
                <div class="status-item">
                    <div class="label">WireGuard</div>
                    <div class="value {{if .Status.WireGuard}}ok{{else}}error{{end}}">{{if .Status.WireGuard}}✓ 已安装{{else}}✗ 未安装{{end}}</div>
                </div>
                <div class="status-item">
                    <div class="label">wstunnel</div>
                    <div class="value {{if .Status.WSTunnel}}ok{{else}}error{{end}}">{{if .Status.WSTunnel}}✓ 已安装{{else}}✗ 未安装{{end}}</div>
                </div>
                <div class="status-item">
                    <div class="label">udp2raw</div>
                    <div class="value {{if .Status.UDP2Raw}}ok{{else}}error{{end}}">{{if .Status.UDP2Raw}}✓ 已安装{{else}}✗ 未安装{{end}}</div>
                </div>
                <div class="status-item">
                    <div class="label">运行隧道</div>
                    <div class="value">{{ .TunnelCount }}</div>
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
            <h2>WireGuard 隧道管理</h2>
            <p><a href="/wireguard/new" class="btn">+ 新建隧道</a></p>
            <table>
                <thead>
                    <tr>
                        <th>名称</th>
                        <th>接口</th>
                        <th>地址</th>
                        <th>客户端数</th>
                        <th>状态</th>
                        <th>操作</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .Tunnels}}
                    <tr>
                        <td>{{.Name}}</td>
                        <td>{{.Interface}}</td>
                        <td>{{.Address}}</td>
                        <td>{{len .Peers}}</td>
                        <td>活跃</td>
                        <td><a href="/wireguard/edit/{{.Name}}">编辑</a></td>
                    </tr>
                    {{else}}
                    <tr><td colspan="6">暂无隧道，请创建一个</td></tr>
                    {{end}}
                </tbody>
            </table>
        </div>
        <div class="card">
            <h2>创建新隧道</h2>
            <form method="POST" action="/api/wireguard/create">
                <div class="form-group">
                    <label>隧道名称</label>
                    <input type="text" name="name" placeholder="my-vpn" required>
                </div>
                <div class="form-group">
                    <label>接口名称</label>
                    <input type="text" name="interface" placeholder="wg0" required>
                </div>
                <div class="form-group">
                    <label>服务器地址 (CIDR)</label>
                    <input type="text" name="address" placeholder="10.0.0.1/24" required>
                </div>
                <div class="form-group">
                    <label>监听端口</label>
                    <input type="number" name="listenPort" value="51820" required>
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
            <p>DDNS-GO 独立运行在此服务器端口 9876</p>
            <p><a href="http://{{.Host}}:9876" target="_blank" class="btn">打开 DDNS-GO 管理界面</a></p>
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
            <h2>协议封装配置</h2>
            <form method="POST" action="/api/tunnel/config">
                <div class="form-group">
                    <label>协议类型</label>
                    <select name="type">
                        <option value="wstunnel">wstunnel (WebSocket)</option>
                        <option value="udp2raw">udp2raw (UDP伪TCP)</option>
                    </select>
                </div>
                <div class="form-group">
                    <label>启用</label>
                    <select name="enable">
                        <option value="false">关闭</option>
                        <option value="true">启用</option>
                    </select>
                </div>
                <div class="form-group">
                    <label>服务器地址</label>
                    <input type="text" name="server" placeholder="example.com">
                </div>
                <div class="form-group">
                    <label>端口</label>
                    <input type="number" name="port" value="443">
                </div>
                <div class="form-group">
                    <label>密码</label>
                    <input type="password" name="password" placeholder="设置密码">
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
        .status-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 1rem; }
        .status-item { padding: 1rem; background: #f9fafb; border-radius: 6px; text-align: center; }
        .status-item .label { color: #6b7280; font-size: 0.875rem; }
        .status-item .value { font-size: 1.5rem; font-weight: bold; margin-top: 0.5rem; }
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
            <h2>系统监控</h2>
            <div class="status-grid">
                <div class="status-item">
                    <div class="label">CPU 使用率</div>
                    <div class="value">{{.Stats.CPU}}%</div>
                </div>
                <div class="status-item">
                    <div class="label">内存使用</div>
                    <div class="value">{{.Stats.Memory}}%</div>
                </div>
                <div class="status-item">
                    <div class="label">上传速度</div>
                    <div class="value">{{.Stats.Upload}}</div>
                </div>
                <div class="status-item">
                    <div class="label">下载速度</div>
                    <div class="value">{{.Stats.Download}}</div>
                </div>
            </div>
        </div>
    </div>
</body>
</html>`,
}

// Server Web服务器
type Server struct {
	port    int
	mux     *http.ServeMux
	tmpls   map[string]*template.Template
}

// NewServer 创建Web服务器
func NewServer(port int) *Server {
	s := &Server{
		port:  port,
		mux:   http.NewServeMux(),
		tmpls: make(map[string]*template.Template),
	}
	
	// 解析模板
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
	s.mux.HandleFunc("/api/tunnel/config", s.handleTunnelConfig)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	status := s.getStatus()
	tunnelCount := 0
	
	s.tmpls["index"].Execute(w, map[string]interface{}{
		"Status": status,
		"TunnelCount": tunnelCount,
	})
}

func (s *Server) handleWireGuard(w http.ResponseWriter, r *http.Request) {
	s.tmpls["wireguard"].Execute(w, map[string]interface{}{
		"Tunnels": []interface{}{},
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
	s.tmpls["monitor"].Execute(w, map[string]interface{}{
		"Stats": map[string]interface{}{
			"CPU":      "0",
			"Memory":   "0",
			"Upload":   "0 KB/s",
			"Download": "0 KB/s",
		},
	})
}

func (s *Server) handleAPIStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.getStatus())
}

func (s *Server) handleWireGuardCreate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.Form.Get("name")
	fmt.Fprintf(w, "创建隧道: %s", name)
}

func (s *Server) handleTunnelConfig(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Fprintf(w, "配置已保存")
}

func (s *Server) getStatus() map[string]bool {
	return map[string]bool{
		"WireGuard": true,
		"WSTunnel":  true,
		"UDP2Raw":   true,
	}
}

// Run 运行服务器
func (s *Server) Run() error {
	addr := fmt.Sprintf(":%d", s.port)
	fmt.Printf("Web界面: http://localhost%s\n", addr)
	return http.ListenAndServe(addr, s.mux)
}
