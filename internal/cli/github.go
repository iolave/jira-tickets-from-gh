package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/iolave/jira-tickets-from-gh/internal/github"
)

type GithubCmd struct {
	ListProject *GithubListProjectCmd `arg:"subcommand:list-projects"`
}

type GithubListProjectCmd struct {
	Org  *string `arg:"--org" help:"GitHub organization"`
	User *string `arg:"-u,--user" help:"GitHub username"`
}

// GithubProjectListAction lists GitHub organization/user projects.
func GithubProjectListAction(args Cmd) {
	if args.Github == nil {
		exitOnInvalidCall("github list-projects")
	}

	if args.GithubToken == nil {
		err := errors.New(`please set the "GITHUB_TOKEN" env variable`)
		exitFromErr(err)
	}

	if args.Github.ListProject.Org == nil && args.Github.ListProject.User == nil {
		exitOnMissingFlags("--org", "--user")
	}

	if args.Github.ListProject.Org != nil && args.Github.ListProject.User != nil {
		exitOnConflictingFlags("--org", "--user")
	}

	gh := github.New(*args.GithubToken)
	if args.Github.ListProject.User != nil {
		result, _, err := gh.ListUserProjects(*args.Github.ListProject.User)
		if err != nil {
			exitFromErr(err)
		}
		b, err := json.Marshal(result.Data.User.Projects.Nodes)
		if err != nil {
			exitFromErr(err)
		}
		fmt.Println(string(b))
		os.Exit(0)
	} else {
		result, _, err := gh.ListOrganizationProjects(*args.Github.ListProject.Org)
		if err != nil {
			exitFromErr(err)
		}
		b, err := json.Marshal(result.Data.Organization.Projects.Nodes)
		if err != nil {
			exitFromErr(err)
		}
		fmt.Println(string(b))
		os.Exit(0)
	}
}
