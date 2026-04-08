package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Abirdcfly/md-mermaid-silent-check-fix/analyzer"
	"github.com/Abirdcfly/md-mermaid-silent-check-fix/fixer"
	"github.com/Abirdcfly/md-mermaid-silent-check-fix/model"
	"github.com/Abirdcfly/md-mermaid-silent-check-fix/parser"
	"github.com/Abirdcfly/md-mermaid-silent-check-fix/reporter"
	"github.com/Abirdcfly/md-mermaid-silent-check-fix/scanner"
)

var _ model.Issue = model.Issue{}

func main() {
	fixFlag := flag.Bool("fix", false, "Automatically fix fixable issues")
	dryRunFlag := flag.Bool("dry-run", false, "Only show changes, do not write to files")
	jsonFlag := flag.Bool("json", false, "Output JSON format")
	strictFlag := flag.Bool("strict", false, "Exit with non-zero code if any issues found")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Usage: mermaid-lint <directory> [--fix] [--dry-run] [--json] [--strict]")
		os.Exit(1)
	}

	rootDir := args[0]
	files, err := scanner.ScanDirectory(rootDir)
	if err != nil {
		fmt.Printf("Error scanning directory: %v\n", err)
		os.Exit(1)
	}

	totalFiles := len(files)
	if totalFiles == 0 {
		fmt.Println("No Markdown files found")
		os.Exit(0)
	}

	totalIssues := 0
	fixedIssues := 0

	for i := range files {
		blocks := parser.ExtractMermaidBlocks(files[i].Content)
		for j := range blocks {
			issues := analyzer.AnalyzeBlock(blocks[j])
			blocks[j].Issues = issues
			totalIssues += len(issues)
			for _, issue := range issues {
				if issue.Fixable {
					fixedIssues++
				}
			}
		}
		files[i].Blocks = blocks
	}

	if *jsonFlag {
		reporter.PrintJSONReport(files)
	} else {
		reporter.PrintTextReport(files)
		if *fixFlag && fixedIssues > 0 {
			fmt.Printf("\n🔧 Found %d fixable issues, applying fixes...\n", fixedIssues)
		}
	}

	if *fixFlag && !*dryRunFlag {
		for _, file := range files {
			updatedContent := file.Content
			for _, block := range file.Blocks {
				hasFixable := false
				for _, issue := range block.Issues {
					if issue.Fixable {
						hasFixable = true
						break
					}
				}
				if hasFixable {
					updatedContent = fixer.ApplyFixes(updatedContent, block)
				}
			}
			if updatedContent != file.Content {
				err := os.WriteFile(file.Path, []byte(updatedContent), 0644)
				if err != nil {
					fmt.Printf("Error writing file %s: %v\n", file.Path, err)
				}
			}
		}
	}

	if *strictFlag && totalIssues > 0 {
		os.Exit(1)
	}
}
