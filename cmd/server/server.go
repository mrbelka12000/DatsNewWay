package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"slices"
	"time"

	"DatsNewWay/entity"
)

func main() {

	body, err := os.ReadFile("check.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/next", ServeNext)
	go StartJob()
	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}

var (
	totalPoints = 0
	now         = time.Now()
)

func StartJob() {
	ticker := time.NewTicker(10 * time.Second)

	for {
		<-ticker.C
		fmt.Printf("total points: %v, time from start: %v\n", totalPoints, time.Since(now))

		for i := 0; i < 3; i++ {
			resp.Food = append(resp.Food, entity.Food{
				C:      []int{rand.Intn(180), rand.Intn(180), rand.Intn(60)},
				Points: rand.Intn(200),
			})
		}
		ticker.Reset(10 * time.Second)
	}
}

func ServeNext(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	req := entity.Payload{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < len(resp.Enemies); i++ {
		resp.Enemies[i].Geometry = append(getRandomDir(resp.Enemies[i].Geometry[0]), resp.Enemies[i].Geometry...)
	}

	for _, snake := range req.Snakes {

		for i := 0; i < len(resp.Snakes); i++ {
			if resp.Snakes[i].Id == snake.Id {

				if len(snake.Direction) > 0 {
					resp.Snakes[i].Geometry[0][0] += snake.Direction[0]
					resp.Snakes[i].Geometry[0][1] += snake.Direction[1]
					resp.Snakes[i].Geometry[0][2] += snake.Direction[2]
					resp.Snakes[i].OldDirection = snake.Direction
				} else {
					resp.Snakes[i].Geometry[0][0] += resp.Snakes[i].OldDirection[0]
					resp.Snakes[i].Geometry[0][1] += resp.Snakes[i].OldDirection[1]
					resp.Snakes[i].Geometry[0][2] += resp.Snakes[i].OldDirection[2]
				}
			}

			ss := resp.Snakes[i]
			head := ss.Geometry[0]
			for _, fence := range resp.Fences {
				if isInSameCells(head, fence) {
					resp.Snakes[i].Status = "dead"
					fmt.Println("suka")
				}
			}

			for _, enemy := range resp.Enemies {
				for _, geo := range enemy.Geometry {
					if isInSameCells(head, geo) {
						resp.Snakes[i].Status = "dead"
						fmt.Println("suka умер от чужой змеи")
					}
				}
			}

			var foodID int = -1
			for i := 0; i < len(resp.Food); i++ {
				food := resp.Food[i]
				if isInSameCells(head, food.C) {
					foodID = i
					totalPoints += food.Points
					break
				}
			}

			if foodID != -1 {
				resp.Points += resp.Food[foodID].Points
				resp.Food = slices.Delete(resp.Food, foodID, foodID+1)
			}
		}
	}

	body, err = json.Marshal(resp)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Write(body)
}

var resp entity.Response

func isInSameCells(x, y []int) bool {
	return x[0] == y[0] && x[1] == y[1] && x[2] == y[2]
}

func getRandomDir(head []int) [][]int {
	dirs := [6][]int{
		{1, 0, 0},
		{-1, 0, 0},
		{0, 1, 0},
		{0, -1, 0},
		{0, 0, 1},
		{0, 0, -1},
	}

	random := dirs[rand.Intn(len(dirs))]
	return [][]int{{head[0] + random[0], head[1] + random[1], head[2] + random[2]}}
}
