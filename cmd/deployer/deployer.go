package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/jsdir/deployer"

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

	app.Usage = "A cli for deployerd"
	app.Author = ""
	app.Email = ""
	app.Commands = []cli.Command{
		{
			Name:      "release",
			ShortName: "r",
			Usage:     "Creates a release",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "environment, e",
					Usage: "Environment to deploy the release to",
				},
			},
			Action: func(c *cli.Context) {
				createRelease(config, c)
				//release := createRelease(config, c)

				// Deploy to environment if it is specified.
				env := c.String("environment")
				if env != "" {
					//createDeploy(config, release, env)
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

func createRelease(config *CliConfig, c *cli.Context) *deployer.Release {
	// Validate arguments
	args := c.Args()
	service := args.Get(0)
	if service == "" {
		log.Fatal("A valid service is required")
	}

	tag := args.Get(1)
	if tag == "" {
		log.Fatal("A valid build tag is required")
	}

	log.Println("Creating release...")

	values := url.Values{}
	values.Set("service", service)
	values.Add("tag", tag)

	res, err := http.PostForm(config.Addr+"/releases", values)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	release := new(deployer.Release)
	json.Unmarshal(data, &release)
	println(string(data[:]))
	println(release.Id)
	return release
}

func createDeploy(config *CliConfig, src string, dest string) {
	log.Println("Creating deploy...")

	values := url.Values{}
	values.Set("src", src)
	values.Add("dest", dest)

	res, err := http.PostForm(config.Addr+"/deploys", values)
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
