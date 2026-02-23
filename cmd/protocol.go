package main

import (
	"flag"
	"fmt"
	"log"
	"wireguard-partner/web"
)

func main() {
	port := flag.Int("port", 8080, "HTTP server port")
	flag.Parse()

	fmt.Println("========================================")
	fmt.Println("  协议转换模块 - WireGuard 伴侣")
	fmt.Println("========================================")
	fmt.Println()

	server := web.NewProtocolServer(*port)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
