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

		return fmt.Errorf("error when retrieving user: %v", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %v", err)
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
	user, err := s.db.GetUser(ctxBackGround, name)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("error when retrieving user: %v", err)
		}
	}

	if user.Name == name {
		return fmt.Errorf("existing user")
	}

	newUser := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}
	createdUser, err := s.db.CreateUser(ctxBackGround, newUser)
	if err != nil {
		return fmt.Errorf("could create user: %v", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %v", err)
	}

	log.Printf("created user: %+v", createdUser)
	return nil
}
