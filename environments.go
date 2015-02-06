package deployer

import (
	"github.com/jsdir/deployer/pkg/environments"
)

var EnvironmentTypes = map[string]EnvironmentType{
	"kubernetes": environments.Kubernetes,
}
