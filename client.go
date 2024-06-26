package dolus

import (
	"context"

	"github.com/DolusMockServer/dolus/internal/api"
	"github.com/DolusMockServer/dolus/pkg/expectation"
)

// Client to the Dolus server
type Client struct {
	client  *api.Client
	context context.Context
	mapper  api.Mapper
}

// AddExpectation to the Dolus server
func (c *Client) AddExpectation(expectation expectation.Expectation) error {
	exp, err := c.mapper.MapToApiExpectation(expectation)
	if err != nil {
		return err
	}
	c.client.CreateExpectation(c.context, *exp)

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
