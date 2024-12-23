package main

import (
	"fmt"
	"log"
	"os"

	arg "github.com/alexflint/go-arg"
	"github.com/iolave/jira-tickets-from-gh/internal/cli"
)

func main() {
	var args cli.Cmd

	p, err := arg.NewParser(arg.Config{Program: cli.NAME}, &args)

	if err != nil {
		log.Fatalf("there was an error in the definition of the Go struct: %v", err)
	}

	err = p.Parse(os.Args[1:])
	switch {
	case err == arg.ErrHelp: // found "--help" on command line
		p.WriteHelp(os.Stdout)
		os.Exit(0)
	case err != nil:
		fmt.Printf("error: %v\n", err)
		p.WriteUsage(os.Stdout)
		os.Exit(1)
	}

	cli.DetectAndRunAction(args)
}
