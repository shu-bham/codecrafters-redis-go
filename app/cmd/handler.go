package cmd

import (
	parser "github.com/codecrafters-io/redis-starter-go/app/resp"
	"strings"
)

func Handle(b []byte) []byte {
	_, resp := parser.Parse(b)
	arr, err := resp.ToStringArr()
	if err == nil {
		c := arr[0]
		if strings.EqualFold(c, "echo") {
			return []byte(arr[1])
		}
	}

	return []byte("+PONG\r\n")
}
