package ctrl

import (
	"fmt"
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
	if err := server.Start(); err != nil {
		fmt.Println("[start server failed] ", err)
	}
}

func (that *XKeeper) PingKeeper() bool {
	xc := utils.NewUClient(that.ksockName)
	if resp, err := xc.GetResp("/ping", map[string]string{}); err == nil {
		return strings.Contains(resp, XtrayOK)
	}
	return false
}

func (that *XKeeper) SendQuitSig() string {
	xc := utils.NewUClient(that.ksockName)
	resp, _ := xc.GetResp("/stop", map[string]string{})
	return resp
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
