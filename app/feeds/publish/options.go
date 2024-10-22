package publish

import (
	"context"
	"flag"
	"fmt"
	"html/template"

	"github.com/sfomuseum/go-activitypub-feeds/templates/html"
	"github.com/sfomuseum/go-activitypub/uris"
	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	AccountsDatabaseURI             string
	ActivitiesDatabaseURI           string
	PostsDatabaseURI                string
	PostTagsDatabaseURI             string
	FollowersDatabaseURI            string
	DeliveriesDatabaseURI           string
	FeedsPublicationLogsDatabaseURI string
	DeliveryQueueURI                string
	AccountName                     string
	Mode                            string
	FeedURIs                        []string
	URIs                            *uris.URIs
	Verbose                         bool
	MaxPostsPerFeed                 int
	MaxAttempts                     int
	Templates                       *template.Template
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

	t, err := html.LoadTemplates(ctx, html.FS)

	if err != nil {
		return nil, fmt.Errorf("Failed to load templates, %w", err)
	}

	opts := &RunOptions{
		AccountsDatabaseURI:             accounts_database_uri,
		ActivitiesDatabaseURI:           activities_database_uri,
		FollowersDatabaseURI:            followers_database_uri,
		DeliveriesDatabaseURI:           deliveries_database_uri,
		PostsDatabaseURI:                posts_database_uri,
		PostTagsDatabaseURI:             post_tags_database_uri,
		FeedsPublicationLogsDatabaseURI: feeds_publication_logs_database_uri,
		DeliveryQueueURI:                delivery_queue_uri,
		AccountName:                     account_name,
		Mode:                            mode,
		FeedURIs:                        feed_uris,
		MaxPostsPerFeed:                 max_posts_per_feed,
		URIs:                            uris_table,
		Verbose:                         verbose,
		Templates:                       t,
	}

	return opts, nil
}
