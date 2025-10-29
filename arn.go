// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package madmin

import (
	"fmt"
	"strings"
)

// ARN is a struct to define arn.
type ARN struct {
	Type     ServiceType
	ID       string
	Region   string
	Resource string
	Bucket   string
}

// Empty returns true if arn struct is empty
func (a ARN) Empty() bool {
	return a.Type == ""
}

func (a ARN) String() string {
	if a.Bucket != "" {
		return fmt.Sprintf("arn:minio:%s:%s:%s:%s", a.Type, a.Region, a.ID, a.Bucket)
	}
	return fmt.Sprintf("arn:minio:%s:%s:%s:%s", a.Type, a.Region, a.ID, a.Resource)
}

// ParseARN return ARN struct from string in arn format.
func ParseARN(s string) (*ARN, error) {
	// ARN must be in the format of arn:minio:<Type>:<REGION>:<ID>:<remote-bucket/remote-resource>
	if !strings.HasPrefix(s, "arn:minio:") {
		return nil, fmt.Errorf("invalid ARN %s", s)
	}

	tokens := strings.Split(s, ":")
	if len(tokens) != 6 {
		return nil, fmt.Errorf("invalid ARN %s", s)
	}

	if tokens[4] == "" || tokens[5] == "" {
		return nil, fmt.Errorf("invalid ARN %s", s)
	}

	return &ARN{
		Type:     ServiceType(tokens[2]),
		Region:   tokens[3],
		ID:       tokens[4],
		Resource: tokens[5],
	}, nil
}

// ServiceType represents service type
type ServiceType string

const (
	// ReplicationService specifies replication service
	ReplicationService ServiceType = "replication"

	// NotificationService specifies notification/lambda service
	NotificationService ServiceType = "sqs"
)

// IsValid returns true if ARN type is set.
func (t1 ServiceType) IsValid(t2 ServiceType) bool {
	return t1 == t2
}
