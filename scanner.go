// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later
//

package madmin

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

//go:generate msgp -d clearomitted -d "timezone utc" -file $GOFILE

// BucketScanInfo contains information of a bucket scan in a given pool/set
type BucketScanInfo struct {
	Pool        int         `msg:"pool"`
	Set         int         `msg:"set"`
	Cycle       uint64      `msg:"cycle"`
	Ongoing     bool        `msg:"ongoing"`
	LastUpdate  time.Time   `msg:"last_update"`
	LastStarted time.Time   `msg:"last_started"`
	LastError   string      `msg:"last_error"`
	Completed   []time.Time `msg:"completed,omitempty"`
}

// BucketScanInfo returns information of a bucket scan in all pools/sets
func (adm *AdminClient) BucketScanInfo(ctx context.Context, bucket string) ([]BucketScanInfo, error) {
	resp, err := adm.executeMethod(ctx,
		http.MethodGet,
		requestData{relPath: adminAPIPrefix + "/scanner/status/" + bucket})
	if err != nil {
		return nil, err
	}
	defer closeResponse(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, httpRespToErrorResponse(resp)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var info []BucketScanInfo
	err = json.Unmarshal(respBytes, &info)
	if err != nil {
		return nil, err
	}

	return info, nil
}
