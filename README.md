# abstract-go-cache
Lib for abstract cache. No more code duplicating.

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
    aCache abstract_go_cache.AbstractCache
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