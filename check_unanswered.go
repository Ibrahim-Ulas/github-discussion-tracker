package main

func checkUnansweredDiscussion(discussion DiscussionNode) bool {
	comments := discussion.Comments.Nodes
	if len(comments) == 0 {
		return true
	}

	for i := range comments {
		for _, reaction := range comments[i].Reactions.Nodes {
			if reaction.Content == "THUMBS_UP" {
				return false
			}
		}
	}

	if comments[len(comments)-1].Author.Username == discussion.Author.Username {
		return true
	}

	return false
}

func markDiscussions(nodes []DiscussionNode, allDiscussions *AllDiscussions) {
	activeDiscussions := make(map[int]struct{}, len(nodes))

	for i := range nodes {
		node := nodes[i]
		discussionID := node.ID

		activeDiscussions[discussionID] = struct{}{}

		if checkUnansweredDiscussion(node) {
			allDiscussions.UnansweredDiscussions[discussionID] = node.Number
		} else {
			delete(allDiscussions.UnansweredDiscussions, discussionID)
		}
	}

	for id := range allDiscussions.UnansweredDiscussions {
		if _, exists := activeDiscussions[id]; !exists {
			delete(allDiscussions.UnansweredDiscussions, id)
		}
	}
}
