// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later
//

package log

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

//msgp:tag json
//go:generate msgp -d clearomitted -d "timezone utc" $GOFILE

// Error represents the error event
type Error struct {
	Version string            `json:"version"`
	Node    string            `json:"node"`
	Time    time.Time         `json:"time"`
	Message string            `json:"message"`
	API     string            `json:"apiName"`
	Trace   *Trace            `json:"trace,omitempty"`
	Tags    map[string]string `json:"tags,omitempty"`
}

// Trace represents the call trace
type Trace struct {
	Source    []string          `json:"source,omitempty"`
	Variables map[string]string `json:"variables,omitempty"`
}

// GetTagValByKey gets the tag value by key
func (e Error) GetTagValByKey(key string) string {
	return e.Tags[key]
}

// String returns the canonical string for Error
func (e Error) String() string {
	values := []string{
		toString("version", e.Version),
		toString("node", e.Node),
		toTime("time", e.Time),
		toString("message", e.Message),
		toString("apiName", e.API),
		toMap("tags", e.Tags),
	}
	if e.Trace != nil {
		values = append(values, fmt.Sprintf("trace={%s}", e.Trace.String()))
	}
	values = filterAndSort(values)
	return strings.Join(values, ",")
}

// String returns the canonical string for Trace
func (t Trace) String() string {
	values := []string{
		toMap("variables", t.Variables),
	}
	if len(t.Source) > 0 {
		src := slices.Clone(t.Source)
		slices.Sort(src)
		values = append(values, fmt.Sprintf("source=[%s]", strings.Join(src, ",")))
	}
	values = filterAndSort(values)
	return strings.Join(values, ",")
}
