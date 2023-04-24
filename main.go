package main

import (
	"fmt"
	"net/url"
)

func main() {
	u := "vless://b1e41627-a3e9-4ebd-9c92-c366dd82b13f@xray.ibgfw.top:2083?encryption=none&security=tls&type=ws&host=&path=/wSXCvstU/#xray.ibgfw.top%3A2083"
	r, _ := url.Parse(u)
	fmt.Println(r.Scheme)
	fmt.Println(r.Hostname())
	fmt.Println(r.Port())
	fmt.Println(r.User.Username())
	fmt.Println(r.User.Password())
	fmt.Println(r.Query().Get("security"))
	fmt.Println(r.Query().Get("encryption"))
	fmt.Println(r.Query().Get("path"))
}
