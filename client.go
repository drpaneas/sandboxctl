package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func getIssues() (Output, error) {
	// Construct the JQL query to fetch all the issues
	listIssueQuery := fmt.Sprintf("/rest/api/2/search?jql=project%%3D%%22%s%%22+ORDER+BY+created+DESC&startAt=%s&maxResults=%s", project, startAt, maxResults)

	// Make a GET request to the JIRA API
	req, err := http.NewRequest("GET", fmt.Sprintf("https://issues.redhat.com%s", listIssueQuery), nil)
	if err != nil {
		fmt.Println(err)
		return Output{}, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jiraAPIToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return Output{}, err
	}

	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return Output{}, err
	}

	// Parse the JSON response
	var allIssues Output
	err = json.Unmarshal(body, &allIssues)
	if err != nil {
		fmt.Println(err)
		return Output{}, err
	}

	return allIssues, nil
}
