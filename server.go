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

	hub := socket.NewHub(conf.Mode)
	go hub.Run()

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        login.Api(data, w, r)
    })

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		socket.ServeWs(hub, w, r)
	})

	http.Handle("/", http.FileServer(http.Dir("./client/")))
    log.Fatal(http.ListenAndServe(conf.Port, nil))
}