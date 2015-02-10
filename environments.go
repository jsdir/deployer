package deployer

import (
	"github.com/jsdir/deployer/pkg/environments"
	"github.com/jsdir/deployer-kubernetes"
)

func GetEnvironmentType(name string) resources.EnvironmentType {
	switch name {
	case "kubernetes":
		return new(kubernetes.Kubernetes)
	}
	return nil
}
