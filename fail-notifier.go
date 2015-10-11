package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/codegangsta/cli"
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
		cli.StringFlag{
			Name:  "command, c",
			Usage: "execute given command",
		},
		cli.BoolFlag{
			Name:  "version, v",
			Usage: "print the version",
		},
	}
	app.Usage = "Send notifications when a given command fails"
	app.Action = func(c *cli.Context) {
		action(c)
	}
	app.Run(os.Args)
}

func action(c *cli.Context) {
	command := c.String("c")
	cmd := exec.Command(command)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if err2, ok := err.(*exec.ExitError); ok {
			if s, ok := err2.Sys().(syscall.WaitStatus); ok {
				println(fmt.Sprintf("command failed. exitStatus=%v stdout=%v stderr=%v", s.ExitStatus(), stdout.String(), stderr.String()))
				return
			} else {
				// Unix や Winodws とは異なり、 exec.ExitError.Sys() が syscall.WaitStatus ではないOSの場合
				println(fmt.Sprintf("command failed. stdout=%v stderr=%v", stdout.String(), stderr.String()))
				return
			}
		} else {
			// may be returned for I/O problems.
			println(fmt.Sprintf("command can't execute. err=%v", err))
			return
		}
	}
	println(stdout.String())
}