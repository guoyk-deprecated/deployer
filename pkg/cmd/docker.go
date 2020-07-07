package cmd

func RunDockerVersion() error {
	return Run("docker", "--version")
}

func RunDockerBuild(dockerFile, imageName string) error {
	return Run("docker", "build", "-t", imageName, "-f", dockerFile, ".")
}

func RunDockerTag(imageName string, imageNameAlt string) error {
	return Run("docker", "tag", imageName, imageNameAlt)
}

func RunDockerPush(imageName string, configDir string) error {
	return Run("docker", "--config", configDir, "push", imageName)
}

func RunDockerRemoveImage(imageName string) error {
	return Run("docker", "rmi", imageName)
}
