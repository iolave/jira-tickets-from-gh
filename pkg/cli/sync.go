package cli

type SyncCmd struct {
	Config string `arg:"--config,-c" help:"path to config file" placeholder:"<PATH>"`
	// TODO: maybe remove the following(?)
	// GhProjectID       string             `arg:"--gh-project-id" help:"GitHub project id" placeholder:"<ID>"`
	// GhAssigneesMap    *map[string]string `arg:"--gh-assignees-map" help:"map of GitHub users to Jira ones (email)" placeholder:"<GH_USER:JIRA_USER,...>"`
	// JiraSubdomain     string             `arg:"--jira-subdomain" help:"Jira subdomain" placeholder:"<STRING>"`
	// JiraProjectKey    string             `arg:"--jira-project-key" help:"Jira project key" placeholder:"<STRING>"`
	// JiraIssuePrefix   *string            `arg:"--jira-issue-prefix" help:"prefix to be added to jira issue title" placeholder:"<STRING>"`
	// JiraEstimateField *string            `arg:"--jira-estimate-field" help:"Jira field name that holds estiamte value" placeholder:"<STRING>"`
	// SleepTime         int                `arg:"--sleep-time" help:"sleep time between executions (if not specified the program will run once)" placeholder:"<MS>"`
	//transitionsToWip?: number[],
	//transitionsToDone?: number[],
}

// SyncCmdAction syncs GitHub projects with Jira cloud boards.
func SyncCmdAction(args Cmd) {
	if args.Sync == nil {
		exitOnInvalidCall("sync")
	}
}
