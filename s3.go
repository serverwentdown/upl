package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type credential struct {
	AccessKey string
	SecretKey string
	// Region is critical when signing requests.
	Region string
	// Endpoint is the base URL of the bucket, including the bucket name (in either the domain or path).
	//
	// Example:
	//   https://bucketname.s3.us-west-2.amazonaws.com
	//   http://my-minio.example.com/bucket-name
	Endpoint string
	// ACL is an optional canned ACL to set on objects
	ACL string
}

func (cred credential) validate() error {
	if strings.HasSuffix(cred.Endpoint, "/") {
		return fmt.Errorf("%w: endpoint should not end with slash", errBadRequest)
	}
	return nil
}

var httpClientS3 *http.Client

func setupS3() {
	httpClientS3 = &http.Client{
		Timeout: 10 * time.Second,
	}
}
