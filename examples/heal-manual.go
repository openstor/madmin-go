//go:build ignore

// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later
//

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

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

	opts := madmin.HealOpts{
		Recursive: true,                  // recursively heal all objects at 'prefix'
		Remove:    true,                  // remove content that has lost quorum and not recoverable
		Recreate:  true,                  // rewrite all old non-inlined xl.meta to new xl.meta
		ScanMode:  madmin.HealNormalScan, // by default do not do 'deep' scanning
	}

	start, _, err := madmClnt.Heal(context.Background(), "healing-rewrite-bucket", "", opts, "", false, false)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Healstart sequence ===")
	enc := json.NewEncoder(os.Stdout)
	if err = enc.Encode(&start); err != nil {
		log.Fatalln(err)
	}

	fmt.Println()
	for {
		_, status, err := madmClnt.Heal(context.Background(), "healing-rewrite-bucket", "", opts, start.ClientToken, false, false)
		if status.Summary == "finished" {
			fmt.Println("Healstatus on items ===")
			for _, item := range status.Items {
				if err = enc.Encode(&item); err != nil {
					log.Fatalln(err)
				}
			}
			break
		}
		if status.Summary == "stopped" {
			fmt.Println("Healstatus on items ===")
			fmt.Println("Heal failed with", status.FailureDetail)
			break
		}

		for _, item := range status.Items {
			if err = enc.Encode(&item); err != nil {
				log.Fatalln(err)
			}
		}

		time.Sleep(time.Second)
	}
}
