package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"DatsNewWay/entity"
)

type Client struct {
	httpClient *http.Client
	token      string
}

const (
	testDomain = "https://games-test.datsteam.dev/play/snake3d/player/move"
	domain     = "https://games.datsteam.dev/play/snake3d"
)

func NewClient(token string) *Client {
	return &Client{
		httpClient: http.DefaultClient,
		token:      token,
	}
}

func (c *Client) Get(ctx context.Context, payload entity.Payload) (entity.Response, error) {
	//fileData, err := os.ReadFile("check.json")
	//if err != nil {
	//	return entity.Response{}, err
	//}
	//
	//obj := entity.Response{}
	//err = json.Unmarshal(fileData, &obj)
	//if err != nil {
	//	return entity.Response{}, err
	//}
	//
	//return obj, nil

	body, err := json.Marshal(payload)
	if err != nil {
		return entity.Response{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, testDomain, bytes.NewReader(body))
	if err != nil {
		return entity.Response{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return entity.Response{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return entity.Response{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return entity.Response{}, err
	}

	go func() {
		file, err := os.Create(fmt.Sprintf("check/%v_%v.json", time.Now().Minute(), time.Now().Unix()))
		if err != nil {
			fmt.Println(err, "create file")
			return
		}
		defer file.Close()

		_, err = file.Write(bytes)
		if err != nil {
			fmt.Println(err, "write file")
		}
	}()

	result := entity.Response{}

	if err = json.Unmarshal(bytes, &result); err != nil {
		return entity.Response{}, err
	}

	return result, nil
}
