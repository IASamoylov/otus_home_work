package executor

import (
	"fmt"
	"os"
	"os/exec"

	envreader "github.com/IASamoylov/otus_home_work/hw08_envdir_tool/env_reader"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func (ctx *Ctx) RunCmd(args []string, env envreader.Environment) (returnCode int) {
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

// setEnvironment prepares runtime environment.
func setEnvironment(env envreader.Environment) error {
	for k, e := range env {
		if err := os.Unsetenv(k); err != nil {
			return err
		}

		if !e.NeedRemove {
			if err := os.Setenv(k, e.Value); err != nil {
				return err
			}
		}
	}

	return nil
}

// newCommand creates new command for running.
func (ctx *Ctx) newCommand(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)

	cmd.Stdin = ctx.stdIn
	cmd.Stdout = ctx.stdOut
	cmd.Stderr = ctx.stdErr
	cmd.Env = os.Environ()

	return cmd
}
