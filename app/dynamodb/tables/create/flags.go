package create

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var refresh bool

var dynamodb_client_uri string
var table_prefix string
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("follow")

	fs.BoolVar(&refresh, "refresh", false, "...")
	fs.StringVar(&dynamodb_client_uri, "dynamodb-client-uri", "", "...")
	fs.StringVar(&table_prefix, "table-prefix", "", "...")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")
	return fs
}
