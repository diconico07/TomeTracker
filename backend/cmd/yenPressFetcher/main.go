package main

import (
	"context"
	"os"

	tomefetcher "github.com/diconico07/TomeTracker/pkg/tomeFetcher"
	yenpressfetcher "github.com/diconico07/TomeTracker/pkg/yenPressFetcher"
	"github.com/rs/zerolog"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	baseUrl := os.Getenv("BASE_URL")
	if baseUrl == "" {
		baseUrl = "http://localhost:8080"
	}
	yenPressFetcher := yenpressfetcher.New()
	fetcher := tomefetcher.New(baseUrl, yenPressFetcher)
	fetcher.ProcessAll(context.Background())
}
