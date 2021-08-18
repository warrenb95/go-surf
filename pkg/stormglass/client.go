package stormglass

import (
	"context"
	"encoding/json"
	"fmt"
	"github/warrenb95/go-surf/internal/config"
	"github/warrenb95/go-surf/internal/gosurf"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	High = "high"
	Low  = "low"
)

type Client struct {
	Config     config.Config
	HttpClient *http.Client
}

type weatherbody struct {
	Hours []struct {
		AirTemperature   []floatvalue `json:"airTemperature"`
		WaterTemperature []floatvalue `json:"waterTemperature"`
		WaveHeight       []floatvalue `json:"waveHeight"`
		WavePeriod       []floatvalue `json:"wavePeriod"`
		WindSpeed        []floatvalue `json:"windSpeed"`
		Time             time.Time    `json:"time"`
	} `json:"hours"`
}

type tidebody struct {
	Data []struct {
		Height float64   `json:"height"`
		Time   time.Time `json:"time"`
		Type   string    `json:"type"`
	} `json:"data"`
}

type floatvalue struct {
	Value float64 `json:"value"`
}

func (c Client) Get(ctx context.Context, spot gosurf.Spot, params string) (map[string]map[int]gosurf.Spot, error) {
	argUrl := fmt.Sprintf("%slat=%v&lng=%v&params=%s", c.Config.StormglassWeatherURL, spot.Lat, spot.Lng, params)

	var wbody weatherbody
	err := c.getObjFromURL(argUrl, &wbody)
	if err != nil {
		log.Printf("error getting weatherbody from %v", argUrl)
		return nil, fmt.Errorf("error getting weatherbody from %v", argUrl)
	}

	argUrl = fmt.Sprintf("%slat=%v&lng=%v", c.Config.StormglassTideURL, spot.Lat, spot.Lng)
	var tbody tidebody
	err = c.getObjFromURL(argUrl, &tbody)
	if err != nil {
		log.Printf("error getting tidebody from %v", argUrl)
		return nil, fmt.Errorf("error getting tidebody from %v", argUrl)
	}

	return processResponse(wbody, tbody, spot)
}

func (c Client) getObjFromURL(url string, obj interface{}) error {
	r, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		log.Printf("error creating request: %v", err)
		return fmt.Errorf("error creating request: %v", err)
	}
	r.Header.Add("Authorization", "40a29832-fba5-11eb-9642-0242ac130002-40a298a0-fba5-11eb-9642-0242ac130002")

	res, err := c.HttpClient.Do(r)
	if err != nil {
		log.Printf("error sending message: %v", err)
		return fmt.Errorf("error sending message: %v", err)
	}
	defer res.Body.Close()

	log.Printf("stormglass res status %v", res.Status)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("can't read body bytes in stormglass response: %v", err)
		return fmt.Errorf("can't read body bytes in stormglass response: %v", err)
	}
	// log.Println(string(body))

	err = json.Unmarshal(body, obj)
	if err != nil {
		log.Printf("error unmarshalling into resbody: %v", err)
		return fmt.Errorf("error unmarshalling into resbody: %v", err)
	}

	return nil
}

// processResponse will process the data and create a slice of Spots.
// this might not belong here as it is specific to the internal pkg but i
// dunno where to put it.
func processResponse(wbody weatherbody, tbody tidebody, spot gosurf.Spot) (map[string]map[int]gosurf.Spot, error) {
	// spotMap use date string as key YYYY-MM-DD
	spotMap := map[string]map[int]gosurf.Spot{}

	// sort out the weather first and then add the tides
	for _, forecast := range wbody.Hours {
		airtemp := sumValueSlice(forecast.AirTemperature) / float64(len(forecast.AirTemperature))
		waterTemp := sumValueSlice(forecast.WaterTemperature) / float64(len(forecast.WaterTemperature))
		waveHeight := sumValueSlice(forecast.WaveHeight) / float64(len(forecast.WaveHeight))
		wavePeriod := sumValueSlice(forecast.WavePeriod) / float64(len(forecast.WavePeriod))
		windSpeed := sumValueSlice(forecast.WindSpeed) / float64(len(forecast.WindSpeed))

		date := forecast.Time.Format("2000-01-01")

		if _, ok := spotMap[date]; !ok {
			spotMap[date] = map[int]gosurf.Spot{}
		}

		hour := forecast.Time.Hour()

		spotMap[date][hour] = gosurf.Spot{
			Name:       spot.Name,
			Lng:        spot.Lng,
			Lat:        spot.Lat,
			Tide:       gosurf.Tide{},
			Period:     wavePeriod,
			WaveHeight: waveHeight,
			Temperature: gosurf.Temperature{
				Air:   airtemp,
				Water: waterTemp,
			},
			WindSpeed: windSpeed,
		}
	}

	// sort out the tides now
	for _, forecast := range tbody.Data {
		date := forecast.Time.Format("2000-01-01")
		if day, ok := spotMap[date]; ok {
			for _, hour := range day {
				if forecast.Type == Low {
					hour.Tide.Low = forecast.Time
				} else {
					hour.Tide.High = forecast.Time
				}
			}
		}
	}

	return spotMap, nil
}

func sumValueSlice(valueSlice []floatvalue) float64 {
	result := 0.0
	for _, tmp := range valueSlice {
		result += tmp.Value
	}
	return result
}
