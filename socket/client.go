package socket

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"../login"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second
	pongWait  = 60 * time.Second

	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		if strings.Split(r.RemoteAddr, ":")[0][:10] == "192.168.5." {
			return true
		}
		return false
	},
}

type client struct {
	hub *Hub

	conn *websocket.Conn

	send chan []byte

	userSafe *login.UserSafe
	user     *login.User
}

func (c *client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("Error: %v", err)
			}
			break
		}
		var pack packet
		pack.ID = c.user.Name
		pack.Data = string(bytes.TrimSpace(bytes.Replace(message, newline, space, -1)))

		if c.hub.mode == "Debug" {
			log.Println("Client socket is sending message to hub:", pack.ID)
		}
		c.hub.request <- &pack
	}
}

func (c *client) write(mt int, payload []byte) error {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteMessage(mt, payload)
}

func (c *client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// ServeWs establishes a socket connection between a client and the hub
func ServeWs(data *login.UserList, hub *Hub, w http.ResponseWriter, r *http.Request, apiKey []byte) {
	tokenString := r.URL.Query().Get("id")
	var idString string
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return apiKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		idString = claims["id"].(string)
	} else {
		log.Println(err)
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	if hub.mode == "Debug" {
		log.Println("Successful socket connection established for user:", idString)
	}
	var userPointer *login.User
	var userSafePointer *login.UserSafe
	for i, v := range data.List {
		if v.Name == idString {
			userPointer = &data.List[i]
		}
	}
	for i, v := range data.SafeList {
		if v.Name == idString {
			userSafePointer = &data.SafeList[i]
		}
	}
	client := &client{hub: hub, conn: conn, send: make(chan []byte, 256), user: userPointer, userSafe: userSafePointer}
	client.hub.register <- client
	go client.writePump()
	client.readPump()
}
