package main

import (
	"fmt"

	"github.com/gen2brain/beeep"
)

func checkUnansweredDiscussion(comments CommentConnection) bool {
	if len(comments.Nodes) == 0 {
		return true
	}
	return false
}

func markDiscussions(nodes []DiscussionNode, allDiscussions *AllDiscussions) {
	for i := range nodes {
		node := nodes[i]
		commentNodes := node.Comments.Nodes
		discussionID := node.ID

		if checkUnansweredDiscussion(node.Comments) {
			allDiscussions.UnansweredDiscussions[discussionID] = node.Number
		} else {
			delete(allDiscussions.UnansweredDiscussions, discussionID)
			if _, ok := allDiscussions.repliedDiscussions[discussionID]; ok {
				if allDiscussions.repliedDiscussions[discussionID] != commentNodes[0].ID {
					fmt.Printf("%v titled discussion has a new comment!\n", node.Title)
					
					err := beeep.Alert("New Comment!",fmt.Sprintf("#%d numbered discussion has a new comment!", node.Number),"")
					if err != nil {
						fmt.Printf("Error sending notification: %v\n", err)
					}
				}
			}
			allDiscussions.repliedDiscussions[discussionID] = commentNodes[0].ID
		}
	}
}
