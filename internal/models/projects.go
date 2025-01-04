package models

import (
	"errors"
	"fmt"
)

type Project struct {
	models *Models
	ID     string // github project id
	Fields struct {
		JiraURL       string
		JiraIssueType string
		Title         string
		Estimate      string
		Status        string
		Assignees     string
		Repository    string
	} // project fields ids within github
}

func (p Project) UpsertIssue(id, title string, status *IssueStatus, jiraUrl, jiraIssueType, repo *string, estimate *int, assignees *[]string) (*Issue, error) {
	return p.models.Issues.Upsert(p.ID, id, title, status, jiraUrl, jiraIssueType, repo, estimate, assignees)
}

func (p Project) GetIssue(id string) (*Issue, error) {
	return p.models.Issues.Get(p.ID, id)
}

type Projects struct {
	models *Models
}

func (p Projects) Upsert(id, jiraUrlId, jiraIssueTypeId, titleId, estimateId, statusId, assigneesId, repositoryId string) (*Project, error) {
	project := new(Project)
	project.models = p.models
	project.ID = id
	project.Fields.Assignees = assigneesId
	project.Fields.JiraURL = jiraUrlId
	project.Fields.JiraIssueType = jiraIssueTypeId
	project.Fields.Status = statusId
	project.Fields.Title = titleId
	project.Fields.Estimate = estimateId
	project.Fields.Repository = repositoryId
	stmt := `INSERT OR REPLACE INTO projects(
			id,
			FID_jiraUrl,
			FID_jiraIssueType,
			FID_title,
			FID_estimate,
			FID_status,
			FID_assignees,
			FID_repository
		) values(?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := p.models.db.Exec(
		stmt,
		project.ID,
		project.Fields.JiraURL,
		project.Fields.JiraIssueType,
		project.Fields.Title,
		project.Fields.Estimate,
		project.Fields.Status,
		project.Fields.Assignees,
		project.Fields.Repository,
	)

	if err != nil {
		return nil, err
	}
	return project, nil
}

// Get gets a project, if no project found *Project will be nil
func (p *Projects) Get(githubId string) (*Project, error) {
	if githubId == "" {
		return nil, errors.New(`please provide a value for "githubId"`)
	}

	stmt := fmt.Sprintf(`SELECT
		id,
		FID_jiraUrl,
		FID_jiraIssueType,
		FID_title,
		FID_estimate,
		FID_status,
		FID_assignees,
		FID_repository
	FROM projects
	WHERE id = "%s"
	`, githubId)
	rows, err := p.models.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	project := new(Project)
	if !rows.Next() {
		return nil, nil
	}
	err = rows.Scan(
		&project.ID,
		&project.Fields.JiraURL,
		&project.Fields.JiraIssueType,
		&project.Fields.Title,
		&project.Fields.Estimate,
		&project.Fields.Status,
		&project.Fields.Assignees,
		&project.Fields.Repository,
	)
	if err != nil {
		return nil, err
	}
	project.models = p.models
	return project, nil
}
