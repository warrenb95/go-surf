package main

import (
	"context"
	"fmt"
	"github/warrenb95/go-surf/internal/command"
	"github/warrenb95/go-surf/internal/config"
	"github/warrenb95/go-surf/internal/gosurf"
	"github/warrenb95/go-surf/pkg/stormglass"
	"github/warrenb95/go-surf/pkg/twilio"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"
)

func main() {
	log.Println("starting go-surf")
	defer log.Println("stopped go-surf")

	ctx := context.Background()
	g, ctx := errgroup.WithContext(ctx)

	// Load the config
	filename, err := filepath.Abs("deployment/go-surf/values.yaml")
	if err != nil {
		log.Fatalf("cannot find config filepath, error: %v", err)
	}
	cfg, err := config.Parse(ctx, filename)
	if err != nil {
		log.Fatal(err)
	}

	twilioToken := os.Getenv("TWILIO_TOKEN")
	if twilioToken == "" {
		log.Fatalln("TWILIO_TOKEN not in env")
	}

	twilioURL := fmt.Sprintf("https://%s:%s@api.twilio.com/2010-04-01/Accounts/AC3596bf01c9b1d9b42a9f4c330a65b4ba/Messages.json",
		"AC3596bf01c9b1d9b42a9f4c330a65b4ba", twilioToken)

	s := command.Server{
		Config: cfg,
		Twilio: twilio.Client{
			Config:             cfg,
			HttpClient:         &http.Client{},
			UrlString:          twilioURL,
			TargetMobileNumber: cfg.TargetMobileNumber,
		},
		SurfGuru: gosurf.SurfGuru{
			Client: stormglass.Client{
				Config:     cfg,
				HttpClient: &http.Client{},
			},
		},
	}

	g.Go(func() error {
		return s.Run(ctx)
	})

	if err := g.Wait(); err != nil {
		log.Printf("error funning server: %v", err.Error())
	}
}
