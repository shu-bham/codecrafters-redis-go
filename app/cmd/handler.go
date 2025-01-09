package cmd

import (
	"fmt"
	parser "github.com/codecrafters-io/redis-starter-go/app/resp_parser"
	"log"
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
	if len(cmd.Args) < 2 {
		return cmd.Error(fmt.Errorf("wrong number of arguments for 'set' command"))
	}
	app.store.Set(cmd.Args[0], cmd.Args[1])
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
