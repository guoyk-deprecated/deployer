package main

import (
	"encoding/json"
	"go.guoyk.net/deployer/pkg/cmd"
	"go.guoyk.net/deployer/pkg/tempfile"
	"log"
	"strings"
	"time"
)

type Patch struct {
	Metadata struct {
		Annotations map[string]string `json:"annotations,omitempty"`
	} `json:"metadata,omitempty"`
	Spec struct {
		Template struct {
			Metadata struct {
				Annotations struct {
					Timestamp string `json:"net.guoyk.deployer/timestamp,omitempty"`
				} `json:"annotations"`
			} `json:"metadata"`
			Spec struct {
				Containers       []PatchContainer       `json:"containers,omitempty"`
				InitContainers   []PatchInitContainer   `json:"initContainers,omitempty"`
				ImagePullSecrets []PatchImagePullSecret `json:"imagePullSecrets,omitempty"`
			} `json:"spec"`
		} `json:"template"`
	} `json:"spec"`
}

type PatchContainer struct {
	Image           string `json:"image"`
	Name            string `json:"name"`
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`
	Resources       struct {
		Limits struct {
			CPU    string `json:"cpu,omitempty"`
			Memory string `json:"memory,omitempty"`
		} `json:"limits,omitempty"`
		Requests struct {
			CPU    string `json:"cpu,omitempty"`
			Memory string `json:"memory,omitempty"`
		} `json:"requests,omitempty"`
	} `json:"resources,omitempty"`
}

type PatchInitContainer struct {
	Image           string `json:"image"`
	Name            string `json:"name"`
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`
}

type PatchImagePullSecret struct {
	Name string `json:"name"`
}

func runDeployStage(opts Options) (err error) {
	log.Println("------------------------ 开始部署 ------------------------")
	defer log.Println("------------------------ 结束部署 ------------------------")

	// build Patch struct
	var p Patch
	p.Metadata.Annotations = opts.ExtraAnnotations
	p.Spec.Template.Metadata.Annotations.Timestamp = time.Now().Format(time.RFC3339)
	for _, name := range opts.ImagePullSecrets {
		secret := PatchImagePullSecret{Name: strings.TrimSpace(name)}
		p.Spec.Template.Spec.ImagePullSecrets = append(p.Spec.Template.Spec.ImagePullSecrets, secret)
	}
	if opts.IsInit {
		container := PatchInitContainer{
			Image:           opts.ImageName,
			Name:            opts.Container,
			ImagePullPolicy: "Always",
		}
		p.Spec.Template.Spec.InitContainers = append(p.Spec.Template.Spec.InitContainers, container)
	} else {
		container := PatchContainer{
			Image:           opts.ImageName,
			Name:            opts.Container,
			ImagePullPolicy: "Always",
		}
		container.Resources.Requests.CPU = opts.RequestsCPU
		container.Resources.Requests.Memory = opts.RequestsMEM
		container.Resources.Limits.CPU = opts.LimitsCPU
		container.Resources.Limits.Memory = opts.LimitsMEM
		p.Spec.Template.Spec.Containers = append(p.Spec.Template.Spec.Containers, container)
	}
	// marshal
	var buf []byte
	if buf, err = json.Marshal(&p); err != nil {
		return
	}
	// generate kubeconfig
	var kubefile string
	if kubefile, err = tempfile.WriteFile(opts.ScriptKubeconfig, "deployer-kubeconfig", ".yaml", false); err != nil {
		return
	}
	log.Printf("生成 Kubeconfig: %s", kubefile)
	// version
	if err = cmd.RunKubectlVersion(kubefile); err != nil {
		return
	}
	// kubectl patch
	if err = cmd.RunKubectlPatch(kubefile, opts.Namespace, opts.Workload, opts.WorkloadType, string(buf)); err != nil {
		return
	}
	return
}
