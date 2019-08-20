package util

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"io"
)

func SHA1Encode(encodeStr string) string {
	h := sha1.New()
	io.WriteString(h, encodeStr)
	s := fmt.Sprintf("%x", h.Sum(nil))
	return  s
}


func HmacSHA1(key, str string) string {
	//hmac ,use sha1
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(str))
	return fmt.Sprintf("%x\n", mac.Sum(nil))
	//return string(mac.Sum(nil))
}
