package analyzer

import (
	"regexp"
	"strings"

	"github.com/Abirdcfly/md-mermaid-silent-check-fix/fixer"
	"github.com/Abirdcfly/md-mermaid-silent-check-fix/model"
)

var (
	nodeIDRegex   = regexp.MustCompile(`^\s*([A-Za-z0-9_-]+)\s*\[`)
	classDefRegex = regexp.MustCompile(`classDef\s+([A-Za-z0-9_-]+)`)
	classUseRegex = regexp.MustCompile(`class\s+([A-Za-z0-9\s,_-]+)\s+([A-Za-z0-9_-]+)`)
	styleRegex    = regexp.MustCompile(`style\s+[A-Za-z0-9_-]+\s+([a-z-]+):`)
	subgraphRegex = regexp.MustCompile(`subgraph\s+([A-Za-z0-9_-]+)`)
	nodeTextRegex = regexp.MustCompile(`\[([^\]]+)\]`)
	htmlTagRegex  = regexp.MustCompile(`<(div|span|p|div\s|span\s|p\s)`)
	// connectionRegex matches any arrow connection and extracts source/target nodes
	// handles all arrow types: -->, ->, ==>, -.->, -.-, ~~~, etc.
	// also handles labeled arrows with spaces in labels
	connectionRegex = regexp.MustCompile(`([A-Za-z0-9_-]+)\s*.*[-=.][-=.]+[>]\s*.*\s+([A-Za-z0-9_-]+)`)

	validStyleProps = map[string]bool{
		"fill": true, "stroke": true, "stroke-width": true, "color": true,
		"background": true, "border": true, "font-size": true, "font-family": true,
		"font-weight": true, "text-align": true, "padding": true, "margin": true,
		"width": true, "height": true, "opacity": true, "rx": true, "ry": true,
	}
)

func AnalyzeBlock(block model.MermaidBlock) []model.Issue {
	var issues []model.Issue

	issues = checkNewlineLiteral(issues, block)
	issues = checkHTMLiteral(issues, block)
	issues = checkDuplicateNode(issues, block)
	issues = checkUndefinedClass(issues, block)
	issues = checkInvalidStyle(issues, block)
	issues = checkUnquotedText(issues, block)
	issues = checkIsolatedNode(issues, block)
	issues = checkDuplicateSubgraph(issues, block)

	return issues
}

func checkNewlineLiteral(issues []model.Issue, block model.MermaidBlock) []model.Issue {
	lines := strings.Split(block.Content, "\n")
	for i, line := range lines {
		if strings.Contains(line, `\n`) {
			issues = append(issues, model.Issue{
				Type:     model.IssueTypeNewlineLiteral,
				Message:  "Found newline literal \\n in text, use <br> instead",
				Line:     block.StartLine + i,
				Fixable:  true,
				Fix:      fixer.FixNewline,
				Severity: model.SeverityError,
			})
		}
	}
	return issues
}

func checkHTMLiteral(issues []model.Issue, block model.MermaidBlock) []model.Issue {
	lines := strings.Split(block.Content, "\n")
	for i, line := range lines {
		if htmlTagRegex.MatchString(strings.ToLower(line)) {
			issues = append(issues, model.Issue{
				Type:     model.IssueTypeHTMLiteral,
				Message:  "Found HTML tag <div>/<span>, may not render correctly in Mermaid",
				Line:     block.StartLine + i,
				Fixable:  false,
				Severity: model.SeverityError,
			})
		}
	}
	return issues
}

func checkDuplicateNode(issues []model.Issue, block model.MermaidBlock) []model.Issue {
	nodes := make(map[string]int)
	lines := strings.Split(block.Content, "\n")

	for i, line := range lines {
		matches := nodeIDRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			nodeID := matches[1]
			if firstLine, exists := nodes[nodeID]; exists {
				issues = append(issues, model.Issue{
					Type:     model.IssueTypeDuplicateNode,
					Message:  "Node " + nodeID + " defined multiple times (first at line " + string(rune(firstLine)) + ")",
					Line:     block.StartLine + i,
					Fixable:  false,
					Severity: model.SeverityError,
				})
			} else {
				nodes[nodeID] = block.StartLine + i
			}
		}
	}
	return issues
}

func checkUndefinedClass(issues []model.Issue, block model.MermaidBlock) []model.Issue {
	classDefs := make(map[string]bool)
	lines := strings.Split(block.Content, "\n")

	for _, line := range lines {
		matches := classDefRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			classDefs[matches[1]] = true
		}
	}

	for i, line := range lines {
		matches := classUseRegex.FindStringSubmatch(line)
		if len(matches) > 2 {
			className := matches[2]
			if !classDefs[className] {
				issues = append(issues, model.Issue{
					Type:     model.IssueTypeUndefinedClass,
					Message:  "Class " + className + " is used but not defined with classDef",
					Line:     block.StartLine + i,
					Fixable:  false,
					Severity: model.SeverityError,
				})
			}
		}
	}
	return issues
}

func checkInvalidStyle(issues []model.Issue, block model.MermaidBlock) []model.Issue {
	lines := strings.Split(block.Content, "\n")
	for i, line := range lines {
		matches := styleRegex.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) > 1 {
				prop := match[1]
				if !validStyleProps[prop] {
					issues = append(issues, model.Issue{
						Type:     model.IssueTypeInvalidStyle,
						Message:  "Invalid style property '" + prop + "', may be a typo",
						Line:     block.StartLine + i,
						Fixable:  false,
						Severity: model.SeverityError,
					})
				}
			}
		}
	}
	return issues
}

func checkUnquotedText(issues []model.Issue, block model.MermaidBlock) []model.Issue {
	lines := strings.Split(block.Content, "\n")
	for i, line := range lines {
		matches := nodeTextRegex.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) > 1 {
				text := match[1]
				if strings.Contains(text, `"`) {
					continue
				}
				// Unwrap wrapping parentheses from node shape syntax
				// Handles [(text)], ((text)), [[text]], etc. that come from different node shapes
				for {
					if len(text) >= 2 && (text[0] == '(' && text[len(text)-1] == ')') ||
						(text[0] == '[' && text[len(text)-1] == ']') {
						text = text[1 : len(text)-1]
					} else {
						break
					}
				}
				// After unwrapping shape syntax, check if actual content contains special characters
				// Only require quotes when there are matching parentheses or curly braces
				// Colons in text content don't require quotes in modern Mermaid
				if (strings.Contains(text, "(") && strings.Contains(text, ")")) ||
					(strings.Contains(text, "{") && strings.Contains(text, "}")) {
					issues = append(issues, model.Issue{
						Type:     model.IssueTypeUnquotedText,
						Message:  "Text contains special characters () : {} and is not quoted",
						Line:     block.StartLine + i,
						Fixable:  true,
						Fix:      fixer.FixUnquotedText,
						Severity: model.SeverityError,
					})
				}
			}
		}
	}
	return issues
}

func checkIsolatedNode(issues []model.Issue, block model.MermaidBlock) []model.Issue {
	nodes := make(map[string]bool)
	connectedNodes := make(map[string]bool)
	lines := strings.Split(block.Content, "\n")

	for _, line := range lines {
		matches := nodeIDRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			nodes[matches[1]] = true
		}

		// Find all connections in this line using regex
		connMatches := connectionRegex.FindAllStringSubmatch(line, -1)
		for _, match := range connMatches {
			if len(match) > 2 {
				// Both source and target nodes are connected
				if len(strings.TrimSpace(match[1])) > 0 {
					connectedNodes[match[1]] = true
				}
				if len(strings.TrimSpace(match[2])) > 0 {
					connectedNodes[match[2]] = true
				}
			}
		}
	}

	for node := range nodes {
		if !connectedNodes[node] && len(nodes) > 1 {
			var line int = 0
			for i, l := range lines {
				if strings.Contains(l, node+"[") {
					line = block.StartLine + i
					break
				}
			}
			issues = append(issues, model.Issue{
				Type:     model.IssueTypeIsolatedNode,
				Message:  "Node " + node + " is isolated (no connections)",
				Line:     line,
				Fixable:  false,
				Severity: model.SeverityWarning,
			})
		}
	}

	return issues
}

func checkDuplicateSubgraph(issues []model.Issue, block model.MermaidBlock) []model.Issue {
	subgraphs := make(map[string]int)
	lines := strings.Split(block.Content, "\n")

	for i, line := range lines {
		matches := subgraphRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			name := matches[1]
			if firstLine, exists := subgraphs[name]; exists {
				issues = append(issues, model.Issue{
					Type:     model.IssueTypeDuplicateSubgraph,
					Message:  "Subgraph " + name + " defined multiple times (first at line " + string(rune(firstLine)) + ")",
					Line:     block.StartLine + i,
					Fixable:  false,
					Severity: model.SeverityError,
				})
			} else {
				subgraphs[name] = block.StartLine + i
			}
		}
	}
	return issues
}
