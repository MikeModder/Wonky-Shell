package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const glURL = "https://gitlab.com/api/v4/projects/7808366/repository/commits"

// GLResp GitLab commits response
type GLResp struct {
	Commits []GLCommit
}

// GLCommit a GitLab commit object
type GLCommit struct {
	CommitID       string   `json:"id"`
	CommitShort    string   `json:"short_id"`
	Title          string   `json:"title"`
	AuthorName     string   `json:"author_name"`
	AuhtorEmail    string   `json:"author_email"`
	CommitDate     string   `json:"authored_date"`
	CommitterName  string   `json:"comitter_name"`
	ComiitterEmail string   `json:"comitter_email"`
	CommittedDate  string   `json:"comitted_date"`
	CreatedAt      string   `json:"created_at"`
	Message        string   `json:"message"`
	ParentIds      []string `json:"parent_ids"`
}

// CheckUpdate Poll GitLab about the latest commit
func CheckUpdate() (bool, bool, string) {
	glGET, e := http.Get(glURL)
	if e != nil {
		fmt.Println("Error while contacting GitLab API!")
		return false, true, ""
	}
	var commits []GLCommit
	e = unmarshal(glGET, &commits)
	if e != nil {
		fmt.Println("Error while decoding GitLab response!")
		return false, true, ""
	}
	latest := commits[0]
	hash := latest.CommitShort[0:7]
	if GitCommit != hash {
		return true, false, hash
	}
	return false, false, ""
}

func unmarshal(b *http.Response, t interface{}) error {
	defer b.Body.Close()
	return json.NewDecoder(b.Body).Decode(t)
}
