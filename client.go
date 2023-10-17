package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func getIssues(project string) (Output, error) {
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

func getLast7dUpdates(project string) (Output, error) {
	//  jira issue list --updated -7d --status="In Progress" --status="Closed" --status="Review" -tStory
	listIssueQuery := fmt.Sprintf("/rest/api/2/search?jql=project%%3D%%22%s%%22+AND+type%%3D%%22Story%%22+AND+updatedDate%%3E%%3D%%22-7d%%22+AND+status+IN+%%28%%22In+Progress%%22%%2C+%%22Closed%%22%%2C+%%22Review%%22%%29+ORDER+BY+updated+DESC&startAt=%s&maxResults=%s", project, startAt, maxResults)

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

func getStoryInfo(storyName string) (IssueOutput, error) {
	//  jira issue view $storyName --plain --debug
	listIssueQuery := fmt.Sprintf("/rest/api/2/issue/%v", storyName)

	// Make a GET request to the JIRA API
	req, err := http.NewRequest("GET", fmt.Sprintf("https://issues.redhat.com%s", listIssueQuery), nil)
	if err != nil {
		fmt.Println(err)
		return IssueOutput{}, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jiraAPIToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return IssueOutput{}, err
	}

	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return IssueOutput{}, err
	}

	// Parse the JSON response
	var reply IssueOutput
	err = json.Unmarshal(body, &reply)
	if err != nil {
		fmt.Println(err)
		return IssueOutput{}, err
	}

	return reply, nil
}
