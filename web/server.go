package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"wireguard-partner/internal/config"
	"wireguard-partner/internal/tunnel"
	"wireguard-partner/internal/updater"
)

// Server Web服务器
type Server struct {
	port    int
	cfg     *config.Config
	tunnelMgr *tunnel.TunnelMgr
	updater  *updater.Checker
	mux     *http.ServeMux
}

// NewServer 创建Web服务器
func NewServer(port int) *Server {
	return &Server{
		port:      port,
		tunnelMgr: tunnel.NewTunnelMgr(),
		updater:   updater.NewChecker(),
		mux:       http.NewServeMux(),
	}
}

// Run 运行服务器
func (s *Server) Run() error {
	s.setupRoutes()
	addr := fmt.Sprintf(":%d", s.port)
	fmt.Printf("Web服务器启动: http://localhost%s\n", addr)
	return http.ListenAndServe(addr, s.mux)
}

func (s *Server) setupRoutes() {
	s.mux.HandleFunc("/", s.handleIndex)
	s.mux.HandleFunc("/api/status", s.handleStatus)
	s.mux.HandleFunc("/api/wireguard", s.handleWireGuard)
	s.mux.HandleFunc("/api/tunnel", s.handleTunnel)
	s.mux.HandleFunc("/api/ddns", s.handleDDNS)
	s.mux.HandleFunc("/api/version", s.handleVersion)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>WireGuard 伴侣</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; max-width: 800px; margin: 50px auto; padding: 20px; }
        h1 { color: #333; }
        .card { border: 1px solid #ddd; border-radius: 8px; padding: 20px; margin: 20px 0; }
        .status { display: inline-block; padding: 4px 12px; border-radius: 4px; font-size: 12px; }
        .status.ok { background: #d4edda; color: #155724; }
        .status.error { background: #f8d7da; color: #721c24; }
        .btn { background: #007bff; color: white; border: none; padding: 8px 16px; border-radius: 4px; cursor: pointer; }
    </style>
</head>
<body>
    <h1>WireGuard 伴侣</h1>
    <div id="app"></div>
    <script>
        fetch('/api/status').then(r=>r.json()).then(data => {
            document.getElementById('app').innerHTML = JSON.stringify(data, null, 2);
        });
    </script>
</body>
</html>`
	w.Write([]byte(html))
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// 检查各组件状态
	wgInstalled, _ := s.updater.CheckWireGuard()
	wsInstalled, _ := s.tunnelMgr.CheckInstalled(tunnel.ProtocolWSTunnel)
	udp2rawInstalled, _ := s.tunnelMgr.CheckInstalled(tunnel.ProtocolUDP2Raw)
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"wireguard": wgInstalled != nil,
		"wstunnel": wsInstalled,
		"udp2raw": udp2rawInstalled,
	})
}

func (s *Server) handleWireGuard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleTunnel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	if r.Method == "GET" {
		// 获取状态
		wsRunning, _ := s.tunnelMgr.GetProcessStatus("wstunnel")
		udpRunning, _ := s.tunnelMgr.GetProcessStatus("udp2raw")
		
		json.NewEncoder(w).Encode(map[string]bool{
			"wstunnel": wsRunning,
			"udp2raw": udpRunning,
		})
	}
}

func (s *Server) handleDDNS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	versions := make(map[string]interface{})
	
	if v, err := s.updater.CheckWireGuard(); err == nil {
		versions["wireguard"] = v
	}
	if v, err := s.updater.CheckWSTunnel(); err == nil {
		versions["wstunnel"] = v
	}
	if v, err := s.updater.CheckUDP2Raw(); err == nil {
		versions["udp2raw"] = v
	}
	
	json.NewEncoder(w).Encode(versions)
}
