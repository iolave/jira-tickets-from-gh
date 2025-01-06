package github

import (
	"fmt"
	"net/http"
)

type ProjectFieldType int

const (
	PROJECT_FIELD_TEXT ProjectFieldType = iota
	PROJECT_FIELD_SINGLE_SELECT
	PROJECT_FIELD_USER
	PROJECT_FIELD_NUMBER
	PROJECT_FIELD_REPO
)

type ListUserProjectsResult struct {
	Errors *[]Error `json:"errors"`
	Data   struct {
		User struct {
			Projects struct {
				Nodes []struct {
					ID    string `json:"id"`
					Title string `json:"title"`
				} `json:"nodes"`
			} `json:"projectsV2"`
		} `json:"user"`
	} `json:"data"`
}

func (c *GitHubClient) ListUserProjects(user string) (ListUserProjectsResult, *http.Response, error) {
	query := fmt.Sprintf(`query{
		user(login:"%s") {
			projectsV2(first:100){ nodes { id title } }
		}
	}`, user)

	var result ListUserProjectsResult

	res, err := c.request(query, &result)
	if err != nil {
		return result, res, err
	}
	err = getErrorFromErrors(result.Errors)

	return result, res, err
}

type ListOrganizationProjectsResult struct {
	Errors *[]Error `json:"errors"`
	Data   struct {
		Organization struct {
			Projects struct {
				Nodes []struct {
					ID    string `json:"id"`
					Title string `json:"title"`
				} `json:"nodes"`
			} `json:"projectsV2"`
		} `json:"organization"`
	} `json:"data"`
}

func (c *GitHubClient) ListOrganizationProjects(org string) (ListOrganizationProjectsResult, *http.Response, error) {
	query := fmt.Sprintf(`query{
		organization(login:"%s") {
			projectsV2(first:100){ nodes { id title } }
		}
	}`, org)

	var result ListOrganizationProjectsResult

	res, err := c.request(query, &result)
	if err != nil {
		return result, res, err
	}
	err = getErrorFromErrors(result.Errors)

	return result, res, err
}

type GetProjectFieldsResult struct {
	Errors *[]Error `json:"errors"`
	Data   struct {
		Node struct {
			Fields struct {
				Nodes []struct {
					ID      string `json:"id"`
					Name    string `json:"name"`
					Options *[]struct {
						ID   string `json:"id"`
						Name string `json:"name"`
					}
				} `json:"nodes"`
			} `json:"fields"`
		} `json:"node"`
	} `json:"data"`
}

func (c *GitHubClient) GetProjectFields(id string) (GetProjectFieldsResult, *http.Response, error) {
	query := fmt.Sprintf(`query{ node(id: "%s") {
		... on ProjectV2 {
			fields(first: 100) {nodes {
				... on ProjectV2Field { id name } 
				... on ProjectV2IterationField { id name }
				... on ProjectV2SingleSelectField { id name options { id name }}
			}}
		}
	}}`, id)

	var result GetProjectFieldsResult

	res, err := c.request(query, &result)
	if err != nil {
		return result, res, err
	}
	err = getErrorFromErrors(result.Errors)

	return result, res, err
}

type ProjectField struct {
	Type       ProjectFieldType
	FieldName  string
	FieldAlias string
}

func (f ProjectField) ToQuery() string {
	switch f.Type {
	case PROJECT_FIELD_TEXT:
		return fmt.Sprintf(`
			%s: fieldValueByName(name: "%s") {
				__typename
				... on ProjectV2ItemFieldTextValue {text}
			}
		`, f.FieldAlias, f.FieldName)
	case PROJECT_FIELD_REPO:
		return fmt.Sprintf(`
			%s: fieldValueByName(name: "%s") {
				__typename
				... on ProjectV2ItemFieldRepositoryValue {repository{nameWithOwner}}
			}
		`, f.FieldAlias, f.FieldName)
	case PROJECT_FIELD_SINGLE_SELECT:
		return fmt.Sprintf(`
			%s: fieldValueByName(name: "%s") {
				__typename
				... on ProjectV2ItemFieldSingleSelectValue {name optionId}
			}
		`, f.FieldAlias, f.FieldName)
	case PROJECT_FIELD_NUMBER:
		return fmt.Sprintf(`
			%s: fieldValueByName(name: "%s") {
				__typename
				... on ProjectV2ItemFieldNumberValue {number}
			}
		`, f.FieldAlias, f.FieldName)
	case PROJECT_FIELD_USER:
		return fmt.Sprintf(`
			%s: fieldValueByName(name: "%s") {
				__typename
				... on ProjectV2ItemFieldUserValue {users(first:100){nodes {login}}}
			}
		`, f.FieldAlias, f.FieldName)
	default:
		return ""
	}
}

type GetProjectItemsResult struct {
	Errors *[]Error `json:"errors"`
	Data   struct {
		Node struct {
			Items struct {
				Nodes    []map[string]any `json:"nodes"` // TODO: Add better way to access items
				PageInfo struct {
					StartCursor string `json:"startCursor"`
					EndCursor   string `json:"endCursor"`
					HasNextPage bool   `json:"hasNextPage"`
					HasPrevPage bool   `json:"hasPreviousPage"`
				} `json:"pageInfo"`
			} `json:"items"`
		} `json:"node"`
	} `json:"data"`
}

// TODO: Add better way to access items
func (c *GitHubClient) GetProjectItems(id string, fields []ProjectField) (GetProjectItemsResult, *http.Response, error) {
	queryFields := ""
	for i := 0; i < len(fields); i++ {
		queryFields = fmt.Sprintf("%s %s", queryFields, fields[i].ToQuery())
	}

	query := fmt.Sprintf(`query{ node(id: "%s") { ... on ProjectV2 {
		items(first: 100) {
			pageInfo{startCursor endCursor hasNextPage hasPreviousPage} 
			nodes{
				id
				content{
					__typename
					... on Issue {comments(first:100) {nodes{body}}}
				}
				%s
			}
		}
	}}}}`, id, queryFields)

	var result GetProjectItemsResult

	res, err := c.request(query, &result)
	if err != nil {
		return result, res, err
	}
	// err = getErrorFromErrors(result.Errors)

	return result, res, err

}
