package publish

import (
	"context"
	"flag"
	"fmt"

	"github.com/sfomuseum/go-activitypub/uris"
	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	AccountsDatabaseURI   string
	FeedsDatabaseURI  string
	PostsDatabaseURI      string
	DeliveryQueueURI      string
	AccountName           string
	FeedURIs []string
	URIs                  *uris.URIs
	Verbose               bool
}

func OptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "ACTIVITYPUB")

	if err != nil {
		return nil, fmt.Errorf("Failed to derive flags from environment variables, %w", err)
	}

	uris_table := uris.DefaultURIs()
	uris_table.Hostname = hostname
	uris_table.Insecure = insecure

	opts := &RunOptions{
		AccountsDatabaseURI:   accounts_database_uri,
		FeedsDatabaseURI:  feeds_database_uri,
		PostsDatabaseURI:      posts_database_uri,
		DeliveryQueueURI:      delivery_queue_uri,
		AccountName:           account_name,
		FeedsURIs: feed_uris,
		URIs:                  uris_table,
		Verbose:               verbose,
	}

	return opts, nil
}
