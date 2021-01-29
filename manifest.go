package main

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	DefaultProfileName = "default"
)

type Manifest struct {
	Build    []string            `yaml:"build"`
	Package  []string            `yaml:"package"`
	CPU      string              `yaml:"cpu"`
	MEM      string              `yaml:"mem"`
	Vars     map[string]string   `yaml:"vars"`
	Profiles map[string]Manifest `yaml:",inline,omitempty"`
}

func loadLegacyManifestLines(filenames ...string) (lines []string) {
	var buf []byte
	var err error
	for _, filename := range filenames {
		buf, err = ioutil.ReadFile(filename)
		if err != nil {
			continue
		}
		blines := bytes.Split(buf, []byte{'\n'})
		for _, bline := range blines {
			lines = append(lines, string(bytes.TrimSpace(bline)))
		}
		log.Printf("加载文件: %s", filepath.Base(filename))
		break
	}
	return
}

func LoadManifest(env string) (mf Manifest, err error) {
	// check version
	if err = checkManifestVersion("deployer.yml"); err != nil {
		return
	}
	// load deployer.yml
	var buf []byte
	if buf, err = ioutil.ReadFile("deployer.yml"); err == nil {
		if err = yaml.Unmarshal(buf, &mf); err != nil {
			return
		}
		mf = mf.Profile(env)
		log.Printf("加载文件: %s", "deployer.yml")
		return
	} else {
		err = nil
	}
	// load docker-build.xxx.sh, Dockerfile.xxx,  docker-build.sh, Dockerfile
	mf.Build = loadLegacyManifestLines("docker-build."+env+".sh", "docker-build.sh")
	mf.Package = loadLegacyManifestLines("Dockerfile."+env, "Dockerfile")
	return
}

func manifestApplyDefault(out *Manifest, s Manifest) {
	if len(out.Build) == 0 {
		out.Build = s.Build
	}
	if len(out.Package) == 0 {
		out.Package = s.Package
	}
	if len(out.CPU) == 0 {
		out.CPU = s.CPU
	}
	if len(out.MEM) == 0 {
		out.MEM = s.MEM
	}
	if out.Vars == nil {
		out.Vars = make(map[string]string)
	}
	for k, v := range s.Vars {
		if out.Vars[k] == "" {
			out.Vars[k] = v
		}
	}
}

func renderLines(lines []string, vars map[string]string) {
	for i, line := range lines {
		for k, v := range vars {
			line = strings.ReplaceAll(line, fmt.Sprintf("{{__%s__}}", k), v)
		}
		lines[i] = line
	}
}

func (s Manifest) Profile(p string) Manifest {
	if len(s.Profiles) == 0 {
		return s
	}
	out := s.Profiles[p]
	dft := s.Profiles[DefaultProfileName]

	manifestApplyDefault(&out, dft)
	manifestApplyDefault(&out, s)

	vars := map[string]string{
		"profile": p,
		"PROFILE": strings.ToUpper(p),
	}
	for k, v := range out.Vars {
		vars[k] = v
		vars[k+"__uppercase"] = strings.ToUpper(v)
		vars[k+"__lowercase"] = strings.ToLower(v)
	}

	renderLines(out.Build, vars)
	renderLines(out.Package, vars)

	return out
}

func (s Manifest) ResourcesCPU() (req, limit int) {
	splits := strings.Split(s.CPU, ":")
	if len(splits) != 2 {
		return
	}
	req, _ = strconv.Atoi(strings.TrimSpace(splits[0]))
	limit, _ = strconv.Atoi(strings.TrimSpace(splits[1]))
	return
}

func (s Manifest) ResourcesMEM() (req, limit int) {
	splits := strings.Split(s.MEM, ":")
	if len(splits) != 2 {
		return
	}
	req, _ = strconv.Atoi(strings.TrimSpace(splits[0]))
	limit, _ = strconv.Atoi(strings.TrimSpace(splits[1]))
	return
}

func (s Manifest) RequestsCPU() string {
	v, _ := s.ResourcesCPU()
	if v == 0 {
		return ""
	}
	return fmt.Sprintf("%dm", v)
}

func (s Manifest) LimitsCPU() string {
	_, v := s.ResourcesCPU()
	if v == 0 {
		return ""
	}
	return fmt.Sprintf("%dm", v)
}

func (s Manifest) RequestsMEM() string {
	v, _ := s.ResourcesMEM()
	if v == 0 {
		return ""
	}
	return fmt.Sprintf("%dMi", v)
}

func (s Manifest) LimitsMEM() string {
	_, v := s.ResourcesMEM()
	if v == 0 {
		return ""
	}
	return fmt.Sprintf("%dMi", v)
}

func (s Manifest) GenerateBuild() []byte {
	buf := &bytes.Buffer{}
	if len(s.Build) == 0 || !strings.HasPrefix(s.Build[0], "#!") {
		buf.WriteString("#!/bin/bash\nset -eux\n")
	}
	for _, l := range s.Build {
		buf.WriteString(l)
		buf.WriteByte('\n')
	}
	return buf.Bytes()
}

func (s Manifest) GenerateDockerfile() []byte {
	buf := &bytes.Buffer{}
	for _, l := range s.Package {
		buf.WriteString(l)
		buf.WriteByte('\n')
	}
	return buf.Bytes()
}
