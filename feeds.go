package feeds

import (
	"context"
	"fmt"
	"time"

	"github.com/sfomuseum/go-activitypub/id"
)

type PublicationLog struct {
	Id        int64  `json:"id"`
	AccountId int64  `json:"account_id"`
	FeedURL   string `json:"feed_url"`
	ItemGUID  string `json:"item_guid"`
	Published int64  `json:"published"`
}

func NewPublicationLog(ctx context.Context, account_id int64, feed_url string, item_guid string) (*PublicationLog, error) {

	id, err := id.NewId()

	if err != nil {
		return nil, fmt.Errorf("Failed to create new ID, %w", err)
	}

	now := time.Now()
	ts := now.Unix()

	log := &PublicationLog{
		Id:        id,
		AccountId: account_id,
		FeedURL:   feed_url,
		ItemGUID:  item_guid,
		// PostId: post_id,
		Published: ts,
	}

	return log, nil
}
