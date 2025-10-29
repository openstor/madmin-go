//go:build ignore

// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later
//

package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/openstor/madmin-go/v4"
	"github.com/openstor/openstor-go/v7/pkg/credentials"
)

func main() {
	// Note: YOUR-ACCESSKEYID, YOUR-SECRETACCESSKEY are
	// dummy values, please replace them with original values.

	// API requests are secure (HTTPS) if secure=true and insecure (HTTP) otherwise.
	// New returns an MinIO Admin client object.
	madmClnt, err := madmin.NewWithOptions("your-minio.example.com:9000", &madmin.Options{
		Creds:  credentials.NewStaticV4("YOUR-ACCESSKEYID", "YOUR-SECRETACCESSKEY", ""),
		Secure: true,
	})
	if err != nil {
		log.Fatalln(err)
	}

	results, err := madmClnt.Speedtest(context.Background(), madmin.SpeedtestOpts{Autotune: true})
	if err != nil {
		log.Fatalln(err)
	}
	for result := range results {
		js, _ := json.MarshalIndent(result, "", "  ")
		log.Printf("Speedtest Result: %s\n", string(js))
	}
}
