package app

import (
	"github.com/AndreyAndreevich/otus_go/hw-01/interfaces"
)

// App is main application struct
type App struct {
	clock  interfaces.Clock
	writer interfaces.Writer
}

// NewApp create new App
func NewApp(clock interfaces.Clock, writer interfaces.Writer) *App {
	return &App{
		clock:  clock,
		writer: writer,
	}
}

// Run write time to writer
func (a *App) Run() error {
	time, err := a.clock.GetCurrentTime()
	if err != nil {
		return err
	}
	return a.writer.Write(time)
}
