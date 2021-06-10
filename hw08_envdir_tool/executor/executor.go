package executor

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/IASamoylov/otus_home_work/hw08_envdir_tool/env_reader"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func (ctx *ExecutorCtx) RunCmd(args []string, env env_reader.Environment) (returnCode int) {
	if err := setEnvironment(env); err != nil {
		fmt.Fprintln(ctx.stdErr, err)
		return 1
	}

	cmd := ctx.newCommand(args[0], args[1:]...)

	if err := cmd.Run(); err != nil {
		fmt.Fprintln(ctx.stdErr, err)
	}
	return cmd.ProcessState.ExitCode()
}

// setEnvironment prepares runtime environment
func setEnvironment(env env_reader.Environment) error {
	for k, e := range env {
		os.Unsetenv(k)

		if !e.NeedRemove {
			if err := os.Setenv(k, e.Value); err != nil {
				return err
			}
		}
	}

	return nil
}

// newCommand creates new command for running
func (ctx *ExecutorCtx) newCommand(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)

	cmd.Stdin = ctx.stdIn
	cmd.Stdout = ctx.stdOut
	cmd.Stderr = ctx.stdErr
	cmd.Env = os.Environ()

	return cmd
}
