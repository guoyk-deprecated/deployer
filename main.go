package main

import (
	"errors"
	"flag"
	"go.guoyk.net/deployer/pkg/tempfile"
	"log"
	"os"
	"path"
	"strings"
)

var (
	opts Options

	optOnlyBuild  bool
	optOnlyDeploy bool
	optSkipDeploy bool

	optKeepGenerated bool
)

func sanitize(strs ...*string) {
	for _, str := range strs {
		val := strings.ReplaceAll(
			strings.TrimSpace(strings.ToLower(*str)),
			"_", "-")
		*str = val
	}
}

func fill(input string, out *string) {
	*out = strings.TrimSpace(*out)
	if len(*out) == 0 {
		*out = input
	}
}

func missing(str *string, err *error, desc string) bool {
	*str = strings.TrimSpace(*str)
	if len(*str) == 0 {
		*err = errors.New(desc)
		return true
	}
	return false
}

func exit(err *error) {
	if *err != nil {
		log.Printf("失败: %s", (*err).Error())
		os.Exit(1)
	} else {
		log.Println("成功")
	}
}

func setup() (err error) {
	log.SetPrefix("[deployer]: ")
	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	var (
		optRegistry      string
		optCluster       string
		optEnv           string
		optBaseImageName string
		optCPU           string
		optMEM           string
		optDeployment    string
	)

	// flags
	flag.StringVar(&optRegistry, "registry", "", "镜像仓库, 指定镜像要推往的仓库地址")
	flag.StringVar(&optCluster, "cluster", "", "集群, 决定使用哪个集群配置文件, 使用 $HOME/.deployer/preset-[CLUSTER].yml")
	flag.StringVar(&opts.Namespace, "namespace", "", "命名空间, 工作负载在集群中的命名空间")
	flag.StringVar(&optDeployment, "deployment", "", "工作负载, 要求和 Kubernetes 上的工作负载名称完全一致（已废弃，使用 -workload）")
	flag.StringVar(&opts.Workload, "workload", "", "工作负载名")
	flag.StringVar(&opts.WorkloadType, "workload-type", "deployment", "工作负载类型，默认为 deployment")
	flag.StringVar(&optEnv, "env", "", "环境名, 决定 deployer.yml 或者 docker-build.XXX.sh,  Dockerfile.XXX 文件的选择, 和构建后的镜像标签")
	flag.StringVar(&optBaseImageName, "image", "", "镜像基础名称, 默认为 '[命名空间]-[工作负载]'")
	flag.StringVar(&opts.Container, "container", "", "容器名称,  默认和 deployment 相同（基于 rancher 习惯）")
	flag.BoolVar(&opts.IsInit, "init", false, "容器为 init 类型的容器")
	flag.BoolVar(&optKeepGenerated, "keep-generated", false, "保存生成的中间文件, 用于调试")
	flag.BoolVar(&opts.KeepImage, "keep-image", false, "在本地 Docker 保留镜像, 便于重复构建, 或者调试")
	flag.StringVar(&opts.LimitsCPU, "limits-cpu", "", "CPU 资源限制, 单位必须为 'm', 千分之一核心（废弃, 使用 --cpu 参数）")
	flag.StringVar(&opts.LimitsMEM, "limits-mem", "", "MEM 资源限制, 单位必须为 'Mi', 兆字节（废弃, 使用 --mem 参数）")
	flag.StringVar(&opts.RequestsCPU, "requests-cpu", "", "CPU 资源请求, 单位必须为 'm', 千分之一核心（废弃, 使用 --cpu 参数）")
	flag.StringVar(&opts.RequestsMEM, "requests-mem", "", "MEM 资源请求, 单位必须为 'Mi', 兆字节（废弃, 使用 --mem 参数）")
	flag.StringVar(&optCPU, "cpu", "", "CPU 资源配置, 格式 '[申请]:[限额]', 如 50:250, 单位 'm' 千分之一核心")
	flag.StringVar(&optMEM, "mem", "", "内存资源配置, 格式 '[申请]:[限额]', 如 64:256, 单位 'Mi' 兆字节")
	flag.BoolVar(&optOnlyBuild, "only-build", false, "只执行 build 步骤, 不进行 package, push, deploy 步骤")
	flag.BoolVar(&optSkipDeploy, "skip-deploy", false, "只执行 build, package, push 步骤, 不进行 deploy 步骤")
	flag.BoolVar(&optOnlyDeploy, "only-deploy", false, "只执行 deploy 步骤, 不进行 build, package, push 步骤")
	flag.Parse()

	// fix deployment
	optDeployment = strings.TrimSpace(optDeployment)
	if len(optDeployment) > 0 {
		opts.Workload = optDeployment
		opts.WorkloadType = "deployment"
	}

	// extract JOB_NAME
	jobNameSplits := strings.Split(os.Getenv("JOB_NAME"), ".")
	if len(jobNameSplits) == 2 {
		fill(jobNameSplits[0], &opts.Workload)
		fill(jobNameSplits[1], &optEnv)
		log.Println("从 Jenkins $JOB_NAME 获取到 工作负载 " + opts.Workload + ", 环境 " + optEnv)
	} else if len(jobNameSplits) == 4 {
		fill(jobNameSplits[0], &optCluster)
		fill(jobNameSplits[1], &opts.Namespace)
		fill(jobNameSplits[2], &opts.Workload)
		fill(jobNameSplits[3], &optEnv)
		log.Println("从 Jenkins $JOB_NAME 获取到 集群 " + optCluster + ", 命名空间 " + opts.Namespace + ", 工作负载 " + opts.Workload + ", 环境 " + optEnv)
	}

	// check target options
	sanitize(&optCluster, &opts.Namespace, &opts.Workload, &optEnv)
	if missing(&optCluster, &err, "错误: 集群未指定, 使用 --cluster 指定集群") {
		return
	}
	if missing(&opts.Namespace, &err, "错误: 命名空间未指定, 使用 --namespace 指定命名空间") {
		return
	}
	if missing(&opts.Workload, &err, "错误: 工作负载未指定, 使用 --workload, 或者 $JOB_NAME 指定工作负载") {
		return
	}
	if missing(&optEnv, &err, "错误: 环境未指定, 使用 --env 或者 $JOB_NAME 指定环境") {
		return
	}
	log.Printf("部署目标: %s -> %s -> %s (ENV: %s)", optCluster, opts.Namespace, opts.Workload, optEnv)
	log.Println("------------------------")

	// container
	fill(opts.Workload, &opts.Container)
	// base image name
	fill(opts.Namespace+"-"+opts.Workload, &optBaseImageName)

	// preset
	var preset Preset
	if preset, err = LoadPreset(optCluster); err != nil {
		return
	}
	opts.ImagePullSecrets = preset.ImagePullSecrets
	opts.ScriptKubeconfig = preset.GenerateKubeconfig()
	opts.ScriptDockerconfig = preset.GenerateDockerconfig()
	fill(preset.Registry, &optRegistry)

	// check registry configuration
	if missing(&optRegistry, &err, "错误: 镜像仓库未指定, 使用 --registry 指定镜像仓库") {
		return
	}

	// generate image name
	buildNumber := os.Getenv("BUILD_NUMBER")
	if buildNumber == "" {
		opts.ImageName = path.Join(optRegistry, optBaseImageName+":"+optEnv)
	} else {
		opts.ImageName = path.Join(optRegistry, optBaseImageName+":"+optEnv+"-build-"+buildNumber)
	}
	log.Printf("镜像完整名称: %s", opts.ImageName)

	// load manifest
	var mf Manifest
	if mf, err = LoadManifest(optEnv); err != nil {
		return
	}
	opts.ScriptBuild = mf.GenerateBuild()
	opts.ScriptPackage = mf.GenerateDockerfile()

	// manifest resources
	if len(optCPU) > 0 {
		mf.CPU = optCPU
	}
	if len(optMEM) > 0 {
		mf.MEM = optMEM
	}
	fill(mf.RequestsCPU(), &opts.RequestsCPU)
	fill(mf.LimitsCPU(), &opts.LimitsCPU)
	fill(mf.RequestsMEM(), &opts.RequestsMEM)
	fill(mf.LimitsMEM(), &opts.LimitsMEM)

	// preset resources
	fill(preset.RequestsCPU, &opts.RequestsCPU)
	fill(preset.LimitsCPU, &opts.LimitsCPU)
	fill(preset.RequestsMEM, &opts.RequestsMEM)
	fill(preset.LimitsMEM, &opts.LimitsMEM)

	return
}

func main() {
	var err error
	defer exit(&err)

	// setup
	if err = setup(); err != nil {
		return
	}

	if !optKeepGenerated {
		defer tempfile.DeleteAll()
	}

	if !optOnlyDeploy {
		if err = runBuildStage(opts); err != nil {
			return
		}

		if optOnlyBuild {
			return
		}

		if err = runPackageStage(opts); err != nil {
			return
		}

		if err = runPushStage(opts); err != nil {
			return
		}

		if optSkipDeploy {
			return
		}
	}

	if err = runDeployStage(opts); err != nil {
		return
	}
}
