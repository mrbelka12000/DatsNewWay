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

	dir := runnerAStar(response, []int{2, 2, 2}, []int{3, 2, 2}, []int{1, 2, 2}, obst)
	fmt.Println(dir)
}

func TestGetProfit(t *testing.T) {

	food := entity.Food{
		C:      []int{1, 3, 2},
		Points: 100,
	}

	fmt.Println(calculateProfit([]int{0, 0, 0}, food, false))
}
