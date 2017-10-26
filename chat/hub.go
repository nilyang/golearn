package main

type Hub struct {
	// Registered clients.
	// 已注册的所有客户端连接
	clients map[*Client]bool

	// Inbound messages from the clients.
	// 用于广播给所有连接客户端的消息
	broadcast chan []byte

	// Register requests from the clients
	register chan *Client

	// Unregister requests from the clients
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

//
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			// 注册客户端连接对象
			h.clients[client] = true
		case client := <-h.unregister:
			// 释放客户端连接（客户端退出）
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			// 广播收到的消息到每个客户端
			for client := range h.clients {
				select {
				case client.send <- message:
					// 广播到每一个客户端的发消息channel
				default:
					// 若广播失败，则说明客户端已退出，应关闭客户端channel，并移除客户端连接
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
