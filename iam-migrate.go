// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package madmin

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

// ImportIAMResult - represents the structure iam import response
type ImportIAMResult struct {
	// Skipped entries while import
	// This could be due to groups, policies etc missing for
	// impoprted entries. We dont fail hard in this case and
	// skip those entries
	Skipped IAMEntities `json:"skipped,omitempty"`

	// Removed entries - this mostly happens for policies
	// where empty might be getting imported and that's invalid
	Removed IAMEntities `json:"removed,omitempty"`

	// Newly added entries
	Added IAMEntities `json:"added,omitempty"`

	// Failed entries while import. This would have details of
	// failed entities with respective errors
	Failed IAMErrEntities `json:"failed,omitempty"`
}

// IAMEntities - represents different IAM entities
type IAMEntities struct {
	// List of policy names
	Policies []string `json:"policies,omitempty"`
	// List of user names
	Users []string `json:"users,omitempty"`
	// List of group names
	Groups []string `json:"groups,omitempty"`
	// List of Service Account names
	ServiceAccounts []string `json:"serviceAccounts,omitempty"`
	// List of user policies, each entry in map represents list of policies
	// applicable to the user
	UserPolicies []map[string][]string `json:"userPolicies,omitempty"`
	// List of group policies, each entry in map represents list of policies
	// applicable to the group
	GroupPolicies []map[string][]string `json:"groupPolicies,omitempty"`
	// List of STS policies, each entry in map represents list of policies
	// applicable to the STS
	STSPolicies []map[string][]string `json:"stsPolicies,omitempty"`
}

// IAMErrEntities - represents errored out IAM entries while import with error
type IAMErrEntities struct {
	// List of errored out policies with errors
	Policies []IAMErrEntity `json:"policies,omitempty"`
	// List of errored out users with errors
	Users []IAMErrEntity `json:"users,omitempty"`
	// List of errored out groups with errors
	Groups []IAMErrEntity `json:"groups,omitempty"`
	// List of errored out service accounts with errors
	ServiceAccounts []IAMErrEntity `json:"serviceAccounts,omitempty"`
	// List of errored out user policies with errors
	UserPolicies []IAMErrPolicyEntity `json:"userPolicies,omitempty"`
	// List of errored out group policies with errors
	GroupPolicies []IAMErrPolicyEntity `json:"groupPolicies,omitempty"`
	// List of errored out STS policies with errors
	STSPolicies []IAMErrPolicyEntity `json:"stsPolicies,omitempty"`
}

// IAMErrEntity - represents errored out IAM entity
type IAMErrEntity struct {
	// Name of the errored IAM entity
	Name string `json:"name,omitempty"`
	// Actual error
	Error error `json:"error,omitempty"`
}

// IAMErrPolicyEntity - represents errored out IAM policies
type IAMErrPolicyEntity struct {
	// Name of entity (user, group, STS)
	Name string `json:"name,omitempty"`
	// List of policies
	Policies []string `json:"policies,omitempty"`
	// Actual error
	Error error `json:"error,omitempty"`
}

// ExportIAM makes an admin call to export IAM data
func (adm *AdminClient) ExportIAM(ctx context.Context) (io.ReadCloser, error) {
	path := adminAPIPrefix + "/export-iam"

	resp, err := adm.executeMethod(ctx,
		http.MethodGet, requestData{
			relPath: path,
		},
	)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		closeResponse(resp)
		return nil, httpRespToErrorResponse(resp)
	}
	return resp.Body, nil
}

// ImportIAM makes an admin call to setup IAM  from imported content
func (adm *AdminClient) ImportIAM(ctx context.Context, contentReader io.ReadCloser) error {
	content, err := io.ReadAll(contentReader)
	if err != nil {
		return err
	}

	path := adminAPIPrefix + "/import-iam"
	resp, err := adm.executeMethod(ctx,
		http.MethodPut, requestData{
			relPath: path,
			content: content,
		},
	)
	defer closeResponse(resp)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return httpRespToErrorResponse(resp)
	}

	return nil
}

// ImportIAMV2 makes an admin call to setup IAM  from imported content
func (adm *AdminClient) ImportIAMV2(ctx context.Context, contentReader io.ReadCloser) (iamr ImportIAMResult, err error) {
	content, err := io.ReadAll(contentReader)
	if err != nil {
		return iamr, err
	}

	path := adminAPIPrefix + "/import-iam-v2"
	resp, err := adm.executeMethod(ctx,
		http.MethodPut, requestData{
			relPath: path,
			content: content,
		},
	)
	defer closeResponse(resp)
	if err != nil {
		return iamr, err
	}

	if resp.StatusCode != http.StatusOK {
		return iamr, httpRespToErrorResponse(resp)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return iamr, err
	}

	if err = json.Unmarshal(b, &iamr); err != nil {
		return iamr, err
	}

	return iamr, nil
}
