package algo

import (
	"fmt"
	"math"
)

type Point struct {
	x, y, z float64
}

func crossProduct(a, b Point) Point {
	return Point{
		x: a.y*b.z - a.z*b.y,
		y: a.z*b.x - a.x*b.z,
		z: a.x*b.y - a.y*b.x,
	}
}

func dotProduct(a, b Point) float64 {
	return a.x*b.x + a.y*b.y + a.z*b.z
}

// Параметрические уравнения отрезков и проверка их пересечения
func areSegmentsIntersecting(p1, p2, q1, q2 Point) bool {
	// Вектор для отрезков
	p2p1 := Point{p2.x - p1.x, p2.y - p1.y, p2.z - p1.z}
	q2q1 := Point{q2.x - q1.x, q2.y - q1.y, q2.z - q1.z}

	// Векторное произведение
	numerator := crossProduct(p2p1, q2q1)

	// Проверка пересечения на уровне 3D
	// Проверка, что пересекаются ли два отрезка
	return math.Abs(numerator.x) > 1e-9 || math.Abs(numerator.y) > 1e-9 || math.Abs(numerator.z) > 1e-9
}

func main() {
	// Пример двух отрезков
	p1 := Point{0, 0, 0}
	p2 := Point{1, 1, 1}
	q1 := Point{0, 1, 0}
	q2 := Point{1, 0, 1}

	if areSegmentsIntersecting(p1, p2, q1, q2) {
		fmt.Println("Отрезки пересекаются")
	} else {
		fmt.Println("Отрезки не пересекаются")
	}
}
