package main

import (
    "log"
    "net/http"
    "encoding/json"
    "os"

    "./socket"
    "./login"
)

type Configuration struct {
    Mode    string
    Port	string
}

var data login.UserList

func main() {
	log.SetFlags(log.Lshortfile)
	
	conf := decodeConfiguration()

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

func decodeConfiguration() Configuration {
	file, _ := os.Open("settings.json")
	decoder := json.NewDecoder(file)
	conf := Configuration{}
	err := decoder.Decode(&conf)
	if err != nil {
	  log.Fatal("Error:", err)
	}
	return conf
}