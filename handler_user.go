package main

import (
	"context"
	"fmt"
	"time"

	"github.com/The-Legolas/aggregator_program/internal/database"
	"github.com/google/uuid"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}

	ctx := context.Background()
	uuid := uuid.New()
	created_at := time.Now()
	updated_at := time.Now()
	name := cmd.args[0]
	args := database.CreateUserParams{
		ID:        uuid,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
		Name:      name,
	}

	user, err := s.db.CreateUser(ctx, args)
	if err != nil {
		return fmt.Errorf("couldn't create user: %s", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %v", err)
	}

	fmt.Println("User created successfully:")
	printUser(user)
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}
	name := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("couldn't find user: %v", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %v", err)
	}

	fmt.Println("User switched successfully!")
	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}
