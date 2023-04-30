package example

import (
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

type Apps struct {
	*cli.App
}

func NewApps() (a *Apps) {
	a = &Apps{
		App: &cli.App{},
	}
	a.initiate()
	return a
}

func (that *Apps) initiate() {
	command := &cli.Command{
		Name:    "shell",
		Aliases: []string{"sh", "s"},
		Usage:   "Start a shell for xtray.",
		Action: func(ctx *cli.Context) error {
			xe := NewXtrayExa()
			xe.Runner.CtrlShell()
			return nil
		},
	}
	that.Commands = append(that.Commands, command)

	command = &cli.Command{
		Name:    "runner",
		Aliases: []string{"run", "r"},
		Usage:   "Start an xtray runner.",
		Action: func(ctx *cli.Context) error {
			xe := NewXtrayExa()
			xe.Runner.Start()
			return nil
		},
	}
	that.Commands = append(that.Commands, command)

	command = &cli.Command{
		Name:    "keeper",
		Aliases: []string{"keep", "k"},
		Usage:   "Start an xtray keeper.",
		Action: func(ctx *cli.Context) error {
			xe := NewXtrayExa()
			xe.Keeper.Run()
			return nil
		},
	}
	that.Commands = append(that.Commands, command)
}

var cmdName string = func() string {
	epath, err := os.Executable()
	if err != nil {
		panic("cannot find executable path")
	}
	return epath
}()

var (
	Starter *exec.Cmd = exec.Command(cmdName, "runner")
	Keeper  *exec.Cmd = exec.Command(cmdName, "keeper")
)
