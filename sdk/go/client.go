package dolus

import (
	"github.com/DolusMockServer/dolus-expectations/pkg/dolus"
	"github.com/DolusMockServer/dolus/internal/client"
)

type Client struct {
	c *client.Client
}

// check we implement the interface

func (c *Client) AddExpectation(expectation dolus.Expectation) error {
	panic("NOT Implemented")
}

// connect to the dolus Mock Server
func Connect() (*Client, error) {
	client, err := client.NewClient("http://localhost:1080")
	if err != nil {
		return nil, err
	}

	return &Client{
		c: client,
	}, nil
}

func test() {
	// c, _ := Connect()
	// c.AddExpectation(expectation.DolusExpectation{}.CueExpectation)
}
