package main

import (
	"log"

	"github.com/AndreyAndreevich/otus_go/hw02/app"
)

func main() {
	unpacker := &app.RepCharUnpacker{}
	application := app.NewApp(unpacker)
	strings := []string{"a4bc2d5e", "abcd", "45", `qwe\4\5`, `qwe\45`, `qwe\\5`}
	err := application.Run(strings)
	if err != nil {
		log.Fatalln(err)
	}
}
