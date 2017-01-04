package imdb

import (
	"github.com/Knetic/govaluate"
)

type evalExpressCache struct {
	imdb IAngelImdb
}

func (f *evalExpressCache) Set(key string, exp *govaluate.EvaluableExpression) {
	f.imdb.Set(key, exp)
}
func (f *evalExpressCache) Get(key string) (exp *govaluate.EvaluableExpression, err error) {
	obj, err := f.imdb.Get(key)
	if err != nil {
		return
	}
	return obj.(*govaluate.EvaluableExpression), nil
}

func (f *evalExpressCache) Remove(key string) error {
	return f.imdb.Remove(key)
}
func (f *evalExpressCache) SafeSet(key string, exp *govaluate.EvaluableExpression) (res *govaluate.EvaluableExpression, err error) {
	obj, err := f.imdb.SafeSet(key, exp)
	if err != nil {
		return
	}
	return obj.(*govaluate.EvaluableExpression), nil
}

func (f *evalExpressCache) ClearAll() {
	f.imdb.ClearAll()
}
