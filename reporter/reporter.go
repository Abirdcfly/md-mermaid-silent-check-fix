package reporter

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Abirdcfly/md-mermaid-silent-check-fix/model"
)

type JSONReport struct {
	Files []FileReport `json:"files"`
	Total int          `json:"total_issues"`
}

type FileReport struct {
	File   string        `json:"file"`
	Issues []IssueReport `json:"issues"`
	Count  int           `json:"issue_count"`
}

type IssueReport struct {
	Type     string `json:"type"`
	Message  string `json:"message"`
	Line     int    `json:"line"`
	Fixable  bool   `json:"fixable"`
	Severity string `json:"severity"`
}

func PrintTextReport(files []model.MarkdownFile) int {
	totalIssues := 0
	for _, file := range files {
		fileIssues := 0
		for _, block := range file.Blocks {
			fileIssues += len(block.Issues)
		}
		totalIssues += fileIssues
		if fileIssues == 0 {
			continue
		}

		fmt.Printf("\nfile: %s\n", file.Path)
		for _, block := range file.Blocks {
			for _, issue := range block.Issues {
				fixMarker := "  "
				if issue.Fixable {
					fixMarker = "✅"
				}
				fmt.Printf("[%s] [%s] %s (line %d)\n  %s\n  fixable: %v\n",
					issue.Severity, issue.Type, fixMarker, issue.Line, issue.Message, issue.Fixable)
			}
		}
	}

	if totalIssues == 0 {
		fmt.Println("No issues found ✓")
	} else {
		fmt.Printf("\n✋ Found total %d issues\n", totalIssues)
	}

	return totalIssues
}

func PrintJSONReport(files []model.MarkdownFile) error {
	report := JSONReport{}
	total := 0

	for _, file := range files {
		fr := FileReport{
			File: file.Path,
		}
		for _, block := range file.Blocks {
			for _, issue := range block.Issues {
				fr.Issues = append(fr.Issues, IssueReport{
					Type:     issue.Type,
					Message:  issue.Message,
					Line:     issue.Line,
					Fixable:  issue.Fixable,
					Severity: issue.Severity,
				})
			}
		}
		fr.Count = len(fr.Issues)
		report.Files = append(report.Files, fr)
		total += fr.Count
	}
	report.Total = total

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}
