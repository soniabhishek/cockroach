package uuid

import (
	"encoding/base32"
	"errors"
	"strings"
)

// Do not change this string
// All api keys will become invalid
// TODO think about how to secure this (most probably config file)
const cEncodeStr = "jkl234GEFVmqrsTCWXdefgYUDnABZabc"

const suffix = "======"

// Our custom base32 encoder
var base32EncoderCustom = base32.NewEncoding(cEncodeStr)

var ErrCEncNotValid = errors.New("Encoded string is not a valid cEncoded string")

// Parses the input cEncoded string to UUID
func FromCEnc(enc string) (UUID, error) {
	enc += suffix
	bty, err := base32EncoderCustom.DecodeString(enc)
	if err != nil {
		return Nil, ErrCEncNotValid
	}

	uuid, err := FromBytes(bty)
	if err != nil {
		return Nil, err
	}
	return uuid, nil
}

// Returns cEncoded string (our custom encoder)
func toCEnc(uuid UUID) string {
	s := base32EncoderCustom.EncodeToString(uuid.Bytes())
	return strings.TrimSuffix(s, suffix)
}
