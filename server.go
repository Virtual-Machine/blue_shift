package main

import (
	"encoding/json"
	"log"
	"net/http"

	"./config"
	"./login"
	"./socket"

	"github.com/boltdb/bolt"
)

var data login.UserList

func main() {
	db, err := bolt.Open("users.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		var user login.User
		b := tx.Bucket([]byte("users"))
		v := b.Get([]byte("ADMIN"))
		err := json.Unmarshal(v, &user)
		if err != nil {
			return err
		}
		data.List = append(data.List, user)
		return nil
	})

	log.SetFlags(log.Lshortfile)

	conf := config.DecodeConfiguration()

	apiKey := []byte(conf.SigningKey)

	hub := socket.NewHub(conf.Mode, &data)
	go hub.Run()

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		login.API(&data, w, r, apiKey, conf.Mode)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		socket.ServeWs(&data, hub, w, r, apiKey)
	})

	http.Handle("/", http.FileServer(http.Dir("./client/")))
	if conf.Mode == "Debug" {
		log.Println("Blue Shift   ---Online---    localhost", conf.Port)
	}
	log.Fatal(http.ListenAndServe(conf.Port, nil))
}
