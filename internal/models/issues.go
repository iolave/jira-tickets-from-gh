package models

import "strings"

type Issues struct {
	models *Models
}

func (service *Issues) Upsert(projectId, id, title string, status *IssueStatus, jiraUrl, jiraIssueType, repo *string, estimate *int, assignees *[]string) (*Issue, error) {
	var assigneesStr *string = nil
	if assignees != nil {
		joined := strings.Join(*assignees, ";")
		assigneesStr = &joined
	}
	stmt := `INSERT OR REPLACE INTO issues(
			projectId,
			id,
			jiraUrl,
			jiraIssueType,
			title,
			estimate,
			status,
			assignees,
			repository
		) values(?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := service.models.db.Exec(
		stmt,
		projectId,
		id,
		jiraUrl,
		jiraIssueType,
		title,
		estimate,
		status,
		assigneesStr,
		repo,
	)

	if err != nil {
		return nil, err
	}

	if assignees == nil {
		assignees = &[]string{}
	}

	issue := new(Issue)
	issue.GitHubProjectID = projectId
	issue.GitHubID = id
	issue.Title = title
	issue.JiraIssueType = jiraIssueType
	issue.JiraURL = jiraUrl
	issue.Assignees = *assignees
	issue.Repository = repo
	issue.Estimate = estimate
	issue.Status = status

	return issue, nil
}

type IssueStatus string

const (
	STATUS_TODO IssueStatus = "Todo"
	STATUS_WIP  IssueStatus = "In Progress"
	STATUS_DONE IssueStatus = "Done"
)

type Issue struct {
	GitHubProjectID string
	GitHubID        string
	Title           string
	JiraURL         *string
	JiraIssueType   *string
	Estimate        *int
	Status          *IssueStatus
	Assignees       []string
	Repository      *string
}
