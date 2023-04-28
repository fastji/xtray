package ctrl

import (
	"os"

	"github.com/moqsien/xtray/pkgs/client"
	"github.com/moqsien/xtray/pkgs/conf"
	"github.com/moqsien/xtray/pkgs/proxy"
	cron "github.com/robfig/cron/v3"
)

var StopChan = make(chan struct{})

type XRunner struct {
	Client   *client.XClient
	Verifier *proxy.Verifier
	Conf     *conf.Conf
	Cron     *cron.Cron
}

func NewXRunner(cnf *conf.Conf) *XRunner {
	return &XRunner{
		Client:   client.NewXClient(),
		Verifier: proxy.NewVerifier(cnf),
		Conf:     cnf,
		Cron:     cron.New(),
	}
}

func (that *XRunner) Start() {
	if !that.Verifier.IsRunning {
		that.Verifier.Run(true)
	}
	that.Cron.AddFunc("@every 2h", func() {
		if !that.Verifier.IsRunning {
			that.Verifier.Run(false)
		}
	})
	that.Cron.Start()
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
