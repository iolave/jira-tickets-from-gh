package models

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

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

// Get gets a project issue, if no issue found *Issue will be nil
func (p *Issues) Get(githubProjectId, githubId string) (*Issue, error) {
	if githubProjectId == "" {
		return nil, errors.New(`please provide a value for "githubProjectId"`)
	}
	if githubId == "" {
		return nil, errors.New(`please provide a value for "githubId"`)
	}

	stmt := fmt.Sprintf(`SELECT
		projectId,
		id,
		jiraUrl,
		jiraIssueType,
		title,
		estimate,
		status,
		assignees,
		repository
	FROM issues
	WHERE id = "%s"
	AND projectId = "%s"
	`, githubId, githubProjectId)
	rows, err := p.models.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	issue := new(Issue)
	var assigneesStr *string
	if !rows.Next() {
		return nil, nil
	}
	err = rows.Scan(
		&issue.GitHubProjectID,
		&issue.GitHubID,
		&issue.JiraURL,
		&issue.JiraIssueType,
		&issue.Title,
		&issue.Estimate,
		&issue.Status,
		&assigneesStr,
		&issue.Repository,
	)
	if err != nil {
		return nil, err
	}
	assignees := []string{}
	if assigneesStr != nil {
		assignees = strings.Split(*assigneesStr, ";")
	}
	issue.Assignees = assignees
	// uncomment when models is avaialbe within an issue
	// issue.models = p.models
	return issue, nil
}

// Get gets a project issue, if no issue found *Issue will be nil
func (p *Issues) GetWithoutUrl(githubProjectId string) ([]*Issue, error) {
	if githubProjectId == "" {
		return nil, errors.New(`please provide a value for "githubProjectId"`)
	}

	stmt := fmt.Sprintf(`SELECT
		projectId,
		id,
		jiraUrl,
		jiraIssueType,
		title,
		estimate,
		status,
		assignees,
		repository
	FROM issues
	WHERE projectId = "%s"
	AND jiraUrl IS NULL
	`, githubProjectId)
	rows, err := p.models.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	issues := []*Issue{}
	for rows.Next() {
		issue := new(Issue)
		var assigneesStr *string

		err = rows.Scan(
			&issue.GitHubProjectID,
			&issue.GitHubID,
			&issue.JiraURL,
			&issue.JiraIssueType,
			&issue.Title,
			&issue.Estimate,
			&issue.Status,
			&assigneesStr,
			&issue.Repository,
		)
		if err != nil {
			return nil, err
		}
		assignees := []string{}
		if assigneesStr != nil && *assigneesStr != "" {
			assignees = strings.Split(*assigneesStr, ";")
		}
		issue.Assignees = assignees
		// uncomment when models is avaialbe within an issue
		// issue.models = p.models
		issues = append(issues, issue)
	}

	return issues, nil
}

// Get gets a project issue, if no issue found *Issue will be nil
func (p *Issues) GetWithUrl(githubProjectId string) ([]*Issue, error) {
	if githubProjectId == "" {
		return nil, errors.New(`please provide a value for "githubProjectId"`)
	}

	stmt := fmt.Sprintf(`SELECT
		projectId,
		id,
		jiraUrl,
		jiraIssueType,
		title,
		estimate,
		status,
		assignees,
		repository
	FROM issues
	WHERE projectId = "%s"
	AND jiraUrl IS NOT NULL
	`, githubProjectId)
	rows, err := p.models.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	issues := []*Issue{}
	for rows.Next() {
		issue := new(Issue)
		var assigneesStr *string

		err = rows.Scan(
			&issue.GitHubProjectID,
			&issue.GitHubID,
			&issue.JiraURL,
			&issue.JiraIssueType,
			&issue.Title,
			&issue.Estimate,
			&issue.Status,
			&assigneesStr,
			&issue.Repository,
		)
		if err != nil {
			return nil, err
		}
		assignees := []string{}
		if assigneesStr != nil && *assigneesStr != "" {
			assignees = strings.Split(*assigneesStr, ";")
		}
		issue.Assignees = assignees
		// uncomment when models is avaialbe within an issue
		// issue.models = p.models
		issues = append(issues, issue)
	}

	return issues, nil
}

// FindThoseThatExist takes a list of issues ids and returns those that does exists.
func (p *Issues) FindThoseThatExist(githubProjectId string, ids []string) ([]string, error) {
	newIds := slices.Clone(ids)
	if githubProjectId == "" {
		return []string{}, errors.New(`please provide a value for "githubProjectId"`)
	}

	for i := 0; i < len(newIds); i++ {
		newIds[i] = fmt.Sprintf(`"%s"`, newIds[i])
	}
	idsQuery := strings.Join(newIds, ",")
	// query to select those ids that exist
	stmt := fmt.Sprintf(`SELECT
		id
	FROM issues
	WHERE projectId = "%s"
	AND id IN (%s)
	`, githubProjectId, idsQuery)
	rows, err := p.models.db.Query(stmt)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()
	idsThatExist := []string{}
	for rows.Next() {
		id := ""

		err = rows.Scan(
			&id,
		)
		if err != nil {
			return []string{}, err
		}
		idsThatExist = append(idsThatExist, id)
	}

	return idsThatExist, nil
}

// FindThoseThatDoesntExist takes a list of issues ids and returns those that does not exists.
func (p *Issues) FindThoseThatDoesntExist(githubProjectId string, ids []string) ([]string, error) {
	idsThatExist, err := p.FindThoseThatExist(githubProjectId, ids)
	if err != nil {
		return []string{}, err
	}

	idsThatDoesntExist := []string{}

	for i := 0; i < len(ids); i++ {
		found := false
		for j := 0; j < len(idsThatExist); j++ {
			if ids[i] == idsThatExist[j] {
				found = true
			}
		}
		if !found {
			idsThatDoesntExist = append(idsThatDoesntExist, ids[i])
		}
	}

	return idsThatDoesntExist, nil
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
