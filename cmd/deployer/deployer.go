package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/jsdir/deployer/pkg/resources"

	"github.com/codegangsta/cli"
)

func main() {
	app := CreateCliApp("http://localhost:7654")
	app.Run(os.Args)
}

func CreateCliApp(addr string) *cli.App {
	app := cli.NewApp()

	app.Usage = "A cli for deployerd"
	app.Author = ""
	app.Email = ""
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "addr",
			Value: addr,
			Usage: "Address and port of the deployerd instance",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "release",
			ShortName: "r",
			Usage:     "Creates a release",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "environment, e",
					Usage: "Environment to deploy the release to (optional)",
				},
			},
			Action: handleRelease,
		},
		{
			Name:      "deploy",
			ShortName: "d",
			Usage:     "Deploys a release to an environment",
			Action:    handleDeploy,
		},
	}

	return app
}

func handleRelease(c *cli.Context) {
	release, err := createRelease(c)
	handleErr(err)

	// Deploy to environment if it is specified.
	env := c.String("environment")
	if env != "" {
		createDeploy(c, string(release.Id), env)
	}
}

func handleDeploy(c *cli.Context) {
	args := c.Args()
	src := args.Get(0)
	if src == "" {
		log.Fatal("A source release or environment is required")
	}

	dest := args.Get(1)
	if dest == "" {
		log.Fatal("A valid destination environment is required")
	}

	err := createDeploy(c, src, dest)
	handleErr(err)
}

func createRelease(c *cli.Context) (*resources.Release, error) {
	// Validate arguments
	args := c.Args()
	if len(args)%2 == 1 {
		return nil, errors.New("Invalid builds")
	}

	// Send request
	fmt.Println("Creating release...")

	//values := url.Values{}
	//values.Set("service", service)
	//values.Add("tag", tag)

	res, err := http.PostForm(c.GlobalString("addr")+"/releases", values)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	release := new(resources.Release)
	json.Unmarshal(data, &release)
	fmt.Println(string(data[:]))
	fmt.Println(release.Id)
	return release, nil
}

func createDeploy(c *cli.Context, src string, dest string) error {
	fmt.Println("Deploying...")

	values := url.Values{}
	values.Set("src", src)
	values.Add("dest", dest)

	res, err := http.PostForm(c.GlobalString("addr")+"/deploys", values)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	//fmt.Printf(string(body[:]))
	return nil
}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
