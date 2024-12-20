package entity

type Snake struct {
	Id             string  `json:"id"`
	Direction      []int   `json:"direction"`
	OldDirection   []int   `json:"oldDirection"`
	Geometry       [][]int `json:"geometry"`
	DeathCount     int     `json:"deathCount"`
	Status         string  `json:"status"`
	ReviveRemainMs int     `json:"reviveRemainMs"`
}

type Enemy struct {
	Geometry [][]int `json:"geometry"`
	Status   string  `json:"status"`
	Kills    int     `json:"kills"`
}

type Food struct {
	C      []int `json:"c"`
	Points int   `json:"points"`
	Type   int   `json:"type"`
}

type SpecialFood struct {
	Golden     [][]int `json:"golden"`
	Suspicious [][]int `json:"suspicious"`
}

type Response struct {
	MapSize          []int         `json:"mapSize"`
	Name             string        `json:"name"`
	Points           int           `json:"points"`
	Fences           [][]int       `json:"fences"`
	Snakes           []Snake       `json:"snakes"`
	Enemies          []Enemy       `json:"enemies"`
	Food             []Food        `json:"food"`
	SpecialFood      SpecialFood   `json:"specialFood"`
	Turn             int           `json:"turn"`
	ReviveTimeoutSec int           `json:"reviveTimeoutSec"`
	TickRemainMs     int           `json:"tickRemainMs"`
	Errors           []interface{} `json:"errors"`
}
