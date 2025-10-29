//go:build ignore

// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later
//

package main

import (
	"context"
	"log"

	"github.com/openstor/madmin-go/v4"
)

func main() {
	// Note: YOUR-ACCESSKEYID, YOUR-SECRETACCESSKEY and my-bucketname are
	// dummy values, please replace them with original values.

	// API requests are secure (HTTPS) if secure=true and insecure (HTTP) otherwise.
	// New returns an MinIO Admin client object.
	client, err := madmin.New("your-minio.example.com:9000", "YOUR-ACCESSKEYID", "YOUR-SECRETACCESSKEY", true)
	if err != nil {
		log.Fatalln(err)
	}
	// return IAM info as zipped content
	r, err := client.ExportIAM(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	// set IAM info from the reader
	if err := client.ImportIAM(context.Background(), r); err != nil {
		log.Fatalln(err)
	}
}
