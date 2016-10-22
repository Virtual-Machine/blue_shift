package socket

type Packet struct {
	Id string `json:id`
	Data string `json:data`
}