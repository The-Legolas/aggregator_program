package main

import (
	"fmt"
)

type command struct {
	name string
	args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
	descriptions       map[string]string
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.registeredCommands[cmd.name]
	if !ok {
		return fmt.Errorf("command not found: %s", cmd.name)
	}

	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error, description string) {
	_, ok := c.registeredCommands[name]
	if ok {
		fmt.Printf("command already in registry: %s\n", name)
		return
	}
	c.registeredCommands[name] = f
	c.descriptions[name] = description
}
