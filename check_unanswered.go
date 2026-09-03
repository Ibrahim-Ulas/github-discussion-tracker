package main

func checkUnansweredDiscussion(comments CommentConnection) bool {
	if len(comments.Nodes) == 0 {
		return true
	}
	return false
}
