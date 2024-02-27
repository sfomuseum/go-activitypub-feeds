package publish

import (
	"context"
	"flag"
	"fmt"
	"log/slog"

	"github.com/mmcdole/gofeed"
	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub-feeds"
	ap_slog "github.com/sfomuseum/go-activitypub/slog"
)

func Run(ctx context.Context, logger *slog.Logger) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs, logger)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet, logger *slog.Logger) error {

	opts, err := OptionsFromFlagSet(ctx, fs)

	if err != nil {
		return fmt.Errorf("Failed to derive options from flagset, %w", err)
	}

	return RunWithOptions(ctx, opts, logger)
}

func RunWithOptions(ctx context.Context, opts *RunOptions, logger *slog.Logger) error {

	ap_slog.ConfigureLogger(logger, opts.Verbose)

	accounts_db, err := activitypub.NewAccountsDatabase(ctx, opts.AccountsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create new database, %w", err)
	}

	defer accounts_db.Close(ctx)

	posts_db, err := activitypub.NewPostsDatabase(ctx, opts.PostsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create instantiate posts database, %w", err)
	}

	defer posts_db.Close(ctx)

	followers_db, err := activitypub.NewFollowersDatabase(ctx, opts.FollowersDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create instantiate followers database, %w", err)
	}

	defer followers_db.Close(ctx)

	deliveries_db, err := activitypub.NewDeliveriesDatabase(ctx, opts.DeliveriesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create instantiate deliveries database, %w", err)
	}

	defer deliveries_db.Close(ctx)

	feeds_db, err := feeds.NewFeedsDatabase(ctx, opts.FeedsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to instantiate feeds database, %w", err)
	}

	defer feeds_db.Close(ctx)

	delivery_q, err := activitypub.NewDeliveryQueue(ctx, opts.DeliveryQueueURI)

	if err != nil {
		return fmt.Errorf("Failed to create new delivery queue, %w", err)
	}

	acct, err := accounts_db.GetAccountWithName(ctx, opts.AccountName)

	if err != nil {
		return fmt.Errorf("Failed to retrieve account %s, %w", opts.AccountName, err)
	}

	fp := gofeed.NewParser()

	for _, feed_url := range opts.FeedURIs {

		// START OF put me in a function, maybe do this concurrently?
		
		feed, err := fp.ParseURL(feed_url)

		if err != nil {
			return fmt.Errorf("Failed to parse URI '%s', %w", feed_url, err)
		}

		for _, item := range feed.Items {

			guid := item.GUID

			is_published, err := feeds_db.IsPublished(ctx, acct.Id, feed_url, guid)

			if err != nil {
				return fmt.Errorf("Failed to determine if %s#%s has been published by %d", feed_url, guid, acct.Id)
			}

			if is_published {
				logger.Debug("Already published", "feed", feed_url, "item", guid)
				continue
			}

			post, err := activitypub.NewPost(ctx, acct, item.Content)

			if err != nil {
				return fmt.Errorf("Failed to create new post, %w", err)
			}

			err = posts_db.AddPost(ctx, post)

			if err != nil {
				return fmt.Errorf("Failed to add post, %w", err)
			}

			deliver_opts := &activitypub.DeliverPostToFollowersOptions{
				AccountsDatabase:   accounts_db,
				FollowersDatabase:  followers_db,
				DeliveriesDatabase: deliveries_db,
				DeliveryQueue:      delivery_q,
				Post:               post,
				URIs:               opts.URIs,
			}

			err = activitypub.DeliverPostToFollowers(ctx, deliver_opts)

			if err != nil {
				return fmt.Errorf("Failed to deliver post, %w", err)
			}

			log, err := feeds.NewPublicationLog(ctx, acct.Id, feed_url, guid)

			if err != nil {
				return fmt.Errorf("Failed to create new publication log, %w", err)
			}

			err = feeds_db.AddPublicationLog(ctx, log)

			if err != nil {
				return fmt.Errorf("Failed to add publication log, %w", err)
			}

			logger.Info("Published feed item", "account", acct.Id, "feed", feed_url, "item", guid, "post", post.Id, "log", log.Id)
			break
		}

		// END OF put me in a function, maybe do this concurrently?		
	}

	return nil
}
