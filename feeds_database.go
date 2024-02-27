package feeds

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-roster"
)

// type GetFeedsCallbackFunc func(context.Context, *Feed) error

type FeedsDatabase interface {
	IsPublished(context.Context, int64, string, string) (bool, error)
	AddPublicationLog(context.Context, *PublicationLog) error
	Close(context.Context) error
}

var feed_database_roster roster.Roster

// FeedsDatabaseInitializationFunc is a function defined by individual feed_database package and used to create
// an instance of that feed_database
type FeedsDatabaseInitializationFunc func(ctx context.Context, uri string) (FeedsDatabase, error)

// RegisterFeedsDatabase registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `FeedsDatabase` instances by the `NewFeedsDatabase` method.
func RegisterFeedsDatabase(ctx context.Context, scheme string, init_func FeedsDatabaseInitializationFunc) error {

	err := ensureFeedsDatabaseRoster()

	if err != nil {
		return err
	}

	return feed_database_roster.Register(ctx, scheme, init_func)
}

func ensureFeedsDatabaseRoster() error {

	if feed_database_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		feed_database_roster = r
	}

	return nil
}

// NewFeedsDatabase returns a new `FeedsDatabase` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `FeedsDatabaseInitializationFunc`
// function used to instantiate the new `FeedsDatabase`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterFeedsDatabase` method.
func NewFeedsDatabase(ctx context.Context, uri string) (FeedsDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := feed_database_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(FeedsDatabaseInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func FeedsDatabaseSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureFeedsDatabaseRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range feed_database_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
