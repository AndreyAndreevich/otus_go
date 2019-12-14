package app

import (
	"fmt"

	"github.com/AndreyAndreevich/otus_go/hw02/interfaces"
)

// App is main application struct
type App struct {
	unpacker interfaces.Unpacker
}

// NewApp create new App
func NewApp(unpacker interfaces.Unpacker) *App {
	return &App{
		unpacker: unpacker,
	}
}

//Run application
func (m *App) Run(strings []string) error {
	for _, str := range strings {
		res, err := m.unpacker.Unpack(str)
		if err != nil {
			fmt.Printf("%s => %s\n", str, err)
			continue
		}
		fmt.Printf("%s => %s\n", str, res)
	}
	return nil
}
