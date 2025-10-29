// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later
//

package madmin

// ClusterRegistrationReq - JSON payload of the subnet api for cluster registration
// Contains a registration token created by base64 encoding  of the registration info
type ClusterRegistrationReq struct {
	Token string `json:"token"`
}

// ClusterRegistrationInfo - Information stored in the cluster registration token
type ClusterRegistrationInfo struct {
	DeploymentID string `json:"deployment_id"`
	LicenseID    string `json:"license_id,omitempty"`
	ClusterName  string `json:"cluster_name"`
	UsedCapacity uint64 `json:"used_capacity"`
	// The "info" sub-node of the cluster registration information struct
	// Intended to be extensible i.e. more fields will be added as and when required
	Info struct {
		MinioVersion    string `json:"minio_version"`
		NoOfServerPools int    `json:"no_of_server_pools"`
		NoOfServers     int    `json:"no_of_servers"`
		NoOfDrives      int    `json:"no_of_drives"`
		NoOfBuckets     uint64 `json:"no_of_buckets"`
		NoOfObjects     uint64 `json:"no_of_objects"`
		TotalDriveSpace uint64 `json:"total_drive_space"`
		UsedDriveSpace  uint64 `json:"used_drive_space"`
		Edition         string `json:"edition"`
	} `json:"info"`
}

// SubnetLoginReq - JSON payload of the SUBNET login api
type SubnetLoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SubnetMFAReq - JSON payload of the SUBNET mfa api
type SubnetMFAReq struct {
	Username string `json:"username"`
	OTP      string `json:"otp"`
	Token    string `json:"token"`
}
