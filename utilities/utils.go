package utilities

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/asaskevich/govalidator"
	"gitlab.com/playment-main/angel/utilities/constants"
)

func TimeInMillis() int64 {
	now := time.Now()
	unixNano := now.UnixNano()
	umillisec := unixNano / 1000000
	return umillisec
}

func TimeDiff(absolute bool, times ...int64) int64 {

	var newTime, oldTime int64

	if len(times) > 1 {
		oldTime = times[0]
		newTime = times[1]
	} else {
		oldTime = times[0]
		newTime = TimeInMillis()
	}

	if absolute {
		return Abs(newTime - oldTime)
	} else {
		return newTime - oldTime
	}
}

func Abs(n int64) int64 {
	if n > 0 {
		return n
	} else {
		return 0 - n
	}
}

func GetRandomID() int64 {
	return time.Now().UnixNano()
}

func IsEmptyOrNil(str string) bool {
	return str == constants.Empty
}

func IsValidError(err error) bool {
	return err != nil
}

func IsInt(val string) bool {
	_, err := strconv.Atoi(val)
	return (err == nil)
}

func GetInt(val string) int {
	ret, _ := strconv.Atoi(val)
	return ret
}

func GetHMAC(val string, hmacKey string) string {
	key := []byte(hmacKey)
	sig := hmac.New(sha256.New, key)
	sig.Write([]byte(val))
	return hex.EncodeToString(sig.Sum(nil))
}

//Some long long key
var encryptionKey string = "29ab862a724bf2c310ba8888ed15965d"

func Encrypt(plainText string) ([]byte, error) {
	text := []byte(plainText)
	block, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

func Decrypt(encryptedText string) ([]byte, error) {
	text := []byte(encryptedText)
	block, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func ValidateUrl(urlStr string) bool {
	var validURL bool
	validURL = govalidator.IsURL(urlStr)
	return validURL
}

func ReplaceEscapeCharacters(bty []byte) []byte {
	bty = bytes.Replace(bty, []byte(`\u003c`), []byte("<"), -1)
	bty = bytes.Replace(bty, []byte(`\u003e`), []byte(">"), -1)
	bty = bytes.Replace(bty, []byte(`\u0026`), []byte("&"), -1)
	return bty
}

//In the upcoming Go1.7 release, you can turn off the escaping in a json.Encoder with SetEscapeHTML(false).
/*
func ReplaceEscapeCharsInGOLAN1_7()  {

	mp := map[string]string{"key": "The Knight & Day"}
	var buf bytes.Buffer

	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.Encode(mp)

	fmt.Println(&buf)
}*/

func IsValidUTF8(rows []string) (int, error) {
	for i, row := range rows {
		row = strings.TrimRight(row, "\n")
		fmt.Println(row, []byte(row))
		if !utf8.ValidString(row) {
			return i, errors.New("!utf8.ValidString")
		}
	}
	return -1, nil
}
