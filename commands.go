package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/scottw0173/aggregator/internal/config"
	"github.com/scottw0173/aggregator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type handlerFunc func(*state, command) error

type commands struct {
	handlers map[string]handlerFunc
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("login requires exactly one argument")
	}

	username := cmd.args[0]
	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("user: %s is not registered with db", username)
	}

	if err := s.cfg.SetUser(username); err != nil {
		return err
	}
	fmt.Printf("username successfully set to %s\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("missing name")
	}

	user := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}
	newUser, err := s.db.CreateUser(context.Background(), user)
	if err != nil {
		return err
	}

	s.cfg.SetUser(cmd.args[0])

	fmt.Printf("user successfully created: %s\n", cmd.args[0])
	fmt.Println(newUser)
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.handlers[cmd.name]
	if !exists {
		return fmt.Errorf("unregistered command")
	}
	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) error {
	if _, exists := c.handlers[name]; exists {
		return fmt.Errorf("duplicate command")
	}
	c.handlers[name] = f
	return nil
}
