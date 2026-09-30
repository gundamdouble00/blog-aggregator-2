package main

import "fmt"

func handlerLogin(s *state, cmd command) error {
	numArgs := len(cmd.Args)
	if numArgs != 1 {
		return fmt.Errorf("\"login\" command should have one argument (currently: %v)", numArgs)
	}

	user := cmd.Args[0]
	err := s.cfg.SetUser(user)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %v", err)
	}

	fmt.Println("User switch successfully!")
	return nil
}
