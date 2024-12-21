package algo

import (
	"DatsNewWay/entity"
	"container/heap"
	"fmt"
	"math"
)

const (
	snakeStatusAlive = "alive"
	snakeStatusDead  = "dead"
)

var center [3]int

func calculateCenter(maxX, maxY, maxZ int) (int, int, int) {
	centerX := maxX / 2
	centerY := maxY / 2
	centerZ := maxZ / 2
	return centerX, centerY, centerZ
}

func distanceToCenter(x, y, z int) float64 {
	return math.Sqrt(math.Pow(float64(x-center[0]), 2) + math.Pow(float64(y-center[1]), 2) + math.Pow(float64(z-center[2]), 2))
}

// Функция, чтобы проверить, является ли точка нехорошей
func isBadPoint(currentPosition, food []int) bool {
	foodDist := distanceToCenter(food[0], food[1], food[2])
	currentDist := distanceToCenter(currentPosition[0], currentPosition[1], currentPosition[2])

	// Проверяем, если точка еды дальше от центра, чем текущая точка
	if foodDist > currentDist {
		return true
	}

	if (food[0] > center[0] && currentPosition[0] < center[0]) || (food[0] < center[0] && currentPosition[0] > center[0]) {
		return true
	}
	if (food[1] > center[1] && currentPosition[1] < center[1]) || (food[1] < center[1] && currentPosition[1] > center[1]) {
		return true
	}
	if (food[2] > center[2] && currentPosition[2] < center[2]) || (food[2] < center[2] && currentPosition[2] > center[2]) {
		return true
	}

	return false
}

func calculatePriority(distance, point, k1, k2 int) int {
	return k1*point + k2*distance
}

var flags = [3]bool{false, false, false}

func GetNextDirection(r entity.Response) (obj entity.Payload) {
	center[0], center[1], center[2] = calculateCenter(r.MapSize[0], r.MapSize[1], r.MapSize[2])
	return bfs(r)
}

var mapping = map[[3]int]bool{}

func bfs(r entity.Response) (obj entity.Payload) {

	obst := make(map[[3]int]bool)

	// fill obstacles with fences
	for _, fence := range r.Fences {
		key := [3]int{fence[0], fence[1], fence[2]}
		obst[key] = true
	}

	// fill obstacles with enemies
	for _, enemy := range r.Enemies {
		for _, coord := range enemy.Geometry {
			key := [3]int{coord[0], coord[1], coord[2]}
			obst[key] = true
		}
	}

	for _, snake := range r.Snakes {
		for _, geo := range snake.Geometry {
			key := [3]int{geo[0], geo[1], geo[2]}
			obst[key] = true
		}
	}

	food := make(map[[3]int]bool)
	// fill food hash table
	for _, f := range r.Food {
		key := [3]int{f.C[0], f.C[1], f.C[2]}
		food[key] = true
	}

	used := make(map[int]bool)
	for idx, snake := range r.Snakes {
		if snake.Status == snakeStatusDead {
			continue
		}

		var (
			minDist     = math.MaxInt32
			minInd      = -1
			maxPriority = math.MinInt32
		)

		for i, f := range r.Food {
			if used[i] {
				continue
			}
			if f.Points < 0 {
				continue
			}
			dist := getManhattanDistance(snake.Geometry[0], f.C)

			if !flags[idx] && isBadPoint(snake.Geometry[0], f.C) {
				continue
			}

			if flags[idx] {
				priority := calculatePriority(dist, f.Points, 2, 1)
				if priority > maxPriority {
					maxPriority = priority
					minInd = i
				}
				continue
			}

			if dist < minDist {
				minDist = dist
				minInd = i
			}
		}

		if minInd == -1 {
			flags[idx] = true
			minInd = 0
		}

		used[minInd] = true
		dir := runnerAStar(r, snake.Geometry[0], r.Food[minInd].C, obst)
		//dir := runner(r, snake.Geometry[0], obst, food, used)
		if _, ok := mapping[[3]int{snake.Geometry[0][0], snake.Geometry[0][1], snake.Geometry[0][2]}]; !ok {
			mapping[[3]int{snake.Geometry[0][0], snake.Geometry[0][1], snake.Geometry[0][2]}] = true
		} else {
			fmt.Printf("индекс змейки: %v, direction: %v\n", r.Snakes[idx], dir)
		}

		obj.Snakes = append(obj.Snakes, entity.Snake{
			Id:        snake.Id,
			Direction: dir,
		})
	}

	return obj
}

type info struct {
	point []int
	path  [][]int // Stores the path as a list of points
	cost  int
	heur  int
}

func runnerAStar(r entity.Response, currPoint, target []int, obst map[[3]int]bool) []int {
	dirs := [6][]int{
		{1, 0, 0},
		{-1, 0, 0},
		{0, 1, 0},
		{0, -1, 0},
		{0, 0, 1},
		{0, 0, -1},
	}

	step := make(map[[3]int]info)

	q := &PQ{}
	heap.Init(q)

	// Start with the current point
	heap.Push(q, info{
		point: currPoint,
		path:  [][]int{}, // Start path with the initial point
		cost:  0,
		heur:  heuristic(currPoint, target),
	})

	for q.Len() > 0 {
		curr := heap.Pop(q).(info)
		cp := curr.point

		// If the target is reached, return the path
		if cp[0] == target[0] && cp[1] == target[1] && cp[2] == target[2] {
			fmt.Println("Target reached:", target, cp, curr.path)
			if len(curr.path) > 0 {
				return curr.path[0]
			}
			continue
		}

		for _, dir := range dirs {
			xx, yy, zz := cp[0]+dir[0], cp[1]+dir[1], cp[2]+dir[2]

			// Check boundaries
			if xx < 0 || xx > r.MapSize[0] || yy < 0 || yy > r.MapSize[1] || zz < 0 || zz > r.MapSize[2] {
				continue
			}

			// Check for obstacles and already visited points
			if obst[[3]int{xx, yy, zz}] {
				continue
			}

			newCost := curr.cost + 1
			heur := heuristic([]int{xx, yy, zz}, target)

			// If a better path is found, update the step map and push the new state into the queue
			if _, ok := step[[3]int{xx, yy, zz}]; !ok || newCost < step[[3]int{xx, yy, zz}].cost {
				// Create a copy of the current path and add the new point
				newPath := make([][]int, len(curr.path))
				copy(newPath, curr.path)
				newPath = append(newPath, dir)

				step[[3]int{xx, yy, zz}] = info{
					point: []int{xx, yy, zz},
					path:  newPath,
					cost:  newCost,
					heur:  heur,
				}

				heap.Push(q, info{
					point: []int{xx, yy, zz},
					path:  newPath,
					cost:  newCost,
					heur:  heur,
				})
			}
		}
	}

	return nil // No path found
}

func heuristic(currPoint []int, target []int) int {
	return abs(currPoint[0]-target[0]) + abs(currPoint[1]-target[1]) + abs(currPoint[2]-target[2])
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
