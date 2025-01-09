package cmd

import (
	"fmt"
	parser "github.com/codecrafters-io/redis-starter-go/app/resp_parser"
	"log"
)

type Command struct {
	Name string
	Args []string
}

func ParseCommand(b []byte) (*Command, error) {
	n, resp := parser.Parse(b)
	if n == 0 {
		return nil, fmt.Errorf("invalid command format")
	}

	arr, err := resp.ToStringArr()
	if err != nil {
		return nil, fmt.Errorf("invalid command format: %v", err)
	}

	if len(arr) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	return &Command{
		Name: arr[0],
		Args: arr[1:],
	}, nil
}

func (c *Command) Error(err error) []byte {
	log.Printf("ERROR: %v", err)
	return parser.AppendError([]byte{}, fmt.Sprintf("ERR %v", err))
}
