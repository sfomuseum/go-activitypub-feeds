package publish

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"math/rand"

	"github.com/aws/aws-lambda-go/lambda"
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

	post_tags_db, err := activitypub.NewPostTagsDatabase(ctx, opts.PostTagsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create instantiate post tags database, %w", err)
	}

	defer post_tags_db.Close(ctx)

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

	feeds_db, err := feeds.NewFeedsPublicationLogsDatabase(ctx, opts.FeedsPublicationLogsDatabaseURI)

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

	// START OF put me... somewhere?
	run := func(ctx context.Context) error {

		fp := gofeed.NewParser()

		for _, feed_url := range opts.FeedURIs {

			// START OF put me in a function, maybe do this concurrently?

			feed, err := fp.ParseURL(feed_url)

			if err != nil {
				return fmt.Errorf("Failed to parse URI '%s', %w", feed_url, err)
			}

			// START OF shuffle items

			items := feed.Items

			for i := range items {
				j := rand.Intn(i + 1)
				items[i], items[j] = items[j], items[i]
			}

			// END OF shuffle items

			published := 0

			for _, item := range items {

				guid := item.GUID

				is_published, err := feeds_db.IsPublished(ctx, acct.Id, feed_url, guid)

				if err != nil {
					return fmt.Errorf("Failed to determine if %s#%s has been published by %d", feed_url, guid, acct.Id)
				}

				if is_published {
					logger.Debug("Already published", "feed", feed_url, "item", guid)
					continue
				}

				// This could be made... better
				// body := item.Content
				// body := fmt.Sprintf(`<div xmlns="http://www.w3.org/1999/xhtml" style="text-align:left;">%s</div>`, item.Content)
				// body := fmt.Sprintf(`<a href="%s">%s</a><br />%s`, item.Link, item.Link, item.Title)

				body := fmt.Sprintf(`%s<br/><br /><a href="%s">%s</a>`, item.Title, item.Link, item.Link)

				post_opts := &activitypub.AddPostOptions{
					URIs:             opts.URIs,
					PostsDatabase:    posts_db,
					PostTagsDatabase: post_tags_db,
				}

				post, post_tags, err := activitypub.AddPost(ctx, post_opts, acct, body)

				if err != nil {
					return fmt.Errorf("Failed to add new post, %w", err)
				}

				deliver_opts := &activitypub.DeliverPostToFollowersOptions{
					AccountsDatabase:   accounts_db,
					FollowersDatabase:  followers_db,
					DeliveriesDatabase: deliveries_db,
					DeliveryQueue:      delivery_q,
					Post:               post,
					PostTags:           post_tags,
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

				published += 1

				if published >= opts.MaxPostsPerFeed {
					break
				}

			}

			// END OF put me in a function, maybe do this concurrently?
		}

		return nil
	}
	// END OF put me... somewhere

	switch mode {
	case "cli":

		return run(ctx)

	case "lambda":

		handler := func(ctx context.Context) error {
			return run(ctx)
		}

		lambda.Start(handler)
		return nil

	default:
		return fmt.Errorf("Invalid or unsupported mode")
	}
}
