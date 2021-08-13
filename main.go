package main

import "log"

func main() {
	log.Println("starting go-surf")
	defer log.Println("stopped go-surf")
}
