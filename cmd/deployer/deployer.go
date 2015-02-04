package main

import (
	"os"
	"log"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/codegangsta/cli"
)

func main() {
	config := loadConfig()
	createCli(config)
}

type CliConfig struct {
	Addr string `json:"addr"`
}

// loadConfig loads config from the .Deployer file in the current directory.
func loadConfig() *CliConfig {
	data, err := ioutil.ReadFile(".Deployer")
	if err != nil {
		log.Fatal(err)
	}

	config := CliConfig{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("Error parsing .Deployer: ", err)
	}

	return &config
}

func createCli(config *CliConfig) {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:      "release",
			ShortName: "r",
			Usage:     "Create a release",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "environment, e",
					Usage: "Environment to deploy the release to",
				},
			},
			Action: func(c *cli.Context) {
				release := createRelease(config, c)

				// Deploy to environment if it is specified.
				env := c.String("environment")
				if env != "" {
					createDeploy(config, release, env)
				}
			},
		},
		{
			Name:      "deploy",
			ShortName: "d",
			Usage:     "Deploys a release to an environment",
			Action: func(c *cli.Context) {
				args := c.Args()
				src := args.Get(0)
				if src == "" {
					log.Fatal("A valid source release or environment is required")
				}

				dest := args.Get(1)
				if dest == "" {
					log.Fatal("A valid destination environment is required")
				}

				createDeploy(config, src, dest)
			},
		},
	}

	app.Run(os.Args)
}

func createRelease(config *CliConfig, c *cli.Context) string {
	// Validate arguments
	args := c.Args()
	service := args.Get(0)
	if service == "" {
		log.Fatal("A valid service is required")
	}

	build := args.Get(1)
	if build == "" {
		log.Fatal("A valid build is required")
	}

	log.Println("Creating release...")

	values := url.Values{}
	values.Set("service", service)
	values.Add("build", build)

	res, err := http.PostForm(config.Addr + "/releases", values)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf(string(body[:]))

	return "1"
}

func createDeploy(config *CliConfig, src string, dest string) {
	log.Println("Creating deploy...")

	values := url.Values{}
    values.Set("src", src)
    values.Add("dest", dest)

	res, err := http.PostForm(config.Addr + "/deploys", values)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf(string(body[:]))
}
