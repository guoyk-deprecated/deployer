package cmd

import (
	"log"
	"time"
)

const KubectlPatchRetries = 3

func RunKubectlVersion(kubeconfig string) error {
	return Run("kubectl", "--kubeconfig", kubeconfig, "version")
}

func RunKubectlPatch(kubeconfig, namespace, workload, workloadType, patch string) (err error) {
	i := 0
	for {
		if err = Run("kubectl",
			"--kubeconfig", kubeconfig,
			"--namespace", namespace,
			"patch", workloadType+"s/"+workload,
			"-p", patch,
		); err == nil {
			return
		}

		i++
		if i > KubectlPatchRetries {
			return
		}

		log.Printf("部署失败，5s 后重试")
		time.Sleep(time.Second * 5)
	}
}
