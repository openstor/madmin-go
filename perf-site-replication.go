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

// SiteNetPerfNodeResult  - stats from each server
type SiteNetPerfNodeResult struct {
	Endpoint        string        `json:"endpoint"`
	TX              uint64        `json:"tx"` // transfer rate in bytes
	TXTotalDuration time.Duration `json:"txTotalDuration"`
	RX              uint64        `json:"rx"` // received rate in bytes
	RXTotalDuration time.Duration `json:"rxTotalDuration"`
	TotalConn       uint64        `json:"totalConn"`
	Error           string        `json:"error,omitempty"`
}

// SiteNetPerfResult  - aggregate results from all servers
type SiteNetPerfResult struct {
	NodeResults []SiteNetPerfNodeResult `json:"nodeResults"`
}

// SiteReplicationPerf - perform site-replication on the MinIO servers
func (adm *AdminClient) SiteReplicationPerf(ctx context.Context, duration time.Duration) (result SiteNetPerfResult, err error) {
	queryVals := make(url.Values)
	queryVals.Set("duration", duration.String())

	resp, err := adm.executeMethod(ctx,
		http.MethodPost, requestData{
			relPath:     adminAPIPrefix + "/speedtest/site",
			queryValues: queryVals,
		})
	if err != nil {
		return result, err
	}
	defer closeResponse(resp)
	if resp.StatusCode != http.StatusOK {
		return result, httpRespToErrorResponse(resp)
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}
