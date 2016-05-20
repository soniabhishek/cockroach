package api

import (
	"encoding/hex"
	"strings"
)

/**
Encodes the hex into upper case string
*/
func encodeHexUpper(bty []byte) string {
	s := hex.EncodeToString(bty)
	return strings.ToUpper(s)
}
