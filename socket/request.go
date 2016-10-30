package socket

type Request struct {
	Type string `json:"type"`
	X	int	`json:"x"`
	Y	int	`json:"y"`
	Message string `json:"message"`
}