package cmd

import (
	"fmt"
	"strings"
)

type App struct {
	store Storage
}

func NewApp() *App {
	return &App{
		store: NewInMemoryStorage(),
	}
}

func (app *App) Handle(b []byte) []byte {
	cmd, err := ParseCommand(b)
	if err != nil {
		return cmd.Error(err)
	}

	handler, exists := commandHandlers[strings.ToUpper(cmd.Name)]
	if !exists {
		return cmd.Error(fmt.Errorf("unknown command '%s'", cmd.Name))
	}

	return handler(app, cmd)
}
