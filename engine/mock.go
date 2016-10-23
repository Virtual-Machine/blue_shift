package engine

import (
	"encoding/json"
	"log"
)

type MapCell struct {
	Background string
	Item 	   string
	Character  string
	Blocked	   bool
	Clickable  bool
}

type MockGameEngine struct {
	Active string
	MapData [60][40]MapCell

}

var GameInstance MockGameEngine

func init() {
	GameInstance.Active = "John"
	GameInstance.MapData[0][0].Background = "Grass"
	GameInstance.MapData[0][0].Blocked = true
	GameInstance.MapData[0][1].Clickable = true
}

func (g *MockGameEngine) ProcessClick(user string, x int, y int) bool {
	if user != g.Active {
		return false
	}
	if g.MapData[x][y].Blocked {
		return false
	}
	if g.MapData[x][y].Clickable {
		return true
	}
	return false
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