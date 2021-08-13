package main

import (
	"github/warrenb95/go-surf/internal/command"
	"log"
)

func main() {
	log.Println("starting go-surf")
	defer log.Println("stopped go-surf")

	s := command.Server{}
	s.Run()
}
