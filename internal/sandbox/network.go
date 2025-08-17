// Copyright 2025 The mdlint Authors
// SPDX-License-Identifier: MIT

package sandbox

import (
	"errors"
	"net/http"
)

// DisableNetwork prevents outbound HTTP requests by replacing the default transport.
// It returns a function that restores the previous transport.
func DisableNetwork() func() {
	prev := http.DefaultTransport
	http.DefaultTransport = roundTripperFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("network access disabled")
	})
	return func() { http.DefaultTransport = prev }
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
