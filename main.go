package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/The-Legolas/aggregator_program/internal/config"
	"github.com/The-Legolas/aggregator_program/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	dbQueries := database.New(db)

	programState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	cmds := &commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerListUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerListFeeds)
	cmds.register("follow", handlerFollow)
	cmds.register("following", handlerFollowing)

	if len(os.Args) < 2 {
		fmt.Println("Usage: gator <command> [args...]")
		os.Exit(1)
	}

	// Map the args to the command struct
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	cmd := command{
		name: cmdName,
		args: cmdArgs,
	}
	err = cmds.run(programState, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
