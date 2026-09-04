package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

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

	checkIntervalString := os.Getenv("CHECK_INTERVAL")
	checkInterval, err := strconv.Atoi(checkIntervalString)
	if err != nil {
		log.Fatalf("Invalid interval: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(time.Duration(checkInterval) * time.Second)

	owner := os.Getenv("OWNER")
	token := os.Getenv("GITHUB_TOKEN")
	repo := os.Getenv("REPO")
	categoryName := os.Getenv("CATEGORY_NAME")
	API_URL := "https://api.github.com/graphql"
	allDiscussions := AllDiscussions{
		repliedDiscussions:    make(map[int]int),
		UnansweredDiscussions: make(map[int]string),
	}

	for {
		select {
		case <-sigChan:
			fmt.Println("\nClosing application...")
			return
		case <-ticker.C:
			fmt.Printf("Getting discussions from: %v/%v\n", owner, repo)

			discussions, err := getDiscussionsFromGitHub(owner, repo, token, API_URL, categoryName)
			if err != nil {
				fmt.Printf("Error getting discussions: %v\n", err)
				return
			}
			nodes := discussions.Data.Repository.Discussions.Nodes
			markDiscussions(nodes, &allDiscussions)
			fmt.Print("Unanswered discussions: \n")
			if len(allDiscussions.UnansweredDiscussions) == 0 {
				fmt.Print("No unanswered discussions\n")
			} else {
				for i := range allDiscussions.UnansweredDiscussions {
					fmt.Print(allDiscussions.UnansweredDiscussions[i] + "\n")
				}
			}
		}
	}
}
