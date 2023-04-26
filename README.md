## abstract-go-cache
# Lib for abstract cache. No more code duplicating.

Usage example:

```go
type YourDataStruct struct {
    // some fields...
}

const (
    cachePrefixSomeUsefulData = "some_useful_data"
)

type SomeInterface interface {
    UsefulProc(ctx context.Context, parameter1 int64, parameter2 string) (*YourDataStruct, error)
}

type ICacheCore interface {
	Put(ctx context.Context, key string, val interface{}) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
}

type YourService struct {
    aCache        *abstract_go_cache.AbstractCache
    someInterface SomeInterface
    //...
}

func NewYourService(someInterface SomeInterface, cacheCore ICacheCore) *YourService {
    return &YourService{
        aCache: abstract_go_cache.NewAbstractCache(cacheCore),
        someInterface: someInterface,
        //...
    }
}

func (s *YourService) GetSomeUsefulData(ctx context.Context, parameter1 int64, parameter2 string) (*YourDataStruct, error) {
    
    // some code here ...
    
    data, err := s.aCache.GetFromCache(ctx, cachePrefixSomeUsefulData,
        func(ctx context.Context, params ...interface{}) (interface{}, error) {
            return s.someInterface.UsefulProc(ctx, (params[0]).(int64), (params[1]).(string))
        },
        parameter1, parameter2)
    if err != nil {
        return nil, fmt.Errorf("get from cache: %w", err)
    }

    var ds YourDataStruct
    if err = json.Unmarshal(data, &ds); err != nil {
        return nil, fmt.Errorf("failed to unmarshal: %w", err)
    }
    return &ds, nil
}
```

How to install:

type `go get -v github.com/rusfort/abstract-go-cache` in your terminal

# Full main.go file - try this out!

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"

	abstract_go_cache "github.com/rusfort/abstract-go-cache/cache"
)

const (
	cachePrefixSomeUsefulData = "some_useful_data"
)

// some struct with useful info

type YourDataStruct struct {
	Value1 int64
	Value2 string
}

// let's try to call the abstract cache

func main() {
	println("hello")
	s := SomeStruct{}
	c := NewCache()
	service := NewYourService(s, c)
	println("1 attempt: empty cache")
	data, err := service.GetSomeUsefulData(context.Background(), 12, "abc")
	if err != nil {
		println(err.Error())
		return
	}
	println(data.Value1, data.Value2)
	println("2 attempt: cache hit")
	data, err = service.GetSomeUsefulData(context.Background(), 12, "abc")
	if err != nil {
		println(err.Error())
		return
	}
	println(data.Value1, data.Value2)
}

// your procedures or rpc

type SomeInterface interface {
	UsefulProc(ctx context.Context, parameter1 int64, parameter2 string) (*YourDataStruct, error)
}

type SomeStruct struct{}

func (SomeStruct) UsefulProc(ctx context.Context, parameter1 int64, parameter2 string) (*YourDataStruct, error) {
	println("proc called")

	// maybe smth good happends here ...

	return &YourDataStruct{
		Value1: parameter1,
		Value2: parameter2,
	}, nil
}

// any cache core implementation

type ICacheCore interface {
	Put(ctx context.Context, key string, val interface{}) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
}

type Cache struct {
	Data map[string]*YourDataStruct
}

func NewCache() *Cache {
	return &Cache{
		Data: make(map[string]*YourDataStruct),
	}
}

func (c *Cache) Put(ctx context.Context, key string, val interface{}) error {
	println("cache put called")
	c.Data[key] = val.(*YourDataStruct)
	return nil
}

func (c *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	println("cache get called")
	data, ok := c.Data[key]
	if !ok {
		println("cache miss")
		return nil, fmt.Errorf("no data")
	}
	res, err := json.Marshal(*data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal")
	}
	println("cache hit")
	return res, nil
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	delete(c.Data, key)
	return nil
}

// your local service

type YourService struct {
	aCache        *abstract_go_cache.AbstractCache
	someInterface SomeInterface
	//...
}

func NewYourService(someInterface SomeInterface, cacheCore ICacheCore) *YourService {
	return &YourService{
		aCache:        abstract_go_cache.NewAbstractCache(cacheCore),
		someInterface: someInterface,
		//...
	}
}

func (s *YourService) GetSomeUsefulData(ctx context.Context, parameter1 int64, parameter2 string) (*YourDataStruct, error) {

	// some code here ...

    // just call GetFromCache method and get the data!
	data, err := s.aCache.GetFromCache(ctx, cachePrefixSomeUsefulData,
		func(ctx context.Context, params ...interface{}) (interface{}, error) {
			return s.someInterface.UsefulProc(ctx, (params[0]).(int64), (params[1]).(string))
		},
		parameter1, parameter2)
	if err != nil {
		return nil, fmt.Errorf("get from cache: %w", err)
	}

    // unmarshal data to your datastruct
	var ds YourDataStruct
	if err = json.Unmarshal(data, &ds); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}
	return &ds, nil
}

```
