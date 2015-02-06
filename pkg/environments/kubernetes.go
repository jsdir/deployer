package environments

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/jsdir/deployer/pkg/resources"
)

type KubernetesConfig struct {
	ManifestGlob string
	Cmd          string
}

type Kubernetes struct{}

func (k *Kubernetes) Deploy(deploy *resources.Deploy) error {
	config := deploy.EnvConfig.(KubernetesConfig)
	// Set default options.
	manifestGlob := config.ManifestGlob
	if manifestGlob == "" {
		manifestGlob = "*.{json,yaml}"
	}

	command := config.Cmd
	if command == "" {
		command = "kubectl"
	}

	// Iterate through and parse templates.
	templates, err := filepath.Glob(manifestGlob)
	if err != nil {
		return err
	}

	results := make(chan bool)

	for _, filename := range templates {
		go func(filename string) {
			data, err := ioutil.ReadFile(filename)
			if err != nil {
				return err
			}

			tmpl, err := template.New("").Parse(string(data))
			if err != nil {
				return err
			}

			cmd := exec.Command(command, "update", "-f", "-")
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

			err = cmd.Wait()
			if err != nil {
				results <- err
			}

			results <- nil

		}(filename)
	}

	for i := 0; i < len(templates); i++ {
		<-results
	}

	return nil
	// TODO: Poll until the release is deployed.
}
