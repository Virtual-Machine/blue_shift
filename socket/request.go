package socket

type Request struct {
	Type string `json:type`
	X	float64	`json:x`
	Y	float64	`json:y`
}