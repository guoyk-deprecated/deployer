package main

import (
	"github.com/guoyk93/deployer/pkg/cmd"
	"github.com/guoyk93/deployer/pkg/tempfile"
	"log"
)

func runPackageStage(opts Options) (err error) {
	log.Println("------------------------ 开始打包镜像 ------------------------")
	defer log.Println("------------------------ 结束打包镜像 ------------------------")

	// version
	if err = cmd.RunDockerVersion(); err != nil {
		return
	}
	// write temp file
	var filename string
	if filename, err = tempfile.WriteFile(opts.ScriptPackage, "deployer-package", ".dockerfile", false); err != nil {
		return
	}
	log.Printf("生成打包文件: %s", filename)
	// execute docker build
	if err = cmd.RunDockerBuild(filename, opts.ImageName); err != nil {
		return
	}
	if opts.ImageNameAlt != "" {
		if err = cmd.RunDockerTag(opts.ImageName, opts.ImageNameAlt); err != nil {
			return
		}
	}
	return
}
