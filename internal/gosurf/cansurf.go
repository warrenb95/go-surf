package gosurf

import (
	"context"
	"log"
)

const params = "wavePeriod,waveHeight,waterTemperature,airTemperature,windSpeed"

type SurfGuru struct {
	Client Fetcher
}

func (s SurfGuru) CanSurf(ctx context.Context, spot Spot) (bool, error) {
	_, err := s.Client.Get(ctx, spot, params)
	if err != nil {
		log.Printf("error calling surfguru fetcher: %v", err)
		return false, err
	}
	return true, nil
}
