// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package madmin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

// BandwidthDetails for the measured bandwidth
type BandwidthDetails struct {
	LimitInBytesPerSecond            int64   `json:"limitInBits"`
	CurrentBandwidthInBytesPerSecond float64 `json:"currentBandwidth"`
}

// BucketBandwidthReport captures the details for all buckets.
type BucketBandwidthReport struct {
	BucketStats map[string]BandwidthDetails `json:"bucketStats,omitempty"`
}

// Report includes the bandwidth report or the error encountered.
type Report struct {
	Report BucketBandwidthReport `json:"report"`
	Err    error                 `json:"error,omitempty"`
}

// GetBucketBandwidth - Gets a channel reporting bandwidth measurements for replication buckets. If no buckets
// generate replication traffic an empty map is returned in the report until traffic is seen.
func (adm *AdminClient) GetBucketBandwidth(ctx context.Context, buckets ...string) <-chan Report {
	queryValues := url.Values{}
	ch := make(chan Report)
	if len(buckets) > 0 {
		queryValues.Set("buckets", strings.Join(buckets, ","))
	}

	reqData := requestData{
		relPath:     adminAPIPrefix + "/bandwidth",
		queryValues: queryValues,
	}
	resp, err := adm.executeMethod(ctx, http.MethodGet, reqData)
	if err != nil {
		defer closeResponse(resp)
		ch <- Report{Err: err}
		return ch
	}
	if resp.StatusCode != http.StatusOK {
		ch <- Report{Err: httpRespToErrorResponse(resp)}
		return ch
	}

	dec := json.NewDecoder(resp.Body)

	go func(ctx context.Context, ch chan<- Report, resp *http.Response) {
		defer func() {
			closeResponse(resp)
			close(ch)
		}()
		for {
			var report BucketBandwidthReport

			if err = dec.Decode(&report); err != nil {
				ch <- Report{Err: err}
				return
			}
			select {
			case <-ctx.Done():
				return
			case ch <- Report{Report: report}:
			}
		}
	}(ctx, ch, resp)
	return ch
}
