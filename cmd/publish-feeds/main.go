package main

import (
	"context"
	"log"

	"github.com/sfomuseum/go-activitypub-feeds/app/feeds/publish"
)

func main() {

	ctx := context.Background()
	err := publish.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to publish feeds, %v", err)
	}
}
