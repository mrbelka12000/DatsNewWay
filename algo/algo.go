package algo

import (
	"container/heap"
	"fmt"
	"math"
	"time"

	"DatsNewWay/entity"
)

type (
	SegmentInfo struct {
		CountSnakes     int
		CountFood       int
		CountGoldenFood int
		TotalFoodPoints int
	}
)

var (
	dirs = [6][3]int{
		{1, 0, 0},
		{-1, 0, 0},
		{0, 1, 0},
		{0, -1, 0},
		{0, 0, 1},
		{0, 0, -1},
	}
	orderX, orderY, orderZ int

	segmentInfo = map[int]*SegmentInfo{
		0: &SegmentInfo{},
		1: &SegmentInfo{},
		2: &SegmentInfo{},
		3: &SegmentInfo{},
		4: &SegmentInfo{},
		5: &SegmentInfo{},
		6: &SegmentInfo{},
		7: &SegmentInfo{},
		8: &SegmentInfo{},
	}
	foodTotalPoints = 0
	avgTotal        = 0

	totalProfit      float64
	totalProfitCount int
	currentProfit    = 5.0

	distAdd = 10.0
)

func init() {
	go runProfitAvg()
}

func runProfitAvg() {
	t := time.NewTicker(10 * time.Second)

	for {
		<-t.C
		oldProfit := currentProfit
		currentProfit = (totalProfit / float64(totalProfitCount)) * 0.75
		fmt.Printf("old profit:%v, new profit:%v, current distAdd:%v\n", oldProfit, currentProfit, distAdd)
		totalProfit = 0
		totalProfitCount = 0
	}
}

const (
	snakeStatusAlive = "alive"
	snakeStatusDead  = "dead"
)

func segmentSnakePriority(point []int, x, y, z int) {
	segmentId := segmentPriority(point, x, y, z)
	segmentInfo[segmentId].CountSnakes++
}

func segmentGoldenFoodPriority(goldenFood []int, x, y, z int) {
	segmentId := segmentPriority(goldenFood, x, y, z)
	segmentInfo[segmentId].CountGoldenFood++
	goldenApproximate := avgTotal * 10
	segmentInfo[segmentId].TotalFoodPoints += goldenApproximate
}

func segmentFoodPriority(food *entity.Food, x, y, z int) {
	food.SegmentInd = segmentPriority(food.C, x, y, z)
	segmentInfo[food.SegmentInd].CountFood++
	foodTotalPoints += food.Points
	segmentInfo[food.SegmentInd].TotalFoodPoints += food.Points
}

func prepareSegmentPriority(r entity.Response) {
	totalCount := 0
	snakeCount := 0
	goldenCount := 0
	segmentInfo = map[int]*SegmentInfo{
		0: &SegmentInfo{},
		1: &SegmentInfo{},
		2: &SegmentInfo{},
		3: &SegmentInfo{},
		4: &SegmentInfo{},
		5: &SegmentInfo{},
		6: &SegmentInfo{},
		7: &SegmentInfo{},
		8: &SegmentInfo{},
	}

	for _, food := range r.Food {
		segmentFoodPriority(&food, r.MapSize[0], r.MapSize[1], r.MapSize[2])
	}
	avgTotal = foodTotalPoints / (len(r.Food) + 1)

	for _, snake := range r.Enemies {
		segmentSnakePriority(snake.Geometry[0], r.MapSize[0], r.MapSize[1], r.MapSize[2])
	}

	for _, goldenFood := range r.SpecialFood.Golden {
		segmentGoldenFoodPriority(goldenFood, r.MapSize[0], r.MapSize[1], r.MapSize[2])
	}

	for _, segment := range segmentInfo {
		totalCount += segment.TotalFoodPoints
		snakeCount += segment.CountSnakes
		goldenCount += segment.CountGoldenFood
	}

	//fmt.Println("total count", totalCount)
	//fmt.Println("snake count", snakeCount)
	//fmt.Println("golden count", goldenCount)

	foodTotalPoints = 0
}

func GetNextDirection(r entity.Response) (obj entity.Payload) {
	prepareSegmentPriority(r)
	orderX = r.MapSize[0]
	orderY = r.MapSize[1]
	orderZ = r.MapSize[2]

	getDistAdd()

	return bfs(r)
}

func calculateProfit(head []int, food entity.Food, isGolden bool, ind int) float64 {
	dist := getManhattanDistance(head, food.C)

	segmentID := segmentPriority(food.C, orderX, orderY, orderZ)

	sInfo := segmentInfo[segmentID]
	if sInfo.TotalFoodPoints == 0 {
		return 0
	}

	// Calculate FoodFactor
	profit := float64(food.Points) / (dist + distAdd)

	if isGolden {
		profit *= 5
	}

	return profit
}

func bfs(r entity.Response) (obj entity.Payload) {
	obst := make(map[[3]int]bool, len(r.Fences)+len(r.Enemies))

	// fill obstacles with fences
	for _, fence := range r.Fences {
		for _, dir := range dirs {
			key := [3]int{fence[0] + dir[0], fence[1] + dir[1], fence[2] + dir[2]}
			obst[key] = true
		}
		key := [3]int{fence[0], fence[1], fence[2]}
		obst[key] = true
	}

	// fill obstacles with enemies
	for _, enemy := range r.Enemies {
		for _, coord := range enemy.Geometry {
			for _, dir := range dirs {
				key := [3]int{coord[0] + dir[0], coord[1] + dir[1], coord[2] + dir[2]}
				obst[key] = true
				kk := key
				for _, dd := range dirs {
					kk[0] += dd[0]
					kk[1] += dd[1]
					kk[2] += dd[2]
					obst[kk] = true
					kk = key
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
			maxProfit       float64
			maxInd          int
			sum             int
			head            = snake.Geometry[0]
			minDist         = math.MaxFloat32
			minInd          int
			totalFoodPoints int
		)

		for i, f := range r.Food {
			if usedIDs[i] || f.Points < 0 {
				continue
			}

			profit := calculateProfit(head, f, false, 10)
			if profit > maxProfit {
				maxProfit = profit
				maxInd = i
			}
			totalFoodPoints += f.Points

			dist := getManhattanDistance(head, f.C)
			if minDist > dist {
				minDist = dist
				minInd = i
			}
			sum += f.Points
		}

		totalProfit += maxProfit
		totalProfitCount++

		if maxProfit > currentProfit {
			// run for profitable mandarin
			fmt.Println("get maximum profit")
			usedIDs[maxInd] = true
			dir := runnerAStar(r, head, getPreviousPoint(snake), r.Food[maxInd].C, obst)
			obj.Snakes = append(obj.Snakes, entity.Snake{
				Id:        snake.Id,
				Direction: dir,
			})
		} else {
			var (
				goldProfit float64
				goldInd    int
			)

			for i, gold := range r.SpecialFood.Golden {
				profit := calculateProfit(head, entity.Food{
					C:      gold,
					Points: totalFoodPoints / len(r.Food),
				}, true, 10)
				if profit > goldProfit {
					goldProfit = profit
					goldInd = i
				}
			}

			if goldProfit > currentProfit {
				fmt.Println("went to golden food")
				usedIDs[goldInd] = true
				dir := runnerAStar(r, head, getPreviousPoint(snake), r.SpecialFood.Golden[goldInd], obst)
				obj.Snakes = append(obj.Snakes, entity.Snake{
					Id:        snake.Id,
					Direction: dir,
				})
			} else {
				// run for minimal distance
				fmt.Println("went to minimal distance food")
				usedIDs[minInd] = true
				dir := runnerAStar(r, head, getPreviousPoint(snake), r.Food[minInd].C, obst)
				obj.Snakes = append(obj.Snakes, entity.Snake{
					Id:        snake.Id,
					Direction: dir,
				})
			}

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
					//fmt.Println("Found direction for: ", currPoint, curr.path[0])
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

func getManhattanDistance(x, y []int) float64 {
	return math.Sqrt(math.Pow(float64(x[0]-y[0]), 2) + math.Pow(float64(x[1]-y[1]), 2) + math.Pow(float64(x[2]-y[2]), 2))
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

func segmentPriority(point []int, x, y, z int) int {
	segmentId := 0
	switch {
	case isCentralized(point, x-x/2, y-y/2, z-z/2):
		segmentId = 1

	case isCentralized(point, x-x/2, y-y/2, z+z/2):
		segmentId = 2

	case isCentralized(point, x-x/2, y+y/2, z-z/2):
		segmentId = 3

	case isCentralized(point, x-x/2, y+y/2, z+z/2):
		segmentId = 4

	case isCentralized(point, x+x/2, y-y/2, z-z/2):
		segmentId = 5

	case isCentralized(point, x+x/2, y-y/2, z+z/2):
		segmentId = 6

	case isCentralized(point, x+x/2, y+y/2, z-z/2):
		segmentId = 7

	case isCentralized(point, x+x/2, y+y/2, z+z/2):
		segmentId = 8
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

func getDistAdd() {
	distAdd = float64(orderX)/50 + float64(orderY)/50 + float64(orderZ)/30
}

/*
180 / 6 = 30
180 / 6 = 30
90 / 6 = 15
*/
