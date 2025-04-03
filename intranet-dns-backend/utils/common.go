package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"math/rand"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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

func RemoveRepeatedElement[T comparable](old []T) (new []T) {
	new = make([]T, 0)
	for i := 0; i < len(old); i++ {
		repeat := false
		for j := i + 1; j < len(old); j++ {
			if old[i] == old[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			new = append(new, old[i])
		}
	}
	return
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

func marshal(src interface{}) string {
	if b, err := json.Marshal(src); err != nil {
		logrus.Error(err)
		return ""
	} else {
		return string(b)
	}
}

func HttpRequest(url string, body interface{}, params map[string]string, header map[string]string, method string, log bool) ([]byte, error) {
	var req *http.Request
	var err error

	if method == http.MethodGet {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, strings.NewReader(marshal(body)))
	}

	if err != nil {
		return nil, err
	}

	if len(params) > 0 {
		q := req.URL.Query()
		for key, val := range params {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}

	for key, val := range header {
		req.Header.Add(key, val)
	}
	if _, ok := header["Content-Type"]; !ok {
		req.Header.Set("Content-Type", "application/json")
	}
	if _, ok := header["Host"]; ok {
		req.Host = header["Host"]
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, syscall.ECONNRESET) {
			rsp, err = http.DefaultClient.Do(req)
		}
		if err != nil {
			return nil, err
		}
	}

	defer rsp.Body.Close()

	rspBody, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	if log {
		logrus.Infof("[http-api] status code: %d, url: %s, body: %s, params: %s, header: %s, method: %s, response: %s\n",
			rsp.StatusCode, url, marshal(body), marshal(params), marshal(header), method, string(rspBody))
		logrus.Infof("[http-api] response header %s \n", marshal(rsp.Header))
	}
	if rsp.StatusCode >= 300 {
		return nil, errors.New("http status not ok: " + strconv.Itoa(rsp.StatusCode))
	}
	return rspBody, nil
}

func Retry(retryInterval time.Duration, times int, mark string, call func() error) error {
	var count int
	var err error
	for {
		count++
		if count > times {
			return err
		}
		err = call()
		if err == nil {
			return nil
		}
		logrus.WithField("mark", mark).Warn(err)
		time.Sleep(retryInterval)
	}
}
