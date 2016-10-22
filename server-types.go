package main

type LoginResponse struct {
	Type 	string
	Message string
	Name 	string
	Password string
}

type User struct {
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