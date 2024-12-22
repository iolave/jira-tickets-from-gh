package cli

import (
	"fmt"
	"os"
)

const (
	VERSION = "v0.5.0-beta.1"
	NAME    = "jira-tickets-from-gh"
)

type Cmd struct {
	Version *bool      `arg:"--version" help:"display the program version"`
	Github  *GithubCmd `arg:"subcommand:github" help:"GitHub utilities"`
}

// PrintVersion prints this program current version.
func PrintVersion() {
	fmt.Println(VERSION)
}

// DetectAndRunAction chooses the proper action to be executed
// based in the given args.
func DetectAndRunAction(args Cmd) {
	if args.Version != nil && *args.Version == true {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	switch {
	case args.Github.ListProject != nil:
		GithubProjectListAction(args)
	}
}
