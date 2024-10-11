package publish

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"log/slog"
	"math/rand"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mmcdole/gofeed"
	// "github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub-feeds"
	"github.com/sfomuseum/go-activitypub/database"
	"github.com/sfomuseum/go-activitypub/posts"
	"github.com/sfomuseum/go-activitypub/queue"
)

type itemTemplateVars struct {
	FeedURL string
	Item    *gofeed.Item
}

func Run(ctx context.Context) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	opts, err := OptionsFromFlagSet(ctx, fs)

	if err != nil {
		return fmt.Errorf("Failed to derive options from flagset, %w", err)
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	if opts.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	// logger := slog.Default()

	accounts_db, err := database.NewAccountsDatabase(ctx, opts.AccountsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create new database, %w", err)
	}

	defer accounts_db.Close(ctx)

	posts_db, err := database.NewPostsDatabase(ctx, opts.PostsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create instantiate posts database, %w", err)
	}

	defer posts_db.Close(ctx)

	post_tags_db, err := database.NewPostTagsDatabase(ctx, opts.PostTagsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create instantiate post tags database, %w", err)
	}

	defer post_tags_db.Close(ctx)

	followers_db, err := database.NewFollowersDatabase(ctx, opts.FollowersDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create instantiate followers database, %w", err)
	}

	defer followers_db.Close(ctx)

	deliveries_db, err := database.NewDeliveriesDatabase(ctx, opts.DeliveriesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create instantiate deliveries database, %w", err)
	}

	defer deliveries_db.Close(ctx)

	feeds_db, err := feeds.NewFeedsPublicationLogsDatabase(ctx, opts.FeedsPublicationLogsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to instantiate feeds database, %w", err)
	}

	defer feeds_db.Close(ctx)

	delivery_q, err := queue.NewDeliveryQueue(ctx, opts.DeliveryQueueURI)

	if err != nil {
		return fmt.Errorf("Failed to create new delivery queue, %w", err)
	}

	acct, err := accounts_db.GetAccountWithName(ctx, opts.AccountName)

	if err != nil {
		return fmt.Errorf("Failed to retrieve account %s, %w", opts.AccountName, err)
	}

	item_t := opts.Templates.Lookup("feed_item")

	if item_t == nil {
		return fmt.Errorf("Failed to load 'feed_item' template, missing")
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

				logger := slog.Default()
				logger = logger.With("feed", feed_url)
				logger = logger.With("item", guid)

				is_published, err := feeds_db.IsPublished(ctx, acct.Id, feed_url, guid)

				if err != nil {
					return fmt.Errorf("Failed to determine if %s#%s has been published by %d", feed_url, guid, acct.Id)
				}

				if is_published {
					logger.Debug("Item already published")
					continue
				}

				vars := itemTemplateVars{
					FeedURL: feed_url,
					Item:    item,
				}

				var buf bytes.Buffer
				wr := bufio.NewWriter(&buf)

				err = item_t.Execute(wr, vars)

				if err != nil {
					return fmt.Errorf("Failed to render template for %s#%s, %w", feed_url, guid, err)
				}

				wr.Flush()

				body := buf.String()

				logger.Debug("Post item", "body", body)

				post_opts := &posts.AddPostOptions{
					URIs:             opts.URIs,
					PostsDatabase:    posts_db,
					PostTagsDatabase: post_tags_db,
				}

				post, mentions, err := posts.AddPost(ctx, post_opts, acct, body)

				if err != nil {
					return fmt.Errorf("Failed to add post, %w", err)
				}

				logger = logger.With("post id", post.Id)

				activity, err := posts.ActivityFromPost(ctx, opts.URIs, acct, post, mentions)

				if err != nil {
					return fmt.Errorf("Failed to create new (create) activity, %w", err)
				}

				logger = logger.With("activity id", activity.Id)

				deliver_opts := &queue.DeliverActivityToFollowersOptions{
					AccountsDatabase:   accounts_db,
					FollowersDatabase:  followers_db,
					DeliveriesDatabase: deliveries_db,
					DeliveryQueue:      delivery_q,
					Activity:           activity,
					PostId:             post.Id,
					Mentions:           mentions,
					URIs:               opts.URIs,
				}

				logger.Debug("Deliver activity")

				err = queue.DeliverActivityToFollowers(ctx, deliver_opts)

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
