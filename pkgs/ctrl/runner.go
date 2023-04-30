package ctrl

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mholt/archiver/v3"
	"github.com/moqsien/free/pkgs/query"
	futils "github.com/moqsien/free/pkgs/utils"
	"github.com/moqsien/goktrl"
	"github.com/moqsien/xtray/pkgs/client"
	"github.com/moqsien/xtray/pkgs/conf"
	"github.com/moqsien/xtray/pkgs/proxy"
	"github.com/moqsien/xtray/pkgs/utils"
	cron "github.com/robfig/cron/v3"
)

var StopChan = make(chan struct{})

const (
	ExtraSocksName = "xtray_runner"
	KtrlSocksName  = "xtray_ktrl"
	XtrayOK        = "ok"
)

type XRunner struct {
	Client     *client.XClient
	Verifier   *proxy.Verifier
	XKeeper    *XKeeper
	Conf       *conf.Conf
	Cron       *cron.Cron
	ExtraSocks string
	KtrlSocks  string
	Ktrl       *goktrl.Ktrl
	starter    *exec.Cmd
	keeper     *exec.Cmd
}

func NewXRunner(cnf *conf.Conf) (r *XRunner) {
	r = &XRunner{
		Client:     client.NewXClient(),
		Verifier:   proxy.NewVerifier(cnf),
		Conf:       cnf,
		Cron:       cron.New(),
		ExtraSocks: ExtraSocksName,
		KtrlSocks:  KtrlSocksName,
		Ktrl:       goktrl.NewKtrl(),
	}
	r.XKeeper = NewXKeeper(cnf, r)
	return r
}

func (that *XRunner) runServer() {
	server := utils.NewUServer(that.ExtraSocks)
	server.AddHandler("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, XtrayOK)
	})
	server.Start()
}

func (that *XRunner) PingXtray() bool {
	xc := utils.NewUClient(that.ExtraSocks)
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
	cronTime := that.Conf.VerifierCron
	if !strings.HasPrefix(cronTime, "@every") {
		cronTime = "@every 2h"
	}
	that.Cron.AddFunc(cronTime, func() {
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
	that.starter = starter
}

func (that *XRunner) RegisterKeeper(keeper *exec.Cmd) {
	that.keeper = keeper
}

// TODO: ctrl shell
func (that *XRunner) initCtrl() {
	that.Ktrl.AddKtrlCommand(&goktrl.KCommand{
		Name: "start",
		Help: "Start an xray-core client.",
		Func: func(c *goktrl.Context) {
			if that.starter == nil {
				fmt.Println("Please register a starter first.")
				return
			}
			if err := that.starter.Run(); err != nil {
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

			if that.keeper != nil {
				if err := that.keeper.Run(); err != nil {
					fmt.Println("Start a keeper failed: ", err)
					return
				} else {
					fmt.Println("Starting a keeper...")
					time.Sleep(time.Second * 3)
					if that.PingXtray() {
						fmt.Println("Start a keeper succeeded.")
						return
					}
					fmt.Println("Please check keeper status.")
				}
			}
		},
		KtrlHandler: func(c *goktrl.Context) {},
		SocketName:  that.KtrlSocks,
	})

	that.Ktrl.AddKtrlCommand(&goktrl.KCommand{
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

	that.Ktrl.AddKtrlCommand(&goktrl.KCommand{
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

	that.Ktrl.AddKtrlCommand(&goktrl.KCommand{
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

	that.Ktrl.AddKtrlCommand(&goktrl.KCommand{
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

	that.Ktrl.AddKtrlCommand(&goktrl.KCommand{
		Name: "status",
		Help: "Show xray-core client running status.",
		Func: func(c *goktrl.Context) {
			if that.PingXtray() {
				fmt.Println("xray-core client is running.")
			} else {
				fmt.Println("xray-core client is stopped.")
			}
			if that.XKeeper.PingKeeper() {
				fmt.Println("client keeper is running.")
			} else {
				fmt.Println("client keeper is stopped.")
			}
		},
		KtrlHandler: func(c *goktrl.Context) {},
		SocketName:  that.KtrlSocks,
	})

	that.Ktrl.AddKtrlCommand(&goktrl.KCommand{
		Name: "omega",
		Help: "Download switchy-omega plugin for Google Chrome Browser.",
		Func: func(c *goktrl.Context) {
			omegaPath := that.SwitchyOmega()
			if ok, _ := futils.PathIsExist(omegaPath); ok {
				fmt.Println("switchy-omega is unarchived in: ", omegaPath)
			} else {
				fmt.Println("download switchy-omega failed.")
			}
		},
		KtrlHandler: func(c *goktrl.Context) {},
		SocketName:  that.KtrlSocks,
	})
}

func (that *XRunner) SwitchyOmega() (omegaPath string) {
	omegaPath = filepath.Join(that.Conf.WorkDir, "switchy_omega")
	if ok, _ := futils.PathIsExist(omegaPath); ok {
		fmt.Println("[Archive Path] ", omegaPath)
		return
	}
	fName := "switchy-omega.zip"
	fpath := filepath.Join(that.Conf.WorkDir, fName)
	d := query.NewDownloader(that.Conf.SwitchyOmegaUrl)
	d.File(fpath)
	if ok, _ := futils.PathIsExist(fpath); ok {
		if err := archiver.Unarchive(fpath, omegaPath); err != nil {
			os.RemoveAll(fpath)
			os.RemoveAll(omegaPath)
			fmt.Println("[Unarchive failed] ", err)
			return
		} else {
			fmt.Println("Swithy-Omega Download Succeeded.")
			fmt.Println("[Archive Path] ", omegaPath)
		}
	}
	return
}

func (that *XRunner) CtrlServer() {
	that.initCtrl()
	that.Ktrl.RunCtrl(that.KtrlSocks)
}

func (that *XRunner) CtrlShell() {
	that.initCtrl()
	that.Ktrl.RunShell(that.KtrlSocks)
}
