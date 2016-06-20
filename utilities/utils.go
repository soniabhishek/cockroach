package utilities

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"strconv"
	"time"
)

const (
	Empty              = ""
	Star               = "*"
	Hyphen             = "-"
	UnderScore         = "_"
	WhiteSpace         = " "
	Colon              = ":"
	Dot                = "."
	Table_Referencer   = "." //Yeah, they both are same.
	Comma              = ","
	Spaced_Comma       = " , "
	Column_Quote       = "\""
	Place_Holder       = "%s"
	Place_Holder_Cover = "(%s)"

	Left_Parentheses  = "("
	Right_Parentheses = ")"
	Parentheses       = "()"
	Left_Braces       = "{"
	Right_Braces      = "{"
	Braces            = "{}"
	Left_Bracket      = "["
	Right_Bracket     = "]"
	Bracket           = "[]"

	StartId = iota
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
	return str == Empty
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
