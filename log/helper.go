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

// stringifyMap sorts and joins a map[string]string as "k1=v1,k2=v2".
func stringifyMap(m map[string]string) string {
	if len(m) == 0 {
		return ""
	}
	pairs := make([]string, 0, len(m))
	for k, v := range m {
		pairs = append(pairs, fmt.Sprintf("%v=%v", k, v))
	}
	slices.Sort(pairs)
	return strings.Join(pairs, ",")
}

// stringifyInterfaceMap sorts and joins a map[string]interface{} as "k1=v1,k2=v2".
// (Uses fmt.Sprint to stringify values.)
func stringifyInterfaceMap(m map[string]interface{}) string {
	if len(m) == 0 {
		return ""
	}
	pairs := make([]string, 0, len(m))
	for k, v := range m {
		pairs = append(pairs, fmt.Sprintf("%v=%v", k, v))
	}
	slices.Sort(pairs)
	return strings.Join(pairs, ",")
}

func toString(key, value string) string {
	if value == "" {
		return ""
	}
	return fmt.Sprintf("%s=%s", key, value)
}

func toInt(key string, value int) string {
	if value == 0 {
		return ""
	}
	return fmt.Sprintf("%s=%d", key, value)
}

func toInt64(key string, value int64) string {
	if value == 0 {
		return ""
	}
	return fmt.Sprintf("%s=%d", key, value)
}

func toTime(key string, t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return fmt.Sprintf("%s=%s", key, t.UTC().Format(time.RFC3339Nano))
}

func toMap(key string, m map[string]string) string {
	if len(m) == 0 {
		return ""
	}
	return fmt.Sprintf("%s={%s}", key, stringifyMap(m))
}

func toInterfaceMap(key string, m map[string]interface{}) string {
	if len(m) == 0 {
		return ""
	}
	return fmt.Sprintf("%s={%s}", key, stringifyInterfaceMap(m))
}

// filterAndSort removes empty entries and sorts.
func filterAndSort(values []string) []string {
	out := values[:0]
	for _, v := range values {
		if v != "" {
			out = append(out, v)
		}
	}
	slices.Sort(out)
	return out
}
