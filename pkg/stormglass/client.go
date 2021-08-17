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

const url = "https://api.stormglass.io/v1/weather/point?"

type Client struct {
	Config     config.Config
	HttpClient *http.Client
}

type resbody struct {
	Hours []struct {
		AirTemperature   []floatvalue `json:"airTemperature"`
		WaterTemperature []floatvalue `json:"waterTemperature"`
		WaveHeight       []floatvalue `json:"waveHeight"`
		WavePeriod       []floatvalue `json:"wavePeriod"`
		WindSpeed        []floatvalue `json:"windSpeed"`
		Time             time.Time    `json:"time"`
	} `json:"hours"`
}

type floatvalue struct {
	Value float64 `json:"value"`
}

func (c Client) Get(ctx context.Context, spot gosurf.Spot, params string) ([]gosurf.Spot, error) {
	argUrl := fmt.Sprintf("%slat=%v&lng=%v&params=%s", c.Config.StormglassURL, spot.Lat, spot.Lng, params)

	r, err := http.NewRequest("GET", argUrl, http.NoBody)
	if err != nil {
		log.Printf("error creating request: %v", err)
		return []gosurf.Spot{}, fmt.Errorf("error creating request: %v", err)
	}
	r.Header.Add("Authorization", "40a29832-fba5-11eb-9642-0242ac130002-40a298a0-fba5-11eb-9642-0242ac130002")

	res, err := c.HttpClient.Do(r)
	if err != nil {
		log.Printf("error sending message: %v", err)
		return []gosurf.Spot{}, fmt.Errorf("error sending message: %v", err)
	}
	defer res.Body.Close()

	log.Printf("stormglass res status %v", res.Status)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("can't read body bytes in stormglass response: %v", err)
		return []gosurf.Spot{}, fmt.Errorf("can't read body bytes in stormglass response: %v", err)
	}
	// log.Println(string(body))

	var resBody resbody
	err = json.Unmarshal(body, &resBody)
	if err != nil {
		log.Printf("error unmarshalling into resbody: %v", err)
		return []gosurf.Spot{}, fmt.Errorf("error unmarshalling into resbody: %v", err)
	}

	log.Println(resBody)

	return processResponse(resBody, spot)
}

func processResponse(response resbody, spot gosurf.Spot) ([]gosurf.Spot, error) {
	spotSlice := []gosurf.Spot{}

	return spotSlice, nil
}
