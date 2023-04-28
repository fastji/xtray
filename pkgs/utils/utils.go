package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os/user"
	"strings"
)

func NormalizeSSR(s string) (r string) {
	r = strings.ReplaceAll(s, "-", "+")
	r = strings.ReplaceAll(r, "_", "/")
	return r
}

func DecodeBase64(str string) (res string) {
	count := len(str) % 4
	for i := 0; i < count; i++ {
		str += "="
	}
	if s, err := base64.StdEncoding.DecodeString(str); err == nil {
		res = string(s)
	}
	return
}

func StringToReader(str string) io.Reader {
	return bytes.NewReader([]byte(str))
}

func GetHomeDir() (homeDir string) {
	u, err := user.Current()
	if err != nil {
		fmt.Println("[CurrentUser]", err)
		return
	}
	return u.HomeDir
}
