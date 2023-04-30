package ctrl

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/moqsien/xtray/pkgs/conf"
	"github.com/moqsien/xtray/pkgs/utils"
	cron "github.com/robfig/cron/v3"
)

const (
	KeeperSocksName = "xtray_keeper"
)

type XKeeper struct {
	runner    *XRunner
	ksockName string
	Conf      *conf.Conf
	Cron      *cron.Cron
}

func NewXKeeper(cnf *conf.Conf, runner *XRunner) *XKeeper {
	return &XKeeper{
		runner:    runner,
		ksockName: KeeperSocksName,
		Conf:      cnf,
		Cron:      cron.New(),
	}
}

func (that *XKeeper) runServer() {
	server := utils.NewUServer(that.ksockName)
	server.AddHandler("/stop", func(c *gin.Context) {
		StopChan <- struct{}{}
		c.String(http.StatusOK, "xtray keeper is stopped.")
	})
	server.AddHandler("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, XtrayOK)
	})
	server.Start()
}

func (that *XKeeper) PingKeeper() bool {
	xc := utils.NewUClient(that.ksockName)
	if resp, _ := xc.GetResp("/ping", map[string]string{}); resp == XtrayOK {
		return true
	}
	return false
}

func (that *XKeeper) checkRunner() {
	if !that.runner.PingXtray() {
		that.runner.starter.Run()
	}
}

func (that *XKeeper) Run() {
	utils.DaemonizeInit()
	go that.runServer()
	cronTime := that.Conf.KeeperCron
	if !strings.HasPrefix(cronTime, "@every") {
		cronTime = "@every 3m"
	}
	that.Cron.AddFunc(cronTime, that.checkRunner)
	that.Cron.Start()
	<-StopChan
}
