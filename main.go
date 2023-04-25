package main

import (
	"fmt"
	"net/url"
)

func main() {
	u := "trojan://b5e4e360-5946-470b-aad0-db98f50faa57@frontend.yijianlian.app:54430?security=tls&type=tcp&headerType=none#%F0%9F%87%BA%F0%9F%87%B8%20Relay%20%F0%9F%87%BA%F0%9F%87%B8%20United%20States%2011%20TG%3A%40SSRSUB"
	r, _ := url.Parse(u)
	if r == nil {
		return
	}
	fmt.Println(r.Scheme)
	fmt.Println(r.Hostname())
	fmt.Println(r.Port())
	fmt.Println(r.User.Username())
	fmt.Println(r.User.Password())
	fmt.Println(r.Query().Get("security"))
	fmt.Println(r.Query().Get("type"))
}
