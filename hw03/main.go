package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/AndreyAndreevich/otus_go/hw03/app"
)

const exampleFileName = "text.txt"

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fileNamePtr := flag.String("filename", path.Join(dir, exampleFileName), "filename")

	flag.Parse()

	fmt.Println(dir)
	fmt.Println(*fileNamePtr)

	text, err := ioutil.ReadFile(*fileNamePtr)
	if err != nil {
		log.Fatal(err)
	}

	top := app.Top10(string(text))
	for i, word := range top {
		fmt.Printf("Top %d => %s\n", i+1, word)
	}
}
