// Copyright (c) 2024

package findings

// Severity indicates rule violation severity.
type Severity string

const (
	Suggestion Severity = "suggestion"
	Warning    Severity = "warning"
	Error      Severity = "error"
)

var severityRank = map[Severity]int{
	Suggestion: 0,
	Warning:    1,
	Error:      2,
}

// String returns the string representation.
func (s Severity) String() string { return string(s) }

// AtLeast reports whether s is greater than or equal to threshold.
func (s Severity) AtLeast(threshold Severity) bool {
	return severityRank[s] >= severityRank[threshold]
}
