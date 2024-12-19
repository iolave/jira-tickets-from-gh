package businesslogic

import (
	jira "github.com/ctreminiom/go-atlassian/jira/v3"
)

type SyncOptions struct {
	JiraHost  string
	JiraEmail string
	JiraToken string
}

func Sync(args SyncOptions) error {
	jc, err := jira.New(nil, args.JiraHost)

	if err != nil {
		return err
	}

	jc.Auth.SetBasicAuth(args.JiraEmail, args.JiraToken)
	return nil
}
