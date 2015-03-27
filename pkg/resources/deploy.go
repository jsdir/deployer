package resources

type Deploy struct {
	Env             *Environment
	LastRelease     *Release
	Release         *Release
	ChangedServices []string
	EnvConfig       interface{}
}

type DeployRequest struct {
	Src  string `form:"src"`
	Dest string `form:"dest"`
}
