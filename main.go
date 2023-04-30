package main

import (
	"time"

	"github.com/moqsien/xtray/example"
	_ "github.com/moqsien/xtray/pkgs/client"
	_ "github.com/moqsien/xtray/pkgs/conf"
	_ "github.com/moqsien/xtray/pkgs/proxy"
)

func init() {
	var cstZone = time.FixedZone("CST", 8*3600)
	time.Local = cstZone
}

func main() {
	// config := conf.NewConf()
	// verifier := proxy.NewVerifier(config)
	// verifier.Run(true)
	// fmt.Println(verifier.VerifiedProxies.VList.List)

	// client.TestTrojan("trojan://b5e4e360-5946-470b-aad0-db98f50faa57@frontend.yijianlian.app:54430?security=tls&type=tcp&headerType=none#%F0%9F%87%BA%F0%9F%87%B8%20Relay%20%F0%9F%87%BA%F0%9F%87%B8%20United%20States%2011%20TG%3A%40SSRSUB")
	// client.TestVless("vless://b1e41627-a3e9-4ebd-9c92-c366dd82b13f@xray.ibgfw.top:2083?encryption=none&security=tls&type=ws&host=&path=/wSXCvstU/#xray.ibgfw.top%3A2083")
	// client.TestVmess("vmess://eyJ2IjogIjIiLCAicHMiOiAiZ2l0aHViLmNvbS9mcmVlZnEgLSBcdTdmOGVcdTU2ZmRDbG91ZEZsYXJlXHU1MTZjXHU1M2Y4Q0ROXHU4MjgyXHU3MGI5IDEiLCAiYWRkIjogIm1pY3Jvc29mdGRlYnVnLmNvbSIsICJwb3J0IjogIjgwIiwgImlkIjogIjEwMTdlZjZhLTY3ZDktNGJiMy1iNjY3LTBkNjdjMWVlNTU0NiIsICJhaWQiOiAiMCIsICJzY3kiOiAiYXV0byIsICJuZXQiOiAid3MiLCAidHlwZSI6ICJub25lIiwgImhvc3QiOiAidjEudXM5Lm1pY3Jvc29mdGRlYnVnLmNvbSIsICJwYXRoIjogIi9zZWN4IiwgInRscyI6ICIiLCAic25pIjogIiJ9")
	// client.TestSS("ss://Y2hhY2hhMjAtaWV0Zi1wb2x5MTMwNTo3MjgyMjliOS0xNjRlLTQ1Y2ItYmZiMy04OTZiM2EwNTZhMTg=@node03.gde52px1vwf5q6301fxn.catapi.management:33907#%F0%9F%87%AC%F0%9F%87%A7%20Relay%20%F0%9F%87%AC%F0%9F%87%A7%20United%20Kingdom%2005%20TG%3A%40SSRSUB")
	// client.TestSSR("ssr://OTEuMjA2LjkyLjIyNzoxNzE5MTpvcmlnaW46cmM0OnBsYWluOmJHNWpiaTV2Y21jZ05ubG8vP29iZnNwYXJhbT0mcmVtYXJrcz01TC1FNTcyWDVwYXZVQSZncm91cD1URzVqYmk1dmNtYw")
	example.Start()
}
