package ctrl

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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
	Client   *client.XClient
	Verifier *proxy.Verifier
	Conf     *conf.Conf
	Cron     *cron.Cron
	SockName string
}

func NewXRunner(cnf *conf.Conf) *XRunner {
	return &XRunner{
		Client:   client.NewXClient(),
		Verifier: proxy.NewVerifier(cnf),
		Conf:     cnf,
		Cron:     cron.New(),
		SockName: "xtray_runner",
	}
}

func (that *XRunner) runServer() {
	server := utils.NewUServer(that.SockName)
	server.AddHandler("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, XtrayOK)
	})
	server.Start()
}

func (that *XRunner) PingXtray() bool {
	xc := utils.NewUClient(that.SockName)
	if resp, _ := xc.GetResp("/ping", map[string]string{}); resp == XtrayOK {
		return true
	}
	return false
}

func (that *XRunner) Start() {
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

func (that *XRunner) Restart(idx int) {
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
	}
}

func (that *XRunner) Stop() {
	StopChan <- struct{}{}
}

// TODO: ctrl shell
func (that *XRunner) initCtrl() {
}

func (that *XRunner) CtrlServer() {
	that.initCtrl()
}

func (that *XRunner) CtrlShell() {
	that.initCtrl()
}
