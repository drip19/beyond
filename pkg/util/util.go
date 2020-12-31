package util

import (
	"encoding/hex"
	"encoding/json"
	"math/rand"
	"net/url"
	"os"
	"time"
)

//CreateDirIfNotExist create given folder
func CreateDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0700)
		return err
	}
	return nil
}

func ToIndentString(v interface{}) string {
	b, err := json.MarshalIndent(&v, "", "\t")
	if err != nil {
		return ""
	}
	return string(b)
}

func ToString(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

// Bytes fills the given byte slice with random bytes.
func Bytes(data []byte) error {
	_, err := rand.Read(data)
	return err
}

func RandomHexString(length int) string {
	if length == 0 {
		return ""
	}
	b := make([]byte, length)
	_ = Bytes(b)
	s := hex.EncodeToString(b)
	return s
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}

func Scheme(endpoint string) (string, string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", "", err
	}
	return u.Scheme, u.Host, nil
}
