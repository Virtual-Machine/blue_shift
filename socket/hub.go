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

func (h *Hub) intakeRequest(packet *Packet){
	if h.mode == "Debug" {
		log.Println("Got packet: " + packet.Data + " From: " + packet.Id)
	}
	var req Request
	if err := json.Unmarshal([]byte(packet.Data), &req); err != nil {
	    log.Println("REQUEST ERROR!!! : ", err)
	    return
    }
    // MARKER Server -> Socket server received data from client.
    if req.Type == "Click" {
		if req.X < 0 || req.Y < 0 || req.X >= 60 || req.Y >= 40 {
			packet.Data = "{\"error\":\"Click is out of bounds\"}"
			h.sendMessage( packet )
			return
		}
		validMove, err := engine.GameInstance.ProcessClick(packet.Id, req.X, req.Y)
		if validMove {
			packet.Data = string(engine.GameInstance.GetData(packet.Id, "MapData"))
			h.sendBroadcast( packet )
		} else {
			packet.Data = "{\"error\":\"" + err.Error() + "\"}"
			h.sendMessage( packet )
		}
	}
	if req.Type == "MapData" {
		packet.Data = string(engine.GameInstance.GetData(packet.Id, "MapData"))
		h.sendBroadcast( packet )
	}
	if req.Type == "ChatMessage" {
		packet.Data = "{\"message\":\"" + req.Message + "\", \"author\": \"" + packet.Id + "\"}"
		h.sendBroadcast( packet )
	}
}

func (h *Hub) sendMessage(packet *Packet){
	if h.mode == "Debug" {
		log.Println("Sending message to : " + packet.Id)
	}
	for client := range h.clients {
		if client.Tag == packet.Id {
			select {
			case client.send <- []byte(packet.Data):
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
}

func (h *Hub) sendBroadcast(packet *Packet){
	if h.mode == "Debug" {
		log.Println("Broadcasting from: " + packet.Id)
	}
	for client := range h.clients {
		select {
		case client.send <- []byte(packet.Data):
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