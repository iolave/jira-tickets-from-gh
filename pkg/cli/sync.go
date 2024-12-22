package cli

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// SyncCmdAction syncs GitHub projects with Jira cloud boards.
func SyncCmdAction(args Cmd) {
	if args.Sync == nil {
		exitOnInvalidCall("sync")
	}

	b, err := os.ReadFile(args.Sync.Config)

	if err != nil {
		exitFromErr(err)
	}

	var config Config

	err = yaml.Unmarshal(b, &config)
	if err != nil {
		exitFromErr(err)
	}

	err = config.validate()
	if err != nil {
		exitFromErr(err)
	}

	// TODO: sync here
	fmt.Println(config)
}

type SyncCmd struct {
	Config string `arg:"required,--config,-c" help:"path to config file" placeholder:"<PATH>"`
}

type Config struct {
	SleepTime *string `yaml:"sleepTime"`
	EnableAPI *bool   `yaml:"enableApi"`
	Projects  []struct {
		Name      string `yaml:"name"`
		Assignees []struct {
			JiraEmail string `yaml:"jiraEmail"`
			GHUser    string `yaml:"ghUser"`
		} `yaml:"assignees"`
		Github struct {
			ProjectID string `yaml:"projectId"`
		}
		Jira struct {
			Subdomain     string  `yaml:"subdomain"`
			ProjectKey    string  `yaml:"projectKey"`
			EstimateField *string `yaml:"estimateField"`
			IssuePrefix   *string `yaml:"issuePrefix"`
			Issues        []struct {
				Type              string `yaml:"type"`
				TransitionsToWIP  []int  `yaml:"transitionsToWip"`
				TransitionsToDone []int  `yaml:"transitionsToDone"`
			} `yaml:"issues"`
		} `yaml:"jira"`
	} `yaml:"sync"`
}

func (c Config) validate() error {
	for i := 0; i < len(c.Projects); i++ {
		proj := c.Projects[i]
		if proj.Name == "" {
			return fmt.Errorf(`"sync[%d].name" property is missing`, i)
		}

		for j := 0; j < len(proj.Assignees); j++ {
			assignee := proj.Assignees[j]

			if assignee.GHUser == "" {
				return fmt.Errorf(`"sync[%d].assignees[%d].ghUser" property is missing`, i, j)
			}
			if assignee.JiraEmail == "" {
				return fmt.Errorf(`"sync[%d].assignees[%d].jiraEmail" property is missing`, i, j)
			}

		}

		if proj.Github.ProjectID == "" {
			return fmt.Errorf(`"sync[%d].github.projectId" property is missing`, i)
		}
		if proj.Jira.Subdomain == "" {
			return fmt.Errorf(`"sync[%d].jira.subdomain" property is missing`, i)
		}
		if proj.Jira.ProjectKey == "" {
			return fmt.Errorf(`"sync[%d].jira.projectKey" property is missing`, i)
		}

		for j := 0; j < len(proj.Jira.Issues); j++ {
			issue := proj.Jira.Issues[j]

			if issue.Type == "" {
				return fmt.Errorf(`"sync[%d].jira.issues[%d].type" property is missing`, i, j)
			}
		}

	}

	return nil
}
