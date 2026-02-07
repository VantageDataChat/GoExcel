package gospreadsheet

// Comment represents a cell comment/annotation.
type Comment struct {
	Author string
	Text   string
}

// NewComment creates a new comment.
func NewComment(author, text string) *Comment {
	return &Comment{Author: author, Text: text}
}
