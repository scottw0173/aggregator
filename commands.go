package main

import (
	"fmt"

	"github.com/scottw0173/aggregator/internal/config"
)

type state struct {
	currentCfg *config.Config
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

	if err := s.currentCfg.SetUser(username); err != nil {
		return err
	}
	fmt.Printf("username successfully set to %s\n", username)
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
