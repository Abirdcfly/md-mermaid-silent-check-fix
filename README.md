# mermaid-lint

> [!WARNING]
> **Archived**: Because [mermaid-js/mermaid#7276](https://github.com/mermaid-js/mermaid/pull/7276) has been merged, `\n` is now natively supported in Mermaid diagrams and there is no need to convert `\n` to `<br>` anymore. The reason why many websites still display `\n` incorrectly is simply because they haven't updated their Mermaid version.

A Go static analysis tool for checking and fixing common issues in Mermaid diagrams embedded in Markdown files.

## Features

- Recursively scans directories for Markdown files
- Detects 8 types of common Mermaid issues
- Automatically fixes fixable issues
- Supports text and JSON output formats
- Can be integrated into CI/CD pipelines with strict mode

## Installation

```bash
go build -o mermaid-lint main.go
```

## Usage

```bash
./mermaid-lint <directory> [flags]

Flags:
  --fix       Automatically fix fixable issues
  --dry-run   Show changes without modifying files
  --json      Output results in JSON format
  --strict    Exit with non-zero code if any issues found
```

## Detected Issues

| Issue | Fixable | Description |
|-------|---------|-------------|
| newline_literal | ✓ | `\n` used instead of `<br>` for line breaks |
| unquoted_text | ✓ | Unquoted text containing special characters `():,` |
| html_literal | | HTML tags used in diagram |
| duplicate_node | | Duplicate node ID definitions |
| undefined_class | | Class used but not defined |
| invalid_style | | Invalid style property names |
| isolated_node | | Node with no connections |
| duplicate_subgraph | | Duplicate subgraph names |

## Example

```bash
# Check current directory
./mermaid-lint .

# Check and fix automatically
./mermaid-lint ./docs --fix

# CI usage - fail if issues found
./mermaid-lint ./docs --strict
```

## License

MIT
