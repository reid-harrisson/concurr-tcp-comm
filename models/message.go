package msg

type Message struct {
	UUID   string `json:"uuid"`
	Prefix string `json:"prefix"`
	Data   string `json:"data"`
}
