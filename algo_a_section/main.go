package algo_a_section

import (
	"DatsNewWay/entity"
	"container/heap"
	"fmt"
)

// Position представляет координаты в 3D
type Position struct {
	X, Y, Z int
	Weight  int
}

// Cell представляет клетку в пространстве
type Cell struct {
	Position Position
	Weight   int
	// Для использования в A* (f(n) = g(n) + h(n))
	G, H, F int
}

// PriorityQueue для приоритетной очереди (используется в A*)
type PriorityQueue []*Cell

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].F < pq[j].F
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}
func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Cell))
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// Дистанция Манхэттена между двумя позициями
func heuristic(a, b Position) int {
	return abs(a.X-b.X) + abs(a.Y-b.Y) + abs(a.Z-b.Z)
}

// Функция для расчета стоимости пути с учетом веса секции
func sectionWeight(pos Position, sectionPriorities map[Position]int) int {
	sectionKey := Position{X: pos.X / 30, Y: pos.Y / 30, Z: pos.Z / 30} // Разделение на секции 30x30x30
	return sectionPriorities[sectionKey]                                // Вес секции
}

// Поиск ближайших высоковесовых клеток
func findHighestWeightCells(weights map[Position]int) []Position {
	var highWeightCells []Position
	// Проходим по всем клеткам и выбираем те, у которых вес максимальный
	maxWeight := -1
	for pos, weight := range weights {
		if weight > maxWeight {
			maxWeight = weight
			highWeightCells = []Position{pos}
		} else if weight == maxWeight {
			highWeightCells = append(highWeightCells, pos)
		}
	}
	return highWeightCells
}

// Алгоритм A* для поиска путей, фокусируясь на высоковесовых клетках
func aStar(start Position, cubeSize [3]int, weights map[Position]int, sectionPriorities map[Position]int) Position {
	pq := &PriorityQueue{}
	heap.Init(pq)

	// Начальная клетка
	heap.Push(pq, &Cell{start, weights[start], 0, 0, 0})

	// Для отслеживания посещенных клеток
	visited := make(map[Position]bool)

	// Направления для движения в 3D (включая все 6 соседей)
	directions := [][3]int{
		{1, 0, 0}, {-1, 0, 0}, {0, 1, 0}, {0, -1, 0}, {0, 0, 1}, {0, 0, -1},
	}

	for pq.Len() > 0 {
		cell := heap.Pop(pq).(*Cell)

		// Отслеживаем клетки с высокими весами
		highWeightCells := findHighestWeightCells(weights)

		// Если текущая клетка является высоковесной целью
		if contains(highWeightCells, cell.Position) {
			// Вычисление пути с учётом веса
			return cell.Position
		}

		if visited[cell.Position] {
			continue
		}
		visited[cell.Position] = true

		// Обходим все соседние клетки
		for _, dir := range directions {
			newPos := Position{X: cell.Position.X + dir[0], Y: cell.Position.Y + dir[1], Z: cell.Position.Z + dir[2]}

			// Проверяем границы
			if newPos.X < 0 || newPos.X >= cubeSize[0] || newPos.Y < 0 || newPos.Y >= cubeSize[1] || newPos.Z < 0 || newPos.Z >= cubeSize[2] {
				continue
			}

			// Добавляем вес клетки и вес секции
			newWeight := cell.G + weights[newPos] + sectionWeight(newPos, sectionPriorities)

			// Вычисляем эвристику (манхэттенское расстояние)
			h := heuristic(newPos, start) // Используем манхэттенское расстояние как эвристику

			// Если клетка еще не была посещена
			if !visited[newPos] {
				heap.Push(pq, &Cell{
					Position: newPos,
					Weight:   newWeight,
					G:        newWeight,
					H:        h,
					F:        newWeight + h,
				})
			}
		}
	}

	return Position{} // Если путь не найден
}

// Утилита для вычисления абсолютного значения
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Функция для проверки, содержит ли список координат целевую позицию
func contains(positions []Position, target Position) bool {
	for _, pos := range positions {
		if pos == target {
			return true
		}
	}
	return false
}

func getSections(head Position, cubeSize [3]int) int {
	centreX, centreY, centreZ := cubeSize[0]/2, cubeSize[1]/2, cubeSize[2]/2

	quadX, quadY, quadZ := centreX/2, centreY/2, centreZ/2

	subquadX, subquadY, subquadZ := quadX/2, quadY/2, quadZ/2

	if head.X >= centreX && head.X <= centreX+quadX && head.Y >= centreY && head.Y <= centreY+quadY && head.Z >= centreZ && head.Z <= centreZ+quadZ {
		return 5
	}

	if head.X >= quadX && head.X <= quadX+subquadX && head.Y >= quadY && head.Y <= quadY+subquadY && head.Z >= quadZ && head.Z <= quadZ+subquadZ {
		return 3
	}

	return 1
}

var depthCounts = make(map[int]int)
var sectionCounts = make(map[int][]int)

// Рекурсивная функция для деления куба на 4 части и подсчета точек в каждой секции
func getSectionsPriorityByPoints(points []*Position, depth int, maxDepth int, prevX, prevY, prevZ int) {
	if depth > maxDepth {
		return
	}

	// Рассчитываем размеры для секций
	midX := prevX / 2
	midY := prevY / 2
	midZ := prevZ / 2

	// Вес зависит от глубины: чем глубже деление, тем меньше вес

	// Инициализация подсчета точек в каждой из 8-х секций
	for i := 0; i < 8; i++ {
		sectionCounts[depth] = []int{0, 0, 0, 0, 0, 0, 0, 0}
	}

	// Подсчитываем точки на текущем уровне
	depthCounts[depth] = 0

	// Обрабатываем точки и назначаем им веса в зависимости от их положения в секции
	for _, point := range points {
		// Определяем, в какую из 8 частей попадает точка
		if point.X >= prevX && point.X < midX && point.Y >= prevY && point.Y < midY && point.Z >= prevZ && point.Z < midZ {
			sectionCounts[depth][0]++
		} else if point.X >= midX && point.X < prevX && point.Y >= prevY && point.Y < midY && point.Z >= prevZ && point.Z < midZ {
			sectionCounts[depth][1]++
		} else if point.X >= prevX && point.X < midX && point.Y >= midY && point.Y < prevY && point.Z >= prevZ && point.Z < midZ {
			sectionCounts[depth][2]++
		} else if point.X >= midX && point.X < prevX && point.Y >= midY && point.Y < prevY && point.Z >= prevZ && point.Z < midZ {
			sectionCounts[depth][3]++
		} else if point.X >= prevX && point.X < midX && point.Y >= prevY && point.Y < midY && point.Z >= midZ && point.Z < prevZ {
			sectionCounts[depth][4]++
		} else if point.X >= midX && point.X < prevX && point.Y >= prevY && point.Y < midY && point.Z >= midZ && point.Z < prevZ {
			sectionCounts[depth][5]++
		} else if point.X >= prevX && point.X < midX && point.Y >= midY && point.Y < prevY && point.Z >= midZ && point.Z < prevZ {
			sectionCounts[depth][6]++
		} else if point.X >= midX && point.X < prevX && point.Y >= midY && point.Y < prevY && point.Z >= midZ && point.Z < prevZ {
			sectionCounts[depth][7]++
		}
	}

	// Дополнительный вес в зависимости от количества точек в секции
	for i := 0; i < 8; i++ {
		// Если в секции много точек, увеличиваем вес
		if sectionCounts[depth][i] > 0 {
			// Увеличиваем вес каждой точки в этой секции пропорционально количеству точек
			for idx, point := range points {
				if point.X >= prevX && point.X < midX && point.Y >= prevY && point.Y < midY && point.Z >= prevZ && point.Z < midZ && i == 0 ||
					point.X >= midX && point.X < prevX && point.Y >= prevY && point.Y < midY && point.Z >= prevZ && point.Z < midZ && i == 1 ||
					point.X >= prevX && point.X < midX && point.Y >= midY && point.Y < prevY && point.Z >= prevZ && point.Z < midZ && i == 2 ||
					point.X >= midX && point.X < prevX && point.Y >= midY && point.Y < prevY && point.Z >= prevZ && point.Z < midZ && i == 3 ||
					point.X >= prevX && point.X < midX && point.Y >= prevY && point.Y < midY && point.Z >= midZ && point.Z < prevZ && i == 4 ||
					point.X >= midX && point.X < prevX && point.Y >= prevY && point.Y < midY && point.Z >= midZ && point.Z < prevZ && i == 5 ||
					point.X >= prevX && point.X < midX && point.Y >= midY && point.Y < prevY && point.Z >= midZ && point.Z < prevZ && i == 6 ||
					point.X >= midX && point.X < prevX && point.Y >= midY && point.Y < prevY && point.Z >= midZ && point.Z < prevZ && i == 7 {
					points[idx].Weight += sectionCounts[depth][i]
				}
			}
		}
	}

	// Рекурсивно делим 8 частей
	getSectionsPriorityByPoints(points, depth+1, maxDepth, midX, midY, midZ)
}

func GetNextDirection(r entity.Response) (obj entity.Payload) {
	// Размер куба
	//cubeSize := [3]int{r.MapSize[0], r.MapSize[1], r.MapSize[2]}

	positions := make([]*Position, len(r.Food))

	for idx, food := range r.Food {
		positions[idx] = &Position{
			X: food.C[0],
			Y: food.C[1],
			Z: food.C[2],
		}
	}

	getSectionsPriorityByPoints(positions, 1, 5, r.MapSize[0], r.MapSize[1], r.MapSize[2])
	fmt.Println(positions)
	return obj
	// Вес клеток (например, случайные или заранее заданные)
	//weights := make(map[Position]int)
	//for x := 0; x < cubeSize[0]; x++ {
	//	for y := 0; y < cubeSize[1]; y++ {
	//		for z := 0; z < cubeSize[2]; z++ {
	//			weights[Position{x, y, z}] = rand.Intn(5) // Пример: случайные веса от 0 до 4
	//		}
	//	}
	//}

	// Вес секций (например, более высокие веса для определённых секций)
	//sectionPriorities := make(map[Position]int)
	//for x := 0; x < cubeSize[0]/30; x++ {
	//	for y := 0; y < cubeSize[1]/30; y++ {
	//		for z := 0; z < cubeSize[2]/30; z++ {
	//			sectionPriorities[Position{x, y, z}] = rand.Intn(10) // Пример: случайный вес для секции
	//		}
	//	}
	//}

	// Начальная позиция
	//start := Position{0, 0, 0}

	// Запуск поиска пути
	//result := aStar(start, cubeSize, weights, sectionPriorities)

	return obj
}
