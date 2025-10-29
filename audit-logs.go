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

// AuditLogOpts represents the options for the audit logs
type AuditLogOpts struct {
	Node     string        `json:"node,omitempty"`
	API      string        `json:"api,omitempty"`
	Bucket   string        `json:"bucket,omitempty"`
	Interval time.Duration `json:"interval,omitempty"`
}

// GetAuditLogs fetches the persisted audit logs from MinIO
func (adm AdminClient) GetAuditLogs(ctx context.Context, opts AuditLogOpts) iter.Seq2[log.Audit, error] {
	return func(yield func(log.Audit, error) bool) {
		auditOpts, err := json.Marshal(opts)
		if err != nil {
			yield(log.Audit{}, err)
			return
		}
		reqData := requestData{
			relPath: adminAPIPrefix + "/logs/audit",
			content: auditOpts,
		}
		resp, err := adm.executeMethod(ctx, http.MethodPost, reqData)
		if err != nil {
			yield(log.Audit{}, err)
			return
		}
		if resp.StatusCode != http.StatusOK {
			yield(log.Audit{}, httpRespToErrorResponse(resp))
			return
		}
		dec := msgp.NewReader(resp.Body)
		for {
			var info log.Audit
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
