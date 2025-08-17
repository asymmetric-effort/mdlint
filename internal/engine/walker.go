//go:build ignore
// +build ignore

// Copyright 2024

package engine

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sort"

	"golang.org/x/sync/errgroup"
)

// Config controls how files are discovered during a run.
type Config struct {
	// Include defines glob patterns of files to include. If empty, all files are
	// included.
	Include []string
	// Exclude defines glob patterns of files or directories to exclude.
	Exclude []string
	// Workers controls the number of concurrent workers used to process files.
	// If zero, runtime.NumCPU is used.
	Workers int
}

// Run walks the file tree rooted at root, applying all registered rules to
// matching files. Findings are returned in a deterministic order.
func Run(ctx context.Context, root string, cfg Config) ([]Finding, error) {
	if root == "" {
		return nil, errors.New("root must not be empty")
	}

	paths := make(chan string)
	g, ctx := errgroup.WithContext(ctx)

	// Walker goroutine.
	g.Go(func() error {
		defer close(paths)
		return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			if d.IsDir() {
				if matchPattern(rel, cfg.Exclude) {
					return filepath.SkipDir
				}
				return nil
			}
			if !shouldInclude(rel, cfg) {
				return nil
			}
			select {
			case paths <- path:
			case <-ctx.Done():
				return ctx.Err()
			}
			return nil
		})
	})

	workerCount := cfg.Workers
	if workerCount <= 0 {
		workerCount = runtime.NumCPU()
	}

	findingsCh := make(chan Finding)

	for i := 0; i < workerCount; i++ {
		g.Go(func() error {
			for path := range paths {
				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				fctx := &Context{FilePath: path}
				for _, r := range Rules() {
					for _, f := range r.Apply(content, fctx) {
						if f.Location == "" {
							f.Location = path
						}
						select {
						case findingsCh <- f:
						case <-ctx.Done():
							return ctx.Err()
						}
					}
				}
			}
			return nil
		})
	}

	var findings []Finding
	collectDone := make(chan struct{})
	go func() {
		for f := range findingsCh {
			findings = append(findings, f)
		}
		close(collectDone)
	}()

	if err := g.Wait(); err != nil {
		close(findingsCh)
		<-collectDone
		return nil, err
	}
	close(findingsCh)
	<-collectDone

	sort.Slice(findings, func(i, j int) bool {
		if findings[i].Location == findings[j].Location {
			return findings[i].RuleID < findings[j].RuleID
		}
		return findings[i].Location < findings[j].Location
	})

	return findings, nil
}

// shouldInclude reports whether the given relative path should be processed.
func shouldInclude(rel string, cfg Config) bool {
	if len(cfg.Include) > 0 && !matchPattern(rel, cfg.Include) {
		return false
	}
	if matchPattern(rel, cfg.Exclude) {
		return false
	}
	return true
}

// matchPattern reports whether name matches any of the provided glob patterns.
func matchPattern(name string, patterns []string) bool {
	for _, p := range patterns {
		if ok, _ := filepath.Match(p, name); ok {
			return true
		}
	}
	return false
}
