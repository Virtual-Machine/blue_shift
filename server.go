package main

import (
    "log"
    "net/http"
    "encoding/json"
    "os"

    "./socket"
)

type LoginResponse struct {
	Type 	string
	Message string
	Name 	string
	Password string
}

type User struct{
    Name      string `json:name`
    Password	string `json:password`
}

type UserList struct {
	List []User
}

type Configuration struct {
    Mode    string
    Port	string
}

var data UserList

func main() {
	log.SetFlags(log.Lshortfile)
	
	file, _ := os.Open("settings.json")
	decoder := json.NewDecoder(file)
	conf := Configuration{}
	err := decoder.Decode(&conf)
	if err != nil {
	  log.Fatal("Error:", err)
	}

	hub := socket.NewHub(conf.Mode)
	go hub.Run()

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        var u User
        if r.Body == nil {
            http.Error(w, "Please send a request body", 400)
            return
        }
        err := json.NewDecoder(r.Body).Decode(&u)
        if err != nil {
            http.Error(w, err.Error(), 400)
            return
        }
		
		for i := range data.List {
			if data.List[i].Name == u.Name {
				if data.List[i].Password != u.Password {
					var err LoginResponse
					err.Type = "Error"
					err.Message = "Password is not correct for this user account"
					json.NewEncoder(w).Encode(err)
					return
				}
			}
		}
		data.List = append(data.List, u)
        var res LoginResponse
        res.Type = "Success"
        res.Message = "Login successful"
        res.Name = u.Name
        res.Password = u.Password
        json.NewEncoder(w).Encode(res)
    })

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		socket.ServeWs(hub, w, r)
	})

	http.Handle("/", http.FileServer(http.Dir("./client/")))
    log.Fatal(http.ListenAndServe(conf.Port, nil))
}