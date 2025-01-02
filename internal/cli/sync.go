package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sync"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
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

	var wg sync.WaitGroup
	for i := 0; i < len(config.Projects); i++ {
		// Increment the wait group counter
		wg.Add(1)
		go func() {
			// Decrement the counter when the go routine completes
			defer wg.Done()
			syncProject(args, config, i)
		}()
	}
	wg.Wait()
}

func syncProject(args Cmd, config Config, projPos int) {
	if args.JiraEmail == nil {
		err := errors.New(`please set the "JIRA_EMAIL" env variable`)
		exitFromErr(err)
	}

	if args.JiraToken == nil {
		err := errors.New(`please set the "JIRA_TOKEN" env variable`)
		exitFromErr(err)
	}

	project := config.Projects[projPos]

	url := fmt.Sprintf("https://%s.atlassian.net", project.Jira.Subdomain)
	jc, err := jira.New(nil, url)
	if err != nil {
		exitFromErr(err)
	}
	token, email, err := getProjectJiraCreds(args, project.Name)
	if err != nil {
		exitFromErr(err)
	}
	jc.Auth.SetBasicAuth(email, token)

	// map to translate github users to jira account ids
	assigneesMap := map[string]string{}

	for i := 0; i < len(project.Assignees); i++ {
		email := project.Assignees[i].JiraEmail
		// FIXME: response returns a 404 when credentials are invalid, fix this
		users, _, err := jc.User.Search.Do(context.Background(), "", email, 0, 2)

		if err != nil {
			exitFromErr(err)
		}

		if len(users) != 1 {
			continue
		}

		assigneesMap[project.Assignees[i].GHUser] = users[0].AccountID
	}
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

		validNamePattern := "^([a-zA-Z0-9_])*$"
		validName, _ := regexp.MatchString(validNamePattern, proj.Name)

		if !validName {
			return fmt.Errorf(`"sync[%d].name" property should match the expression "%s"`, i, validNamePattern)
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

func getProjectJiraCreds(args Cmd, projectName string) (token string, email string, err error) {
	envEmail := fmt.Sprintf("JIRA_EMAIL_%s", projectName)
	envToken := fmt.Sprintf("JIRA_TOKEN_%s", projectName)

	email = os.Getenv(envEmail)
	token = os.Getenv(envToken)

	if email != "" && token != "" {
		return token, email, nil
	}

	if args.JiraEmail == nil {
		err := errors.New(`please set the "JIRA_EMAIL" env variable`)
		return "", "", err
	}

	if args.JiraToken == nil {
		err := errors.New(`please set the "JIRA_TOKEN" env variable`)
		return "", "", err
	}

	return *args.JiraToken, *args.JiraEmail, nil
}
