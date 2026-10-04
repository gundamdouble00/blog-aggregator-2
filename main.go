package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gundamdouble00/blog-aggregator-2/internal/config"
	"github.com/gundamdouble00/blog-aggregator-2/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error connecting to the database: %v", err)
	}

	dbQueries := database.New(db)
	newState := &state{
		cfg: &cfg,
		db:  dbQueries,
	}

	cmdArgs := os.Args
	if len(cmdArgs) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	newCommands := commands{
		cmdHandlers: make(map[string]func(*state, command) error),
	}
	newCommands.register("login", handlerLogin)
	newCommands.register("register", handlerRegister)
	newCommands.register("reset", handlerReset)
	cmd := command{
		Name: cmdArgs[1],
		Args: cmdArgs[2:],
	}
	err = newCommands.run(newState, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
