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

var client *madmin.AdminClient

func main() {
	verifyTLS := true
	var err error
	client, err = madmin.New("your-minio.example.com:9000", "YOUR-ACCESSKEYID", "YOUR-SECRETACCESSKEY", verifyTLS)
	if err != nil {
		log.Fatalln(err)
	}

	getClusterSummary()
	getPools()
	getSinglePool()
	getErasureSetsForSinglePool()
	getDrivesForSinglePool()
}

func getClusterSummary() {
	resp, xerr := client.ClusterSummaryQuery(context.Background(), madmin.ClusterSummaryResourceOpts{})
	if xerr != nil {
		log.Fatalln(xerr)
	}

	fmt.Printf("%+v", resp)
}

func getPools() {
	resp, xerr := client.PoolsQuery(context.Background(), &madmin.PoolsResourceOpts{
		Offset: 0,
		Limit:  1000,
		Filter: "",
		Sort:   "PoolIndex",
	})
	if xerr != nil {
		log.Fatalln(xerr)
	}

	for _, v := range resp.Results {
		fmt.Printf("%+v\n", v)
	}
}

func getSinglePool() {
	resp, xerr := client.PoolsQuery(context.Background(), &madmin.PoolsResourceOpts{
		Offset: 0,
		Limit:  1,
		Filter: "PoolIndex = 1",
	})
	if xerr != nil {
		log.Fatalln(xerr)
	}

	for _, v := range resp.Results {
		fmt.Printf("%+v\n", v)
	}
}

func getErasureSetsForSinglePool() {
	resp, xerr := client.ErasureSetsQuery(context.Background(), &madmin.ErasureSetsResourceOpts{
		Offset:       0,
		Limit:        1000,
		Filter:       "PoolIndex = 1",
		Sort:         "SetIndex",
		SortReversed: false,
	})
	if xerr != nil {
		log.Fatalln(xerr)
	}

	for _, v := range resp.Results {
		fmt.Printf("%+v\n", v)
	}
}

func getDrivesForSinglePool() {
	resp, xerr := client.DrivesQuery(context.Background(), &madmin.DrivesResourceOpts{
		Offset:       0,
		Limit:        1000,
		Filter:       "PoolIndex = 1",
		Sort:         "SetIndex",
		SortReversed: false,
	})
	if xerr != nil {
		log.Fatalln(xerr)
	}

	for _, v := range resp.Results {
		fmt.Printf("%+v\n", v)
	}
}
