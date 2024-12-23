package cli

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

	if args.Github.ListProject.Org == nil && args.Github.ListProject.User == nil {
		exitOnMissingFlags("--org", "--user")
	}

	if args.Github.ListProject.Org != nil && args.Github.ListProject.User != nil {
		exitOnConflictingFlags("--org", "--user")
	}

	// TODO: logic goes here
}
