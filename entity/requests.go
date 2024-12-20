package entity

type Payload struct {
	Snakes []Snake `json:"snakes"`
}

type PayloadSnake struct {
	Id        string `json:"id"`
	Direction []int  `json:"direction"`
}
