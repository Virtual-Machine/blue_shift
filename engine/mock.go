package engine

import (
	"encoding/json"
	"errors"
	"log"
)

type mapCell struct {
	Background string
	Item       string
	Character  string
	Blocked    bool
	Clickable  bool
}

type mockGameEngine struct {
	Active  string
	Players []string
	MapData [60][40]mapCell
}

// GameInstance is the interface point of the socket server
var GameInstance mockGameEngine

func init() {
	GameInstance.MapData[0][0].Background = "Grass"
	GameInstance.MapData[0][0].Blocked = true
	GameInstance.MapData[0][1].Clickable = true
}

func (g *mockGameEngine) StartGame(players []string) {
	g.Players = players
	g.Active = players[0]
}

func (g *mockGameEngine) ProcessClick(user string, x int, y int) (bool, error) {
	if user != g.Active {
		return false, errors.New("User is not active")
	}
	if g.MapData[x][y].Blocked {
		return false, errors.New("Cell is blocked")
	}
	if g.MapData[x][y].Clickable {
		return true, nil
	}
	return false, errors.New("Cell is not clickable")
}

func (g *mockGameEngine) GetData() []byte {
	blob, err := json.Marshal(g)
	if err != nil {
		log.Fatal(err)
	}
	return blob
}
