package models

import (
	"bytes"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const (
	// 将信息写入peer端所允许的最大时间
	writeWait = 10 * time.Second
	// 读取peer端返回信息所允许的最大时间
	pongWait = 60 * time.Second
	// 发送pings到peer端所允许的时间，必须比小于pongWait
	pingPeriod = (pongWait * 9) / 10

	// 消息体最大字节数
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// 将普通http升级成websocket协议的必备神器
var upgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: true,
}

type Client struct {
	hub *Hub

	// websocket 连接
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	// 发消息缓冲
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
// readPump函数将websocket连接中的消息抽取到hub中
//
// The application runs readPump in a per-connection goroutine.
// 应用程序在每一个连接goroutine中运行readPump。
// The application ensures that there is at most one reader on a connection by executing all reads from this goroutine.
// 应用程序通过执行当前goroutine上的所有读取器（reads），以保证每个连接上至少有一个读取器（reader）。
func (c *Client) readPump() {
	defer func() {
		// goroutine 结束时，将当前客户端对象注销掉，并关闭当前连接
		c.hub.unregister <- c
		c.conn.Close()
	}()

	// 设置每次读取消息的大小，超过大小则会给客户端报错，并断开连接
	c.conn.SetReadLimit(maxMessageSize)
	// 设置读取超时时间
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	// 设置响应回调函数
	c.conn.SetPongHandler(func(appData string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// 无限循环 读取连接消息，广播到每个客户端
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			// 出错跳出循环(会关闭当前连接goroutine)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		// 格式化消息，去掉回车换行符，以空格替换
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		// 通过channel 广播消息（channel chan []byte）
		c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
// writePump 将hub中的信息抽出并写入ebsocket连接
//
// A goroutine running writePump is started for each connection.
// The application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send: // 读取广播消息
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// NextWriter 的行为：刷新之前的数据到连接客户端，重新开启新的连接写数据
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				// 失败则返回？
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				// 失败则返回？
				return
			}

		case <-ticker.C: // 读取ticker 时钟io
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// 每个请求处理一次，hub是全局的，会有内部修改所以传入指针，请求是当前请求
func ServeWebsocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

func SendMessage(hub *Hub, msg string) {
	msgBytes := []byte(msg)
	hub.broadcast <- msgBytes
}