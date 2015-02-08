package deployer

import (
	"github.com/jsdir/deployer/pkg/environments"
	"github.com/jsdir/deployer/pkg/resources"
)

func GetEnvironmentType(name string) resources.EnvironmentType {
	switch name {
	case "kubernetes":
		return new(environments.Kubernetes)
	}
	return nil
}
