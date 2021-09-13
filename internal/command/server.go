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
	// for _, spot := range s.Config.Spots {
	// 	// cansurf spot??
	// canSurf, surfMap, err := s.SurfGuru.CanSurf(ctx, spot)
	// if err != nil {
	// 	log.Println("error handlling spot %s, %v", spot.Name, err)
	// 	return err
	// }

	// 	// if can surf then send report to user
	// 	if canSurf {
	// 		g, _ := errgroup.WithContext(ctx)

	// 		for _, s2 := range surfMap {
	// 			g.Go(func() error {
	// 				return s.Twilio.SendAlert(s2.String())
	// 			})
	// 		}

	// 		if err := g.Wait(); err != nil {
	// 			log.Printf("error sending alert: %v", err.Error())
	// 		}
	// 	}
	// }

	genChan := generatePipeline(s.Config.Spots)
	surfChan := s.SurfGuru.CanSurf(ctx, genChan)

	go func() {
		for spotMap := range surfChan {
			for _, spot := range spotMap {
				err := s.Twilio.SendAlert(spot.String())
				if err != nil {
					log.Printf("error send twillio alert of spot: %v error: %v", spot.String(), err)
				}
			}

		}
	}()

	return nil
}

func generatePipeline(spots []gosurf.Spot) <-chan gosurf.Spot {
	outchan := make(chan gosurf.Spot, 10)

	go func() {
		defer close(outchan)
		for _, spot := range spots {
			outchan <- spot
		}
	}()

	return outchan
}
