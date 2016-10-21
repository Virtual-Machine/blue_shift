package socket

import "log"

type Hub struct {
	clients map[*Client]bool
	broadcast chan []byte
	request chan []byte
	register chan *Client
	unregister chan *Client
	mode string
}

func NewHub(modeStr string) *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		request:  	make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		mode: 		modeStr,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.connect(client)
		case client := <-h.unregister:
			h.disconnect(client)
		case request := <-h.request:
			h.intakeRequest(request)
		case message := <-h.broadcast:
			h.sendBroadcast(message)
		}
	}
}

func (h *Hub) connect(client *Client){
	if h.mode == "Debug" {
		log.Println("Connecting socket @", client.conn.RemoteAddr())
	}
	h.clients[client] = true
}

func (h *Hub) disconnect(client *Client){
	if h.mode == "Debug" {
		log.Println("Disconnecting socket @", client.conn.RemoteAddr())
	}
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)
	}
}

func (h *Hub) intakeRequest(request []byte){
	if h.mode == "Debug" {
		log.Println("Got request: " + string(request))
	}
	// TODO process intake request properly
	if string(request) == "BROADCAST" {
		h.sendBroadcast(request)
	}
}

func (h *Hub) sendBroadcast(message []byte){
	if h.mode == "Debug" {
		log.Println("Broadcasting: " + string(message))
	}
	for client := range h.clients {
		select {
		case client.send <- message:
			if h.mode == "Debug" {
				log.Println(client.conn.RemoteAddr(), "received: ", string(message))
			}
		default:
			if h.mode == "Debug" {
				log.Println("Default Condition -- No message into client.send -- see line 54 hub.go")
			}
			close(client.send)
			delete(h.clients, client)
		}
	}
}