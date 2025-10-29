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
	"iter"
	"net/http"
	"time"

	"github.com/openstor/madmin-go/v4/log"
	"github.com/tinylib/msgp/msgp"
)

// ErrorLogOpts represents the options for the ErrorLogs
type ErrorLogOpts struct {
	Node     string        `json:"node,omitempty"`
	API      string        `json:"api,omitempty"`
	Bucket   string        `json:"bucket,omitempty"`
	Prefix   string        `json:"prefix,omitempty"`
	Interval time.Duration `json:"interval,omitempty"`
}

// GetErrorLogs fetches the persisted error logs from MinIO
func (adm AdminClient) GetErrorLogs(ctx context.Context, opts ErrorLogOpts) iter.Seq2[log.Error, error] {
	return func(yield func(log.Error, error) bool) {
		errOpts, err := json.Marshal(opts)
		if err != nil {
			yield(log.Error{}, err)
			return
		}
		reqData := requestData{
			relPath: adminAPIPrefix + "/logs/error",
			content: errOpts,
		}
		resp, err := adm.executeMethod(ctx, http.MethodPost, reqData)
		if err != nil {
			yield(log.Error{}, err)
			return
		}
		if resp.StatusCode != http.StatusOK {
			yield(log.Error{}, httpRespToErrorResponse(resp))
			return
		}
		dec := msgp.NewReader(resp.Body)
		for {
			var info log.Error
			if err = info.DecodeMsg(dec); err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				continue
			}
			select {
			case <-ctx.Done():
				return
			default:
				yield(info, nil)
			}
		}
	}
}
