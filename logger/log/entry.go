// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later
//

package log

import (
	"strings"
	"time"

	"github.com/openstor/madmin-go/v4"
)

//msgp:tag json
//go:generate msgp -d clearomitted -d "timezone utc" -file $GOFILE

// ObjectVersion object version key/versionId
type ObjectVersion struct {
	ObjectName string `json:"objectName"`
	VersionID  string `json:"versionId,omitempty"`
}

// Args - defines the arguments for the API.
type Args struct {
	Bucket    string            `json:"bucket,omitempty"`
	Object    string            `json:"object,omitempty"`
	VersionID string            `json:"versionId,omitempty"`
	Objects   []ObjectVersion   `json:"objects,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// Trace - defines the trace.
type Trace struct {
	Message   string                 `json:"message,omitempty"`
	Source    []string               `json:"source,omitempty"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// API - defines the api type and its args.
type API struct {
	Name string `json:"name,omitempty"`
	Args *Args  `json:"args,omitempty"`
}

//msgp:replace madmin.LogKind with:string

// Entry - defines fields and values of each log entry.
type Entry struct {
	Site         string         `json:"site,omitempty"`
	DeploymentID string         `json:"deploymentid,omitempty"`
	Level        madmin.LogKind `json:"level"`
	LogKind      madmin.LogKind `json:"errKind,omitempty"` // Deprecated Jan 2024
	Time         time.Time      `json:"time"`
	API          *API           `json:"api,omitempty"`
	RemoteHost   string         `json:"remotehost,omitempty"`
	Host         string         `json:"host,omitempty"` // Deprecated Apr 2025
	RequestHost  string         `json:"requestHost,omitempty"`
	RequestNode  string         `json:"requestNode,omitempty"`
	RequestID    string         `json:"requestID,omitempty"`
	UserAgent    string         `json:"userAgent,omitempty"`
	Message      string         `json:"message,omitempty"`
	Trace        *Trace         `json:"error,omitempty"`
}

// Info holds console log messages
type Info struct {
	Entry
	NodeName string `json:"node"`
	Err      error  `json:"-"`
}

// Mask returns the mask based on the error level.
func (l Info) Mask() uint64 {
	return l.Level.LogMask().Mask()
}

// SendLog returns true if log pertains to node specified in args.
func (l Info) SendLog(node string, logKind madmin.LogMask) bool {
	if logKind.Contains(l.Level.LogMask()) {
		return node == "" || strings.EqualFold(node, l.NodeName) && !l.Time.IsZero()
	}
	return false
}
