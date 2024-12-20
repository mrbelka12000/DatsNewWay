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

type node struct {
	point []int
	next  *node
}

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
		)

		for i, f := range r.Food {
			if f.Points < 0 {
				break
			}
			dist := getManhattanDistance(snake.Geometry[0], f.C)
			if dist < minDist {
				minDist = dist
				minInd = i
			}
		}

		dir := runnerAStar(r, snake.Geometry[0], r.Food[minInd].C, obst, used)
		//dir := runner(r, snake.Geometry[0], obst, food, used)
		obj.Snakes = append(obj.Snakes, entity.Snake{
			Id:        snake.Id,
			Direction: dir,
		})
	}

	return obj
}

type info struct {
	point []int
	steps []int
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

	heap.Push(q, info{
		point: currPoint,
		cost:  0,
		heur:  heuristic(currPoint, target),
	})

	for q.Len() > 0 {

		curr := heap.Pop(q).(info)

		cp := curr.point
		if cp[0] == target[0] && cp[1] == target[1] && cp[2] == target[2] {
			fmt.Println("Target reached:", currPoint, curr.steps)
			return curr.steps
		}

		used[[3]int{cp[0], cp[1], cp[2]}] = true

		for _, dir := range dirs {
			xx, yy, zz := cp[0]+dir[0], cp[1]+dir[1], cp[2]+dir[2]

			if xx < 0 || xx > r.MapSize[0] || yy < 0 || yy > r.MapSize[1] || zz < 0 || zz > r.MapSize[2] {
				continue
			}

			if obst[[3]int{xx, yy, zz}] || used[[3]int{xx, yy, zz}] {
				continue
			}

			steps := dir
			if len(curr.steps) != 0 {
				steps = curr.steps
			}

			newCost := curr.cost + 1
			heur := heuristic([]int{xx, yy, zz}, target)

			if _, ok := step[[3]int{xx, yy, zz}]; !ok || newCost < step[[3]int{xx, yy, zz}].cost {
				step[[3]int{xx, yy, zz}] = info{
					point: []int{xx, yy, zz},
					steps: steps,
					cost:  newCost,
					heur:  heur,
				}

				heap.Push(q, info{
					point: []int{xx, yy, zz},
					steps: steps,
					cost:  newCost,
					heur:  heur,
				})
			}
		}
	}

	return nil
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
