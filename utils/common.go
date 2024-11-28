package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"math/rand"

	"github.com/google/uuid"
)

var (
	DefaultTimeFormat = fmt.Sprintf("%s %s", time.DateOnly, time.TimeOnly)
	alphanumeric      = []byte("0123456789abcdefghijklmnopqrstuvwxyz")
	JwtSecret         = []byte("10ad0120-c1d9-4015-8a34-d8d7ba9e4c48")
	UserPasswdSalt    = "03eb7df1-137b-44f3-a946-6413bd2808a6"
)

func GenUUID() string {
	return uuid.NewString()
}

// a random string of arbitrary length
func GenRandStr(n int) string {
	b := make([]byte, n)
	end := len(alphanumeric)
	for i := 0; i < n; i++ {
		b[i] = alphanumeric[rand.Intn(end)]
	}
	return string(b)
}

// subitem is in slice
func Contains[T comparable](as []T, sub T) bool {
	for _, v := range as {
		if v == sub {
			return true
		}
	}
	return false
}

// hash string
func Sha256Hash(origin string) (hash string) {
	h := sha256.New()
	h.Write([]byte(origin))
	return hex.EncodeToString(h.Sum(nil))
}

// simple check whether the string is in sha256 format
func IsSha256(hash string) bool {
	re := regexp.MustCompile(`^[A-Fa-f0-9]{64}$`)
	return re.MatchString(hash)
}

// classic sha256 encode: api auth
// check if appid in db, get appkey
// check if the timestamp is within three minutes of the current time
// check if sha256(appkey+random+timestamp) is equal
func GenSha256AuthHeaders(appid, appkey string) (headers map[string]string) {
	unixNano := time.Now().UnixNano()
	currentTime := strconv.FormatInt(unixNano/int64(time.Second), 10)
	random := strconv.Itoa(rand.New(rand.NewSource(unixNano)).Int())
	h := sha256.New()
	h.Write([]byte(appkey + random + currentTime))
	signature := hex.EncodeToString(h.Sum(nil))
	headers = map[string]string{
		"AppId":     appid,
		"Random":    random,
		"TimeStamp": currentTime,
		"Signature": signature,
	}
	return
}
