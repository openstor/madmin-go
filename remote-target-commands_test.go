// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later
//

package madmin

import (
	"net/url"
	"testing"
)

func isOpsEqual(op1 []TargetUpdateType, op2 []TargetUpdateType) bool {
	if len(op1) != len(op2) {
		return false
	}
	for _, o1 := range op1 {
		found := false
		for _, o2 := range op2 {
			if o2 == o1 {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// TestGetTargetUpdateOps tests GetTargetUpdateOps
func TestGetTargetUpdateOps(t *testing.T) {
	testCases := []struct {
		values      url.Values
		expectedOps []TargetUpdateType
	}{
		{
			values: url.Values{
				"update": []string{"true"},
			},
			expectedOps: []TargetUpdateType{},
		},
		{
			values: url.Values{
				"update": []string{"false"},
				"path":   []string{"true"},
			},
			expectedOps: []TargetUpdateType{},
		},
		{
			values: url.Values{
				"update": []string{"true"},
				"path":   []string{""},
			},
			expectedOps: []TargetUpdateType{},
		},
		{
			values: url.Values{
				"update": []string{"true"},
				"path":   []string{"true"},
				"bzzzz":  []string{"true"},
			},
			expectedOps: []TargetUpdateType{PathUpdateType},
		},

		{
			values: url.Values{
				"update":      []string{"true"},
				"path":        []string{"true"},
				"creds":       []string{"true"},
				"sync":        []string{"true"},
				"proxy":       []string{"true"},
				"bandwidth":   []string{"true"},
				"healthcheck": []string{"true"},
			},
			expectedOps: []TargetUpdateType{
				PathUpdateType, CredentialsUpdateType, SyncUpdateType, ProxyUpdateType, BandwidthLimitUpdateType, HealthCheckDurationUpdateType,
			},
		},
	}
	for i, test := range testCases {
		gotOps := GetTargetUpdateOps(test.values)
		if !isOpsEqual(gotOps, test.expectedOps) {
			t.Fatalf("test %d: expected %v got %v", i+1, test.expectedOps, gotOps)
		}
	}
}
