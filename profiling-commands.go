// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later
//

package madmin

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

// ProfilerType represents the profiler type
// passed to the profiler subsystem.
type ProfilerType string

// Different supported profiler types.
const (
	ProfilerCPU        ProfilerType = "cpu"        // represents CPU profiler type
	ProfilerCPUIO      ProfilerType = "cpuio"      // represents CPU with IO (fgprof) profiler type
	ProfilerMEM        ProfilerType = "mem"        // represents MEM profiler type
	ProfilerBlock      ProfilerType = "block"      // represents Block profiler type
	ProfilerMutex      ProfilerType = "mutex"      // represents Mutex profiler type
	ProfilerTrace      ProfilerType = "trace"      // represents Trace profiler type
	ProfilerThreads    ProfilerType = "threads"    // represents ThreadCreate profiler type
	ProfilerGoroutines ProfilerType = "goroutines" // represents Goroutine dumps.
	ProfilerRuntime    ProfilerType = "runtime"    // Include runtime metrics
)

// StartProfilingResult holds the result of starting
// profiler result in a given node.
type StartProfilingResult struct {
	NodeName string `json:"nodeName"`
	Success  bool   `json:"success"`
	Error    string `json:"error"`
}

// StartProfiling makes an admin call to remotely start profiling on a
// standalone server or the whole cluster in case of a distributed setup.
// NOTE: For simpler use cases use Profile() API.
func (adm *AdminClient) StartProfiling(ctx context.Context, profiler ProfilerType) ([]StartProfilingResult, error) {
	v := url.Values{}
	v.Set("profilerType", string(profiler))
	resp, err := adm.executeMethod(ctx,
		http.MethodPost, requestData{
			relPath:     adminAPIPrefix + "/profiling/start",
			queryValues: v,
		},
	)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, httpRespToErrorResponse(resp)
	}

	jsonResult, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var startResults []StartProfilingResult
	err = json.Unmarshal(jsonResult, &startResults)
	if err != nil {
		return nil, err
	}

	return startResults, nil
}

// DownloadProfilingData makes an admin call to download profiling data of a
// standalone server or of the whole cluster in case of a distributed setup.
// NOTE: For simpler use cases use Profile() API, must be
func (adm *AdminClient) DownloadProfilingData(ctx context.Context) (io.ReadCloser, error) {
	path := adminAPIPrefix + "/profiling/download"
	resp, err := adm.executeMethod(ctx,
		http.MethodGet, requestData{
			relPath: path,
		},
	)
	if err != nil {
		closeResponse(resp)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, httpRespToErrorResponse(resp)
	}

	if resp.Body == nil {
		return nil, errors.New("body is nil")
	}

	return resp.Body, nil
}

// Profile makes an admin call to remotely start profiling on a standalone
// server or the whole cluster in  case of a distributed setup for a specified duration.
func (adm *AdminClient) Profile(ctx context.Context, profiler ProfilerType, duration time.Duration) (io.ReadCloser, error) {
	v := url.Values{}
	v.Set("profilerType", string(profiler))
	v.Set("duration", duration.String())
	resp, err := adm.executeMethod(ctx,
		http.MethodPost, requestData{
			relPath:     adminAPIPrefix + "/profile",
			queryValues: v,
		},
	)
	if err != nil {
		closeResponse(resp)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, httpRespToErrorResponse(resp)
	}

	if resp.Body == nil {
		return nil, errors.New("body is nil")
	}
	return resp.Body, nil
}
