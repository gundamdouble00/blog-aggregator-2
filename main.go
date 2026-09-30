package main

import (
	"log"
	"os"

	"github.com/gundamdouble00/blog-aggregator-2/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	newState := &state{
		cfg: &cfg,
	}
	newCommands := commands{
		cmdHandlers: make(map[string]func(*state, command) error),
	}
	newCommands.register("login", handlerLogin)
	cmdArgs := os.Args
	if len(cmdArgs) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmd := command{
		Name: cmdArgs[1],
		Args: cmdArgs[2:],
	}
	err = newCommands.run(newState, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
