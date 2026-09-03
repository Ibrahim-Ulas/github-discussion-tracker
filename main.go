package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AllDiscussions struct {
	repliedDiscussions    map[int]int
	UnansweredDiscussions map[int]string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading or reading .env file")
	}

	owner := os.Getenv("OWNER")
	token := os.Getenv("GITHUB_TOKEN")
	repo := os.Getenv("REPO")
	API_URL := "https://api.github.com/graphql"
	allDiscussions := AllDiscussions{
		repliedDiscussions:    make(map[int]int),
		UnansweredDiscussions: make(map[int]string),
	}
	fmt.Printf("Getting discussions from: %v/%v\n", owner, repo)

	discussions, err := getDiscussionsFromGitHub(owner, repo, token, API_URL)
	if err != nil {
		fmt.Printf("Error getting discussions: %v\n", err)
		return
	}
	nodes := discussions.Data.Repository.Discussions.Nodes
	for i := range nodes {
		if checkUnansweredDiscussion(nodes[i].Comments) {
			allDiscussions.UnansweredDiscussions[nodes[i].ID] = nodes[i].Title
		}
	}

	fmt.Print("Cevaplanmamış discussionlar: \n")
	for i := range allDiscussions.UnansweredDiscussions {
		fmt.Print(allDiscussions.UnansweredDiscussions[i] + "\n")
	}
}
