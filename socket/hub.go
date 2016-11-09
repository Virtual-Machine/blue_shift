package socket

import (
	"encoding/json"
	"log"
	"strings"

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
	started    bool
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
		started:    false,
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
		log.Println("Connecting socket @", client.conn.RemoteAddr(), client.user.Name)
	}
	h.clients[client] = true

	if client.user.Admin {
		var pack packet
		pack.ID = client.user.Name

		data, _ := json.Marshal(h.users.SafeList)
		pack.Data = "{\"display_admin_panel\":\"true\", \"user_list\": " + string(data) + "}"

		h.sendMessage(&pack)
		return
	}

	if client.userSafe != nil {
		// MARKER Server -> Client connected.
		client.userSafe.Status = "Online"
		client.userSafe.Connections++
		var pack packet
		pack.ID = client.user.Name
		data, _ := json.Marshal(h.users.SafeList)
		pack.Data = "{\"user_list\": " + string(data) + "}"
		h.sendBroadcast(&pack)
		return
	}
}

func (h *Hub) disconnect(client *client) {
	if h.mode == "Debug" {
		log.Println("Disconnecting socket @", client.conn.RemoteAddr(), client.user.Name)
	}
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)
	}
	if client.userSafe != nil {
		// MARKER Server -> Client disconnected.
		client.userSafe.Connections--
		if client.userSafe.Connections <= 0 {
			client.userSafe.Connections = 0
			client.userSafe.Status = "Offline"
		}
		var pack packet
		pack.ID = client.user.Name
		data, _ := json.Marshal(h.users.SafeList)
		pack.Data = "{\"user_list\": " + string(data) + "}"
		h.sendBroadcast(&pack)
		return
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

	h.processRequest(packet, req)
}

func (h *Hub) processRequest(packet *packet, req request) {
	// MARKER Server -> Socket server received data from client.
	if req.Type == "Click" {
		h.processClick(packet, req)
		return
	}
	if req.Type == "StartGame" {
		if h.started {
			packet.Data = "{\"admin_error\":\"This server is already running a game\"}"
			h.sendMessage(packet)
			return
		}
		h.processStart(packet, req)
		return
	}
	if req.Type == "GameData" {
		packet.Data = string(engine.GameInstance.GetData())
		h.sendBroadcast(packet)
		return
	}
	if req.Type == "ChatMessage" {
		packet.Data = "{\"message\":\"" + req.Message + "\", \"author\": \"" + packet.ID + "\"}"
		h.sendBroadcast(packet)
		return
	}
}

func (h *Hub) sendMessage(packet *packet) {
	if h.mode == "Debug" {
		log.Println("Sending message to : " + packet.ID)
	}
	for client := range h.clients {
		if client.user.Name == packet.ID {
			select {
			case client.send <- []byte(packet.Data):
				if h.mode == "Debug" {
					log.Println(client.conn.RemoteAddr(), client.user.Name, "received data")
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
				log.Println(client.conn.RemoteAddr(), client.user.Name, "received data")
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

func (h *Hub) processClick(p *packet, r request) {
	if r.X < 0 || r.Y < 0 || r.X >= 60 || r.Y >= 40 {
		p.Data = "{\"error\":\"Click is out of bounds\"}"
		h.sendMessage(p)
		return
	}
	validMove, err := engine.GameInstance.ProcessClick(p.ID, r.X, r.Y)
	if validMove {
		p.Data = string(engine.GameInstance.GetData())
		h.sendBroadcast(p)
		return
	}
	p.Data = "{\"error\":\"" + err.Error() + "\"}"
	h.sendMessage(p)
}

func (h *Hub) processStart(p *packet, r request) {
	for _, v := range h.users.List {
		if v.Name == p.ID && v.Admin == true {
			names := strings.Split(r.Message, ";")
			count := len(names)
			if count < 2 || count > 4 {
				p.Data = "{\"admin_error\":\"This server is setup to only support 2-4 players\"}"
				h.sendMessage(p)
				return
			}
			for _, name := range names {
				found := false
				for _, v2 := range h.users.List {
					if v2.Name == name {
						found = true
					}
				}
				if !found {
					p.Data = "{\"admin_error\":\"Submitted name: " + name + " not found on server\"}"
					h.sendMessage(p)
					return
				}
			}
			engine.GameInstance.StartGame(names)
			h.started = true
			jnames, _ := json.Marshal(names)
			p.Data = "{\"success\":\"Game Started\", \"players\":" + string(jnames) + "}"
			h.sendBroadcast(p)
			return
		}
	}
}
