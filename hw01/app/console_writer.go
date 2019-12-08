package app

import "fmt"

// ConsoleWriter write info to console
type ConsoleWriter struct {
}

func (c *ConsoleWriter) Write(args ...interface{}) error {
	_, err := fmt.Println(args...)
	return err
}
