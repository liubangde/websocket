package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

// 统计接收到的是第几条数据
var count int64

// websocket 连接小工具
func main() {
	var host string

	var headerString string

	flag.StringVar(&host, "h", "127.0.0.1:8000", "ws server address:port")
	flag.StringVar(&headerString, "r", `{}`, "set request header, Multiple parameters in the request header can be separated by commas")
	flag.Parse()
	log.Printf("start con.. %s", host)

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

	// 终止连接依据,接收信号
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan struct{})

	// 控制连接粒度
	buf := make(chan struct{}, 1)
	buf <- struct{}{}

	//初始化尝试次数
	count = 0

	// 进行连接
	for {
		if count < 5 {
			<-buf
			go con(host, headerH, interrupt, buf, done)
		} else {
			os.Exit(0)
		}
	}

}

func con(host string, header http.Header, interrupt chan os.Signal, buf chan struct{}, done chan struct{}) {
	conn, _, err := websocket.DefaultDialer.Dial(host, header)
	if err != nil {
		log.Println("conn fail: ", err)
		atomic.AddInt64(&count, 1) // 失败后连接次数+1
		buf <- struct{}{}
		return
	}
	defer func() {
		conn.Close()
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	}()

	log.Println("conn successfully！")

	// 读取
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read server message fail:", err)
				os.Exit(1)
			}
			log.Printf("Read the message returned by the server is : %s", string(message))
		}
	}()

	// 写入
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		// 交互模式，控制台输入
		log.Println("interactive mode，Please enter a message to send :")
		in := bufio.NewReader(os.Stdin)
		str, _, err := in.ReadLine()
		if err != nil {
			log.Println("输入错误：", err)
			return
		}

		select {
		case <-done:
			os.Exit(0)
			return
		case <-ticker.C:
			err := conn.WriteMessage(websocket.TextMessage, str)
			if err != nil {
				log.Println("write server fail: ", err)
				os.Exit(0)
				return
			}
		case <-interrupt:
			log.Println("interrupt")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				os.Exit(0)
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}

}
