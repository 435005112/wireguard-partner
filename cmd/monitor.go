package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"text/template"
)

type Status struct {
	Docker      bool
	WG Easy    bool
	DDNSGo     bool
	WSTunnel   bool
	UDP2Raw    bool
	WireGuard  bool
}

var tmpl = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>WireGuard伴侣 - 统一监控</title>
    <meta http-equiv="refresh" content="30">
    <style>
        body { font-family: -apple-system, sans-serif; background: #f5f5f5; margin: 0; }
        .header { background: #2563eb; color: white; padding: 1rem 2rem; }
        .container { padding: 2rem; max-width: 800px; margin: 0 auto; }
        .card { background: white; border-radius: 8px; padding: 1.5rem; margin-bottom: 1rem; }
        .card h2 { margin: 0 0 1rem; color: #1f2937; }
        .status-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 1rem; }
        .item { padding: 1rem; background: #f9fafb; border-radius: 6px; }
        .item .name { color: #6b7280; font-size: 0.875rem; }
        .item .value { font-size: 1.5rem; font-weight: bold; margin-top: 0.5rem; }
        .item .value.ok { color: #10b981; }
        .item .value.error { color: #ef4444; }
        .btn { display: inline-block; padding: 0.5rem 1rem; background: #2563eb; color: white; text-decoration: none; border-radius: 6px; }
    </style>
</head>
<body>
    <div class="header">
        <h1>WireGuard伴侣 - 统一监控</h1>
    </div>
    <div class="container">
        <div class="card">
            <h2>服务状态</h2>
            <div class="status-grid">
                <div class="item">
                    <div class="name">Docker</div>
                    <div class="value {{if .Status.Docker}}ok{{else}}error{{end}}">
                        {{if .Status.Docker}}✓ 运行中{{else}}✗ 未安装{{end}}
                    </div>
                </div>
                <div class="item">
                    <div class="name">wg-easy</div>
                    <div class="value {{if .Status.WG-Easy}}ok{{else}}error{{end}}">
                        {{if .Status.WG-Easy}}✓ 运行中{{else}}✗ 未运行{{end}}
                    </div>
                </div>
                <div class="item">
                    <div class="name">ddns-go</div>
                    <div class="value {{if .Status.DDNSGo}}ok{{else}}error{{end}}">
                        {{if .Status.DDNSGo}}✓ 运行中{{else}}✗ 未运行{{end}}
                    </div>
                </div>
                <div class="item">
                    <div class="name">wstunnel</div>
                    <div class="value {{if .Status.WSTunnel}}ok{{else}}error{{end}}">
                        {{if .Status.WSTunnel}}✓ 运行中{{else}}✗ 未运行{{end}}
                    </div>
                </div>
                <div class="item">
                    <div class="name">udp2raw</div>
                    <div class="value {{if .Status.UDP2Raw}}ok{{else}}error{{end}}">
                        {{if .Status.UDP2Raw}}✓ 运行中{{else}}✗ 未运行{{end}}
                    </div>
                </div>
                <div class="item">
                    <div class="name">WireGuard</div>
                    <div class="value {{if .Status.WireGuard}}ok{{else}}error{{end}}">
                        {{if .Status.WireGuard}}✓ 有隧道{{else}}✗ 无隧道{{end}}
                    </div>
                </div>
            </div>
        </div>
        <div class="card">
            <h2>快速访问</h2>
            <p>
                <a href="http://{{.Host}}:51821" target="_blank" class="btn">wg-easy 管理</a>
                <a href="http://{{.Host}}:9876" target="_blank" class="btn">ddns-go 管理</a>
            </p>
        </div>
        <div class="card">
            <h2>重新安装服务</h2>
            <p><a href="/install" class="btn">运行安装向导</a></p>
        </div>
    </div>
</body>
</html>`

func checkDocker() bool {
	cmd := exec.Command("docker", "info")
	return cmd.Run() == nil
}

func checkWG-Easy() bool {
	cmd := exec.Command("docker", "ps", "--filter", "name=wg-easy", "--format", "{{.Names}}")
	output, _ := cmd.Output()
	return strings.Contains(string(output), "wg-easy")
}

func checkDDNSGo() bool {
	cmd := exec.Command("pgrep", "-f", "ddns-go")
	return cmd.Run() == nil
}

func checkWSTunnel() bool {
	cmd := exec.Command("pgrep", "-f", "wstunnel")
	return cmd.Run() == nil
}

func checkUDP2Raw() bool {
	cmd := exec.Command("pgrep", "-f", "udp2raw")
	return cmd.Run() == nil
}

func checkWireGuard() bool {
	cmd := exec.Command("wg", "show", "interfaces")
	output, _ := cmd.Output()
	return strings.TrimSpace(string(output)) != ""
}

func getStatus() Status {
	return Status{
		Docker:    checkDocker(),
		WG-Easy:    checkWG-Easy(),
		DDNSGo:   checkDDNSGo(),
		WSTunnel:  checkWSTunnel(),
		UDP2Raw:  checkUDP2Raw(),
		WireGuard: checkWireGuard(),
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	status := getStatus()
	host := strings.Split(r.Host, ":")[0]
	
	t, _ := template.New("status").Parse(tmpl)
	t.Execute(w, map[string]interface{}{
		"Status": status,
		"Host":   host,
	})
}

func main() {
	fmt.Println("WireGuard伴侣 - 统一监控")
	fmt.Println("========================================")
	fmt.Println("访问: http://<ip>:8080")
	fmt.Println("========================================")
	
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
