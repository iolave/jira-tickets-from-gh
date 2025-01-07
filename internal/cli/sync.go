package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"sync"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/iolave/jira-tickets-from-gh/internal/github"
	"github.com/iolave/jira-tickets-from-gh/internal/models"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// SyncCmdAction syncs GitHub projects with Jira cloud boards.
func SyncCmdAction(args Cmd) {
	level := logrus.InfoLevel
	if args.Debug != nil && *args.Debug == true {
		level = logrus.DebugLevel
	}
	log := newLogger(level)
	log.Debugln("initializing db models")
	m, err := models.Initialize()

	if err != nil {
		log.WithFields(logrus.Fields{"err": err}).Errorln("db models initialization failed")
		exitFromErr(err)
	}

	if args.Sync == nil {
		log.WithFields(logrus.Fields{"err": err}).Errorln("call to sync action args.Sync property is nil")
		exitOnInvalidCall("sync")
	}

	log.WithFields(logrus.Fields{"config": args.Sync.Config}).Debugln("reading config file from location")
	b, err := os.ReadFile(args.Sync.Config)

	if err != nil {
		log.WithFields(logrus.Fields{"config": args.Sync.Config, "err": err}).Errorln("reading config file from location failed")
		exitFromErr(err)
	}

	log.Debugln("parsing config content")
	var config Config
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		log.WithFields(logrus.Fields{"err": err}).Errorln("parsing config content failed")
		exitFromErr(err)
	}

	log.Debugln("validating config properties")
	err = config.validate()
	if err != nil {
		log.WithFields(logrus.Fields{"err": err}).Errorln("config properties validation failed")
		exitFromErr(err)
	}

	if args.GithubToken == nil {
		err := errors.New(`please set the "GITHUB_TOKEN" env variable`)
		log.WithFields(logrus.Fields{"err": err}).Errorln("GithubToken property is nil")
		exitFromErr(err)
	}
	gh := github.New(*args.GithubToken)

	var wg sync.WaitGroup
	for i := 0; i < len(config.Projects); i++ {
		// Increment the wait group counter
		wg.Add(1)
		go func() {
			// Decrement the counter when the go routine completes
			defer wg.Done()
			syncProject(args, config, i, m, gh, log)
		}()
	}
	wg.Wait()
}

func syncProject(args Cmd, config Config, projPos int, m *models.Models, gh *github.GitHubClient, log *logrus.Logger) {
	projectCfg := config.Projects[projPos]
	log.WithFields(logrus.Fields{"project": projectCfg.Name}).Debugln("syncing project")

	// creates new jira client
	log.WithFields(logrus.Fields{"project": projectCfg.Name}).Debugln("creating new jira client")
	url := fmt.Sprintf("https://%s.atlassian.net", projectCfg.Jira.Subdomain)
	jc, err := jira.New(nil, url)
	if err != nil {
		log.WithFields(logrus.Fields{"err": err, "project": projectCfg.Name}).Errorln("failed creating jira client")
		exitFromErr(err)
	}
	log.WithFields(logrus.Fields{"project": projectCfg.Name}).Debugln("getting jira creds")
	token, email, err := getProjectJiraCreds(args, projectCfg.Name)
	if err != nil {
		log.WithFields(logrus.Fields{"err": err, "project": projectCfg.Name}).Errorln("failed to retrieve jira creds")
		exitFromErr(err)
	}
	jc.Auth.SetBasicAuth(email, token)

	// get and set required github project fields into the model
	log.WithFields(logrus.Fields{"project": projectCfg.Name}).Debugln("retrieving github project fields")
	fieldsResult, _, err := gh.GetProjectFields(projectCfg.Github.ProjectID)
	if err != nil {
		log.WithFields(logrus.Fields{"err": err, "project": projectCfg.Name}).Errorln("failed retrieving github project fields")
		exitFromErr(err)
	}
	fieldsIds := struct{ JiraUrl, JiraIssueType, Title, Estimate, Status, Repo, Assignees string }{}
	for _, v := range fieldsResult.Data.Node.Fields.Nodes {
		switch v.Name {
		case models.FIELD_NAME_JIRA_URL:
			fieldsIds.JiraUrl = v.ID
		case models.FIELD_NAME_JIRA_ISSUE_TYPE:
			fieldsIds.JiraIssueType = v.ID
		case models.FIELD_NAME_TITLE:
			fieldsIds.Title = v.ID
		case models.FIELD_NAME_ESTIMATE:
			fieldsIds.Estimate = v.ID
		case models.FIELD_NAME_STATUS:
			fieldsIds.Status = v.ID
		case models.FIELD_NAME_ASSIGNEES:
			fieldsIds.Assignees = v.ID
		case models.FIELD_NAME_REPO:
			fieldsIds.Repo = v.ID
		}
	}
	v := reflect.ValueOf(fieldsIds)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).String() == "" {
			err := fmt.Errorf(`error: missing field in project "%s", make sure it have the following fields [%s, %s, %s, %s, %s, %s, %s]`,
				projectCfg.Name,
				models.FIELD_NAME_JIRA_URL,
				models.FIELD_NAME_JIRA_ISSUE_TYPE,
				models.FIELD_NAME_TITLE,
				models.FIELD_NAME_ESTIMATE,
				models.FIELD_NAME_STATUS,
				models.FIELD_NAME_ASSIGNEES,
				models.FIELD_NAME_REPO,
			)
			log.WithFields(logrus.Fields{"err": err, "project": projectCfg.Name}).Errorln("some fields are not present in github project")
			exitFromErr(err)
		}

	}
	// TODO: maybe is not necesesary to store the project fields id as they can be accessed from variables
	log.WithFields(logrus.Fields{"project": projectCfg.Name, "fields": fieldsIds}).Debugln("upserting project fields ids")
	_, err = m.Projects.Upsert(projectCfg.Github.ProjectID, fieldsIds.JiraUrl, fieldsIds.JiraIssueType, fieldsIds.Title, fieldsIds.Estimate, fieldsIds.Status, fieldsIds.Assignees, fieldsIds.Repo)
	if err != nil {
		log.WithFields(logrus.Fields{"err": err, "project": projectCfg.Name, "fields": fieldsIds}).Errorln("upserting project fields ids failed")
		exitFromErr(err)
	}

	// translates github users to jira account ids
	log.WithFields(logrus.Fields{"project": projectCfg.Name, "assignees": projectCfg.Assignees}).Debugln("translating jira emails to github users")
	assigneesMap := map[string]string{}

	for i := 0; i < len(projectCfg.Assignees); i++ {
		email := projectCfg.Assignees[i].JiraEmail
		// FIXME: response returns a 404 when credentials are invalid, fix this
		users, _, err := jc.User.Search.Do(context.Background(), "", email, 0, 2)

		if err != nil {
			log.WithFields(logrus.Fields{"err": err, "project": projectCfg.Name, "assignee": projectCfg.Assignees[i].JiraEmail}).Errorln("translating jira emails to github user failed")
			exitFromErr(err)
		}

		if len(users) != 1 {
			log.WithFields(logrus.Fields{"project": projectCfg.Name, "assignee": projectCfg.Assignees[i].JiraEmail}).Warnln("found more than one match while translating jira email to github user")
			continue
		}
		assigneesMap[projectCfg.Assignees[i].GHUser] = users[0].AccountID
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

		validNamePattern := "^[a-zA-Z0-9_]*$"
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

		// checking if project id is duped
		for j := 0; j < len(c.Projects); j++ {
			if j == i {
				continue
			}

			if proj.Github.ProjectID == c.Projects[j].Github.ProjectID {
				return fmt.Errorf(`"sync[%d].github.projectId" property is duplicated at "sync[%d].github.projectId"`, i, j)
			}

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
