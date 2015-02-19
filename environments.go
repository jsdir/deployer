package deployer

import (
	"github.com/jsdir/deployer-kubernetes"
	"github.com/jsdir/deployer/pkg/resources"
)

func GetEnvironmentType(name string) resources.EnvironmentType {
	switch name {
	case "kubernetes":
		return new(kubernetes.Kubernetes)
	}
	return nil
}
