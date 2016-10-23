package socket

import (
	"log"
	"encoding/json"

	"../engine"
)

type Hub struct {
	clients map[*Client]bool
	broadcast chan *Packet
	request chan *Packet
	register chan *Client
	unregister chan *Client
	mode string
}

func NewHub(modeStr string) *Hub {
	return &Hub{
		broadcast:  make(chan *Packet),
		request:  	make(chan *Packet),
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
		log.Println("Connecting socket @", client.conn.RemoteAddr(), client.Tag)
	}
	h.clients[client] = true
}

func (h *Hub) disconnect(client *Client){
	if h.mode == "Debug" {
		log.Println("Disconnecting socket @", client.conn.RemoteAddr(), client.Tag)
	}
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)
	}
}

func (h *Hub) intakeRequest(request *Packet){
	if h.mode == "Debug" {
		log.Println("Got request: " + request.Data + " From: " + request.Id)
	}
	var req Request
	if err := json.Unmarshal([]byte(request.Data), &req); err != nil {
	    log.Println("REQUEST ERROR!!! : ", err)
	    return
    }
    if req.Type == "Click" {
    	validMove := engine.GameInstance.ProcessClick(request.Id, req.X, req.Y)
    	if validMove {
    		// 	h.sendBroadcast( newMapData )
		} else {
			// notify client that their selected move was rejected
		}
    }
	if req.Type == "MapData" {
		mapData := engine.GameInstance.GetData(request.Id, req.Type)
		log.Println(mapData[0])
	}
}

func (h *Hub) sendBroadcast(message *Packet){
	if h.mode == "Debug" {
		log.Println("Broadcasting: " + message.Data + " From: " + message.Id)
	}
	for client := range h.clients {
		select {
		case client.send <- []byte(message.Data):
			if h.mode == "Debug" {
				log.Println(client.conn.RemoteAddr(), "received: ", message.Data)
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