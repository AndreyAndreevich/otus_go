package main

import (
	"log"

	"github.com/AndreyAndreevich/otus_go/hw01/app"
)

func main() {
	clock := app.NewNtpClock("0.beevik-ntp.pool.ntp.org")
	writer := new(app.ConsoleWriter)
	application := app.NewApp(clock, writer)

	err := application.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
