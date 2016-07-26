package auther

import (
	"crypto/hmac"
	"errors"
	"strconv"

	"crypto/sha256"
	"encoding/base64"

	"github.com/crowdflux/angel/app/config"
	"github.com/crowdflux/angel/app/models/uuid"
)

var minKeyLength = 20

type auther struct {
	key []byte
}

func New(key string) (a auther, err error) {
	if len(key) < minKeyLength {
		err = errors.New("Key should be more than " + strconv.Itoa(minKeyLength) + " characters")
		return
	}
	return auther{[]byte(key)}, nil
}

func (a auther) GetAPIKey(id uuid.UUID) string {
	mac := hmac.New(sha256.New, a.key)

	mac.Write(id.Bytes())

	sha := base64.RawStdEncoding.EncodeToString(mac.Sum(nil))

	return sha
}

func (a auther) Check(id uuid.UUID, key string) bool {

	bty, err := base64.RawStdEncoding.DecodeString(key)
	if err != nil {
		return false
	}

	mac := hmac.New(sha256.New, a.key)
	mac.Write(id.Bytes())
	return hmac.Equal(mac.Sum(nil), bty)
}

//--------------------------------------------------------------------------------//

var StdProdAuther = auther{[]byte(config.AUTHER_PLAYMENT_SECRET.Get())}
