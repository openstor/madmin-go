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
	var kiB uint64 = 1 << 10
	ctx := context.Background()
	quota := &madmin.BucketQuota{
		Size: 32 * kiB,
		Type: madmin.HardQuota,
	}
	// set bucket quota config
	if err := madmClnt.SetBucketQuota(ctx, "bucket-name", quota); err != nil {
		log.Fatalln(err)
	}
	// gets bucket quota config
	quotaCfg, err := madmClnt.GetBucketQuota(ctx, "bucket-name")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(quotaCfg)
}
