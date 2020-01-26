package main

import (
	"log"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/calendar"
	"github.com/AndreyAndreevich/otus_go/calendar/internal/pkg/memorystorage"
)

func main() {
	storage := memorystorage.New()
	currentCalendar := calendar.New(storage)

	if err := currentCalendar.Run(); err != nil {
		log.Fatal(err)
	}
}
