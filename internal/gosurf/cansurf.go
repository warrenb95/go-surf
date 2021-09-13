package gosurf

import (
	"context"
	"log"
	"time"
)

const params = "wavePeriod,waveHeight,waterTemperature,airTemperature,windSpeed"

type SurfGuru struct {
	Client Fetcher
}

func (s SurfGuru) CanSurf(ctx context.Context, inChan <-chan Spot) <-chan map[string]*Spot {
	outchan := make(chan map[string]*Spot)

	go func() {
		defer close(outchan)
		for spot := range inChan {
			ctx, cancel := context.WithTimeout(ctx, time.Second*5)
			defer cancel()
			spotMap, err := s.Client.Get(ctx, spot, params)
			if err != nil {
				log.Printf("error calling surfguru fetcher: %v", err)
			}
			outchan <- calculate(spot, spotMap)
		}

	}()

	return outchan
}

func calculate(spot Spot, spotMap map[string]map[int]*Spot) map[string]*Spot {
	canSurfMap := map[string]*Spot{}

	for date, day := range spotMap {
		log.Printf("spot date key: %v", date)
		// calculate the wave period and tide
		for hour, sp := range day {
			if sp.Period < spot.Period {
				break
			}

			// is the tide ok?
			if spot.Tide.Position == AllTide {
				canSurfMap[date] = sp
				break
			} else if spot.Tide.Position == LowTide {
				tideHour := sp.Tide.Low.Hour()
				if hour < tideHour-sp.Tide.TimeBefore || hour > tideHour+sp.Tide.TimeAfter {
					break
				}
			} else if spot.Tide.Position == HighTide {
				tideHour := sp.Tide.High.Hour()
				if hour < tideHour-sp.Tide.TimeBefore || hour > tideHour+sp.Tide.TimeAfter {
					break
				}
			}

			canSurfMap[date] = sp
		}
	}

	return canSurfMap
}
