// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package madmin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

// PoolDecommissionInfo currently draining information
type PoolDecommissionInfo struct {
	StartTime   time.Time `json:"startTime"`
	StartSize   int64     `json:"startSize"`
	TotalSize   int64     `json:"totalSize"`
	CurrentSize int64     `json:"currentSize"`
	Complete    bool      `json:"complete"`
	Failed      bool      `json:"failed"`
	Canceled    bool      `json:"canceled"`

	ObjectsDecommissioned     int64 `json:"objectsDecommissioned"`
	ObjectsDecommissionFailed int64 `json:"objectsDecommissionedFailed"`
	BytesDone                 int64 `json:"bytesDecommissioned"`
	BytesFailed               int64 `json:"bytesDecommissionedFailed"`
}

// PoolStatus captures current pool status
type PoolStatus struct {
	ID           int                   `json:"id"`
	CmdLine      string                `json:"cmdline"`
	LastUpdate   time.Time             `json:"lastUpdate"`
	Decommission *PoolDecommissionInfo `json:"decommissionInfo,omitempty"`
}

// DecommissionPool - starts moving data from specified pool to all other existing pools.
// Decommissioning if successfully started this function will return `nil`, to check
// for on-going draining cycle use StatusPool.
func (adm *AdminClient) DecommissionPool(ctx context.Context, pool string) error {
	values := url.Values{}
	values.Set("pool", pool)
	resp, err := adm.executeMethod(ctx, http.MethodPost, requestData{
		// POST <endpoint>/<admin-API>/pools/decommission?pool=http://server{1...4}/disk{1...4}
		relPath:     adminAPIPrefix + "/pools/decommission",
		queryValues: values,
	})
	if err != nil {
		return err
	}
	defer closeResponse(resp)
	if resp.StatusCode != http.StatusOK {
		return httpRespToErrorResponse(resp)
	}
	return nil
}

// CancelDecommissionPool - cancels an on-going decommissioning process,
// this automatically makes the pool available for writing once canceled.
func (adm *AdminClient) CancelDecommissionPool(ctx context.Context, pool string) error {
	values := url.Values{}
	values.Set("pool", pool)
	resp, err := adm.executeMethod(ctx, http.MethodPost, requestData{
		// POST <endpoint>/<admin-API>/pools/cancel?pool=http://server{1...4}/disk{1...4}
		relPath:     adminAPIPrefix + "/pools/cancel",
		queryValues: values,
	})
	if err != nil {
		return err
	}
	defer closeResponse(resp)
	if resp.StatusCode != http.StatusOK {
		return httpRespToErrorResponse(resp)
	}
	return nil
}

// StatusPool return current status about pool, reports any draining activity in progress
// and elapsed time.
func (adm *AdminClient) StatusPool(ctx context.Context, pool string) (PoolStatus, error) {
	values := url.Values{}
	values.Set("pool", pool)
	resp, err := adm.executeMethod(ctx, http.MethodGet, requestData{
		// GET <endpoint>/<admin-API>/pools/status?pool=http://server{1...4}/disk{1...4}
		relPath:     adminAPIPrefix + "/pools/status",
		queryValues: values,
	})
	if err != nil {
		return PoolStatus{}, err
	}
	defer closeResponse(resp)

	if resp.StatusCode != http.StatusOK {
		return PoolStatus{}, httpRespToErrorResponse(resp)
	}

	var info PoolStatus
	if err = json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return PoolStatus{}, err
	}

	return info, nil
}

// ListPoolsStatus returns list of pools currently configured and being used
// on the cluster.
func (adm *AdminClient) ListPoolsStatus(ctx context.Context) ([]PoolStatus, error) {
	resp, err := adm.executeMethod(ctx, http.MethodGet, requestData{
		relPath: adminAPIPrefix + "/pools/list", // GET <endpoint>/<admin-API>/pools/list
	})
	if err != nil {
		return nil, err
	}
	defer closeResponse(resp)
	if resp.StatusCode != http.StatusOK {
		return nil, httpRespToErrorResponse(resp)
	}
	var pools []PoolStatus
	if err = json.NewDecoder(resp.Body).Decode(&pools); err != nil {
		return nil, err
	}
	return pools, nil
}
