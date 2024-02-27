package publish

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"	
)

var accounts_database_uri string
var posts_database_uri string
var feeds_database_uri string

var delivery_queue_uri string

var feed_uris multi.MultiCSVString

var account_name string

var hostname string
var insecure bool
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("follow")

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "", "...")
	fs.StringVar(&posts_database_uri, "posts-database-uri", "", "...")
	fs.StringVar(&feeds_database_uri, "feeds-database-uri", "", "...")

	fs.StringVar(&delivery_queue_uri, "delivery-queue-uri", "synchronous://", "...")

	fs.StringVar(&account_name, "account-name", "", "...")
	
	fs.StringVar(&hostname, "hostname", "localhost:8080", "...")
	fs.BoolVar(&insecure, "insecure", false, "...")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	return fs
}
