package expectation

import "github.com/MartinSimango/dstruct"

type Response struct {
	Body          any                     `json:"body"`
	GeneratedBody dstruct.GeneratedStruct `json:"-"`
	Status        int                     `json:"status"`
	Headers       *map[string][]string    `json:"headers"`
	Cookies       *[]Cookie               `json:"cookies"`
}
