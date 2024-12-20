package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"

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
	resp.MapSize[0] = 30
	resp.MapSize[1] = 30

	for i := 0; i < len(resp.Food); i++ {
		if resp.Food[i].C[0] > 30 || resp.Food[i].C[1] > 30 {
			resp.Food = slices.Delete(resp.Food, i, i+1)
		}
	}

	for i := 0; i < len(resp.Fences); i++ {
		if resp.Fences[i][0] > 30 || resp.Fences[i][1] > 30 {
			resp.Fences = slices.Delete(resp.Fences, i, i+1)
		}
	}

	http.HandleFunc("/next", ServeNext)
	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
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
			if resp.Snakes[i].Id == snake.Id && len(snake.Direction) > 0 {
				resp.Snakes[i].Geometry[0][0] += snake.Direction[0]
				resp.Snakes[i].Geometry[0][1] += snake.Direction[1]
				resp.Snakes[i].Geometry[0][2] += snake.Direction[2]
			}

			ss := resp.Snakes[i]
			for _, fence := range resp.Fences {
				if ss.Geometry[0][0] == fence[0] && ss.Geometry[0][1] == fence[1] && ss.Geometry[0][2] == fence[2] && ss.Geometry[0][3] == fence[2] {
					resp.Snakes[i].Status = "dead"
					fmt.Println("Snake dead", snake.Id)
				}
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
