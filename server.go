package main

import (
    "log"
    "net/http"
    "encoding/json"
    "os"

    "./socket"
)

type Configuration struct {
    Mode    string
}

func main() {
	log.SetFlags(log.Lshortfile)
	
	file, _ := os.Open("settings.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
	  log.Fatal("Error:", err)
	}

	hub := socket.NewHub(configuration.Mode)
	go hub.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		socket.ServeWs(hub, w, r)
	})

	http.Handle("/", http.FileServer(http.Dir("./client/")))
    log.Fatal(http.ListenAndServe(":8090", nil))
}