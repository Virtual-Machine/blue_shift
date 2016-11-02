package engine

import (
	"encoding/json"
	"errors"
	"log"
)

type MapCell struct {
	Background string
	Item       string
	Character  string
	Blocked    bool
	Clickable  bool
}

type MockGameEngine struct {
	Active  string
	Players []string
	MapData [60][40]MapCell
}

var GameInstance MockGameEngine

func init() {
	GameInstance.MapData[0][0].Background = "Grass"
	GameInstance.MapData[0][0].Blocked = true
	GameInstance.MapData[0][1].Clickable = true
}

func (g *MockGameEngine) StartGame(players []string) {
	g.Players = players
	g.Active = players[0]
}

func (g *MockGameEngine) ProcessClick(user string, x int, y int) (bool, error) {
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

func (g *MockGameEngine) GetData(user string, request string) []byte {
	if request == "MapData" {
		blob, err := json.Marshal(g.MapData)
		if err != nil {
			log.Fatal(err)
		}
		return blob
	}
	return []byte("Unknown request string")
}
