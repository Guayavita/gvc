package diag

import (
	"fmt"
	"strings"
)

type Severity string

type Position struct {
	File   string // filename
	Line   int
	Column int
}

type Span struct {
	Start Position
	End   Position
}

const (
	Error   Severity = "error"
	Warning Severity = "warning"
	Note    Severity = "note"
)

type Diagnostic struct {
	Severity Severity
	Message  string
	Span     Span
	Notes    []string // optional "help" or "hint" messages
}

func (d Diagnostic) Render(source string) string {
	lines := strings.Split(source, "\n")
	startLine := d.Span.Start.Line
	col := d.Span.Start.Column

	// Get the line of code
	codeLine := ""
	if startLine-1 < len(lines) {
		codeLine = lines[startLine-1]
	}
	underline := strings.Repeat(" ", col-1) + "^"
	return fmt.Sprintf(
		"%s: %s\n --> %s:%d:%d\n    |\n %2d | %s\n    | %s\n",
		d.Severity, d.Message,
		d.Span.Start.File, d.Span.Start.Line, d.Span.Start.Column,
		d.Span.Start.Line, codeLine, underline,
	)
}
