package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/The-Legolas/aggregator_program/internal/config"
	"github.com/The-Legolas/aggregator_program/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

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
		descriptions:       make(map[string]string),
	}

	cmds.register("login", handlerLogin, "Log in as a user: <username>")
	cmds.register("register", handlerRegister, "Register a new user: <username>")
	cmds.register("reset", handlerReset, "Deletes all users and data from the database")
	cmds.register("users", handlerListUsers, "List all users currently in the Database")
	cmds.register("agg", handlerAgg, "Start aggregating feeds on an interval: <time_between_reqs> (e.g. 1m, 30s)")
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed), "Subscribe to a new feed: <name> <url>")
	cmds.register("feeds", handlerListFeeds, "List all feeds in the database")
	cmds.register("follow", middlewareLoggedIn(handlerFollow), "Lets the current user follow a specific feed: <feed_url>")
	cmds.register("following", middlewareLoggedIn(handlerFollowing), "List all feeds the current user follows")
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow), "Unfollow a feed: <feed_url>")
	cmds.register("browse", middlewareLoggedIn(handlerBrowse), "Browse posts: [limit]")
	cmds.register("help", func(s *state, cmd command) error {
		return handlerHelp(s, cmd, cmds)
	}, "List all available commands and their usage")

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

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("Error trying to fetch user: %w", err)
		}
		return handler(s, cmd, user)
	}
}
