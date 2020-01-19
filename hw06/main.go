package main

import (
	"flag"
	"log"
	"os"

	copy2 "github.com/AndreyAndreevich/otus_go/hw06/copy"
)

func main() {
	from := flag.String("from", "", "from file path")
	to := flag.String("to", "", "to file path")
	offset := flag.Int("offset", 0, "file offset (default 0)")
	limit := flag.Int("limit", 0, "file limit (default all file size")

	flag.Parse()
	if *from == "" || *to == "" {
		flag.Usage()
		os.Exit(1)
	}

	err := copy2.Copy(*from, *to, *limit, *offset)
	if err != nil {
		log.Fatalln("File doesn't copy. Error: ", err)
	}
}
