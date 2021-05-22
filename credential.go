package main

import (
	"fmt"
	"strings"
)

/* types */

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
	// Prefix is a string to prepend to object keys
	Prefix string
}

func (cred credential) validate() error {
	if strings.HasSuffix(cred.Endpoint, "/") {
		return fmt.Errorf("%w: endpoint should not end with slash", errBadRequest)
	}
	if strings.HasPrefix(cred.Prefix, "/") {
		return fmt.Errorf("%w: prefix should not start with slash", errBadRequest)
	}
	return nil
}
