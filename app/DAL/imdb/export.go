package imdb

type IAngelImdb interface {
	Set(key string, value interface{})
	Get(key string) (value interface{}, err error)
	Remove(key string) error
	SafeSet(key string, value interface{}) (existingValue interface{}, err error)
	ClearAll()
}

func new() IAngelImdb {
	return &imdb{
		db: make(map[string]interface{}),
	}
}

func newFluUploadCache() *fluUploadCache {
	return &fluUploadCache{
		imdb: new(),
	}
}

var FluUploadCache = newFluUploadCache()
