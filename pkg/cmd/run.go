package cmd

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func Run(name string, args ...string) (err error) {
	log.Printf("执行: %s %s", name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if ee, ok := err.(*exec.ExitError); ok {
		log.Printf("执行完成: %d", ee.ExitCode())
	}
	return
}
