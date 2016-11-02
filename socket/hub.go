package socket

import (
	"encoding/json"
	"log"

	"../engine"
	"../login"
)

// Hub is the default socker server struct
type Hub struct {
	clients    map[*client]bool
	broadcast  chan *packet
	request    chan *packet
	register   chan *client
	unregister chan *client
	mode       string
	users      *login.UserList
}

// NewHub provides a pointer to a new default socket server
func NewHub(modeStr string, userList *login.UserList) *Hub {
	return &Hub{
		broadcast:  make(chan *packet),
		request:    make(chan *packet),
		register:   make(chan *client),
		unregister: make(chan *client),
		clients:    make(map[*client]bool),
		mode:       modeStr,
		users:      userList,
	}
}

// Run initiates channels for actions
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

func (h *Hub) connect(client *client) {
	if h.mode == "Debug" {
		log.Println("Connecting socket @", client.conn.RemoteAddr(), client.Tag)
	}
	h.clients[client] = true
	for i, v := range h.users.SafeList {
		if v.Name == client.Tag {
			// MARKER Server -> Client connected.
			// TODO implement smart client connections to establish first player
			h.users.SafeList[i].Status = "Online"
			h.users.SafeList[i].Connections++
			var pack packet
			pack.ID = client.Tag
			data, _ := json.Marshal(h.users.SafeList)
			pack.Data = "{\"user_list\": " + string(data) + "}"
			h.sendBroadcast(&pack)
			return
		}
	}
}

func (h *Hub) disconnect(client *client) {
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
			// TODO implement smart client connections to establish first player on disconnect as well
			h.users.SafeList[i].Connections--
			if h.users.SafeList[i].Connections <= 0 {
				h.users.SafeList[i].Connections = 0
				h.users.SafeList[i].Status = "Offline"
			}
			var pack packet
			pack.ID = client.Tag
			data, _ := json.Marshal(h.users.SafeList)
			pack.Data = "{\"user_list\": " + string(data) + "}"
			h.sendBroadcast(&pack)
			return
		}
	}
}

func (h *Hub) intakeRequest(packet *packet) {
	if h.mode == "Debug" {
		log.Println("Got packet: " + packet.Data + " From: " + packet.ID)
	}
	var req request
	if err := json.Unmarshal([]byte(packet.Data), &req); err != nil {
		log.Println("REQUEST ERROR!!! : ", err)
		return
	}
	// MARKER Server -> Socket server received data from client.
	if req.Type == "Click" {
		if req.X < 0 || req.Y < 0 || req.X >= 60 || req.Y >= 40 {
			packet.Data = "{\"error\":\"Click is out of bounds\"}"
			h.sendMessage(packet)
			return
		}
		validMove, err := engine.GameInstance.ProcessClick(packet.ID, req.X, req.Y)
		if validMove {
			packet.Data = string(engine.GameInstance.GetData(packet.ID, "MapData"))
			h.sendBroadcast(packet)
		} else {
			packet.Data = "{\"error\":\"" + err.Error() + "\"}"
			h.sendMessage(packet)
		}
	}
	if req.Type == "MapData" {
		packet.Data = string(engine.GameInstance.GetData(packet.ID, "MapData"))
		h.sendBroadcast(packet)
	}
	if req.Type == "ChatMessage" {
		packet.Data = "{\"message\":\"" + req.Message + "\", \"author\": \"" + packet.ID + "\"}"
		h.sendBroadcast(packet)
	}
}

func (h *Hub) sendMessage(packet *packet) {
	if h.mode == "Debug" {
		log.Println("Sending message to : " + packet.ID)
	}
	for client := range h.clients {
		if client.Tag == packet.ID {
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

func (h *Hub) sendBroadcast(packet *packet) {
	if h.mode == "Debug" {
		log.Println("Broadcasting from: " + packet.ID)
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
