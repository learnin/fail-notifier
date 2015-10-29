package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/codegangsta/cli"

	"github.com/learnin/fail-notifier/plugin"
)

func main() {
	app := cli.NewApp()
	app.Name = "fail-notifier"
	app.Version = "0.0.1"
	app.Author = "Manabu Inoue"
	app.Email = ""
	app.HideVersion = true
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "help, h",
			Usage: "show help",
		},
		cli.BoolFlag{
			Name:  "version, v",
			Usage: "show the version",
		},
	}
	app.HideHelp = true
	app.ArgsUsage = "command"
	app.Usage = "Send notifications when a given command fails"
	app.Action = func(c *cli.Context) {
		action(c)
	}
	app.Run(os.Args)
}

func action(c *cli.Context) {
	if !(c.Args().Present()) {
		cli.ShowAppHelp(c)
		return
	}

	// FIXME 使用する plugin は設定ファイルで定義するようにする
	typeName := "plugin.Stdout"
	p, _ := plugin.CreatePlugin(typeName)

	var cmd *exec.Cmd
	if len(c.Args()) == 1 {
		cmd = exec.Command(c.Args().First())
	} else {
		cmd = exec.Command(c.Args().First(), c.Args().Tail()...)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if err2, ok := err.(*exec.ExitError); ok {
			if s, ok := err2.Sys().(syscall.WaitStatus); ok {
				// FIXME 文字列渡しではなく、exitStatus/stdout/stderrをもった構造体を渡すようにする
				p.Notice(fmt.Sprintf("command failed. exitStatus=%v stdout=%v stderr=%v", s.ExitStatus(), stdout.String(), stderr.String()))
				return
			} else {
				// Unix や Winodws とは異なり、 exec.ExitError.Sys() が syscall.WaitStatus ではないOSの場合
				p.Notice(fmt.Sprintf("command failed. stdout=%v stderr=%v", stdout.String(), stderr.String()))
				return
			}
		} else {
			// may be returned for I/O problems.
			p.Notice(fmt.Sprintf("command can't execute. err=%v", err))
			return
		}
	}
	p.Notice(stdout.String())
}
