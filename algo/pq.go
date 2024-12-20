package algo

type PQ []info

func (pq PQ) Len() int {
	return len(pq)
}

func (pq PQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

// min--heap
func (pq PQ) Less(i, j int) bool {
	return pq[i].cost+pq[i].heur < pq[j].cost+pq[j].heur
}

func (pq *PQ) Push(x interface{}) {
	tmp := x.(info)
	*pq = append(*pq, tmp)
}

func (pq *PQ) Pop() interface{} {
	n := len(*pq)
	tmp := (*pq)[n-1]
	*pq = (*pq)[:n-1]
	return tmp
}
