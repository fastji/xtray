package ctrl

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/moqsien/goktrl"
	"github.com/moqsien/xtray/pkgs/client"
	"github.com/moqsien/xtray/pkgs/conf"
	"github.com/moqsien/xtray/pkgs/proxy"
	"github.com/moqsien/xtray/pkgs/utils"
	cron "github.com/robfig/cron/v3"
)

var StopChan = make(chan struct{})

const (
	XtrayOK = "ok"
)

type XRunner struct {
	Client    *client.XClient
	Verifier  *proxy.Verifier
	Conf      *conf.Conf
	Cron      *cron.Cron
	AddSocks  string
	KtrlSocks string
	Ktrcl     *goktrl.Ktrl
	starters  []*exec.Cmd
}

func NewXRunner(cnf *conf.Conf) *XRunner {
	return &XRunner{
		Client:    client.NewXClient(),
		Verifier:  proxy.NewVerifier(cnf),
		Conf:      cnf,
		Cron:      cron.New(),
		AddSocks:  "xtray_runner",
		KtrlSocks: "xtray_ktrl",
		Ktrcl:     goktrl.NewKtrl(),
		starters:  []*exec.Cmd{},
	}
}

func (that *XRunner) runServer() {
	server := utils.NewUServer(that.AddSocks)
	server.AddHandler("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, XtrayOK)
	})
	server.Start()
}

func (that *XRunner) PingXtray() bool {
	xc := utils.NewUClient(that.AddSocks)
	if resp, _ := xc.GetResp("/ping", map[string]string{}); resp == XtrayOK {
		return true
	}
	return false
}

func (that *XRunner) Start() {
	if that.PingXtray() {
		fmt.Println("xtray is already running.")
		return
	}
	utils.DaemonizeInit()
	if !that.Verifier.IsRunning {
		that.Verifier.Run(true)
	}
	that.Cron.AddFunc("@every 2h", func() {
		if !that.Verifier.IsRunning {
			that.Verifier.Run(false)
		}
	})
	that.Cron.Start()
	go that.runServer()
	that.Restart(0)
	<-StopChan
	os.Exit(0)
}

func (that *XRunner) Restart(idx int) (result string) {
	if that.Client == nil {
		that.Client = client.NewXClient()
	}
	that.Client.Close()
	rawUri := that.Verifier.VerifiedProxies.GetByIndex(idx)
	if rawUri != "" {
		that.Client.Start(&client.ClientParams{
			RawUri: rawUri,
			InPort: that.Conf.Port,
		})
		result = fmt.Sprintf("%d.%s", idx, that.Client.Out.GetString())
	}
	return
}

func (that *XRunner) Stop() {
	StopChan <- struct{}{}
}

func (that *XRunner) RegisterStarter(starter *exec.Cmd) {
	that.starters = append(that.starters, starter)
}

// TODO: ctrl shell
func (that *XRunner) initCtrl() {
	that.Ktrcl.AddKtrlCommand(&goktrl.KCommand{
		Name: "start",
		Help: "Start an xray-core client.",
		Func: func(c *goktrl.Context) {
			if len(that.starters) == 0 {
				fmt.Println("No [starter] has been registered.")
				return
			}

			for _, starter := range that.starters {
				if err := starter.Run(); err != nil {
					fmt.Println("Start a client failed: ", err)
					return
				} else {
					fmt.Println("Starting a client...")
					time.Sleep(time.Second * 3)
					if that.PingXtray() {
						fmt.Println("Start a client succeeded.")
						return
					}
					fmt.Println("Please check client status.")
				}
			}
		},
		KtrlHandler: func(c *goktrl.Context) {},
		SocketName:  that.KtrlSocks,
	})

	that.Ktrcl.AddKtrlCommand(&goktrl.KCommand{
		Name: "stop",
		Help: "Stop the running xray-core client.",
		Func: func(c *goktrl.Context) {
			result, _ := c.GetResult()
			if len(result) > 0 {
				fmt.Println(string(result))
				// TODO: stop keeper
				// that.sendQuitSignal()
			}
		},
		KtrlHandler: func(c *goktrl.Context) {
			that.Stop()
			c.Send("xray-core client stopped.", 200)
		},
		SocketName: that.KtrlSocks,
	})

	that.Ktrcl.AddKtrlCommand(&goktrl.KCommand{
		Name: "restart",
		Help: "Restart the running xray-core client.",
		Func: func(c *goktrl.Context) {
			result, _ := c.GetResult()
			if len(result) > 0 {
				fmt.Println(string(result))
			}
		},
		ArgsDescription: "choose a specified proxy by index.",
		KtrlHandler: func(c *goktrl.Context) {
			idx := 0
			if len(c.Args) > 0 {
				idx, _ = strconv.Atoi(c.Args[0])
			}
			r := that.Restart(idx)
			c.Send(fmt.Sprintf("Restart client using [%s]", r), 200)
		},
		SocketName: that.KtrlSocks,
	})

	that.Ktrcl.AddKtrlCommand(&goktrl.KCommand{
		Name: "show",
		Help: "Show vpn list info.",
		Func: func(c *goktrl.Context) {
			that.Verifier.Reload(false)
			fmt.Println("Raw free vpn list statistics: ")
			fmt.Printf("vmess: %d, vless: %d, ss: %d, ssr: %d, trojan: %d, updated_at: %s",
				that.Verifier.RawProxies.VmessList.Total,
				that.Verifier.RawProxies.VlessList.Total,
				that.Verifier.RawProxies.SSList.Total,
				that.Verifier.RawProxies.SSRList.Total,
				that.Verifier.RawProxies.Trojan.Total,
				that.Verifier.RawProxies.UpdateTime)
			verifiedList := proxy.NewVerifiedList(that.Conf.PorxyFile)
			verifiedList.Load()
			fmt.Printf("verifed vpn list(@%s): ", verifiedList.VList.UpdateTime)
			for idx, v := range verifiedList.VList.List {
				fmt.Printf("%d. %s| rtt: %dms", idx, client.ParseRawUri(v.RawUri), v.RTT)
			}
		},
		KtrlHandler: func(c *goktrl.Context) {},
		SocketName:  that.KtrlSocks,
	})

	type filterOpts struct {
		Force bool `alias:"f" descr:"Force to get new raw vpn list."`
	}

	that.Ktrcl.AddKtrlCommand(&goktrl.KCommand{
		Name: "filter",
		Help: "Filter vpns by verifier.",
		Opts: &filterOpts{},
		Func: func(c *goktrl.Context) {
			result, _ := c.GetResult()
			if len(result) > 0 {
				fmt.Println(string(result))
			}
		},
		KtrlHandler: func(c *goktrl.Context) {
			if that.Verifier.IsRunning {
				c.Send("verifier is already running.", 200)
				return
			}
			opts := c.Options.(*filterOpts)
			go that.Verifier.Run(opts.Force)
			c.Send("verifier starts running.", 200)
		},
		SocketName: that.KtrlSocks,
	})

	that.Ktrcl.AddKtrlCommand(&goktrl.KCommand{
		Name: "status",
		Help: "Show xray-core client running status.",
		Func: func(c *goktrl.Context) {
			if that.PingXtray() {
				fmt.Println("xray-core client is running.")
				return
			}
			fmt.Println("xray-core client is stopped.")
		},
		KtrlHandler: func(c *goktrl.Context) {},
		SocketName:  that.KtrlSocks,
	})

	that.Ktrcl.AddKtrlCommand(&goktrl.KCommand{
		Name:        "omega",
		Help:        "Download switchy-omega plugin for Google Chrome Browser.",
		Func:        func(c *goktrl.Context) {},
		KtrlHandler: func(c *goktrl.Context) {},
		SocketName:  that.KtrlSocks,
	})
}

func (that *XRunner) CtrlServer() {
	that.initCtrl()
}

func (that *XRunner) CtrlShell() {
	that.initCtrl()
}
