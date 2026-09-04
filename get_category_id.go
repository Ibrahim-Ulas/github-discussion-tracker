package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CategoriesResponse struct {
	Data struct {
		Repository struct {
			DiscussionCategories struct {
				Nodes []DiscussionCategory `json:"nodes"`
			} `json:"discussionCategories"`
		} `json:"repository"`
	} `json:"data"`
}

type DiscussionCategory struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func getCategoryID(owner, repoName, token, apiURL, categoryName string) (string, error) {
	categoryQuery := fmt.Sprintf(`query{
		repository(owner:"%s", name:"%s") {
			discussionCategories(first:10) {
				nodes{
					id
					name
				}
			}
		}
	}`, owner, repoName)

	client := http.Client{}

	payload := map[string]string{"query": categoryQuery}
	requestBody, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshaling payload: %v", err)
		return "", err
	}
	requestBodyReader := bytes.NewReader(requestBody)
	req, err := http.NewRequest("POST", apiURL, requestBodyReader)
	if err != nil {
		fmt.Printf("Error making the request: %v\n", err)
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Github-Discussion-Tracker")

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("API responded with: %v\n", err)
		return "", err
	}
	responseBodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return "", err
	}
	responseBody := CategoriesResponse{}
	err = json.Unmarshal(responseBodyBytes, &responseBody)
	if err != nil {
		fmt.Printf("Error unmarshaling response body: %v\n", err)
		return "", err
	}

	for _, category := range responseBody.Data.Repository.DiscussionCategories.Nodes {
		if category.Name == categoryName {
			return category.ID, nil
		}
	}
	return "", fmt.Errorf("No categories found by the name %v", categoryName)
}
