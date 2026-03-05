package main

import (
	"fmt"
	"log"
	"os"

	"github.com/The-Legolas/aggregator_program/internal/config"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	programState := &state{cfg: &cfg}

	cmds := &commands{
		names: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("Usage: gator <command> [args...]")
		os.Exit(1)
	}

	// Step 10: Map the args to your command struct
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
