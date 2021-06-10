package main

import (
	"fmt"
	"os"

	envreader "github.com/IASamoylov/otus_home_work/hw08_envdir_tool/env_reader"
	"github.com/IASamoylov/otus_home_work/hw08_envdir_tool/executor"
)

func main() {
	if len(os.Args) <= 2 {
		fmt.Print("incorrectly configured script")
		os.Exit(1)
	}

	readerCtx := envreader.NewOSContext()
	env, err := readerCtx.ReadDir(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	executorCtx := executor.NewExecutorCtx(os.Stdin, os.Stdout, os.Stderr)
	os.Exit(executorCtx.RunCmd(os.Args[2:], env))
}
