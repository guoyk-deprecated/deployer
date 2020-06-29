package cmd

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
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

func RunRetries(retry int, name string, args ...string) (err error) {
	if retry < 1 {
		retry = 1
	}
	for {
		if err = Run(name, args...); err == nil {
			return
		}

		retry--
		if retry == 0 {
			return
		}
		time.Sleep(time.Second * 5)
		log.Printf("5s 后重试, 剩余 %d", retry)
	}
}
