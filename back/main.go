package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"mc-bedrock-server-dashboard/back/server"
)

func main() {

	// 命令行参数 Path, 缺省值为当前工作目录
	path := flag.String("path", ".", "path to the minecraft server executable directory")
	flag.Parse()

	cmd, stdout, stderr, err := server.StartServer(*path)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
		return
	}

	defer func() {
		// 在程序退出时关闭服务器进程
		if cmd != nil && cmd.Process != nil {
			fmt.Println("Stopping Minecraft server...")
			if err := cmd.Process.Kill(); err != nil {
				log.Println("Failed to kill server:", err)
			} else {
				cmd.Wait() // 等待进程结束
				fmt.Println("Minecraft server stopped.")
			}
		}
		if stdout != nil {
			stdout.Close()
		}
		if stderr != nil {
			stderr.Close()
		}
	}()

	r := gin.Default()
	r.GET("/", func(context *gin.Context) {
		server.StreamLogs(context, stdout)
	})
	log.Println("已于localhost:8080上启动服务")
	err = r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
