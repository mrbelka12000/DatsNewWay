package algo

import (
	"container/heap"
	"fmt"
	"math"
	"sort"

	"DatsNewWay/entity"
)

const (
	snakeStatusAlive = "alive"
	snakeStatusDead  = "dead"
)

func GetNextDirection(r entity.Response) (obj entity.Payload) {
	return bfs(r)
}

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

	sort.Slice(r.Food, func(i, j int) bool {
		return r.Food[i].Points > r.Food[j].Points
	})

	used := make(map[[3]int]bool)
	for _, snake := range r.Snakes {
		if snake.Status == snakeStatusDead {
			continue
		}

		var (
			minDist = math.MaxInt32
			minInd  int
			head    = snake.Geometry[0]
		)

		for i, f := range r.Food {
			if f.Points < 0 {
				break
			}
			if isCentralized(f.C, r.MapSize[0], r.MapSize[1], r.MapSize[2]) {
				dist := getManhattanDistance(head, f.C)
				if dist < minDist {
					minDist = dist
					minInd = i
				}
			}
		}

		if !isCentralized(head, r.MapSize[0], r.MapSize[1], r.MapSize[2]) && minDist > 5 {
			dir := runnerAStar(r, head, []int{r.MapSize[0] / 2, r.MapSize[1] / 2, r.MapSize[2] / 2}, obst, used)
			obj.Snakes = append(obj.Snakes, entity.Snake{
				Id:        snake.Id,
				Direction: dir,
			})
			fmt.Println(head, " not in centre")
			continue
		}
		fmt.Println(head, " in centre")

		dir := runnerAStar(r, head, r.Food[minInd].C, obst, used)
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

func runnerAStar(r entity.Response, currPoint, target []int, obst, used map[[3]int]bool) []int {
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
	var deep int

	for q.Len() > 0 {
		if deep > 10 {
			break
		}
		size := q.Len()

		for i := 0; i < size; i++ {
			curr := heap.Pop(q).(info)
			cp := curr.point

			// If the target is reached, return the path
			if cp[0] == target[0] && cp[1] == target[1] && cp[2] == target[2] {
				if len(curr.path) > 0 {
					fmt.Println("Target reached:", target, cp, curr.path[0])
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

/*
180 / 6 = 30
180 / 6 = 30
90 / 6 = 15
*/
