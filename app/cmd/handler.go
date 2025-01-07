package cmd

import (
	parser "github.com/codecrafters-io/redis-starter-go/app/resp_parser"
	"log"
	"strings"
)

type App struct {
	Store map[string]string
}

// CommandError represents an error that occurred while processing a command
type CommandError struct {
	Command string
	Err     error
}

func (e *CommandError) Error() string {
	log.Printf("command '%s' failed: %v", e.Command, e.Err)
	return e.Err.Error()
}

// Handle processes incoming Redis commands and returns the appropriate response
func (app *App) Handle(b []byte) []byte {
	var out []byte
	n, resp := parser.Parse(b)
	if n == 0 {
		log.Printf("ERROR: Failed to parse command: invalid format")
		return parser.AppendError(out, "ERR invalid command format")
	}

	arr, err := resp.ToStringArr()
	if err != nil {
		log.Printf("ERROR: Failed to convert command to string array: %v", err)
		return parser.AppendError(out, "ERR invalid command format")
	}

	if len(arr) == 0 {
		log.Printf("ERROR: Received empty command")
		return parser.AppendError(out, "ERR empty command")
	}

	command := strings.ToUpper(arr[0])
	log.Printf("INFO: Processing command: %s", command)

	switch command {
	case "PING":
		return parser.AppendString(out, "PONG")

	case "ECHO":
		if len(arr) < 2 {
			log.Printf("ERROR: ECHO command received without argument")
			return parser.AppendError(out, "ERR wrong number of arguments for 'echo' command")
		}
		return parser.AppendBulk(out, []byte(arr[1]))

	case "SET":
		app.Store[strings.ToLower(arr[1])] = arr[2]
		log.Printf("DEBUG: SET key='%s', value='%s'", arr[1], arr[2])
		return parser.AppendString(out, "OK")

	case "GET":
		s, ok := app.Store[strings.ToLower(arr[1])]
		if !ok {
			log.Printf("DEBUG: GET key='%s', Value not found", arr[1])
			return parser.AppendNull(out)
		}
		log.Printf("DEBUG: GET key='%s', value='%s'", arr[1], s)
		return parser.AppendBulk(out, []byte(s))

	default:
		log.Printf("WARN: Unknown command received: %s", command)
		return parser.AppendError(out, "ERR unknown command '"+command+"'")
	}
}
