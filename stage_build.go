package main

import (
	"go.guoyk.net/deployer/pkg/cmd"
	"go.guoyk.net/deployer/pkg/tempfile"
	"log"
)

func runBuildStage(opts Options) (err error) {
	log.Println("------------------------ 开始构建 ------------------------")
	defer log.Println("------------------------ 结束构建 ------------------------")
	// write temp file
	var buildfile string
	if buildfile, err = tempfile.WriteFile(opts.ScriptBuild, "deployer-build", ".sh", true); err != nil {
		return
	}
	log.Printf("生成构建文件: %s", buildfile)
	// execute the script
	if err = cmd.Run(buildfile); err != nil {
		return
	}
	return
}
