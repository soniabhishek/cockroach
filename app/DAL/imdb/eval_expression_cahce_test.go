package imdb

import (
	"testing"
	"github.com/Knetic/govaluate"
	"github.com/stretchr/testify/assert"
)

func TestEvalExpressCache_Set(t *testing.T) {
	idb := newEvalExpressionCache()
	exp, err := govaluate.NewEvaluableExpression("1>3")
	assert.NoError(t, err)
	idb.Set("hello", exp)
	val, err := idb.Get("hello")
	assert.NoError(t, err)
	assert.Equal(t, val, exp)
}
