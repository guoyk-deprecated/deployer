package main

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type manifestVersionCheck struct {
	Version int `yaml:"version"`
}

func checkManifestVersion(filename string) (err error) {
	var buf []byte
	if buf, err = ioutil.ReadFile(filename); err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return
	}
	var m manifestVersionCheck
	if err = yaml.Unmarshal(buf, &m); err != nil {
		return
	}
	if m.Version != 0 {
		err = errors.New("deployer 不兼容 version: 2 的 deployer.yml 文件，请使用 deployer2")
		return
	}
	return
}
