// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later
//

package madmin

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/openstor/madmin-go/v4/estream"
)

//msgp:replace TraceType with:uint64
//go:generate msgp -d clearomitted -d "timezone utc" $GOFILE

// HTTPFilter defines parameters for filtering traces based on incoming http request properties
type HTTPFilter struct {
	Func      string            `json:"funcFilter"`
	UserAgent string            `json:"userAgent"`
	Header    map[string]string `json:"header"`
}

// ServiceTelemetryOpts is a request to add following types to tracing.
type ServiceTelemetryOpts struct {
	// Types to add to tracing.
	Types TraceType `json:"types"`

	// Public cert to encrypt stream.
	PubCert []byte

	// Sample rate to set for this filter.
	// If <=0 or >=1 no sampling will be performed
	// and all hits will be traced.
	SampleRate float64 `json:"sampleRate"`

	// Disable sampling and only do tracing when a trace id is set on incoming request.
	ParentOnly bool `json:"parentOnly"`

	// Tag adds a `custom.tag` field to all traces triggered by this.
	TagKV map[string]string `json:"tags"`

	// On incoming HTTP types, only trigger if substring is in request.
	HTTPFilter HTTPFilter `json:"httpFilter"`
}

//msgp:ignore ServiceTelemetry

// ServiceTelemetry holds http telemetry spans, serialized and compressed.
type ServiceTelemetry struct {
	SpanMZ []byte // Serialized and Compressed spans.
	Err    error  // Any error that occurred
}

// ServiceTelemetryStream - gets raw stream for service telemetry.
func (adm AdminClient) ServiceTelemetryStream(ctx context.Context, opts ServiceTelemetryOpts) (io.ReadCloser, error) {
	bopts, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}
	reqData := requestData{
		relPath: adminAPIPrefix + "/telemetry",
		content: bopts,
	}
	// Execute GET to call trace handler
	resp, err := adm.executeMethod(ctx, http.MethodPost, reqData)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		closeResponse(resp)
		return nil, httpRespToErrorResponse(resp)
	}

	return resp.Body, nil
}

// ServiceTelemetry - perform trace request and return individual packages.
// If options contains a public key the private key must be provided.
// If context is canceled the function will return.
func (adm AdminClient) ServiceTelemetry(ctx context.Context, opts ServiceTelemetryOpts, dst chan<- ServiceTelemetry, pk *rsa.PrivateKey) {
	defer close(dst)
	resp, err := adm.ServiceTelemetryStream(ctx, opts)
	if err != nil {
		dst <- ServiceTelemetry{Err: err}
		return
	}
	dec, err := estream.NewReader(resp)
	if err != nil {
		dst <- ServiceTelemetry{Err: err}
		return
	}
	if pk != nil {
		dec.SetPrivateKey(pk)
	}
	for {
		st, err := dec.NextStream()
		if err != nil {
			dst <- ServiceTelemetry{Err: err}
			return
		}
		if ctx.Err() != nil {
			return
		}
		block, err := io.ReadAll(st)
		if err == nil && len(block) == 0 {
			// Ignore 0 sized blocks.
			continue
		}
		if ctx.Err() != nil {
			return
		}
		select {
		case <-ctx.Done():
			return
		case dst <- ServiceTelemetry{SpanMZ: block, Err: err}:
			if err != nil {
				return
			}
		}
	}
}

// ParseTraceType - given a slice of trace types, returns OR'd tracetypes counterpart
func ParseTraceType(typeSlice []string) TraceType {
	var traceType TraceType
	for _, t := range typeSlice {
		x := TraceAll
		v := TraceType(1)
		for x != 0 {
			if strings.EqualFold(v.String(), t) {
				traceType |= v
			}
			v <<= 1
			x >>= 1
		}
	}

	if traceType == 0 {
		traceType = TraceS3
	}
	return traceType
}

// ParseSampleRate converts a sample rate from string format to float64.
func ParseSampleRate(s string) (float64, error) {
	// Parse "x/y" entries.
	if strings.ContainsRune(s, '/') {
		split := strings.Split(s, "/")
		if len(split) != 2 {
			return 0, fmt.Errorf("invalid sample rate (%s)", s)
		}
		nom, err := strconv.ParseFloat(strings.TrimSpace(split[0]), 64)
		if err != nil {
			return 0, fmt.Errorf("invalid sample rate (%s)", s)
		}
		det, err := strconv.ParseFloat(strings.TrimSpace(split[1]), 64)
		if err != nil || det <= 0 {
			return 0, fmt.Errorf("invalid sample rate (%s)", s)
		}
		return nom / det, nil
	}
	// Parse percentage
	if strings.ContainsRune(s, '%') {
		split := strings.Split(strings.TrimSpace(s), "%")
		if len(split) != 2 {
			return 0, fmt.Errorf("invalid sample rate (%s)", s)
		}
		pct, err := strconv.ParseFloat(strings.TrimSpace(split[0]), 64)
		if err != nil || pct <= 0 {
			return 0, fmt.Errorf("invalid sample rate (%s)", s)
		}
		return pct / 100, nil
	}
	// Neither a fraction nor a percentage. So we treat s as a floating-point number.
	return strconv.ParseFloat(strings.TrimSpace(s), 64)
}
