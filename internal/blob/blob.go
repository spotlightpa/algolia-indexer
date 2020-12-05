package blob

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"path"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/memblob"
	_ "gocloud.dev/blob/s3blob"
	"gocloud.dev/gcerrors"
)

var ErrNotFound = errors.New("not found")

type Store struct {
	b    *blob.Bucket
	path string
	l    Logger
}

// New opens a Bucket at the provided URL
func New(ctx context.Context, urlstr string, l Logger) (*Store, error) {
	u, err := url.Parse(urlstr)
	if err != nil {
		return nil, err
	}

	store := Store{
		path: u.Path,
		l:    l,
	}
	if store.b, err = blob.OpenBucket(ctx, urlstr); err != nil {
		return nil, err
	}

	return &store, nil
}

func (store *Store) printf(format string, v ...interface{}) {
	if store.l != nil {
		store.l.Printf(format, v...)
	}
}

func (store *Store) name(key string) string {
	return path.Join(store.path, key+".json")
}

// Get calls GET in Redis and converts values from JSON bytes
func (store *Store) Get(ctx context.Context, key string, getv interface{}) (err error) {
	key = store.name(key)
	store.printf("get %q", key)
	getb, err := store.b.ReadAll(ctx, key)
	if err != nil {
		if gcerrors.Code(err) == gcerrors.NotFound {
			return ErrNotFound
		}
		return err
	}
	return json.Unmarshal(getb, getv)
}

// Set converts values to JSON bytes and calls SET in Redis
func (store *Store) Set(ctx context.Context, key string, setv interface{}) (err error) {
	key = store.name(key)
	store.printf("set %q", key)

	setb, err := json.Marshal(setv)
	if err != nil {
		return err
	}
	opts := blob.WriterOptions{
		ContentType: "application/json",
	}
	return store.b.WriteAll(ctx, key, setb, &opts)
}
