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
	"net/url"
)

// QuotaType represents bucket quota type
type QuotaType string

const (
	// HardQuota specifies a hard quota of usage for bucket
	HardQuota QuotaType = "hard"
)

// IsValid returns true if quota type is one of Hard
func (t QuotaType) IsValid() bool {
	return t == HardQuota
}

// BucketQuota holds bucket quota restrictions
type BucketQuota struct {
	Size     uint64    `json:"size"`     // Indicates maximum size allowed per bucket
	Rate     uint64    `json:"rate"`     // Indicates bandwidth rate allocated per bucket
	Requests uint64    `json:"requests"` // Indicates number of requests allocated per bucket
	Type     QuotaType `json:"quotatype,omitempty"`
}

// GetBucketQuota - get info on a user
func (adm *AdminClient) GetBucketQuota(ctx context.Context, bucket string) (q BucketQuota, err error) {
	queryValues := url.Values{}
	queryValues.Set("bucket", bucket)

	reqData := requestData{
		relPath:     adminAPIPrefix + "/get-bucket-quota",
		queryValues: queryValues,
	}

	// Execute GET on /minio/admin/v4/get-quota
	resp, err := adm.executeMethod(ctx, http.MethodGet, reqData)

	defer closeResponse(resp)
	if err != nil {
		return q, err
	}

	if resp.StatusCode != http.StatusOK {
		return q, httpRespToErrorResponse(resp)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return q, err
	}
	if err = json.Unmarshal(b, &q); err != nil {
		return q, err
	}

	return q, nil
}

// SetBucketQuota - sets a bucket's quota, if quota is set to '0'
// quota is disabled.
func (adm *AdminClient) SetBucketQuota(ctx context.Context, bucket string, quota *BucketQuota) error {
	data, err := json.Marshal(quota)
	if err != nil {
		return err
	}

	queryValues := url.Values{}
	queryValues.Set("bucket", bucket)

	reqData := requestData{
		relPath:     adminAPIPrefix + "/set-bucket-quota",
		queryValues: queryValues,
		content:     data,
	}

	// Execute PUT on /minio/admin/v4/set-bucket-quota to set quota for a bucket.
	resp, err := adm.executeMethod(ctx, http.MethodPut, reqData)

	defer closeResponse(resp)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return httpRespToErrorResponse(resp)
	}

	return nil
}
