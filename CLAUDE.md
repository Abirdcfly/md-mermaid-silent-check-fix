# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

`mermaid-lint` is a Go-based static analysis tool that detects common issues in Mermaid diagrams embedded within Markdown files. It can automatically fix certain issues and integrates into CI/CD pipelines.

## Build Commands

```bash
# Run from source
go run main.go <directory> [--fix] [--dry-run] [--json] [--strict]

# Build executable
go build -o mermaid-lint main.go

# Check Go formatting
go fmt ./...

# Run Go vet
go vet ./...
```

## Test Commands

```bash
# Test on the included test file
./mermaid-lint .

# Test with JSON output
./mermaid-lint . --json

# Test fixes with dry run (no changes written)
./mermaid-lint . --fix --dry-run

# Test with strict mode (exits non-zero if issues found)
./mermaid-lint . --strict
```

## Command Line Flags

- `--fix`: Automatically fix fixable issues (newline_literal and unquoted_text)
- `--dry-run`: Show what would be fixed without modifying files
- `--json`: Output results in JSON format for machine processing
- `--strict`: Exit with non-zero status code if any issues are found (for CI)

## Architecture

The project follows a modular, clean architecture with separate packages for each responsibility:

- **main.go**: Entry point, parses flags, orchestrates the pipeline
- **scanner/**: Recursively scans directories for Markdown files, reads file content
- **parser/**: Extracts Mermaid code blocks from Markdown content
- **analyzer/**: Analyzes Mermaid blocks to detect issues
- **fixer/**: Applies automated fixes to fixable issues
- **reporter/**: Generates output reports (text or JSON format)
- **model/**: Shared data structures and issue type constants

## Issue Types Detected

| Issue Type | Description | Fixable |
|------------|-------------|---------|
| `newline_literal` | Uses `\n` instead of `<br>` for line breaks | ✓ Yes |
| `unquoted_text` | Text containing special characters (`()`, `:`, `{}`) not quoted | ✓ Yes |
| `html_literal` | HTML tags (like `<div>`) used in diagrams | No |
| `duplicate_node` | Same node ID defined multiple times | No |
| `undefined_class` | Class referenced but not defined with `classDef` | No |
| `invalid_style` | Invalid CSS style property names | No |
| `isolated_node` | Node with no connections to other nodes | No |
| `duplicate_subgraph` | Same subgraph name used multiple times | No |

## Dependencies

- Go 1.26+
- No external dependencies - pure Go standard library only
