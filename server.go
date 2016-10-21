package main

import (
    "fmt"
    "net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./client/")))
    http.ListenAndServe(":8090", nil)
    fmt.Println("HTTP Server listening on localhost:8090")
}