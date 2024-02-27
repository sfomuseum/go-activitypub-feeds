package activitypub

import (
	"context"
	"fmt"
	"io"

	aa_docstore "github.com/aaronland/gocloud-docstore"
	gc_docstore "gocloud.dev/docstore"
)

type DocstorePostsDatabase struct {
	PostsDatabase
	collection *gc_docstore.Collection
}

func init() {

	ctx := context.Background()

	RegisterPostsDatabase(ctx, "awsdynamodb", NewDocstorePostsDatabase)

	for _, scheme := range gc_docstore.DefaultURLMux().CollectionSchemes() {
		RegisterPostsDatabase(ctx, scheme, NewDocstorePostsDatabase)
	}
}

func NewDocstorePostsDatabase(ctx context.Context, uri string) (PostsDatabase, error) {

	col, err := aa_docstore.OpenCollection(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to open collection, %w", err)
	}

	db := &DocstorePostsDatabase{
		collection: col,
	}

	return db, nil
}

func (db *DocstorePostsDatabase) AddPost(ctx context.Context, p *Post) error {

	return db.collection.Put(ctx, p)
}

func (db *DocstorePostsDatabase) GetPostWithId(ctx context.Context, id int64) (*Post, error) {

	q := db.collection.Query()
	q = q.Where("Id", "=", id)

	return db.getPost(ctx, q)
}

func (db *DocstorePostsDatabase) GetPostWithUUID(ctx context.Context, uuid string) (*Post, error) {

	q := db.collection.Query()
	q = q.Where("UUID", "=", uuid)

	return db.getPost(ctx, q)
}

func (db *DocstorePostsDatabase) getPost(ctx context.Context, q *gc_docstore.Query) (*Post, error) {

	iter := q.Get(ctx)
	defer iter.Stop()

	for {

		var p Post
		err := iter.Next(ctx, &p)

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("Failed to interate, %w", err)
		} else {
			return &p, nil
		}
	}

	return nil, ErrNotFound
}

func (db *DocstorePostsDatabase) Close(ctx context.Context) error {
	return db.collection.Close()
}
