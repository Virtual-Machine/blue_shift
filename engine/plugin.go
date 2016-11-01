package engine

type Plugin interface {
	ProcessClick(user string, x int, y int) (bool, error)
	GetData(user string, request string) []byte
	StartGame(players []string)
}

// If a game engine implements the above interface, it may be plugged into the server
// There are currently restrictions on the map data based on the client side rendering
// Future releases may allow the game engine to mandate the client side settings

// This package contains a mock implementation to satisfy development needs
// Simply import and use the desired game engine in place of the mock