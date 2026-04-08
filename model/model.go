package model

type MarkdownFile struct {
	Path    string
	Content string
	Blocks  []MermaidBlock
}

type MermaidBlock struct {
	StartLine int
	EndLine   int
	Content   string
	Issues    []Issue
}

type Issue struct {
	Type    string
	Message string
	Line    int
	Fixable bool
	Fix     func(string) string
}

const (
	IssueTypeNewlineLiteral   = "newline_literal"
	IssueTypeHTMLiteral       = "html_literal"
	IssueTypeDuplicateNode    = "duplicate_node"
	IssueTypeUndefinedClass   = "undefined_class"
	IssueTypeInvalidStyle     = "invalid_style"
	IssueTypeUnquotedText     = "unquoted_text"
	IssueTypeIsolatedNode     = "isolated_node"
	IssueTypeDuplicateSubgraph = "duplicate_subgraph"
)
