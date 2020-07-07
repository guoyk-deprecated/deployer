package main

type Options struct {
	ScriptBuild        []byte
	ScriptPackage      []byte
	ScriptKubeconfig   []byte
	ScriptDockerconfig []byte

	ImageName    string
	ImageNameAlt string
	KeepImage    bool

	Namespace        string
	Workload         string
	WorkloadType     string
	Container        string
	IsInit           bool
	ImagePullSecrets []string
	RequestsCPU      string
	RequestsMEM      string
	LimitsCPU        string
	LimitsMEM        string
}
