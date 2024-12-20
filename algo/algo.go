package algo

import "DatsNewWay/entity"

func GetNextDirection(r entity.Response) {

}

func getDirection(head, target [3]int) [3]int {
	if head[0] != target[0] {
		if head[0] < target[0] {
			return [3]int{1, 0, 0}
		}
		return [3]int{-1, 0, 0}
	}

	if head[1] != target[1] {
		if head[1] < target[1] {
			return [3]int{0, 1, 0}
		}
		return [3]int{0, -1, 0}
	}

	if head[2] != target[2] {
		if head[2] < target[2] {
			return [3]int{0, 0, 1}
		}
		return [3]int{0, 0, -1}
	}

	return [3]int{0, 0, 0} // Return this if head equals target
}

func getManhattanDistance(x, y [3]int) int {
	return abs(x[0]-y[0]) + abs(x[1]-y[1]) + abs(x[2]-y[2])
}

func abs(a int) int {
	if a < 0 {
		return a * -1
	}
	return a
}
