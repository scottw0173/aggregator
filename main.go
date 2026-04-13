package main

import (
	"fmt"
	"os"

	"github.com/scottw0173/aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("error creating config")
	}

	newState := state{
		currentCfg: &cfg,
	}

	cmds := commands{
		handlers: make(map[string]handlerFunc),
	}
	cmds.register("login", handlerLogin)

	args := os.Args
	//for _, arg := range args {
	//	fmt.Println(arg)
	//}

	if len(args) < 2 {
		fmt.Println("too few arguments")
		os.Exit(1)
	}
	newCommand := command{
		name: args[1],
		args: args[2:],
	}
	if err := cmds.run(&newState, newCommand); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
