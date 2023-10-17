package main

import (
	"fmt"
	"os"
)

// Set the default query parameters
var (
	projects   = []string{"SVPI", "SANDBOX", "ASC", "APPSVC", "PLNSRVCE", "SRVKP", "TKNECO", "GITOPSRVCE", "GITOPS", "GITOPSVCCR", "DEVHAS", "RHIDP"}
	maxResults = "1000"
	startAt    = "0"
)

var jiraAPIToken string

func init() {
	jiraAPIToken = os.Getenv("JIRA_API_TOKEN")
	if jiraAPIToken == "" {
		fmt.Println("JIRA_API_TOKEN is empty. Please set it to a valid JIRA API token.")
		return
	}
}
