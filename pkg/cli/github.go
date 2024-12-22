package cli

import (
	"fmt"
	"os"
)

type GithubCmd struct {
	ListProject *GithubListProjectCmd `arg:"subcommand:list-projects"`
}

type GithubListProjectCmd struct {
	Org  *string `arg:"--org" help:"GitHub organization"`
	User *string `arg:"-u,--user" help:"GitHub username"`
}

func GithubProjectListAction(args Cmd) {
	if args.Github == nil {
		fmt.Println(`error: probablly not a "github list-projects" call?`)
		os.Exit(1)
	}

	if args.Github.ListProject.Org == nil && args.Github.ListProject.User == nil {
		fmt.Println(`error: please provide one of the following flags "--org,--user"`)
		os.Exit(1)
	}

	if args.Github.ListProject.Org != nil && args.Github.ListProject.User != nil {
		fmt.Println(`error: flags "--org,--user" conflicts with each other`)
		os.Exit(1)
	}

	// TODO: logic goes here
}
