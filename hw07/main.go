package main

import (
	"log"
	"os"

	"github.com/AndreyAndreevich/otus_go/hw07/envdir"
)

func main() {

	if len(os.Args) < 3 {
		log.Fatal("pls write to correct line: hw07 /path/to/envdir command args...")
	}

	env, err := envdir.ReadDir(os.Args[1])
	if err != nil {
		log.Fatal("incorrect envdir")
	}

	cmd := os.Args[2:]

	code := envdir.RunCmd(cmd, env)
	os.Exit(code)
}
