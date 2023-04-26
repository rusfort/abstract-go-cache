package cache

import (
	"context"
	"encoding/json"
	"fmt"
)

// ICache implements any of your local cache or Memcache.
type ICache interface {
	Put(ctx context.Context, key string, val interface{}) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
}

// AbstractCache is an abstraction for your cache
type AbstractCache struct {
	core ICache
}

// NewAbstractCache constructs new AbstractCache based on your cache core
func NewAbstractCache(core ICache) *AbstractCache {
	return &AbstractCache{core: core}
}

// AbstractProc is any of your procedures or rpc calls wrapped in an abstract function
type AbstractProc func(ctx context.Context, params ...interface{}) (interface{}, error)

// GetFromCache gets data from your cache if it exist (and didn't expire) and calls your proc if not (and then saves the result in your cache)
func (c *AbstractCache) GetFromCache(
	ctx context.Context,
	keyPrefix string,
	proc AbstractProc,
	params ...interface{},
) ([]byte, error) {
	var (
		err  error
		got  interface{}
		data []byte
	)

	key := keyPrefix
	for _, p := range params {
		key += fmt.Sprintf("-%v", p)
	}

	data, err = c.core.Get(ctx, key)
	if err != nil {
		got, err = proc(ctx, params...)
		if err != nil {
			return nil, fmt.Errorf("error on proc: %w", err)
		}

		data, err = json.Marshal(got)
		if err != nil {
			return nil, fmt.Errorf("error on marshal: %w", err)
		}

		_ = c.core.Put(ctx, key, got)
	}

	return data, nil
}
