//
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

	"github.com/openstor/madmin-go/v4"
	"github.com/openstor/openstor-go/v7/pkg/credentials"
)

func main() {
	c, err := madmin.NewWithOptions("127.0.0.1:9001", &madmin.Options{
		Creds:     credentials.NewStaticV4("minio", "minio123", ""),
		Secure:    false,
		Transport: nil,
	})
	fatalErr(err)

	stats, err := c.ClusterAPIStats(context.Background())
	fatalErr(err)

	b, err := json.MarshalIndent(stats, "", "  ")
	fatalErr(err)
	fmt.Println(string(b))
}

func fatalErr(err error) {
	if err != nil {
		panic(err)
	}
}
