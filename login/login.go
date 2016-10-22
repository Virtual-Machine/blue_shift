package login

import (
	"net/http"
	"encoding/json"
)

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

func Api(data UserList, w http.ResponseWriter, r *http.Request){
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
				sendErrorResponse(w)
				return
			}
		}
	}
	sendSuccessResponse(data, w, u)
}

func sendSuccessResponse(data UserList, w http.ResponseWriter, u User){
	data.List = append(data.List, u)
    var res LoginResponse
    res.Type = "Success"
    res.Message = "Login successful"
    res.Name = u.Name
    res.Password = u.Password
    json.NewEncoder(w).Encode(res)
}

func sendErrorResponse(w http.ResponseWriter){
	var err LoginResponse
	err.Type = "Error"
	err.Message = "Password is not correct for this user account"
	json.NewEncoder(w).Encode(err)
}