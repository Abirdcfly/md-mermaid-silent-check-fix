package analyzer

import (
	"testing"

	"github.com/Abirdcfly/md-mermaid-silent-check-fix/model"
	"github.com/stretchr/testify/assert"
)

func TestCheckIsolatedNode_HyphenatedNodeIDs(t *testing.T) {
	// Test that hyphenated node IDs are detected correctly and connections are found
	content := `flowchart LR
    node-1[First Node] --> node-2[Second Node]
    node-2 --> node-3[Third Node]
`
	block := model.MermaidBlock{
		StartLine: 1,
		EndLine:   4,
		Content:   content,
	}

	var issues []model.Issue
	issues = checkIsolatedNode(issues, block)

	// No isolated nodes should be found because all are connected
	assert.Len(t, issues, 0)
}

func TestCheckIsolatedNode_DottedArrowWithLabel(t *testing.T) {
	// Test the original issue from the user that OBS is not incorrectly detected as isolated
	content := `flowchart LR
  MS["ModelServing CR"]
  subgraph Controller[ModelServing Controller]
    PG["Pod Generator<br/>GenerateEntry/Worker"]
    PM["Plugin Manager<br/>(run in order)"]
  end
  PL["Plugins<br/>BuiltIn/Webhook"]
  POD["Pod<br/>(mutated spec)"]
  OBS["Events / Status<br/>Annotations"]

  MS -->|reconcile| PG;
  MS -->|spec.plugins config| PM;
  PG -->|OnPodCreate| PM;
  PM -->|mutate Pod spec| POD;
  POD -.->|OnPodReady| PM;
  PM -.->|events/conditions| OBS;
`
	block := model.MermaidBlock{
		StartLine: 1,
		EndLine:   15,
		Content:   content,
	}

	var issues []model.Issue
	issues = checkIsolatedNode(issues, block)

	// PL is actually isolated (no connections) should be reported, OBS should not
	// Expect exactly 1 issue (PL only, not OBS)
	assert.Len(t, issues, 1)
	assert.Equal(t, model.IssueTypeIsolatedNode, issues[0].Type)
	assert.Contains(t, issues[0].Message, "PL")
	assert.Equal(t, model.SeverityWarning, issues[0].Severity)
}

func TestCheckUnquotedText_CylinderShape(t *testing.T) {
	// Test that OS[(Radix Tree State Snapshot)] doesn't get incorrectly flagged
	content := `flowchart LR
    OS[(Radix Tree State Snapshot)]
`
	block := model.MermaidBlock{
		StartLine: 1,
		EndLine:   3,
		Content:   content,
	}

	var issues []model.Issue
	issues = checkUnquotedText(issues, block)

	// After unwrapping, the text has no special characters, shouldn't be reported
	assert.Len(t, issues, 0)
}

func TestCheckUnquotedText_ColonInText(t *testing.T) {
	// Test that nodes with only colons in text don't need to be reported (colons don't need quotes in modern Mermaid)
	content := `flowchart LR
    NATS_PROM_EXP[nats-prom-exp :7777 /metrics] --> NATS_SERVER[nats-server :4222]
`
	block := model.MermaidBlock{
		StartLine: 1,
		EndLine:   3,
		Content:   content,
	}

	var issues []model.Issue
	issues = checkUnquotedText(issues, block)

	// Only parentheses and curly braces need quotes, colons don't
	assert.Len(t, issues, 0)
}

func TestCheckUnquotedText_ColonWithArrow(t *testing.T) {
	// Test the user's example that shouldn't be flagged: load[Load: CPU → GPU<br/>skip encoder]
	content := `flowchart LR
    cpu -- yes --> load[Load: CPU → GPU<br/>skip encoder]
`
	block := model.MermaidBlock{
		StartLine: 1,
		EndLine:   3,
		Content:   content,
	}

	var issues []model.Issue
	issues = checkUnquotedText(issues, block)

	// Doesn't need quotes because only has colon, no parentheses or curly braces
	assert.Len(t, issues, 0)
}

func TestCheckUnquotedText_ParenWithColon(t *testing.T) {
	// Test that if we have both parentheses and colon, it's still reported
	content := `flowchart LR
    A[Node with (parens) and: colon]
`
	block := model.MermaidBlock{
		StartLine: 1,
		EndLine:   3,
		Content:   content,
	}

	var issues []model.Issue
	issues = checkUnquotedText(issues, block)

	// Has parentheses so should still be reported
	assert.Len(t, issues, 1)
	assert.Equal(t, model.IssueTypeUnquotedText, issues[0].Type)
}

func TestCheckIsolatedNode_HyphenNodeConnection(t *testing.T) {
	// Test the original user example with colons and hyphen-ready node IDs
	content := `flowchart LR
    NATS_PROM_EXP[nats-prom-exp :7777 /metrics] -->|:8222/varz| NATS_SERVER[nats-server :4222, :6222, :8222]
`
	block := model.MermaidBlock{
		StartLine: 1,
		EndLine:   3,
		Content:   content,
	}

	var issues []model.Issue
	issues = checkIsolatedNode(issues, block)

	// Both nodes are connected - no isolated node issues
	assert.Len(t, issues, 0)
}
