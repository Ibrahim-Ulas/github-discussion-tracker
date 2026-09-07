package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
)

func sendNotification(lastNotificationTime time.Time, allDiscussions *AllDiscussions, notificationCooldownMultiplier int) {
	if len(allDiscussions.UnansweredDiscussions) > 0 {
		if time.Since(lastNotificationTime) >= time.Duration(notificationCooldownMultiplier)*time.Second-1*time.Second {
			var sb strings.Builder
			sb.WriteString("Unanswered discussions:\n")
			for i := range allDiscussions.UnansweredDiscussions {
				sb.WriteString(fmt.Sprintf("- #%d\n", allDiscussions.UnansweredDiscussions[i]))
			}
			err := beeep.Alert("Unanswered Discussions", sb.String(), "")
			if err != nil {
				fmt.Printf("Error sending notification: %v\n", err)
			} else {
				lastNotificationTime = time.Now()
			}
		}
	}
}
