package socket

import (
	"log"
	"encoding/json"

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
	for i, v := range h.users.SafeList {
		if v.Name == client.Tag {
			// MARKER Server -> Client connected.
			h.users.SafeList[i].Status = "Online"
			h.users.SafeList[i].Connections++
			var pack Packet
			pack.Id = client.Tag
			data, _ := json.Marshal(h.users.SafeList)
			pack.Data = "{\"user_list\": " + string(data) + "}"
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
	for i, v := range h.users.SafeList {
		if v.Name == client.Tag {
			// MARKER Server -> Client disconnected.
			h.users.SafeList[i].Connections--
			if h.users.SafeList[i].Connections <= 0 {
				h.users.SafeList[i].Connections = 0
				h.users.SafeList[i].Status = "Offline"
			}
			var pack Packet
			pack.Id = client.Tag
			data, _ := json.Marshal(h.users.SafeList)
			pack.Data = "{\"user_list\": " + string(data) + "}"
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
    // MARKER Server -> Socket server received data from client.
    if req.Type == "Click" {
    	validMove := engine.GameInstance.ProcessClick(request.Id, req.X, req.Y)
    	if validMove {
    		request.Data = string(engine.GameInstance.GetData(request.Id, "MapData"))
    		h.sendBroadcast( request )
		} else {
			// TODO Notify client that their selected move was rejected by engine
		}
    }
	if req.Type == "MapData" {
		request.Data = string(engine.GameInstance.GetData(request.Id, "MapData"))
		h.sendBroadcast( request )
	}
	if req.Type == "ChatMessage" {
		request.Data = "{\"message\":\"" + req.Message + "\", \"author\": \"" + request.Id + "\"}"
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