package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	sproxy "github.com/Asutorufa/yuhaiin/pkg/net/interfaces/proxy"
	"github.com/Asutorufa/yuhaiin/pkg/node/register"
	"github.com/Asutorufa/yuhaiin/pkg/protos/node/point"
	"github.com/Asutorufa/yuhaiin/pkg/protos/node/protocol"
	"github.com/moqsien/free/pkgs/runner"
	futils "github.com/moqsien/free/pkgs/utils"
	"github.com/moqsien/xtray/pkgs/client"
	"github.com/moqsien/xtray/pkgs/conf"
	"github.com/moqsien/xtray/pkgs/utils"
)

type VList struct {
	List       []*Proxy `koanf,json:"list"`
	UpdateTime string   `koanf,json:"update_time"`
}

type VerifiedList struct {
	VList *VList
	*utils.Koanfer
}

func NewVerifiedList(p string) *VerifiedList {
	return &VerifiedList{
		VList:   &VList{List: []*Proxy{}},
		Koanfer: utils.NewKoanfer(p),
	}
}

func (that *VerifiedList) Save() {
	that.VList.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	that.Koanfer.Save(that.VList)
}

func (that *VerifiedList) Load() {
	that.Koanfer.Load(that.VList)
}

func (that *VerifiedList) GetByIndex(idx int) (p string) {
	that.Load()
	if len(that.VList.List) == 0 {
		return
	}
	if idx < 0 || idx >= len(that.VList.List) {
		idx = 0
	}
	return that.VList.List[idx].RawUri
}

type Verifier struct {
	conf            *conf.Conf
	RawProxies      *runner.Result
	VerifiedProxies *VerifiedList
	ProxyChan       chan IProxy
	IsRunning       bool
	fetcher         *Fetcher
	wg              *sync.WaitGroup
	lock            *sync.RWMutex
}

func NewVerifier(conf *conf.Conf) *Verifier {
	v := &Verifier{
		conf:    conf,
		fetcher: NewFetcher(conf),
		RawProxies: &runner.Result{
			VmessList: &runner.VList{List: []string{}},
			VlessList: &runner.VList{List: []string{}},
			SSRList:   &runner.VList{List: []string{}},
			SSList:    &runner.VList{List: []string{}},
			Trojan:    &runner.VList{List: []string{}},
		},
		wg:   &sync.WaitGroup{},
		lock: &sync.RWMutex{},
	}
	v.VerifiedProxies = NewVerifiedList(v.conf.PorxyFile)
	return v
}

func (that *Verifier) Reload(force bool) {
	if ok, _ := futils.PathIsExist(that.conf.RawProxyFile); !ok || force {
		that.fetcher.GetFile()
	}
	if rawProxy, err := os.ReadFile(that.conf.RawProxyFile); err == nil {
		json.Unmarshal(rawProxy, that.RawProxies)
	}
}

func (that *Verifier) getRawList() (result []string) {
	result = append(result, that.RawProxies.VmessList.List...)
	result = append(result, that.RawProxies.VlessList.List...)
	result = append(result, that.RawProxies.SSList.List...)
	// result = append(result, that.RawProxies.SSRList.List...)
	result = append(result, that.RawProxies.Trojan.List...)
	return
}

func (that *Verifier) dispatchProxies() {
	that.ProxyChan = make(chan IProxy, 30)
	for _, v := range that.getRawList() {
		p := &Proxy{RawUri: v}
		that.ProxyChan <- p
	}
	close(that.ProxyChan)
}

func (that *Verifier) stopXClient(wait chan struct{}) {
	wait <- struct{}{}
}

func (that *Verifier) sendReq(param *client.ClientParams, wait chan struct{}) {
	node := &point.Point{
		Protocols: []*protocol.Protocol{
			{
				Protocol: &protocol.Protocol_Simple{
					Simple: &protocol.Simple{
						Host:             "127.0.0.1",
						Port:             int32(param.InPort),
						PacketConnDirect: true,
					},
				},
			},
			{
				Protocol: &protocol.Protocol_Socks5{
					Socks5: &protocol.Socks5{},
				},
			},
		},
	}

	pro, err := register.Dialer(node)
	if err != nil {
		fmt.Println("[Dialer error] ", err)
		that.stopXClient(wait)
		return
	}
	if that.conf.Timeout == 0 {
		that.conf.Timeout = 3
	}
	c := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				add, err := sproxy.ParseAddress(sproxy.PaseNetwork(network), addr)
				if err != nil {
					return nil, fmt.Errorf("parse address failed: %w", err)
				}
				return pro.Conn(ctx, add)
			}},
		Timeout: time.Duration(that.conf.Timeout) * time.Second,
	}
	if that.conf.TestUrl == "" {
		that.conf.TestUrl = "https://www.google.com"
	}
	startTime := time.Now()
	resp, err := c.Get(that.conf.TestUrl)
	timeLag := time.Since(startTime).Milliseconds()
	if err != nil {
		fmt.Println("[Verify url failed] ", param.RawUri, " ", err)
		that.stopXClient(wait)
		return
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	defer resp.Body.Close()
	if strings.Contains(buf.String(), "</html>") {
		p := &Proxy{RawUri: param.RawUri, RTT: int(timeLag)}
		that.lock.Lock()
		that.VerifiedProxies.VList.List = append(that.VerifiedProxies.VList.List, p)
		that.lock.Unlock()
	}
	that.stopXClient(wait)
}

func (that *Verifier) RunXClient(port int) {
	if that.ProxyChan == nil {
		return
	}
	xclient := client.NewXClient()
	that.wg.Add(1)
	for {
		select {
		case pxy, ok := <-that.ProxyChan:
			if pxy == nil && !ok {
				that.wg.Done()
				return
			}
			param := &client.ClientParams{
				RawUri: pxy.GetRawUri(),
				InPort: port,
			}
			xclient.Start(param)
			wait := make(chan struct{})
			go that.sendReq(param, wait)
			<-wait
			xclient.Close()
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (that *Verifier) Run(force bool) {
	that.IsRunning = true
	that.VerifiedProxies = NewVerifiedList(that.conf.PorxyFile)
	that.Reload(force)
	go that.dispatchProxies()
	time.Sleep(time.Millisecond * 50)
	start := that.conf.PortRange.Start
	end := that.conf.PortRange.End
	if start > end {
		start, end = end, start
	}
	for port := start; port <= end; port++ {
		go that.RunXClient(port)
	}
	that.wg.Wait()
	that.VerifiedProxies.Save()
	that.IsRunning = false
}
