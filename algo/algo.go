package algo

import (
	"container/heap"
	"fmt"
	"math"

	"DatsNewWay/entity"
)

var dirs = [6][3]int{
	{1, 0, 0},
	{-1, 0, 0},
	{0, 1, 0},
	{0, -1, 0},
	{0, 0, 1},
	{0, 0, -1},
}

const (
	snakeStatusAlive = "alive"
	snakeStatusDead  = "dead"
)

func GetNextDirection(r entity.Response) (obj entity.Payload) {
	for _, food := range r.Food {
		segmentPriority(food.C, r.MapSize[0], r.MapSize[1], r.MapSize[2], "FOOD")
	}

	for _, snake := range r.Enemies {
		segmentPriority(snake.Geometry[0], r.MapSize[0], r.MapSize[1], r.MapSize[2], "SNAKE")
	}

	return bfs(r)
}

func calculateProfit(head []int, food entity.Food, isGolden bool) float64 {
	dist := getManhattanDistance(head, food.C)

	profit := float64(food.Points) / (float64(dist) + 10)

	if isGolden {
		profit *= 10
	}

	return profit
}

func bfs(r entity.Response) (obj entity.Payload) {

	obst := make(map[[3]int]bool, len(r.Fences)+len(r.Enemies))

	// fill obstacles with fences
	for _, fence := range r.Fences {
		key := [3]int{fence[0], fence[1], fence[2]}
		obst[key] = true
	}

	// fill obstacles with enemies
	for _, enemy := range r.Enemies {
		for i, coord := range enemy.Geometry {
			if i == 0 {
				// handle + 2 cell after enemies head
				for _, dir := range dirs {
					key := [3]int{coord[0] + dir[0], coord[1] + dir[1], coord[2] + dir[2]}
					obst[key] = true
					for _, dd := range dirs {
						key[0] += dd[0]
						key[1] += dd[1]
						key[2] += dd[2]
						obst[key] = true
					}
				}
			}

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

	usedIDs := make(map[int]bool)
	for _, snake := range r.Snakes {
		if snake.Status == snakeStatusDead {
			continue
		}

		var (
			maxProfit float64
			maxInd    int
			sum       int
			head      = snake.Geometry[0]
			minDist   = math.MaxInt32
			minInd    int
		)

		for i, f := range r.Food {
			if usedIDs[i] {
				continue
			}
			if f.Points < 0 {
				continue
			}

			profit := calculateProfit(head, f, false)
			if profit > maxProfit {
				maxProfit = profit
				maxInd = i
			}

			dist := getManhattanDistance(head, f.C)
			if minDist > dist {
				minDist = dist
				minInd = i
			}
			sum += f.Points
		}

		usedIDs[maxInd] = true
		fmt.Println(maxProfit)
		if !isCentralized(head, r.MapSize[0], r.MapSize[1], r.MapSize[2]) && maxProfit < 2 {
			fmt.Println("Идем в центр: ", snake.Id, maxProfit)
			dir := runnerAStar(r, head, getPreviousPoint(snake), []int{r.MapSize[0] / 2, r.MapSize[1] / 2, r.MapSize[2] / 2}, obst)
			obj.Snakes = append(obj.Snakes, entity.Snake{
				Id:        snake.Id,
				Direction: dir,
			})
		} else if maxProfit > 6 {
			// run for profitable mandarin
			dir := runnerAStar(r, head, getPreviousPoint(snake), r.Food[maxInd].C, obst)
			obj.Snakes = append(obj.Snakes, entity.Snake{
				Id:        snake.Id,
				Direction: dir,
			})
		} else {
			// run for minimal distance
			dir := runnerAStar(r, head, getPreviousPoint(snake), r.Food[minInd].C, obst)
			obj.Snakes = append(obj.Snakes, entity.Snake{
				Id:        snake.Id,
				Direction: dir,
			})
		}
	}

	return obj
}

type info struct {
	point []int
	path  [][]int // Stores the path as a list of points
	cost  int
	heur  int
}

func runnerAStar(r entity.Response, currPoint, prevPoint, target []int, obst map[[3]int]bool) []int {
	prevDir := getOpositeDir([3]int{currPoint[0] - prevPoint[0], currPoint[1] - prevPoint[1], currPoint[2] - prevPoint[2]})
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
	var deep int

	for q.Len() > 0 {
		if deep > 5 {
			break
		}
		size := q.Len()

		for i := 0; i < size; i++ {
			curr := heap.Pop(q).(info)
			cp := curr.point

			// If the target is reached, return the path
			if cp[0] == target[0] && cp[1] == target[1] && cp[2] == target[2] {
				if len(curr.path) > 0 {
					fmt.Println("Found direction for: ", currPoint, curr.path[0])
					return curr.path[0]
				}
				continue
			}

			for _, dir := range dirs {
				if dir == prevDir {
					continue
				}
				xx, yy, zz := cp[0]+dir[0], cp[1]+dir[1], cp[2]+dir[2]

				// Check boundaries
				if xx < 0 || xx >= r.MapSize[0] || yy < 0 || yy >= r.MapSize[1] || zz < 0 || zz >= r.MapSize[2] {
					continue
				}

				// Check for obstacles
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
					newPath = append(newPath, dir[:])

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
		prevDir = [3]int{}
		deep++
	}

	if q.Len() == 0 {
		return nil
	}

	return heap.Pop(q).(info).path[0] // No path found
}

func heuristic(currPoint []int, target []int) int {
	return abs(currPoint[0]-target[0]) + abs(currPoint[1]-target[1]) + abs(currPoint[2]-target[2])
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

func isCentralized(head []int, x, y, z int) bool {
	centreX := x / 2
	centreY := y / 2
	centreZ := z / 2

	//dist := getManhattanDistance(head, []int{centreX, centreY, centreZ})

	quadX := centreX / 2
	quadY := centreY / 2
	quadZ := centreZ / 2
	return centreX-quadX < head[0] && centreX+quadX > head[0] &&
		centreY-quadY < head[1] && centreY+quadY > head[1] &&
		centreZ-quadZ < head[2] && centreZ+quadZ > head[2]
}

var (
	segmentFoodInfo  = make(map[int]int)
	segmentSnakeInfo = make(map[int]int)
)

func segmentPriority(point []int, x, y, z int, t string) int {
	segmentId := 0

	if isCentralized(point, x-x/2, y-y/2, z-z/2) {
		segmentId = 1
	}

	if isCentralized(point, x-x/2, y-y/2, z+z/2) {
		segmentId = 2
	}

	if isCentralized(point, x-x/2, y+y/2, z-z/2) {
		segmentId = 3
	}

	if isCentralized(point, x-x/2, y+y/2, z+z/2) {
		segmentId = 4
	}

	if isCentralized(point, x+x/2, y-y/2, z-z/2) {
		segmentId = 5
	}

	if isCentralized(point, x+x/2, y-y/2, z+z/2) {
		segmentId = 6
	}

	if isCentralized(point, x+x/2, y+y/2, z-z/2) {
		segmentId = 7
	}

	if isCentralized(point, x+x/2, y+y/2, z+z/2) {
		segmentId = 8
	}

	if t == "SNAKE" {
		segmentSnakeInfo[segmentId]++
	} else if t == "FOOD" {
		segmentFoodInfo[segmentId]++
	}

	return segmentId
}

func getPreviousPoint(snake entity.Snake) []int {
	if len(snake.Geometry) == 1 {
		return snake.Geometry[0]
	}
	return snake.Geometry[1]
}

func getOpositeDir(dir [3]int) [3]int {
	switch dir {
	case [3]int{1, 0, 0}:
		return [3]int{-1, 0, 0}
	case [3]int{-1, 0, 0}:
		return [3]int{1, 0, 0}
	case [3]int{0, 1, 0}:
		return [3]int{0, -1, 0}
	case [3]int{0, -1, 0}:
		return [3]int{0, 1, 0}
	case [3]int{0, 0, -1}:
		return [3]int{0, 0, 1}
	case [3]int{0, 0, 1}:
		return [3]int{0, 0, -1}
	}

	return [3]int{0, 0, 0}
}

/*
180 / 6 = 30
180 / 6 = 30
90 / 6 = 15
*/
