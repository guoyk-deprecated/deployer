package main

import (
	"github.com/guoyk93/deployer/pkg/cmd"
	"github.com/guoyk93/deployer/pkg/tempfile"
	"log"
)

func runPushStage(opts Options) (err error) {
	log.Println("------------------------ 开始推送镜像 ------------------------")
	defer log.Println("------------------------ 结束推送镜像 ------------------------")
	var dcdir, dcfile string
	if dcdir, dcfile, err = tempfile.WriteDirFile(
		opts.ScriptDockerconfig,
		"deployer-dockerconfig",
		"config.json",
		false,
	); err != nil {
		return
	}
	log.Printf("生成 dockerconfig 文件: %s", dcfile)
	if err = cmd.RunDockerPush(opts.ImageName, dcdir); err != nil {
		return
	}
	if opts.ImageNameAlt != "" {
		if err = cmd.RunDockerPush(opts.ImageNameAlt, dcdir); err != nil {
			return
		}
	}
	if !opts.KeepImage {
		if err = cmd.RunDockerRemoveImage(opts.ImageName); err != nil {
			return
		}
		if opts.ImageNameAlt != "" {
			if err = cmd.RunDockerRemoveImage(opts.ImageNameAlt); err != nil {
				return
			}
		}
	}
	return
}
