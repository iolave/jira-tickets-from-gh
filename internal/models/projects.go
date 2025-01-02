package models

type Project struct {
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

type Projects struct {
	models *Models
}

func (p *Projects) Upsert(projects ...Project) error {
	for i := 0; i < len(projects); i++ {
		stmt := `
		INSERT OR REPLACE INTO projects(
			id,
			FID_jiraUrl,
			FID_jiraIssueType,
			FID_title,
			FID_estimate,
			FID_status,
			FID_assignees,
			FID_repository
		) values(?, ?, ?, ?, ?, ?, ?, ?)
		`
		_, err := p.models.db.Exec(
			stmt,
			projects[i].ID,
			projects[i].Fields.JiraURL,
			projects[i].Fields.JiraIssueType,
			projects[i].Fields.Title,
			projects[i].Fields.Estimate,
			projects[i].Fields.Status,
			projects[i].Fields.Assignees,
			projects[i].Fields.Repository,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
