package publish

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

var accounts_database_uri string
var activities_database_uri string
var followers_database_uri string
var deliveries_database_uri string
var posts_database_uri string
var post_tags_database_uri string
var feeds_publication_logs_database_uri string

var delivery_queue_uri string

var feed_uris multi.MultiCSVString
var max_posts_per_feed int

var account_name string

var mode string

var hostname string
var insecure bool
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("publish")

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "", "...")
	fs.StringVar(&activities_database_uri, "activities-database-uri", "", "...")
	fs.StringVar(&posts_database_uri, "posts-database-uri", "", "...")
	fs.StringVar(&post_tags_database_uri, "post-tags-database-uri", "", "...")
	fs.StringVar(&deliveries_database_uri, "deliveries-database-uri", "", "...")
	fs.StringVar(&followers_database_uri, "followers-database-uri", "", "...")

	fs.StringVar(&feeds_publication_logs_database_uri, "feeds-publication-logs-database-uri", "", "...")

	fs.Var(&feed_uris, "feed-uri", "...")

	fs.StringVar(&delivery_queue_uri, "delivery-queue-uri", "synchronous://", "...")

	fs.StringVar(&account_name, "account-name", "", "...")

	fs.StringVar(&mode, "mode", "cli", "...")

	fs.IntVar(&max_posts_per_feed, "max-posts-per-feed", 10, "...")

	fs.StringVar(&hostname, "hostname", "localhost:8080", "...")
	fs.BoolVar(&insecure, "insecure", false, "...")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	return fs
}
