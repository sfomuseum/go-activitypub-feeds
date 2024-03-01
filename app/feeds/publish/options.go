package publish

import (
	"context"
	"flag"
	"fmt"

	"github.com/sfomuseum/go-activitypub/uris"
	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	AccountsDatabaseURI             string
	FeedsPublicationLogsDatabaseURI string
	PostsDatabaseURI                string
	FollowersDatabaseURI            string
	DeliveriesDatabaseURI           string
	DeliveryQueueURI                string
	AccountName                     string
	Mode                            string
	FeedURIs                        []string
	URIs                            *uris.URIs
	Verbose                         bool
	MaxPostsPerFeed                 int
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
		AccountsDatabaseURI:             accounts_database_uri,
		FeedsPublicationLogsDatabaseURI: feeds_publication_logs_database_uri,
		FollowersDatabaseURI:            followers_database_uri,
		DeliveriesDatabaseURI:           deliveries_database_uri,
		PostsDatabaseURI:                posts_database_uri,
		DeliveryQueueURI:                delivery_queue_uri,
		AccountName:                     account_name,
		Mode:                            mode,
		FeedURIs:                        feed_uris,
		MaxPostsPerFeed:                 max_posts_per_feed,
		URIs:                            uris_table,
		Verbose:                         verbose,
	}

	return opts, nil
}
