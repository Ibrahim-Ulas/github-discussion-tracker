package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ReactionNode struct {
	Content string `json:"content"`
}

type ReactionConnection struct {
	Nodes []ReactionNode `json:"nodes"`
}
type Actor struct {
	Username string `json:"login"`
}
type CommentNode struct {
	Author    Actor              `json:"author"`
	ID        int                `json:"databaseId"`
	Reactions ReactionConnection `json:"reactions"`
}

type CommentConnection struct {
	Nodes []CommentNode `json:"nodes"`
}

type DiscussionNode struct {
	Author   Actor             `json:"author"`
	ID       int               `json:"databaseId"`
	Number   int               `json:"number"`
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

func getDiscussionsFromGitHub(owner, repoName, token, apiUrl, categoryName string) (GraphQLResponse, error) {
	var categoryId string
	var err error
	if categoryName != "" {
		categoryId, err = getCategoryID(owner, repoName, token, apiUrl, categoryName)
		if err != nil {
			fmt.Print(err)
			return GraphQLResponse{}, err
		}
	}

	categoryIdQuery := ""
	if categoryId != "" {
		categoryIdQuery = fmt.Sprintf(`categoryId: "%s",`, categoryId)
	}

	discussionsQuery := fmt.Sprintf(`query {
		repository(owner:"%s", name:"%s"){
			discussions(first:20, %s states:OPEN) {
				nodes {
					author{
						login
					}
					title
					number
					databaseId
					comments(last:10) {
						nodes{
							author {
								login
							}
							reactions(first:5) {
								nodes{
									content
								}
								
							}
							databaseId
						}
					}
				}
			}
		}
	}`, owner, repoName, categoryIdQuery)

	client := http.Client{Timeout: 60 * time.Second}

	payload := map[string]string{"query": discussionsQuery}
	requestBody, err := json.Marshal(payload)
	if err != nil {
		return GraphQLResponse{}, err
	}

	requestBodyReader := bytes.NewReader(requestBody)
	req, err := http.NewRequest("POST", apiUrl, requestBodyReader)
	if err != nil {
		fmt.Printf("Error creating new request: %v\n", err)
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
	// fmt.Print(string(resBytes))
	defer res.Body.Close()
	fmt.Printf("Status Code: %v\n", res.StatusCode)
	discussions := GraphQLResponse{}
	err = json.Unmarshal(resBytes, &discussions)
	if err != nil {
		return GraphQLResponse{}, err
	}
	return discussions, nil
}
