package main

import (
	"encoding/base32"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

// Benchmark functions

const n = 16

func BenchmarkRunes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringRunes(n)
	}
}

func BenchmarkBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringBytes(n)
	}
}

func BenchmarkBytesRmndr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringBytesRmndr(n)
	}
}

func BenchmarkBytesMask(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringBytesMask(n)
	}
}

func BenchmarkBytesMaskImpr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringBytesMaskImpr(n)
	}
}

func BenchmarkBytesMaskImprSrc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringBytesMaskImprSrc(n)
	}
}

func TestRandomFromString(t *testing.T) {
	id := uuid.NewV4()
	fmt.Println(id)

	a := base32.StdEncoding.EncodeToString(id.Bytes())
	fmt.Println(a)

	b, err := base32.StdEncoding.DecodeString(a)
	assert.NoError(t, err)

	c, err := uuid.FromBytes(b)
	assert.NoError(t, err)

	fmt.Println(c.String())
}
