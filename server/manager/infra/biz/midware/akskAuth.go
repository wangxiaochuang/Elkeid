package midware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/levigross/grequests"
	"github.com/wangxiaochuang/Elkeid/server/manager/infra"
)

var (
	AK = infra.AccessKey
	SK = infra.SecretKey
)

func sha256byteArr(in []byte) string {
	if in == nil || len(in) == 0 {
		return ""
	}
	h := sha256.New()
	h.Write(in)
	return hex.EncodeToString(h.Sum(nil))
}

func generateSign(method, url, query, ak, timestamp, sk string, requestBody []byte) string {
	return hmacSha256(fmt.Sprintf(`%s\n%s\n%s\n%s\n%s\n%s`, method, url, query, ak, timestamp, sha256byteArr(requestBody)), sk)
}

func hmacSha256(data string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func formatURLPath(in string) string {
	in = strings.TrimSpace(in)
	if strings.HasSuffix(in, "/") {
		return in[:len(in)-1]
	}
	return in
}

func beforeRequestFuncWithKey(req *http.Request, ak, sk string) error {
	var (
		timestamp   = fmt.Sprintf(`%d`, time.Now().Unix())
		err         error
		requestBody []byte
	)

	if req.Body != nil {
		requestBody, err = ioutil.ReadAll(req.Body)
		if err != nil {
			return err
		}
		//Reset after reading
		req.Body.Close()
		req.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
	} else {
		requestBody = []byte{}
	}
	sign := generateSign(req.Method, formatURLPath(req.URL.Path), req.URL.RawQuery, ak, timestamp, sk, requestBody)
	req.Header.Add("AccessKey", ak)
	req.Header.Add("Signature", sign)
	req.Header.Add("TimeStamp", timestamp)
	return nil
}

func sdBeforeRequestFunc(req *http.Request) error {
	beforeRequestFuncWithKey(req, infra.SdAK, infra.SdSK)
	return nil
}

func SdAuthRequestOption() *grequests.RequestOptions {
	option := &grequests.RequestOptions{
		InsecureSkipVerify: true,
		BeforeRequest:      sdBeforeRequestFunc,
	}
	return option
}
