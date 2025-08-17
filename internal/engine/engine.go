// Copyright 2024 MdLint Authors

package engine

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/asymmetric-effort/mdlint/internal/config"
	"github.com/asymmetric-effort/mdlint/internal/findings"
	"github.com/asymmetric-effort/mdlint/internal/sandbox"
)

// Engine executes lint rules.
// Engine executes lint rules using configured limits.
type Engine struct {
	Limits config.LimitsConfig
}

// Run processes files and returns findings.
func (e Engine) Run(paths []string) ([]findings.Finding, error) {
	restore := sandbox.DisableNetwork()
	defer restore()

	if e.Limits.MaxFiles > 0 && len(paths) > e.Limits.MaxFiles {
		return nil, fmt.Errorf("too many files: %d > %d", len(paths), e.Limits.MaxFiles)
	}
	conc := e.Limits.Concurrency
	if conc <= 0 {
		conc = 1
	}
	sem := make(chan struct{}, conc)
	var wg sync.WaitGroup
	result := []findings.Finding{}
	var mu sync.Mutex
	errCh := make(chan error, len(paths))

	for _, p := range paths {
		p := p
		sem <- struct{}{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			fi, err := os.Stat(p)
			if err != nil {
				errCh <- err
				return
			}
			if e.Limits.MaxFileSize > 0 && fi.Size() > e.Limits.MaxFileSize {
				errCh <- fmt.Errorf("file %s exceeds size limit", p)
				return
			}
			file, err := os.Open(p)
			if err != nil {
				errCh <- err
				return
			}
			scanner := bufio.NewScanner(file)
			line := 1
			for scanner.Scan() {
				text := scanner.Text()
				if idx := strings.Index(text, "TODO"); idx >= 0 {
					mu.Lock()
					result = append(result, findings.Finding{
						Rule:    "MD9000",
						Message: "TODO found",
						File:    p,
						Line:    line,
						Column:  idx + 1,
					})
					mu.Unlock()
				}
				line++
			}
			if err := scanner.Err(); err != nil {
				_ = file.Close()
				errCh <- err
				return
			}
			_ = file.Close()
		}()
	}

	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
