// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later
//

package madmin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

// NetperfNodeResult - stats from each server
type NetperfNodeResult struct {
	Endpoint string `json:"endpoint"`
	TX       uint64 `json:"tx"`
	RX       uint64 `json:"rx"`
	Error    string `json:"error,omitempty"`
}

// NetperfResult - aggregate results from all servers
type NetperfResult struct {
	NodeResults []NetperfNodeResult `json:"nodeResults"`
}

// Netperf - perform netperf on the MinIO servers
func (adm *AdminClient) Netperf(ctx context.Context, duration time.Duration) (result NetperfResult, err error) {
	queryVals := make(url.Values)
	queryVals.Set("duration", duration.String())

	resp, err := adm.executeMethod(ctx,
		http.MethodPost, requestData{
			relPath:     adminAPIPrefix + "/speedtest/net",
			queryValues: queryVals,
		})
	if err != nil {
		return result, err
	}
	if resp.StatusCode != http.StatusOK {
		return result, httpRespToErrorResponse(resp)
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}
