package algo

import (
	"fmt"
	"testing"

	"DatsNewWay/entity"
)

func TestRunnerAStar(t *testing.T) {

	response := entity.Response{
		MapSize: []int{120, 120, 120},
		Fences: [][]int{
			{1, 1, 1},
			{0, 1, 1},
			{0, 2, 1},
		},
		Snakes: []entity.Snake{
			{
				Id: "test",
				Geometry: [][]int{
					{0, 0, 0},
				},
				Status: snakeStatusAlive,
			},
		},
		Food: []entity.Food{
			{
				C: []int{
					3, 3, 3,
				},
			},
		},
	}

	obst := make(map[[3]int]bool)
	for _, v := range response.Fences {
		obst[[3]int{v[0], v[1], v[2]}] = true
	}

	dir := runnerAStar(response, []int{0, 0, 0}, []int{4, 4, 4}, obst)
	fmt.Println(dir)
}
