package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type T struct {
	MapSize []int   `json:"mapSize"`
	Name    string  `json:"name"`
	Points  int     `json:"points"`
	Fences  [][]int `json:"fences"`
	Snakes  []struct {
		Id             string  `json:"id"`
		Direction      []int   `json:"direction"`
		OldDirection   []int   `json:"oldDirection"`
		Geometry       [][]int `json:"geometry"`
		DeathCount     int     `json:"deathCount"`
		Status         string  `json:"status"`
		ReviveRemainMs int     `json:"reviveRemainMs"`
	} `json:"snakes"`
	Enemies []struct {
		Geometry [][]int `json:"geometry"`
		Status   string  `json:"status"`
		Kills    int     `json:"kills"`
	} `json:"enemies"`
	Food []struct {
		C      []int `json:"c"`
		Points int   `json:"points"`
		Type   int   `json:"type"`
	} `json:"food"`
	SpecialFood struct {
		Golden     [][]int `json:"golden"`
		Suspicious [][]int `json:"suspicious"`
	} `json:"specialFood"`
	Turn             int           `json:"turn"`
	ReviveTimeoutSec int           `json:"reviveTimeoutSec"`
	TickRemainMs     int           `json:"tickRemainMs"`
	Errors           []interface{} `json:"errors"`
}

func main() {
	body, err := os.ReadFile("check.json")
	if err != nil {
		panic(err)
	}

	resp := T{}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		panic(err)
	}

	p := 0

	for _, r := range resp.Food {
		p += r.Points
	}

	fmt.Println(p)
}
