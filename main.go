package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/joho/godotenv"
)

type AllDiscussions struct {
	UnansweredDiscussions map[int]int
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading or reading .env file")
	}
	fmt.Print("Starting Discussion Tracker...\n")
	beeep.AppName = "GitHub Discussion Tracker"
	checkIntervalString := os.Getenv("CHECK_INTERVAL_IN_SECONDS")
	checkInterval, err := strconv.Atoi(checkIntervalString)
	if err != nil {
		log.Fatalf("Invalid interval: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(1 * time.Second)
	remainingSeconds := checkInterval

	var lastNotificationTime time.Time

	owner := os.Getenv("OWNER")
	token := os.Getenv("GITHUB_TOKEN")
	repo := os.Getenv("REPO")
	notificationCooldownMultiplier, err := strconv.Atoi(os.Getenv("NOTIFICATION_COOLDOWN_MULTIPLIER"))
	if err != nil {
		log.Fatalf("Invalid notification cooldown multiplier: %v", err)
	}
	categoryName := os.Getenv("CATEGORY_NAME")
	API_URL := "https://api.github.com/graphql"
	allDiscussions := AllDiscussions{
		UnansweredDiscussions: make(map[int]int),
	}

	check := func() {
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
				fmt.Printf("Discussion %d: %d\n", i, allDiscussions.UnansweredDiscussions[i])
			}
			sendNotification(&lastNotificationTime, &allDiscussions, notificationCooldownMultiplier)
		}
	}

	check()
	for {
		select {
		case <-sigChan:
			fmt.Println("\nClosing application...")
			return
		case <-ticker.C:
			remainingSeconds--
			if remainingSeconds >= 0 {
				fmt.Printf("\r Next check in: %2d seconds...", remainingSeconds)
			} else {
				fmt.Printf("\r\n")
				check()
				remainingSeconds = checkInterval
			}
		}
	}
}
