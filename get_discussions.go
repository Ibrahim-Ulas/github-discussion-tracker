package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CommentNode struct {
	ID int `json:"databaseId"`
}

type CommentConnection struct {
	Nodes []CommentNode `json:"nodes"`
}

type DiscussionNode struct {
	ID       int               `json:"databaseId"`
	Title    string            `json:"title"`
	Comments CommentConnection `json:"comments"`
}

type GraphQLResponse struct {
	Data struct {
		Repository struct {
			Discussions struct {
				Nodes []DiscussionNode `json:"nodes"`
			} `json:"discussions"`
		} `json:"repository"`
	} `json:"data"`
}

func getDiscussionsFromGitHub(owner, repoName, token, apiUrl string) (GraphQLResponse, error) {
	discussionsQuery := fmt.Sprintf(`query {
		repository(owner:"%s", name:"%s"){
			discussions(first:5, states:OPEN) {
				nodes {
					title
					databaseId
					comments(first:10) {
						nodes{
							databaseId
						}
					}
				}
			}
		}
	}`, owner, repoName)
	client := http.Client{}

	payload := map[string]string{"query": discussionsQuery}
	requestBody, err := json.Marshal(payload)
	if err != nil {
		return GraphQLResponse{}, err
	}

	requestBodyReader := bytes.NewReader(requestBody)
	req, err := http.NewRequest("POST", apiUrl, requestBodyReader)
	if err != nil {
		fmt.Printf("Error creating new request: %v", err)
		return GraphQLResponse{}, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Github-Discussion-Tracker")

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error requesting: %v", err)
		return GraphQLResponse{}, err
	}
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return GraphQLResponse{}, err
	}

	discussions := GraphQLResponse{}
	err = json.Unmarshal(resBytes, &discussions)
	if err != nil {
		return GraphQLResponse{}, err
	}

	return discussions, nil
}
