package dolus

import (
	"context"

	"github.com/DolusMockServer/dolus/internal/api"
	"github.com/DolusMockServer/dolus/pkg/expectation/models"
)

// Client to the Dolus server
type Client struct {
	client  *api.Client
	context context.Context
	mapper  api.Mapper
}

// AddExpectation to the Dolus server
func (c *Client) AddExpectation(expectation models.Expectation) error {
	exp, err := c.mapper.MapCueExpectation(expectation)
	if err != nil {
		return err
	}
	c.client.PostV1DolusExpectations(c.context, *exp)

	return nil
}

// Connect to the Dolus server
func Connect(ctx context.Context, address string) (*Client, error) {
	client, err := api.NewClient(address)
	if err != nil {
		return nil, err
	}
	return &Client{
		client:  client,
		context: ctx,
		mapper:  api.NewMapper(),
	}, nil
}
