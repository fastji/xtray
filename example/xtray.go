package example

import (
	"os"

	"github.com/moqsien/xtray/pkgs/conf"
	"github.com/moqsien/xtray/pkgs/ctrl"
)

type XtrayExa struct {
	Conf   *conf.Conf
	Runner *ctrl.XRunner
	Keeper *ctrl.XKeeper
}

func NewXtrayExa() *XtrayExa {
	xe := &XtrayExa{
		Conf: conf.NewConf(),
	}
	xe.Runner = ctrl.NewXRunner(xe.Conf)
	xe.Runner.RegisterStarter(Starter)
	xe.Runner.RegisterKeeper(Keeper)
	xe.Keeper = ctrl.NewXKeeper(xe.Conf, xe.Runner)
	return xe
}

func Start() {
	app := NewApps()
	app.Run(os.Args)
}
