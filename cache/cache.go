package cache

import (
	"fmt"
	"time"
    "context"
    "encoding/json"
)

type ICache interface {
	Put(ctx context.Context, key string, val interface{}) error
	Get(ctx context.Context, key string) interface{}
	Delete(ctx context.Context, key string) error
}

type AbstractCache struct {
	core ICache
}

func NewAbstractCache(core ICache) *AbstractCache {
	return AbstractCache{core: core}
}

type AbstractProc func(ctx context.Context, params ...interface{}) (interface{}, error)

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

        _ = c.core.Put(setCtx, key, got)
    }

    return data, nil
}

