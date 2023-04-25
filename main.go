package main

import (
	"fmt"
	"strings"
)

func main() {
	u := "188.119.65.143:12337:origin:rc4:plain:bG5jbi5vcmcgbjh0/?obfsparam=&remarks=5L-E572X5pavTw&group=TG5jbi5vcmc"
	// r, _ := url.Parse(u)
	// if r == nil {
	// 	return
	// }
	// fmt.Println(r.Scheme)
	// fmt.Println(r.Hostname())
	// fmt.Println(r.Port())
	// fmt.Println(r.User.Username())
	// fmt.Println(r.User.Password())
	// fmt.Println(r.Query().Get("remarks"))
	// fmt.Println(r.Query().Get("group"))
	v := strings.SplitN(u, ":", 3)
	fmt.Println(v[0])
	fmt.Println(v[1])
	fmt.Println(v[2])
}
