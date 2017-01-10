package imdb

type IAngelImdb interface {
	Set(key string, value interface{})
	Get(key string) (value interface{}, err error)
	Remove(key string) error
	SafeSet(key string, value interface{}) (existingValue interface{}, err error)
	ClearAll()
}

type ICmap interface {
	get(key interface{}) (val interface{},ok bool)
	set(key, value interface{})
	delete(key interface{})
	reset()
}

func NewCmap() ICmap {
	return &cmap{
		cmap: make(map[interface{}]interface{}),
	}
}
func new() IAngelImdb {
	return &imdb{
		db: cmap{
			cmap: make(map[interface{}]interface{}),
		},
	}
}

func newFluUploadCache() *fluUploadCache {
	return &fluUploadCache{
		imdb: new(),
	}
}

func newEvalExpressionCache() *evalExpressCache {
	return &evalExpressCache{
		imdb: new(),
	}
}

func NewFluValidateCache() *FluValidateCache {
	return &FluValidateCache{
		imdb: new(),
	}
}

var EvalExpressionCache = newEvalExpressionCache()
var FluUploadCache = newFluUploadCache()
