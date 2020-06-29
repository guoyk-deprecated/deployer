package cmd

const KubectlPatchRetries = 3

func RunKubectlVersion(kubeconfig string) error {
	return RunRetries(KubectlPatchRetries, "kubectl", "--kubeconfig", kubeconfig,
		"version")
}

func RunKubectlPatch(kubeconfig, namespace, workload, workloadType, patch string) error {
	return RunRetries(KubectlPatchRetries, "kubectl", "--kubeconfig", kubeconfig,
		"--namespace", namespace, "patch", workloadType+"s/"+workload, "-p", patch)
}
