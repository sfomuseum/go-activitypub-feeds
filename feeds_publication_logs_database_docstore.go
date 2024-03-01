package feeds

import (
	"context"
	"fmt"
	"io"

	aa_docstore "github.com/aaronland/gocloud-docstore"
	gc_docstore "gocloud.dev/docstore"
)

type DocstoreFeedsPublicationLogsDatabase struct {
	FeedsPublicationLogsDatabase
	collection *gc_docstore.Collection
}

func init() {

	ctx := context.Background()

	RegisterFeedsPublicationLogsDatabase(ctx, "awsdynamodb", NewDocstoreFeedsPublicationLogsDatabase)

	for _, scheme := range gc_docstore.DefaultURLMux().CollectionSchemes() {
		RegisterFeedsPublicationLogsDatabase(ctx, scheme, NewDocstoreFeedsPublicationLogsDatabase)
	}
}

func NewDocstoreFeedsPublicationLogsDatabase(ctx context.Context, uri string) (FeedsPublicationLogsDatabase, error) {

	col, err := aa_docstore.OpenCollection(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to open collection, %w", err)
	}

	db := &DocstoreFeedsPublicationLogsDatabase{
		collection: col,
	}

	return db, nil
}

func (db *DocstoreFeedsPublicationLogsDatabase) IsPublished(ctx context.Context, account_id int64, feed_url string, item_guid string) (bool, error) {

	q := db.collection.Query()
	q = q.Where("AccountId", "=", account_id)
	q = q.Where("FeedURL", "=", feed_url)
	q = q.Where("ItemGUID", "=", item_guid)

	iter := q.Get(ctx)
	defer iter.Stop()

	var l PublicationLog
	err := iter.Next(ctx, &l)

	if err == io.EOF {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("Failed to interate, %w", err)
	} else {
		return true, nil
	}
}

func (db *DocstoreFeedsPublicationLogsDatabase) AddPublicationLog(ctx context.Context, log *PublicationLog) error {
	return db.collection.Put(ctx, log)
}

func (db *DocstoreFeedsPublicationLogsDatabase) Close(ctx context.Context) error {
	return db.collection.Close()
}
