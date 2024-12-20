package main

import (
	"encoding/json"
	"fmt"
	"io"
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
				if head[0] == fence[0] && head[1] == fence[1] && head[2] == fence[2] && head[3] == fence[2] {
					resp.Snakes[i].Status = "dead"
					fmt.Println("suka")
				}
			}

			var foodID int = -1
			for i := 0; i < len(resp.Food); i++ {
				food := resp.Food[i]
				if head[0] == food.C[0] && head[1] == food.C[1] && head[2] == food.C[2] {
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
