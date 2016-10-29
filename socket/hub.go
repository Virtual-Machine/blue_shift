package socket

import (
	"log"
	"encoding/json"
	"strconv"

	"../engine"
	"../login"
)

type Hub struct {
	clients map[*Client]bool
	broadcast chan *Packet
	request chan *Packet
	register chan *Client
	unregister chan *Client
	mode string
	users *login.UserList
	count int
}

func NewHub(modeStr string, userList *login.UserList) *Hub {
	return &Hub{
		broadcast:  make(chan *Packet),
		request:  	make(chan *Packet),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		mode: 		modeStr,
		users:		userList,
		count: 		0,
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
	h.count++
	for i, v := range h.users.List {
		if v.Name == client.Tag {
			h.users.List[i].Status = "Online"
			var pack Packet
			pack.Id = client.Tag
			pack.Data = "{\"count\":" + strconv.Itoa(h.count) + "}"
			h.sendBroadcast(&pack)
			return
		}
	}
}

func (h *Hub) disconnect(client *Client){
	if h.mode == "Debug" {
		log.Println("Disconnecting socket @", client.conn.RemoteAddr(), client.Tag)
	}
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)
	}
	h.count--
	for i, v := range h.users.List {
		if v.Name == client.Tag {
			h.users.List[i].Status = "Offline"
			var pack Packet
			pack.Id = client.Tag
			pack.Data = "{\"count\": " + strconv.Itoa(h.count) + "}"
			h.sendBroadcast(&pack)
			return
		}
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
    		request.Data = string(engine.GameInstance.GetData(request.Id, "MapData"))
    		h.sendBroadcast( request )
		} else {
			// notify client that their selected move was rejected
		}
    }
	if req.Type == "MapData" {
		request.Data = string(engine.GameInstance.GetData(request.Id, "MapData"))
		h.sendBroadcast( request )
	}
}

func (h *Hub) sendBroadcast(message *Packet){
	if h.mode == "Debug" {
		log.Println("Broadcasting from: " + message.Id)
	}
	for client := range h.clients {
		select {
		case client.send <- []byte(message.Data):
			if h.mode == "Debug" {
				log.Println(client.conn.RemoteAddr(), client.Tag, "received data")
			}
		default:
			if h.mode == "Debug" {
				log.Println("Default Condition -- No message into client.send")
			}
			close(client.send)
			delete(h.clients, client)
		}
	}
}