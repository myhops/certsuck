package main

import (
	"context"
	"log"
	"os"

	"github.com/myhops/certsuck/command"
)

func run(args []string) error {
	cmd := command.NewCommand()
	return cmd.Run(context.Background(), args)
}

func main() {
	if err := run(os.Args); err != nil {
		log.Fatalf("run failed: %s", err.Error())
	}
}
