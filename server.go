package main

import (
	"log"
	"net/http"

	"./config"
	"./login"
	"./socket"
)

var data login.UserList

func init() {
	var u login.User
	u.Name = "ADMIN"
	u.Password = "Th3_D00d_1928!^^"
	u.Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IkFETUlOIn0.9Z-Ce-q0y7SkJBlBCb6Ghepz9hsUx6PuXBO6-3X565o"
	u.Admin = true
	data.List = append(data.List, u)
}

func main() {
	log.SetFlags(log.Lshortfile)

	conf := config.DecodeConfiguration()

	api_key := []byte(conf.SigningKey)

	hub := socket.NewHub(conf.Mode, &data)
	go hub.Run()

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		login.Api(&data, w, r, api_key, conf.Mode)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		socket.ServeWs(hub, w, r, api_key)
	})

	http.Handle("/", http.FileServer(http.Dir("./client/")))
	if conf.Mode == "Debug" {
		log.Println("Blue Shift   ---Online---    localhost", conf.Port)
	}
	log.Fatal(http.ListenAndServe(conf.Port, nil))
}
