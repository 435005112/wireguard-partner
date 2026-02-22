package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
)

// Builder 打包构建器
type Builder struct {
	platform string
	arch    string
}

// NewBuilder 创建构建器
func NewBuilder() *Builder {
	return &Builder{
		platform: runtime.GOOS,
		arch:    runtime.GOARCH,
	}
}

// BuildWireGuard 编译WireGuard (需要目标平台交叉编译)
func (b *Builder) BuildWireGuard() error {
	fmt.Println("WireGuard 需要目标平台单独安装")
	return nil
}

// InstallDDNSGo 安装DDNS-GO
func (b *Builder) InstallDDNSGo() error {
	// 下载DDNS-GO
	url := "https://github.com/jeessy2/ddns-go/releases/latest"
	fmt.Println("请手动安装DDNS-GO:", url)
	return nil
}

// BuildMain 编译主程序
func (b *Builder) BuildMain(output string) error {
	cmd := exec.Command("go", "build", "-o", output, "./cmd/main.go")
	cmd.Dir = "/root/.openclaw/workspace/wireguard-partner"
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("编译失败: %v\n%s", err, output)
	}
	
	fmt.Printf("编译成功: %s\n", output)
	return nil
}

// PackageForN1 为N1打包
func (b *Builder) PackageForN1() error {
	// N1 是 ARM64 Linux
	log.Println("为斐讯N1 (ARM64) 打包...")
	
	// 设置目标平台
	os.Setenv("GOOS", "linux")
	os.Setenv("GOARCH", "arm64")
	
	// 编译
	return b.BuildMain("wireguard-partner-n1")
}

// PackageForPi 为树莓派打包
func (b *Builder) PackageForPi() error {
	// 树莓派 3/4 是 ARM64
	log.Println("为树莓派 (ARM64) 打包...")
	
	os.Setenv("GOOS", "linux")
	os.Setenv("GOARCH", "arm64")
	
	return b.BuildMain("wireguard-partner-pi")
}

// PackageForX86 为x86打包
func (b *Builder) PackageForX86() error {
	log.Println("为x86_64 Linux 打包...")
	
	os.Setenv("GOOS", "linux")
	os.Setenv("GOARCH", "amd64")
	
	return b.BuildMain("wireguard-partner-x86")
}
