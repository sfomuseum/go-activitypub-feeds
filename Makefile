GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

TABLE_PREFIX=

ACCOUNTS_DB_URI=awsdynamodb://$(TABLE_PREFIX)accounts?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:
ACTIVITIES_DB_URI=awsdynamodb://$(TABLE_PREFIX)activities?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:
FOLLOWING_DB_URI=awsdynamodb://$(TABLE_PREFIX)following?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:
FOLLOWERS_DB_URI=awsdynamodb://$(TABLE_PREFIX)followers?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:
BLOCKS_DB_URI=awsdynamodb://$(TABLE_PREFIX)blocks?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:
NOTES_DB_URI=awsdynamodb://$(TABLE_PREFIX)notes?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:
POSTS_DB_URI=awsdynamodb://$(TABLE_PREFIX)posts?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:
MESSAGES_DB_URI=awsdynamodb://$(TABLE_PREFIX)messages?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:
DELIVERIES_DB_URI=awsdynamodb://$(TABLE_PREFIX)deliveries?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:
FEEDS_DB_URI=awsdynamodb://$(TABLE_PREFIX)feeds_publication_logs?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:

cli:
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/publish-feeds cmd/publish-feeds/main.go

publish:
	go run cmd/publish-feeds/main.go \
		-accounts-database-uri '$(ACCOUNTS_DB_URI)' \
		-activities-database-uri '$(ACTIVITIES_DB_URI)' \
		-followers-database-uri '$(FOLLOWERS_DB_URI)' \
		-posts-database-uri '$(POSTS_DB_URI)' \
		-deliveries-database-uri '$(DELIVERIES_DB_URI)' \
		-feeds-database-uri '$(FEEDS_DB_URI)' \
		-account-name alice \
		-feed-uri $(FEED) \
		-hostname localhost:8080 \
		-insecure \
		-verbose

dynamo-tables-local:
	go run -mod vendor cmd/create-dynamodb-tables/main.go \
		-refresh \
		-table-prefix '$(TABLE_PREFIX)' \
		-dynamodb-client-uri 'awsdynamodb://?local=true&region=localhost&credentials=anon:'

lambda:
	if test -f bootstrap; then rm -f bootstrap; fi
	if test -f publish.zip; then rm -f publish.zip; fi
	GOARCH=arm64 GOOS=linux go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -tags lambda.norpc -o bootstrap cmd/publish-feeds/main.go
	zip publish.zip bootstrap
	rm -f bootstrap
