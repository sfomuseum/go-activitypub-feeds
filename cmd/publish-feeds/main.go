package main

import (
	"context"
	"os"

	"github.com/sfomuseum/go-activitypub-feeds/app/feeds/publish"
	"github.com/sfomuseum/go-activitypub/slog"
)

func main() {

	ctx := context.Background()
	logger := slog.Default()

	err := publish.Run(ctx, logger)

	if err != nil {
		logger.Error("Failed to publish feeds", "error", err)
		os.Exit(1)
	}
}
