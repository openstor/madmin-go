//
//go:build ignore

// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/openstor/madmin-go/v4"
	"github.com/openstor/openstor-go/v7/pkg/credentials"
)

func main() {
	// Note: YOUR-ACCESSKEYID, YOUR-SECRETACCESSKEY and my-bucketname are
	// dummy values, please replace them with original values.

	// API requests are secure (HTTPS) if secure=true and insecure (HTTP) otherwise.
	// New returns an MinIO Admin client object.
	madmClnt, err := madmin.NewWithOptions("localhost:9000", &madmin.Options{
		Creds:  credentials.NewStaticV4("minio", "minio123", ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}
	for event, err := range madmClnt.GetErrorLogs(context.Background(), madmin.ErrorLogOpts{}) {
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("Event: %+v\n", event)
	}
}
