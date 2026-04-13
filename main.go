package main

import (
	"fmt"

	"github.com/scottw0173/aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("error creating config")
	}
	cfg.SetUser("KingOfAnts")

	newCfg, err := config.Read()
	if err != nil {
		fmt.Println("error reading newCfg")
	}

	fmt.Println(newCfg)
}
