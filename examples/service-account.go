//
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
	"time"

	"github.com/openstor/madmin-go/v4"
)

func main() {
	// Note: YOUR-ACCESSKEYID, YOUR-SECRETACCESSKEY and my-bucketname are
	// dummy values, please replace them with original values.

	// API requests are secure (HTTPS) if secure=true and insecure (HTTP) otherwise.
	// New returns an MinIO Admin client object.
	madminClient, err := madmin.New("your-minio.example.com:9000", "YOUR-ACCESSKEYID", "YOUR-SECRETACCESSKEY", true)
	if err != nil {
		log.Fatalln(err)
	}
	ctx := context.Background()

	// add service account
	expiration := time.Now().Add(30 * time.Minute)
	addReq := madmin.AddServiceAccountReq{
		TargetUser: "my-username",
		Expiration: &expiration,
	}
	addRes, err := madminClient.AddServiceAccount(context.Background(), addReq)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(addRes)

	// update service account
	newExpiration := time.Now().Add(45 * time.Minute)
	updateReq := madmin.UpdateServiceAccountReq{
		NewStatus:     "my-status",
		NewExpiration: &newExpiration,
	}
	if err := madminClient.UpdateServiceAccount(ctx, "my-accesskey", updateReq); err != nil {
		log.Fatalln(err)
	}

	// get service account
	listRes, err := madminClient.ListServiceAccounts(ctx, "my-username")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(listRes)

	// delete service account
	if err := madminClient.DeleteServiceAccount(ctx, "my-accesskey"); err != nil {
		log.Fatalln(err)
	}
}
