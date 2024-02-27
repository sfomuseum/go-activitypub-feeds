package feeds

import (
	"context"
)

type NullFeedsDatabase struct {
	FeedsDatabase
}

func init() {
	ctx := context.Background()
	RegisterFeedsDatabase(ctx, "null", NewNullFeedsDatabase)
}

func NewNullFeedsDatabase(ctx context.Context, uri string) (FeedsDatabase, error) {
	db := &NullFeedsDatabase{}
	return db, nil
}

func (db *NullFeedsDatabase) IsPublished(ctx context.Context, account_id int64, feed_url string, item_guid string) (bool, error) {
	return false, nil
}

func (db *NullFeedsDatabase) AddPublicationLog(ctx context.Context, log *PublicationLog) error {
	return nil
}

func (db *NullFeedsDatabase) Close(ctx context.Context) error {
	return nil
}
