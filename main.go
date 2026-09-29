package main

import (
	"fmt"
	"log"

	"github.com/gundamdouble00/blog-aggregator-2/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	fmt.Printf("Read config: %+v\n", cfg)
	err = cfg.SetUser("yuuki")
	if err != nil {
		log.Fatalf("couldn't set current user: %v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	fmt.Printf("%+v\n", cfg)
}
