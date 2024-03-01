package feeds

import (
	"context"
)

type NullFeedsPublicationLogsDatabase struct {
	FeedsPublicationLogsDatabase
}

func init() {
	ctx := context.Background()
	RegisterFeedsPublicationLogsDatabase(ctx, "null", NewNullFeedsPublicationLogsDatabase)
}

func NewNullFeedsPublicationLogsDatabase(ctx context.Context, uri string) (FeedsPublicationLogsDatabase, error) {
	db := &NullFeedsPublicationLogsDatabase{}
	return db, nil
}

func (db *NullFeedsPublicationLogsDatabase) IsPublished(ctx context.Context, account_id int64, feed_url string, item_guid string) (bool, error) {
	return false, nil
}

func (db *NullFeedsPublicationLogsDatabase) AddPublicationLog(ctx context.Context, log *PublicationLog) error {
	return nil
}

func (db *NullFeedsPublicationLogsDatabase) Close(ctx context.Context) error {
	return nil
}
