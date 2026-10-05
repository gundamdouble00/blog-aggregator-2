package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gundamdouble00/blog-aggregator-2/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	numArgs := len(cmd.Args)
	if numArgs != 1 {
		return fmt.Errorf("\"login\" command should have one argument (currently: %v)", numArgs)
	}

	name := cmd.Args[0]
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("invalid user: %v", name)
		}

		return fmt.Errorf("error when retrieving user: %w", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User switch successfully!")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	numArgs := len(cmd.Args)
	if numArgs != 1 {
		return fmt.Errorf("\"register\" command should have one argument (currently: %v)", numArgs)
	}

	name := cmd.Args[0]
	ctxBackGround := context.Background()
	_, err := s.db.GetUser(ctxBackGround, name)
	if err == nil {
		return fmt.Errorf("user already exists")
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("error when retrieving user: %w", err)
	}

	newUser := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}
	createdUser, err := s.db.CreateUser(ctxBackGround, newUser)
	if err != nil {
		return fmt.Errorf("could create user: %w", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User created successfullly!")
	log.Printf("created user: %+v", createdUser)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("\"reset\" command doesn't have any arguments")
	}

	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error when deleting all user: %w", err)
	}

	fmt.Println("Database reset successfully!")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error retrieving users: %w", err)
	}

	for _, user := range users {
		fmt.Printf("* %s", user.Name)
		if user.Name == s.cfg.CurrentUserName {
			fmt.Print(" (current)")
		}
		fmt.Println()
	}
	return nil
}
