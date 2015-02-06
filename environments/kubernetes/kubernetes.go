package kubernetes

import (
	"text/template"
)

type Kubernetes struct {
}

func (k *Kubernetes) Deploy(release Release, config map) {
	config.manifestDir

	config.map
	// Run manifest templates through parser with release and map
	// Upload the new manifests to kubernetes.
	// Poll until the release is deployed.
}
