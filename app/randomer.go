package app

import (
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var letters string
var url_len int

const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

var src = rand.NewSource(time.Now().UnixNano())

func GetRandomString() string {
	sb := strings.Builder{}
	sb.Grow(url_len)
	for i, cache, remain := url_len-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letters) {
			sb.WriteByte(letters[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return sb.String()
}

func init() {
	rand.Seed(time.Now().UnixNano())
	godotenv.Load()
	letters = os.Getenv("LETTERS")
	if letters == "" {
		letters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"
	}
	url_len, _ = strconv.Atoi(os.Getenv("URLLEN"))
	if url_len < 1 {
		url_len = 10
	}
}

// faster than simple random:
// 1.0 version b := make([]byte, n)
// b[i] = letters[rand.Int63()%int64(len(letters))]
