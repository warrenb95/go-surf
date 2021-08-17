package command

import (
	"context"
	"github/warrenb95/go-surf/internal/config"
	"github/warrenb95/go-surf/internal/gosurf"
	"github/warrenb95/go-surf/pkg/twilio"
	"log"
)

type Server struct {
	Config   config.Config
	Twilio   twilio.Client
	SurfGuru gosurf.SurfGuru
}

func (s Server) Run(ctx context.Context) error {
	log.Println("starting server")
	defer log.Println("stopped server")

	// loop over spots
	for _, spot := range s.Config.Spots {
		// cansurf spot??
		_, err := s.SurfGuru.CanSurf(ctx, spot)
		if err != nil {
			log.Println("error handlling spot %s, %v", spot.Name, err)
			return err
		}

		// if can surf then send report to user
		// if canSurf {
		// 	s.Twilio.SendAlert(spot.String())
		// }
	}

	return nil
}
