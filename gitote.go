package main

import (
	"os"

	"github.com/urfave/cli"
	"gitote.com/gitote/gitote/cmd"
	"gitote.com/gitote/gitote/pkg/setting"
)

const APP_VER = "0.0.2-alpha-rc7"

func init() {
	setting.AppVer = APP_VER
}

func main() {
	app := cli.NewApp()
	app.Name = "Gitote"
	app.Usage = "Software version control made simple"
	app.Version = APP_VER
	app.Commands = []cli.Command{
		cmd.Web,
		cmd.Serv,
		cmd.Hook,
		cmd.Cert,
		cmd.Admin,
		cmd.Import,
		cmd.Backup,
		cmd.Restore,
	}
	app.Run(os.Args)
}
