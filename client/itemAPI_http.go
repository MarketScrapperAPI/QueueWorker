package client

import (
	"bytes"
	"net/http"

	"github.com/MarketScrapperAPI/QueueWorker/models"
)

type MarketAPIHTTPClient struct {
	client *http.Client
	apiUrl string
}

func NewMarketAPIHTTPClient(apiUrl string) MarketAPIHTTPClient {

	return MarketAPIHTTPClient{
		client: &http.Client{},
		apiUrl: apiUrl,
	}
}

func (c *MarketAPIHTTPClient) AddItem(message models.Message) (int, error) {
	myJson := bytes.NewBuffer([]byte(""))
	resp, err := c.client.Post(c.apiUrl+"item", "application/json", myJson)
	if err != nil {
		return resp.StatusCode, err
	}

	return resp.StatusCode, nil
}
