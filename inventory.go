// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package madmin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// InventoryJobControlStatus represents the status of an inventory job control operation
type InventoryJobControlStatus string

const (
	InventoryJobStatusCanceled  InventoryJobControlStatus = "canceled"
	InventoryJobStatusSuspended InventoryJobControlStatus = "suspended"
	InventoryJobStatusResumed   InventoryJobControlStatus = "resumed"
)

// InventoryJobControlResponse represents the response from inventory job control operations
type InventoryJobControlResponse struct {
	Status      InventoryJobControlStatus `json:"status"`
	Bucket      string                    `json:"bucket"`
	InventoryID string                    `json:"inventoryId"`
	Message     string                    `json:"message,omitempty"`
}

// CancelInventoryJob cancels a running inventory job
func (adm *AdminClient) CancelInventoryJob(ctx context.Context, bucket, inventoryID string) (*InventoryJobControlResponse, error) {
	return adm.controlInventoryJob(ctx, bucket, inventoryID, "cancel")
}

// SuspendInventoryJob suspends a running inventory job
func (adm *AdminClient) SuspendInventoryJob(ctx context.Context, bucket, inventoryID string) (*InventoryJobControlResponse, error) {
	return adm.controlInventoryJob(ctx, bucket, inventoryID, "suspend")
}

// ResumeInventoryJob resumes a suspended inventory job
func (adm *AdminClient) ResumeInventoryJob(ctx context.Context, bucket, inventoryID string) (*InventoryJobControlResponse, error) {
	return adm.controlInventoryJob(ctx, bucket, inventoryID, "resume")
}

// controlInventoryJob is the common function for all inventory job control operations
func (adm *AdminClient) controlInventoryJob(ctx context.Context, bucket, inventoryID, action string) (*InventoryJobControlResponse, error) {
	if bucket == "" {
		return nil, fmt.Errorf("bucket name cannot be empty")
	}
	if inventoryID == "" {
		return nil, fmt.Errorf("inventory ID cannot be empty")
	}

	resp, err := adm.executeMethod(ctx, http.MethodPost,
		requestData{
			relPath: adminAPIPrefix + fmt.Sprintf("/inventory/%s/%s/%s", bucket, inventoryID, action),
		},
	)
	if err != nil {
		return nil, err
	}
	defer closeResponse(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, httpRespToErrorResponse(resp)
	}

	// Decode the response
	var result InventoryJobControlResponse
	dec := json.NewDecoder(resp.Body)
	if err = dec.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
