package cache

import (
	"fmt"
	"time"
)

type ICache interface {
	Put(key string, val interface{}) error
	Get(key string) interface{}
	Delete(key string) error
}

type AbstractCache struct {
	core ICache
}

func NewAbstractCache(core ICache) *AbstractCache {
	return AbstractCache{core: core}
}

func (a *AbstractCache) ()  {
	
}


func (c *CacheManager) getFromCache(
    ctx context.Context,
 keyPrefix string,
 getter func(ctx context.Context, params ...interface{}) (interface{}, error),
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

 data, err = c.cache.GetData(ctx, key)
 if err != nil {
  logger.WarnKV(ctx, "error when getting from cache", "error", err)
  if len(params) > 0 {
   got, err = getter(ctx, params...)
  } else {
   got, err = getter(ctx)
  }
  if err != nil {
   return nil, err
  }
  data, _ = json.Marshal(got)
  key := key
  got := got
  go func() {
   traceCtx := ctxutil.WithValues(context.Background(), ctx)
   setCtx, setCnl := context.WithTimeout(traceCtx, 20*time.Millisecond)
   defer setCnl()
   if putErr := c.cache.Put(setCtx, key, got); putErr != nil {
    logger.WarnKV(traceCtx, "error when putting in cache", "error", putErr)
   }
  }()
 }

 return data, nil
}

