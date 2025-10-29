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
	// API requests are secure (HTTPS) if secure=true and insecure (HTTPS) otherwise.
	// NewAnonymousClient returns an anonymous MinIO Admin client object.
	// Anonymous client doesn't require any credentials
	madmAnonClnt, err := madmin.NewAnonymousClient("your-minio.example.com:9000", true)
	if err != nil {
		log.Fatalln(err)
	}
	// To enable trace :-
	// madmAnonClnt.TraceOn(os.Stderr)
	opts := madmin.HealthOpts{
		ClusterRead: false, // set to "true" to check if the cluster has read quorum
		Maintenance: false, // set to "true" to check if the cluster is taken down for maintenance
	}
	healthResult, err := madmAnonClnt.Healthy(context.Background(), opts)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(healthResult)
}
