package main

import "errors"

type command struct {
	Name string
	Args []string
}

type commands struct {
	cmdHandlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	cmdHandler, ok := c.cmdHandlers[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}

	return cmdHandler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmdHandlers[name] = f
}
