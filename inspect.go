// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package madmin

import (
	"bufio"
	"context"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"net/url"
)

// InspectOptions provides options to Inspect.
type InspectOptions struct {
	Volume, File string
	PublicKey    []byte // PublicKey to use for inspected data.
}

// Inspect makes an admin call to download a raw files from disk.
// If inspect is called with a public key no key will be returned
// and the data is returned encrypted with the public key.
func (adm *AdminClient) Inspect(ctx context.Context, d InspectOptions) (key []byte, c io.ReadCloser, err error) {
	// Add form key/values in the body
	form := make(url.Values)
	form.Set("volume", d.Volume)
	form.Set("file", d.File)
	if d.PublicKey != nil {
		form.Set("public-key", base64.StdEncoding.EncodeToString(d.PublicKey))
	}

	method := ""
	reqData := requestData{
		relPath: adminAPIPrefix + "/inspect-data",
	}

	// If the public-key is specified, create a POST request and send
	// parameters as multipart-form instead of query values
	if d.PublicKey != nil {
		method = http.MethodPost
		reqData.customHeaders = make(http.Header)
		reqData.customHeaders.Set("Content-Type", "application/x-www-form-urlencoded")
		reqData.content = []byte(form.Encode())
	} else {
		method = http.MethodGet
		reqData.queryValues = form
	}

	resp, err := adm.executeMethod(ctx, method, reqData)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != http.StatusOK {
		closeResponse(resp)
		return nil, nil, httpRespToErrorResponse(resp)
	}

	bior := bufio.NewReaderSize(resp.Body, 4<<10)
	format, err := bior.ReadByte()
	if err != nil {
		closeResponse(resp)
		return nil, nil, err
	}

	switch format {
	case 1:
		key = make([]byte, 32)
		// Read key...
		_, err = io.ReadFull(bior, key[:])
		if err != nil {
			closeResponse(resp)
			return nil, nil, err
		}
	case 2:
		if err := bior.UnreadByte(); err != nil {
			return nil, nil, err
		}
	default:
		closeResponse(resp)
		return nil, nil, errors.New("unknown data version")
	}

	// Return body
	return key, &closeWrapper{Reader: bior, Closer: resp.Body}, nil
}

type closeWrapper struct {
	io.Reader
	io.Closer
}
