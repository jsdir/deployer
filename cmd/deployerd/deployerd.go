package main

import (
	"os"

	"github.com/jsdir/deployer"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()

	app.Usage = "The deployer server"
	app.Author = ""
	app.Email = ""
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "config.json",
			Usage: "Location of the config file",
		},
		cli.IntFlag{
			Name:  "port, p",
			Value: 3000,
			Usage: "The port to listen on",
		},
	}
	app.Action = func(c *cli.Context) {
		config := deployer.ServerConfig{Port: c.Int("port")}
		db := deployer.StartDb()
		defer db.Close()

		deployer.LoadConfig(c.String("config"), &config)
		deployer.CreateServer(&deployer.Context{
			Db:        db,
			Config:    &config,
			NewBuilds: deployer.NewBroadcaster(),
		})
	}

	app.Run(os.Args)
}
