package cmd

import (
	"fmt"
	parser "github.com/codecrafters-io/redis-starter-go/app/resp_parser"
	"log"
	"strconv"
	"strings"
	"time"
)

type CommandHandler func(*App, *Command) []byte

var commandHandlers = map[string]CommandHandler{
	"PING": handlePing,
	"ECHO": handleEcho,
	"SET":  handleSet,
	"GET":  handleGet,
}

func handlePing(app *App, cmd *Command) []byte {
	return parser.AppendString([]byte{}, "PONG")
}

func handleEcho(app *App, cmd *Command) []byte {
	if len(cmd.Args) < 1 {
		return cmd.Error(fmt.Errorf("wrong number of arguments for 'echo' command"))
	}
	return parser.AppendBulk([]byte{}, []byte(cmd.Args[0]))
}

func handleSet(app *App, cmd *Command) []byte {
	if len(cmd.Args) != 2 && len(cmd.Args) != 4 {
		return cmd.Error(fmt.Errorf("wrong number of arguments for 'set' command"))
	}

	if len(cmd.Args) == 4 {
		if strings.EqualFold(cmd.Args[2], "px") {
			// milliseconds
			atoi, err := strconv.Atoi(cmd.Args[3])
			if err != nil {
				return cmd.Error(fmt.Errorf("invalid value provided as expiry for 'set' command"))
			}
			expireAt := time.Now().Add(time.Millisecond * time.Duration(atoi))
			app.store.Set(cmd.Args[0], cmd.Args[1], &expireAt)
			return parser.AppendString([]byte{}, "OK")
		}
	}
	app.store.Set(cmd.Args[0], cmd.Args[1], nil)
	log.Printf("DEBUG: SET key='%s', value='%s'", cmd.Args[0], cmd.Args[1])
	return parser.AppendString([]byte{}, "OK")
}

func handleGet(app *App, cmd *Command) []byte {
	if len(cmd.Args) < 1 {
		return cmd.Error(fmt.Errorf("wrong number of arguments for 'get' command"))
	}

	value, exists := app.store.Get(cmd.Args[0])
	if !exists {
		log.Printf("DEBUG: GET key='%s', Value not found", cmd.Args[0])
		return parser.AppendNull([]byte{})
	}

	log.Printf("DEBUG: GET key='%s', value='%s'", cmd.Args[0], value)
	return parser.AppendBulk([]byte{}, []byte(value))
}
