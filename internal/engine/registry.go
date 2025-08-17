// (c) 2025 Asymmetric Effort, LLC
//
// Package engine provides a simple rule registry and finding types.

package engine

import (
	"sort"
	"sync"
)

// registry holds all registered rules keyed by their identifier.
var (
	registry   = map[string]Rule{}
	registryMu sync.RWMutex
)

// Register adds the given rule to the global registry. It panics if a rule with
// the same ID has already been registered. Registration is typically performed
// in the rule's init function.
func Register(rule Rule) {
	registryMu.Lock()
	defer registryMu.Unlock()
	id := rule.ID()
	if _, exists := registry[id]; exists {
		panic("rule already registered: " + id)
	}
	registry[id] = rule
}

// GetRule returns the rule with the given identifier if registered.
func GetRule(id string) (Rule, bool) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	r, ok := registry[id]
	return r, ok
}

// Rules returns all registered rules sorted by their identifiers. The returned
// slice is a copy and modifications to it do not affect the registry.
func Rules() []Rule {
	registryMu.RLock()
	defer registryMu.RUnlock()
	rules := make([]Rule, 0, len(registry))
	for _, r := range registry {
		rules = append(rules, r)
	}
	sort.Slice(rules, func(i, j int) bool { return rules[i].ID() < rules[j].ID() })
	return rules
}
