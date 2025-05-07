package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

func StreamLogs(c *gin.Context) {
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

	for {
		time.Sleep(1 * time.Second)
		err := connection.WriteMessage(websocket.TextMessage, []byte("测试日志消息"))
		if err != nil {
			log.Println(err)
			return
		}
	}
}
