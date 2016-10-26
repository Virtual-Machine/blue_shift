package main

import (
    "log"
    "net/http"

    "./socket"
    "./login"
    "./config"
)

var data login.UserList

func main() {
	log.SetFlags(log.Lshortfile)
	
	conf := config.DecodeConfiguration()

	api_key := []byte(conf.SigningKey)

	hub := socket.NewHub(conf.Mode)
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