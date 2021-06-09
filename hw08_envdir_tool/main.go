package main

import (
	"fmt"

	"github.com/IASamoylov/otus_home_work/hw08_envdir_tool/env_reader"
)

func main() {
	// Place your code here.

	ctx := env_reader.NewOSContext()
	env, err := ctx.ReadDir("testdata/env")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(env)
}
