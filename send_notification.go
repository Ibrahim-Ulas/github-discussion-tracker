package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
)

type notification struct {
	title   string
	message string
}

var notifyCh = make(chan notification, 5)

func init() {
	go func() {
		for i := range notifyCh {
			i := i
			go func() {
				err := beeep.Alert(i.title, i.message, "")
				if err != nil {
					fmt.Printf("Error sending notification %v\n", err)
				}
			}()
		}
	}()
}

func queueNotification(title, message string) {
	select {
	case notifyCh <- notification{title: title, message: message}:
	default:
		fmt.Println("Notification queue full, skipping")
	}
}

func sendNotification(lastNotificationTime *time.Time, allDiscussions *AllDiscussions, notificationCooldownMultiplier int) {
	if len(allDiscussions.UnansweredDiscussions) > 0 {
		if time.Since(*lastNotificationTime) >= time.Duration(notificationCooldownMultiplier)*time.Second-1*time.Second {
			*lastNotificationTime = time.Now()
			var sb strings.Builder
			sb.WriteString("Unanswered discussions:\n")
			for i := range allDiscussions.UnansweredDiscussions {
				sb.WriteString(fmt.Sprintf("- #%d\n", allDiscussions.UnansweredDiscussions[i]))
			}
			queueNotification("Unanswered Discussions", sb.String())
		}
	}
}
