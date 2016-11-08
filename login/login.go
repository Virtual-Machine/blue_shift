package login

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/scrypt"
)

type loginResponse struct {
	Type    string
	Message string
	Name    string
}

// UserSafe is the client safe data of the user
type UserSafe struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	Connections int    `json:"connections"`
}

// User stores the sensitive records of a user account
type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Admin    bool   `json:"admin"`
	Hashed   []byte `json:"hashed"`
}

// UserList is the servers array of user data
type UserList struct {
	List     []User
	SafeList []UserSafe
}

// API allows users to connect via an post submittal
func API(data *UserList, w http.ResponseWriter, r *http.Request, mySigningKey []byte, mode string, db *bolt.DB) {
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
	if strings.TrimSpace(u.Name) == "" || strings.TrimSpace(u.Password) == "" {
		if mode == "Debug" {
			log.Println("User submitted only whitespace for username and/or password")
		}
		sendErrorResponse(w, "Username / Password cannot consist only of whitespace characters")
		return
	}

	dk, err := scrypt.Key([]byte(u.Password), []byte("all-your-hash-belong-to-dr-seuss"), 16384, 8, 1, 32)

	if err != nil {
		sendErrorResponse(w, "There was an unexpected error, try again.")
		return
	}
	u.Password = ""
	u.Hashed = dk
	for i := range data.List {
		if data.List[i].Name == u.Name {
			if !bytes.Equal(data.List[i].Hashed, u.Hashed) {
				if mode == "Debug" {
					log.Println("Invalid submission attempt for account:", data.List[i].Name)
				}
				sendErrorResponse(w, "Password is not correct for this user account")
				return
			}
			if mode == "Debug" {
				log.Println("Successful login via API for account:", u.Name)
			}
			sendSuccessResponse(w, data.List[i])
			return
		}
	}
	// MARKER Server -> A new client is registering with server
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": u.Name,
	})

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		log.Fatal(err)
	}
	u.Token = tokenString
	u.Admin = false
	data.List = append(data.List, u)
	var uSafe UserSafe
	uSafe.Name = u.Name
	uSafe.Status = "Offline"
	uSafe.Connections = 0
	data.SafeList = append(data.SafeList, uSafe)
	if mode == "Debug" {
		log.Println("Successful login via API for account:", u.Name)
	}
	sendSuccessResponse(w, u)
}

func sendSuccessResponse(w http.ResponseWriter, u User) {
	var res loginResponse
	res.Type = "Success"
	res.Name = u.Name
	res.Message = u.Token
	json.NewEncoder(w).Encode(res)
}

func sendErrorResponse(w http.ResponseWriter, message string) {
	var err loginResponse
	err.Type = "Error"
	err.Message = message
	json.NewEncoder(w).Encode(err)
}
