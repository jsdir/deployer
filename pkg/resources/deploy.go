package resources

import (
	"github.com/mholt/binding"
)

type Deploy struct {
	Env             *Environment
	LastRelease     *Release
	Release         *Release
	ChangedServices []string
	EnvConfig       interface{}
}

type DeployRequest struct {
	Src  string
	Dest string
}

func (d *DeployRequest) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&d.Src: binding.Field{
			Form:     "src",
			Required: true,
		},
		&d.Dest: binding.Field{
			Form:     "dest",
			Required: true,
		},
	}
}
