package environments

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/jsdir/deployer/pkg/resources"

	"github.com/mitchellh/mapstructure"
)

type KubernetesConfig struct {
	ManifestGlob string
	Cmd          string
}

type Kubernetes struct{}

func (k *Kubernetes) Deploy(deploy *resources.Deploy) error {
	// Get environment config.
	config := KubernetesConfig{
		ManifestGlob: "manifests/*.json",
		Cmd:          "kubectl",
	}
	err := mapstructure.Decode(deploy.EnvConfig, &config)
	if err != nil {
		return err
	}

	log.Println(config)

	// Iterate through and parse templates.
	templates, err := filepath.Glob(config.ManifestGlob)
	if err != nil {
		return err
	}

	// Execute commands in parallel.
	count := len(templates)
	if count == 0 {
		log.Println("info: no templates found")
		return nil
	}

	sem := make(chan error, count)
	for _, filename := range templates {
		go func(filename string) {
			sem <- updateManifest(filename, deploy, &config)
			return
		}(filename)
	}

	for i := 0; i < count; i++ {
		err := <-sem
		if err != nil {
			log.Println(err)
			return errors.New("could not run command")
		}
	}

	return nil
}

func updateManifest(filename string, deploy *resources.Deploy, config *KubernetesConfig) error {
	log.Printf("Uploading %s", filename)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	tmpl, err := template.New("").Parse(string(data))
	if err != nil {
		return err
	}

	cmd := exec.Command(config.Cmd, "update", "-f", "-")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	err = tmpl.Execute(stdin, deploy)
	if err != nil {
		return err
	}

	stdin.Close()
	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}
