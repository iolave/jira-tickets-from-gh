package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	VERSION = "v0.4.0-beta.4"
	NAME    = "jira-tickets-from-gh"
)

type Cmd struct {
	Version     *bool      `arg:"--version" help:"display the program version"`
	GithubToken *string    `arg:"env:GITHUB_TOKEN,--gh-token" help:"GitHub token" placeholder:"<STRING>"`
	JiraEmail   *string    `arg:"env:JIRA_EMAIL,--jira-email" help:"Jira email used for basic auth" placeholder:"<STRING>"`
	Debug       *bool      `arg:"--debug" help:"enables debug mode"`
	JiraToken   *string    `arg:"env:JIRA_TOKEN,--jira-token" help:"Jira api token used for basic auth" placeholder:"<STRING>"`
	Github      *GithubCmd `arg:"subcommand:github" help:"GitHub utilities" `
	Sync        *SyncCmd   `arg:"subcommand:sync" help:"sync GitHub project tickets with Jira"`
}

func newLogger(level logrus.Level) *logrus.Logger {
	log := logrus.New()
	log.SetLevel(level)
	//log.SetReportCaller(true)
	log.SetFormatter(&logrus.JSONFormatter{})
	return log
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
	case args.Github != nil:
		switch {
		case args.Github.ListProject != nil:
			GithubProjectListAction(args)
		}
	case args.Sync != nil:
		SyncCmdAction(args)
	}
}

func exitOnInvalidCall(cmd string) {
	msg := fmt.Sprintf(`error: probablly not a "%s" cmd call?`, cmd)

	fmt.Println(msg)
	os.Exit(1)
}

func exitOnMissingFlags(flags ...string) {
	joinedFlags := strings.Join(flags, ",")
	msg := fmt.Sprintf(`error: please provide one of the following flags "%s"`, joinedFlags)

	fmt.Println(msg)
	os.Exit(1)
}

func exitOnConflictingFlags(flags ...string) {
	joinedFlags := strings.Join(flags, ",")
	msg := fmt.Sprintf(`error: flags "%s" conflicts with each other`, joinedFlags)

	fmt.Println(msg)
	os.Exit(1)
}

func exitFromErr(err error) {
	msg := fmt.Sprintf("error: %s", err.Error())
	fmt.Println(msg)
	os.Exit(1)
}
