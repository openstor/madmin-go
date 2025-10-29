// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package madmin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

// LogMask is a bit mask for log types.
type LogMask uint64

const (
	LogMaskFatal LogMask = 1 << iota
	LogMaskWarning
	LogMaskError
	LogMaskEvent
	LogMaskInfo

	// LogMaskAll must be the last.
	LogMaskAll LogMask = (1 << iota) - 1
)

// Mask returns the LogMask as uint64
func (m LogMask) Mask() uint64 {
	return uint64(m)
}

// Contains returns whether all flags in other is present in t.
func (m LogMask) Contains(other LogMask) bool {
	return m&other == other
}

// LogKind specifies the kind of error log
type LogKind string

const (
	LogKindFatal   LogKind = "FATAL"
	LogKindWarning LogKind = "WARNING"
	LogKindError   LogKind = "ERROR"
	LogKindEvent   LogKind = "EVENT"
	LogKindInfo    LogKind = "INFO"
)

// LogMask returns the mask based on the kind.
func (l LogKind) LogMask() LogMask {
	switch l {
	case LogKindFatal:
		return LogMaskFatal
	case LogKindWarning:
		return LogMaskWarning
	case LogKindError:
		return LogMaskError
	case LogKindEvent:
		return LogMaskEvent
	case LogKindInfo:
		return LogMaskInfo
	}
	return LogMaskAll
}

func (l LogKind) String() string {
	return string(l)
}

// LogInfo holds console log messages
type LogInfo struct {
	logEntry
	NodeName string `json:"node"`
	Err      error  `json:"-"`
}

// GetLogs - listen on console log messages.
func (adm AdminClient) GetLogs(ctx context.Context, node string, lineCnt int, logKind string) <-chan LogInfo {
	logCh := make(chan LogInfo, 1)

	// Only success, start a routine to start reading line by line.
	go func(logCh chan<- LogInfo) {
		defer close(logCh)
		urlValues := make(url.Values)
		urlValues.Set("node", node)
		urlValues.Set("limit", strconv.Itoa(lineCnt))
		urlValues.Set("logType", logKind)
		for {
			reqData := requestData{
				relPath:     adminAPIPrefix + "/log",
				queryValues: urlValues,
			}
			// Execute GET to call log handler
			resp, err := adm.executeMethod(ctx, http.MethodGet, reqData)
			if err != nil {
				closeResponse(resp)
				return
			}

			if resp.StatusCode != http.StatusOK {
				logCh <- LogInfo{Err: httpRespToErrorResponse(resp)}
				return
			}
			dec := json.NewDecoder(resp.Body)
			for {
				var info LogInfo
				if err = dec.Decode(&info); err != nil {
					break
				}
				select {
				case <-ctx.Done():
					return
				case logCh <- info:
				}
			}

		}
	}(logCh)

	// Returns the log info channel, for caller to start reading from.
	return logCh
}

// Mask returns the mask based on the error level.
func (l LogInfo) Mask() uint64 {
	return l.LogKind.LogMask().Mask()
}
