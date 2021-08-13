package command

import "log"

type Server struct{}

func (s Server) Run() {
	log.Println("starting server")
	defer log.Println("stopped server")
}
