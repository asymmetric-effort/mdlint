// Copyright 2024 MdLint Authors

package formatter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/asymmetric-effort/mdlint/internal/findings"
)

// Format returns formatted findings.
func Format(fs []findings.Finding, format string) (string, error) {
	switch format {
	case "json":
		b, err := json.Marshal(fs)
		if err != nil {
			return "", err
		}
		return string(b), nil
	case "text":
		var sb strings.Builder
		for _, f := range fs {
			fmt.Fprintf(&sb, "%s:%d:%d %s %s\n", f.File, f.Line, f.Column, f.Rule, f.Message)
		}
		return sb.String(), nil
	default:
		return "", fmt.Errorf("unknown format %q", format)
	}
}
