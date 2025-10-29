//go:build ignore

// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later
//

package main

import (
	"context"
	"fmt"
	"log"

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
	// Print a diff of unreplicated objects in a particular prefix in my-bucketname for the remote target specified.
	// Leaving out the ARN returns
	diffCh := madmClnt.BucketReplicationDiff(context.Background(), "my-bucketname", madmin.ReplDiffOpts{
		ARN:    "<remote-arn>",
		Prefix: "prefix/path",
	})
	for diff := range diffCh {
		if diff.Err != nil {
			log.Fatalln(diff.Err)
		}
		fmt.Println(diff)
	}
}
