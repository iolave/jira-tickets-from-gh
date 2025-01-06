package github

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type GitHubClient struct {
	client *http.Client
	token  string
}

type Error struct {
	Message string  `json:"message"`
	Type    *string `json:"type"`
}

func (e Error) ToError() error {
	if e.Type != nil {
		return fmt.Errorf("%s: %s", *e.Type, e.Message)
	}

	return fmt.Errorf("%s: %s", *e.Type, e.Message)
}

func getErrorFromErrors(errs *[]Error) error {
	if errs == nil {
		return nil
	}

	if len(*errs) == 0 {
		return errors.New("unknown graphql error")
	}

	err := (*errs)[0]

	return err.ToError()
}

func New(token string) *GitHubClient {
	transport := http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{
		Transport: &transport,
	}

	return &GitHubClient{
		client: client,
		token:  token,
	}
}

func (c *GitHubClient) request(query string, result any) (*http.Response, error) {
	var requestBody bytes.Buffer
	requestBodyObj := struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables"`
	}{
		Query:     query,
		Variables: map[string]interface{}{},
	}

	if err := json.NewEncoder(&requestBody).Encode(requestBodyObj); err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, "https://api.github.com/graphql", &requestBody)
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", c.token))
	if err != nil {
		return nil, err
	}
	res, err := c.client.Do(req)
	if err != nil {
		return res, err
	}

	if res.StatusCode != http.StatusOK {
		return res, errors.New("failed to send github graphql request")
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(b, &result)
	if err != nil {
		return res, err
	}

	return res, nil
}
