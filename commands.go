package main

import (
	"fmt"

	"github.com/The-Legolas/aggregator_program/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	names map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	arg := cmd.args
	if len(arg) == 0 {
		return fmt.Errorf("no arguments")
	}

	err := s.cfg.SetUser(arg[0])
	if err != nil {
		return fmt.Errorf("error setting user: %v", err)
	}

	fmt.Printf("The user has been succesfully set to: %v\n", arg[0])

	return nil
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.names[cmd.name]
	if !ok {
		return fmt.Errorf("command not found: %s", cmd.name)
	}

	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	_, ok := c.names[name]
	if ok {
		fmt.Printf("command already in registry: %s\n", name)
		return
	}

	c.names[name] = f
}
