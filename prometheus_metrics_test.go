// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later
//

package madmin

import (
	"strings"
	"testing"

	"github.com/prometheus/prom2json"
)

func TestParsePrometheusResultsReturnsPrometheusObjectsFromStringReader(t *testing.T) {
	prometheusResults := `# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
		# TYPE go_gc_duration_seconds summary
		go_gc_duration_seconds_sum 0.248349766
		go_gc_duration_seconds_count 397
	`
	myReader := strings.NewReader(prometheusResults)
	results, err := ParsePrometheusResults(myReader)
	if err != nil {
		t.Errorf("error not expected, got: %v", err)
	}

	expectedResults := []*prom2json.Family{
		{
			Name: "go_gc_duration_seconds",
			Type: "SUMMARY",
			Help: "A summary of the pause duration of garbage collection cycles.",
			Metrics: []interface{}{
				prom2json.Summary{}, // We just verify length, not content
			},
		},
	}

	if len(results) != len(expectedResults) {
		t.Errorf("len(results): %d  not equal to len(expectedResults): %d", len(results), len(expectedResults))
	}

	for i, result := range results {
		if result.Name != expectedResults[i].Name {
			t.Errorf("result.Name: %v  not equal to expectedResults[i].Name: %v", result.Name, expectedResults[i].Name)
		}
		if len(result.Metrics) != len(expectedResults[i].Metrics) {
			t.Errorf("len(result.Metrics): %d  not equal to len(expectedResults[i].Metrics): %d", len(result.Metrics), len(expectedResults[i].Metrics))
		}
	}
}
