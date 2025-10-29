//go:build ignore

// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later
//

package main

import (
	"context"
	"log"
	"os"

	"github.com/openstor/madmin-go/v4"
)

func main() {
	// Note: YOUR-ACCESSKEYID, YOUR-SECRETACCESSKEY and my-bucketname are
	// dummy values, please replace them with original values.

	// API requests are secure (HTTPS) if secure=true and insecure (HTTP) otherwise.
	// New returns an MinIO Admin client object.
	madmClnt, err := madmin.New("your-minio.example.com:9000", "YOUR-ACCESSKEYID", "YOUR-SECRETACCESSKEY", true)
	if err != nil {
		log.Fatalln(err)
	}
	madmClnt.TraceOn(os.Stderr)

	job := `
replicate:
  flags:
    # (optional)
    name: "weekly-replication-job"
  
  target:
    type: "minio"
    bucket: "testbucket"
    endpoint: "https://play.min.io"
    credentials:
      accessKey: "minioadmin"
      secretKey: "minioadmin"
      sessionToken: ""
      
  source:
    type: "minio"
    bucket: "testbucket"
    prefix: ""
`
	if _, err = madmClnt.StartBatchJob(context.Background(), job); err != nil {
		log.Fatalln(err)
	}
}
