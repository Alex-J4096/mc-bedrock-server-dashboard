package server

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
)

func StartServer(path string) (*exec.Cmd, io.ReadCloser, io.ReadCloser, error) {
	var cmd *exec.Cmd
	absPath, err := filepath.Abs(path)

	fmt.Println("path:", "Starting server at ", absPath)
	if runtime.GOOS == "windows" {
		cmd = exec.Command(absPath + "\\bedrock_server.exe")
	} else if runtime.GOOS == "linux" {

	} else {
		log.Fatal("Unsupported platform", runtime.GOOS)
	}

	// 创建输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create stdout pipe: %v", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create stderr pipe: %v", err)
	}

	err = cmd.Start()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to start server: %v", err)
	}

	// 异步读取错误输出并打印到控制台
	go func() {
		errReader := bufio.NewReader(stderr)
		for {
			line, err := errReader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					log.Println("Error reading stderr:", err)
				}
				return
			}
			log.Println("[Minecraft Error]", line)
		}
	}()

	fmt.Println("Minecraft server started")
	return cmd, stdout, stderr, nil
}

func StreamLogs(c *gin.Context, stdout io.ReadCloser) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// 允许所有来源
			return true
		},
	}

	// 升级http请求为webSocket
	connection, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer connection.Close()
	defer stdout.Close() // 关闭管道

	reader := bufio.NewReader(stdout)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading stdout:", err)
			}
			break
		}
		err = connection.WriteMessage(websocket.TextMessage, []byte(line))
		if err != nil {
			log.Println("Error writing to websocket:", err)
			break
		}
	}
}
