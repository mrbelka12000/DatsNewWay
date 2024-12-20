package algo

import (
	"math"

	"DatsNewWay/entity"
)

const (
	snakeStatusAlive = "alive"
	snakeStatusDead  = "dead"
)

func GetNextDirection(r entity.Response) (obj entity.Payload) {
	used := make(map[int]bool)

	for _, snake := range r.Snakes {
		if snake.Status == snakeStatusDead {
			continue
		}
		var (
			direction []int
			minDist   = math.MaxInt32
			minInd    int
		)

		for i, food := range r.Food {
			dist := getManhattanDistance(snake.Geometry[0], food.C)
			if dist < minDist && !used[i] {
				minDist = dist
				minInd = i
				direction = getDirection(snake.Geometry[0], food.C)
			}
		}

		used[minInd] = true

		obj.Snakes = append(obj.Snakes, entity.Snake{
			Id:        snake.Id,
			Direction: direction,
		})
	}

	return obj
}

func getDirection(head, target []int) []int {
	if head[0] != target[0] {
		if head[0] < target[0] {
			return []int{1, 0, 0}
		}
		return []int{-1, 0, 0}
	}

	if head[1] != target[1] {
		if head[1] < target[1] {
			return []int{0, 1, 0}
		}
		return []int{0, -1, 0}
	}

	if head[2] != target[2] {
		if head[2] < target[2] {
			return []int{0, 0, 1}
		}
		return []int{0, 0, -1}
	}

	return []int{0, 0, 0} // Return this if head equals target
}

func getManhattanDistance(x, y []int) int {
	return abs(x[0]-y[0]) + abs(x[1]-y[1]) + abs(x[2]-y[2])
}

func abs(a int) int {
	if a < 0 {
		return a * -1
	}
	return a
}
