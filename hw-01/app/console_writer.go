package app

import "fmt"

type ConsoleWriter struct {

}

func (c *ConsoleWriter)Write(args ... interface{}) error {
	_, err := fmt.Println(args ...)
	return err
}
