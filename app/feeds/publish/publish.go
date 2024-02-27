package publish

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/mmcdole/gofeed"
	"github.com/sfomuseum/go-activitypub"
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

	feeds_db, err := activitypub.NewFeedsDatabase(ctx, opts.FeedsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to instantiate feeds database, %w", err)
	}

	defer followers_db.Close(ctx)

	posts_db, err := activitypub.NewPostsDatabase(ctx, opts.PostsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create instantiate posts database, %w", err)
	}

	defer posts_db.Close(ctx)


	delivery_q, err := activitypub.NewDeliveryQueue(ctx, opts.DeliveryQueueURI)

	if err != nil {
		return fmt.Errorf("Failed to create new delivery queue, %w", err)
	}

	acct, err := accounts_db.GetAccountWithName(ctx, opts.AccountName)

	if err != nil {
		return fmt.Errorf("Failed to retrieve account %s, %w", opts.AccountName, err)
	}

	//
	
	fp := gofeed.NewParser()
	
	for _, uri := range opts.FeedURIs {

		feed, err := fp.ParseURL(uri)
	
		if err != nil {
			return fmt.Errorf("Failed to parse URI '%s', %w", uri, err)
		}

		for _, item := range feed.Items {
		
			guid := item.GUID

			is_published, err := feeds_db.IsPublished(ctx, acct.Id, uri, guid)

			if err != nil {
				return fmt.Errorf("Failed to determine if %s#%s has been published by %d", uri, guid, acct.Id)
			}

			if is_published {
				continue
			}

			slog.Info(item.Content)
		}
	}

	//
	
	/*
	p, err := activitypub.NewPost(ctx, acct, opts.Message)

	if err != nil {
		return fmt.Errorf("Failed to create new post, %w", err)
	}

	err = posts_db.AddPost(ctx, p)

	if err != nil {
		return fmt.Errorf("Failed to add post, %w", err)
	}

	deliver_opts := &activitypub.DeliverPostToFollowersOptions{
		AccountsDatabase:   accounts_db,
		FollowersDatabase:  followers_db,
		DeliveriesDatabase: deliveries_db,
		DeliveryQueue:      delivery_q,
		Post:               p,
		URIs:               opts.URIs,
	}

	err = activitypub.DeliverPostToFollowers(ctx, deliver_opts)

	if err != nil {
		return fmt.Errorf("Failed to deliver post, %w", err)
	}

	*/
	
	return nil
}
