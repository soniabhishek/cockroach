package api

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"gitlab.com/playment-main/angel/app/models/integer"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

var minKeyLength = integer.Int(20)

type author struct {
	key []byte
}

func New(key string) (a author, err error) {
	if len(key) < minKeyLength.I() {
		err = errors.New("Key should be more than " + minKeyLength.String() + " characters")
		return
	}
	return author{[]byte(key)}, nil
}

func (a author) GetAPIKey(id uuid.UUID) string {
	mac := hmac.New(sha512.New512_256, a.key)

	mac.Write(id.Bytes())

	sha := encodeHexUpper(mac.Sum(nil))

	return sha
}

func (a author) Check(id uuid.UUID, key string) bool {

	bty, err := hex.DecodeString(key)
	if err != nil {
		return false
	}

	mac := hmac.New(sha512.New512_256, a.key)
	mac.Write(id.Bytes())
	return hmac.Equal(mac.Sum(nil), bty)
}
