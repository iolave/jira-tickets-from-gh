package models

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Models struct {
	db       *sql.DB
	Projects Projects
}

func (m *Models) Close() error {
	return m.db.Close()
}

func Initialize() (*Models, error) {
	models := new(Models)

	db, err := sql.Open("sqlite3", "./data/storage.db")
	if err != nil {
		return nil, err
	}
	models.db = db

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS projects (
		id string not null,
		FID_jiraUrl string not null,
		FID_jiraIssueType string not null,
		FID_title string not null,
		FID_estimate string not null,
		FID_status string not null,
		FID_assignees string not null,
		FID_repository string not null,
		primary key (id)
	)`)
	if err != nil {
		return nil, err
	}

	models.Projects = Projects{models: models}
	return models, nil
}
