package envdir

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// RunCmd run command with flags and env
func RunCmd(cmd []string, env map[string]string) int {
	name := cmd[0]
	args := cmd[1:]

	command := exec.Command(name, args...)
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin

	if env != nil {
		envSlice := make([]string, 0, len(env))
		for k, v := range env {
			envSlice = append(envSlice, fmt.Sprintf("%s=%s", k, v))
		}
		command.Env = envSlice
	}

	if err := command.Run(); err != nil {
		log.Panic(err)
	}

	return 0
}
