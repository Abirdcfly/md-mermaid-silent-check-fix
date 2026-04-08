package fixer

import (
	"regexp"
	"strings"

	"github.com/Abirdcfly/md-mermaid-silent-check-fix/model"
)

func FixNewline(content string) string {
	return strings.ReplaceAll(content, `\n`, "<br>")
}

func FixUnquotedText(content string) string {
	re := regexp.MustCompile(`\[([^\]]*[():{}][^\]]*)\]`)
	return re.ReplaceAllStringFunc(content, func(match string) string {
		if strings.Contains(match, `"`) {
			return match
		}
		return re.ReplaceAllString(match, `["$1"]`)
	})
}

func ApplyFixes(content string, block model.MermaidBlock) string {
	fixedContent := block.Content
	for _, issue := range block.Issues {
		if issue.Fixable && issue.Fix != nil {
			fixedContent = issue.Fix(fixedContent)
		}
	}
	if fixedContent == block.Content {
		return content
	}

	lines := strings.Split(content, "\n")
	var newLines []string

	inBlock := false
	for i, line := range lines {
		lineNum := i + 1
		if lineNum == block.StartLine {
			inBlock = true
			newLines = append(newLines, line)
			fixedLines := strings.Split(fixedContent, "\n")
			newLines = append(newLines, fixedLines...)
			continue
		}
		if lineNum == block.EndLine {
			inBlock = false
			newLines = append(newLines, line)
			continue
		}
		if !inBlock {
			newLines = append(newLines, line)
		}
	}

	return strings.Join(newLines, "\n")
}
