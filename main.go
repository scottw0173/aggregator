package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/scottw0173/aggregator/internal/config"
	"github.com/scottw0173/aggregator/internal/database"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("error creating config")
	}
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		fmt.Println(err)
	}
	dbQueries := database.New(db)

	newState := state{
		db:  dbQueries,
		cfg: &cfg,
	}

	cmds := commands{
		handlers: make(map[string]handlerFunc),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)

	args := os.Args

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
