// Copyright 2024 MdLint Authors

package engine

import (
	"bufio"
	"os"
	"strings"

	"github.com/asymmetric-effort/mdlint/internal/findings"
)

// Engine executes lint rules.
type Engine struct{}

// Run processes files and returns findings.
func (Engine) Run(paths []string) ([]findings.Finding, error) {
	var result []findings.Finding
	for _, p := range paths {
		file, err := os.Open(p)
		if err != nil {
			return nil, err
		}
		scanner := bufio.NewScanner(file)
		line := 1
		for scanner.Scan() {
			text := scanner.Text()
			if idx := strings.Index(text, "TODO"); idx >= 0 {
				result = append(result, findings.Finding{
					Rule:    "MD9000",
					Message: "TODO found",
					File:    p,
					Line:    line,
					Column:  idx + 1,
				})
			}
			line++
		}
		if err := scanner.Err(); err != nil {
			_ = file.Close()
			return nil, err
		}
		_ = file.Close()
	}
	return result, nil
}
