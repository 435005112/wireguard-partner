package web

import (
	"fmt"
	"html/template"
	"net/http"
	"wireguard-partner/internal/protocol"
)

// ProtocolServer 协议转换Web服务器
type ProtocolServer struct {
	port      int
	mux       *http.ServeMux
	mgr       *protocol.Manager
	templates map[string]*template.Template
}

// NewProtocolServer 创建协议Web服务器
func NewProtocolServer(port int) *ProtocolServer {
	mgr := protocol.NewManager()
	mgr.Load()

	s := &ProtocolServer{
		port:      port,
		mux:       http.NewServeMux(),
		mgr:       mgr,
		templates: make(map[string]*template.Template),
	}

	s.setupTemplates()
	s.setupRoutes()
	return s
}

func (s *ProtocolServer) setupTemplates() {
	// 加载模板文件
	templates := []string{"list", "edit", "log"}
	for _, name := range templates {
		tmpl := template.Must(template.ParseFiles(fmt.Sprintf("web/templates/%s.html", name)))
		s.templates[name] = tmpl
	}
}

func (s *ProtocolServer) setupRoutes() {
	s.mux.HandleFunc("/", s.handleList)
	s.mux.HandleFunc("/new", s.handleNew)
	s.mux.HandleFunc("/edit", s.handleEdit)
	s.mux.HandleFunc("/save", s.handleSave)
	s.mux.HandleFunc("/delete", s.handleDelete)
	s.mux.HandleFunc("/start", s.handleStart)
	s.mux.HandleFunc("/stop", s.handleStop)
	s.mux.HandleFunc("/restart", s.handleRestart)
	s.mux.HandleFunc("/log", s.handleLog)
}

// =====================================================
// 页面处理函数
// =====================================================

func (s *ProtocolServer) handleList(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	configs := s.mgr.List()
	toolsStatus := s.mgr.InstallTools()

	s.templates["list"].Execute(w, map[string]interface{}{
		"Configs":     configs,
		"ToolsStatus": toolsStatus,
	})
}

func (s *ProtocolServer) handleNew(w http.ResponseWriter, r *http.Request) {
	s.templates["edit"].Execute(w, map[string]interface{}{
		"EditMode":   false,
		"Config":     protocol.ProtocolConfig{},
		"ToolStatus": s.mgr.InstallTools(),
	})
}

func (s *ProtocolServer) handleEdit(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	cfg := s.mgr.Get(name)

	if cfg == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	s.templates["edit"].Execute(w, map[string]interface{}{
		"EditMode":   true,
		"Config":     cfg,
		"ToolStatus": s.mgr.InstallTools(),
	})
}

func (s *ProtocolServer) handleSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	r.ParseForm()

	var cfg protocol.ProtocolConfig
	cfg.Type = protocol.ProtocolType(r.Form.Get("type"))
	cfg.Role = protocol.RoleType(r.Form.Get("role"))

	if cfg.Type == protocol.ProtocolWSTunnel {
		var wsCfg protocol.WSTunnelConfig
		wsCfg.Name = r.Form.Get("name")
		wsCfg.Server = r.Form.Get("server")
		fmt.Sscanf(r.Form.Get("port"), "%d", &wsCfg.Port)
		wsCfg.LocalTun = r.Form.Get("localTun")
		wsCfg.UseTLS = r.Form.Get("useTLS") == "on"
		wsCfg.SNI = r.Form.Get("sni")
		wsCfg.PathPrefix = r.Form.Get("pathPrefix")
		wsCfg.LogLevel = r.Form.Get("logLevel")
		wsCfg.HTTPProxy = r.Form.Get("httpProxy")
		wsCfg.HTTPProxyLogin = r.Form.Get("httpProxyLogin")
		wsCfg.HTTPProxyPass = r.Form.Get("httpProxyPass")
		wsCfg.TLSConf.CertFile = r.Form.Get("tlsCert")
		wsCfg.TLSConf.KeyFile = r.Form.Get("tlsKey")
		cfg.WSTunnel = &wsCfg
	} else if cfg.Type == protocol.ProtocolUDP2Raw {
		var udpCfg protocol.UDP2RawConfig
		udpCfg.Name = r.Form.Get("name")
		udpCfg.Server = r.Form.Get("server")
		fmt.Sscanf(r.Form.Get("serverPort"), "%d", &udpCfg.ServerPort)
		udpCfg.LocalAddr = r.Form.Get("localAddr")
		udpCfg.RemoteAddr = r.Form.Get("remoteAddr")
		udpCfg.Key = r.Form.Get("key")
		udpCfg.Mode = r.Form.Get("mode")
		udpCfg.Cipher = r.Form.Get("cipher")
		udpCfg.Auth = r.Form.Get("auth")
		udpCfg.AutoIPT = r.Form.Get("autoIPT") == "on"
		fmt.Sscanf(r.Form.Get("logLevel"), "%d", &udpCfg.LogLevel)
		udpCfg.SourceIP = r.Form.Get("sourceIP")
		udpCfg.DisableReplay = r.Form.Get("disableReplay") == "on"
		udpCfg.Device = r.Form.Get("device")
		fmt.Sscanf(r.Form.Get("sockBuf"), "%d", &udpCfg.SockBuf)
		cfg.UDP2Raw = &udpCfg
	}

	// 检查是否已存在
	existing := s.mgr.Get(cfg.GetName())
	if existing != nil {
		// 更新时保留状态
		cfg.PID = existing.PID
		cfg.Status = existing.Status
	}

	if err := s.mgr.Save(&cfg); err != nil {
		fmt.Fprintf(w, "保存失败: %v", err)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *ProtocolServer) handleDelete(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if err := s.mgr.Delete(name); err != nil {
		fmt.Fprintf(w, "删除失败: %v", err)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *ProtocolServer) handleStart(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if err := s.mgr.Start(name); err != nil {
		fmt.Fprintf(w, "启动失败: %v", err)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *ProtocolServer) handleStop(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if err := s.mgr.Stop(name); err != nil {
		fmt.Fprintf(w, "停止失败: %v", err)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *ProtocolServer) handleRestart(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if err := s.mgr.Restart(name); err != nil {
		fmt.Fprintf(w, "重启失败: %v", err)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *ProtocolServer) handleLog(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	cfg := s.mgr.Get(name)

	if cfg == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	lines := 100
	fmt.Sscanf(r.URL.Query().Get("lines"), "%d", &lines)

	log, _ := s.mgr.GetLog(name, lines)

	s.templates["log"].Execute(w, map[string]interface{}{
		"Name":   name,
		"Log":    log,
		"Config": cfg,
	})
}

// Run 运行服务器
func (s *ProtocolServer) Run() error {
	addr := fmt.Sprintf(":%d", s.port)
	fmt.Printf("协议转换模块: http://localhost%s\n", addr)
	return http.ListenAndServe(addr, s.mux)
}
