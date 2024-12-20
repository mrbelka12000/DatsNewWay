package client

import (
	"DatsNewWay/entity"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type Client struct {
	httpClient *http.Client
	token      string
}

const (
	testDomain = "https://games-test.datsteam.dev/play/snake3d"
	domain     = "https://games.datsteam.dev/play/snake3d"
)

func NewClient(token string) *Client {
	return &Client{
		httpClient: http.DefaultClient,
		token:      token,
	}
}

func (c *Client) Get(ctx context.Context, payload entity.Payload) (entity.Response, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return entity.Response{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, testDomain, bytes.NewReader(body))
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

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return entity.Response{}, err
	}

	result := entity.Response{}

	if err = json.Unmarshal(bytes, &result); err != nil {
		return entity.Response{}, err
	}

	return result, nil
}
