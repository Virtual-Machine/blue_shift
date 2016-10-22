package login

import (
	"net/http"
	"encoding/json"
	"log"

	"github.com/dgrijalva/jwt-go"
)

type LoginResponse struct {
	Type 	string
	Message string
	Name 	string
}

type User struct {
    Name      string `json:name`
    Password	string `json:password`
    Token		string `json:token`
}

type UserList struct {
	List []User
}

func Api(data *UserList, w http.ResponseWriter, r *http.Request, mySigningKey []byte){
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	    "id": u.Name,
	})

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		log.Fatal(err)
	}
	u.Token = tokenString
	data.List = append(data.List, u)
	sendSuccessResponse(w, u)
}

func sendSuccessResponse(w http.ResponseWriter, u User){
	var res LoginResponse
	res.Type = "Success"
	res.Name = u.Name
	res.Message = u.Token
	json.NewEncoder(w).Encode(res)
}

func sendErrorResponse(w http.ResponseWriter){
	var err LoginResponse
	err.Type = "Error"
	err.Message = "Password is not correct for this user account"
	json.NewEncoder(w).Encode(err)
}