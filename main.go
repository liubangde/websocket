package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// 统计接收到的是第几条数据
var count int64

// wbsoket 连接小工具
func main() {
	var host string

	var headerString string

	flag.StringVar(&host, "h", "127.0.0.1:8000", "ws server address:port")
	flag.StringVar(&headerString, "r", `{"Upgrade":"websocket", "Accept-Language": "zh-CN,zh;q=0.9"}`, "set request header, Multiple parameters in the request header can be separated by commas")
	flag.Parse()

	// 进行socket 连接
	log.Printf("start con.. %s", host)

	// 获取上下文
	ctx, cancel := context.WithCancel(context.Background())

	// 转换header 格式
	var header map[string]interface{}
	err := json.Unmarshal([]byte(headerString), &header)
	if err != nil {
		panic(err)
	}
	headerH := make(http.Header)
	for k, v := range header {
		headerH[k] = []string{v.(string)}
	}

	// 进行连接
	con(ctx, host, headerH)

	// 进程结束
	cancel()

}

func con(ctx context.Context, host string, header http.Header) {
	conn, _, err := websocket.DefaultDialer.Dial(host, header)
	if err != nil {
		log.Println("连接失败", err)
		return
	}
	defer func() {
		conn.Close()
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	}()
	// 监听读取过来的消息
	go func() {
		for {
			_, p, err := conn.ReadMessage()
			log.Println("read server message is：", bytes.NewBuffer(p).String())
			if err != nil {
				<-ctx.Done()
				log.Println("read server message err: ", err)
				return
			}
		}
	}()

}
