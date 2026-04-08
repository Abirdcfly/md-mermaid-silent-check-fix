package parser

import (
	"strings"

	"github.com/Abirdcfly/md-mermaid-silent-check-fix/model"
)

func ExtractMermaidBlocks(content string) []model.MermaidBlock {
	var blocks []model.MermaidBlock
	lines := strings.Split(content, "\n")

	inMermaid := false
	startLine := 0
	var contentLines []string

	for i, line := range lines {
		lineNum := i + 1
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "```mermaid") {
			inMermaid = true
			startLine = lineNum
			contentLines = []string{}
			continue
		}

		if inMermaid && trimmed == "```" {
			inMermaid = false
			blocks = append(blocks, model.MermaidBlock{
				StartLine: startLine,
				EndLine:   lineNum,
				Content:   strings.Join(contentLines, "\n"),
			})
			continue
		}

		if inMermaid {
			contentLines = append(contentLines, line)
		}
	}

	return blocks
}
