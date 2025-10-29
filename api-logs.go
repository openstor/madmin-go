// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

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

// APILogOpts represents the options for the APILogOpts
type APILogOpts struct {
	Node       string        `json:"node,omitempty"`
	API        string        `json:"api,omitempty"`
	Bucket     string        `json:"bucket,omitempty"`
	Prefix     string        `json:"prefix,omitempty"`
	StatusCode int           `json:"statusCode,omitempty"`
	Interval   time.Duration `json:"interval,omitempty"`
	Origin     log.Origin    `json:"origin,omitempty"`
	Type       log.APIType   `json:"type,omitempty"`
}

// GetAPILogs fetches the persisted API logs from MinIO
func (adm AdminClient) GetAPILogs(ctx context.Context, opts APILogOpts) iter.Seq2[log.API, error] {
	return func(yield func(log.API, error) bool) {
		apiOpts, err := json.Marshal(opts)
		if err != nil {
			yield(log.API{}, err)
			return
		}
		reqData := requestData{
			relPath: adminAPIPrefix + "/logs/api",
			content: apiOpts,
		}
		resp, err := adm.executeMethod(ctx, http.MethodPost, reqData)
		if err != nil {
			yield(log.API{}, err)
			return
		}
		if resp.StatusCode != http.StatusOK {
			yield(log.API{}, httpRespToErrorResponse(resp))
			return
		}
		dec := msgp.NewReader(resp.Body)
		for {
			var info log.API
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
