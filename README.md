# Golang Admin Client API Reference [![Slack](https://slack.min.io/slack?type=svg)](https://slack.min.io)
The MinIO Admin Golang Client SDK provides APIs to manage MinIO services.

This document assumes that you have a working [Golang setup](https://golang.org/doc/install).

## Initialize MinIO Admin Client object.

##  MinIO

```go

package main

import (
    "fmt"

    "github.com/openstor/madmin-go/v4"
    "github.com/openstor/openstor-go/v7/pkg/credentials"
)

func main() {
    // Initialize minio client object.
    mdmClnt, err := madmin.NewWithOptions("your-minio.example.com:9000", &madmin.Options{
        Creds:  credentials.NewStaticV4("YOUR-ACCESSKEYID", "YOUR-SECRETKEY", ""),
        Secure: true,
    })
    if err != nil {
        fmt.Println(err)
        return
    }

    // Fetch service status.
	info, err := mdmClnt.ClusterInfo(context.Background())
    if err != nil {
        fmt.Println(err)
        return
    }
	fmt.Printf("%#v\n", info)
}
```

## Documentation
All documentation is available [here](https://pkg.go.dev/github.com/openstor/madmin-go/v4)

## License
This SDK is licensed under [GNU AGPLv3](https://github.com/minio/madmin-go/blob/master/LICENSE).
