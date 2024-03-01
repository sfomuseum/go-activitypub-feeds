package feeds

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-roster"
)

type FeedsPublicationLogsDatabase interface {
	IsPublished(context.Context, int64, string, string) (bool, error)
	AddPublicationLog(context.Context, *PublicationLog) error
	Close(context.Context) error
}

var feeds_publication_logs_database_roster roster.Roster

// FeedsPublicationLogsDatabaseInitializationFunc is a function defined by individual feeds_publication_logs_database package and used to create
// an instance of that feeds_publication_logs_database
type FeedsPublicationLogsDatabaseInitializationFunc func(ctx context.Context, uri string) (FeedsPublicationLogsDatabase, error)

// RegisterFeedsPublicationLogsDatabase registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `FeedsPublicationLogsDatabase` instances by the `NewFeedsPublicationLogsDatabase` method.
func RegisterFeedsPublicationLogsDatabase(ctx context.Context, scheme string, init_func FeedsPublicationLogsDatabaseInitializationFunc) error {

	err := ensureFeedsPublicationLogsDatabaseRoster()

	if err != nil {
		return err
	}

	return feeds_publication_logs_database_roster.Register(ctx, scheme, init_func)
}

func ensureFeedsPublicationLogsDatabaseRoster() error {

	if feeds_publication_logs_database_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		feeds_publication_logs_database_roster = r
	}

	return nil
}

// NewFeedsPublicationLogsDatabase returns a new `FeedsPublicationLogsDatabase` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `FeedsPublicationLogsDatabaseInitializationFunc`
// function used to instantiate the new `FeedsPublicationLogsDatabase`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterFeedsPublicationLogsDatabase` method.
func NewFeedsPublicationLogsDatabase(ctx context.Context, uri string) (FeedsPublicationLogsDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := feeds_publication_logs_database_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(FeedsPublicationLogsDatabaseInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func FeedsPublicationLogsDatabaseSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureFeedsPublicationLogsDatabaseRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range feeds_publication_logs_database_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
