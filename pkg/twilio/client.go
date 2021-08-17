package twilio

import (
	"github/warrenb95/go-surf/internal/config"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Client struct {
	Config             config.Config
	HttpClient         *http.Client
	UrlString          string
	TargetMobileNumber string
}

func (c *Client) SendAlert(alertBody string) error {
	payload := url.Values{}
	payload.Set("To", c.TargetMobileNumber)
	payload.Set("MessagingServiceSid", c.Config.MessagingServiceSid)
	payload.Set("Body", alertBody)

	r, err := http.NewRequest("POST", c.UrlString, strings.NewReader(payload.Encode())) // URL-encoded payload
	if err != nil {
		log.Fatal(err)
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(payload.Encode())))

	res, err := c.HttpClient.Do(r)
	if err != nil {
		log.Printf("error sending message: %v", err)
	}
	defer res.Body.Close()

	log.Printf("send alert %v\nstatus %v", alertBody, res.Status)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(body))

	return nil
}
