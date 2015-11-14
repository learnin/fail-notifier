package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"

	"github.com/codegangsta/cli"

	"github.com/learnin/fail-notifier/plugin"
)

type config struct {
	Plugins []interface{}
}

var plugins = map[string]*plugin.Plugin{}

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
		cli.StringFlag{
			Name:  "config, c",
			Value: "config.json",
			Usage: "specify the path to the configuration file",
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

	cfg, err := loadConfig(c.String("config"))
	if err != nil {
		panic(err)
	}
	if len(cfg.Plugins) == 0 {
		// FIXME
		println("config error. no plugins")
		return
	}
	if !setupPlugins(&cfg) {
		// FIXME
		println("config error. the plugin does not exist.")
		return
	}

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
				notice(fmt.Sprintf("command failed. exitStatus=%v stdout=%v stderr=%v", s.ExitStatus(), stdout.String(), stderr.String()))
				return
			} else {
				// Unix や Winodws とは異なり、 exec.ExitError.Sys() が syscall.WaitStatus ではないOSの場合
				notice(fmt.Sprintf("command failed. stdout=%v stderr=%v", stdout.String(), stderr.String()))
				return
			}
		} else {
			// may be returned for I/O problems.
			notice(fmt.Sprintf("command can't execute. err=%v", err))
			return
		}
	}
	notice(stdout.String())
}

func loadConfig(path string) (config, error) {
	var cfg config
	cfgJson, err := ioutil.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(cfgJson, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func setupPlugins(cfg *config) bool {
	for _, p := range cfg.Plugins {
		p, ok := p.(map[string]interface{})
		if !ok {
			return false
		}
		typeName, ok := p["name"].(string)
		if !ok {
			return false
		}
		typeName = "plugin." + typeName
		pluginInstance, ok := plugin.CreatePlugin(typeName, p)
		if !ok {
			return false
		}
		plugins[typeName] = pluginInstance
	}
	return true
}

func notice(s string) {
	for _, p := range plugins {
		(*p).Notice(s)
	}
}
