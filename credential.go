package main

import (
	"fmt"
	"net/url"
	"strings"
)

/* types */

type credential struct {
	AccessKey string
	SecretKey string
	// Region is critical when signing requests.
	Region string
	// Endpoint is the base URL of the bucket, including the bucket name (in
	// either the domain or path).
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

func newCredential(endpoint, region, accessKey, secretKey, prefix, acl string) credential {
	parsedEndpoint, _ := url.Parse(endpoint)
	return credential{
		Endpoint:  parsedEndpoint.String(),
		Region:    region,
		AccessKey: accessKey,
		SecretKey: secretKey,
		Prefix:    prefix,
		ACL:       acl,
	}
}

func (cred credential) validate() error {
	parsedEndpoint, err := url.Parse(cred.Endpoint)
	if err != nil {
		return fmt.Errorf("%w: endpoint must be a URL and not empty", errBadRequest)
	} else if parsedEndpoint.Host == "" {
		return fmt.Errorf("%w: endpoint must have a valid host", errBadRequest)
	} else if parsedEndpoint.User != nil {
		return fmt.Errorf("%w: endpoint must not have user credentials", errBadRequest)
	} else if parsedEndpoint.RawQuery != "" {
		return fmt.Errorf("%w: endpoint must not have query parameters", errBadRequest)
	} else if parsedEndpoint.RawFragment != "" {
		return fmt.Errorf("%w: endpoint must not have fragment", errBadRequest)
	} else if parsedEndpoint.Scheme != "http" && parsedEndpoint.Scheme != "https" {
		return fmt.Errorf("%w: endpoint must be http(s)", errBadRequest)
	}

	if cred.Region == "" {
		return fmt.Errorf("%w: region must not be empty", errBadRequest)
	}

	if strings.HasSuffix(cred.Endpoint, "/") {
		return fmt.Errorf("%w: endpoint should not end with slash", errBadRequest)
	}

	if strings.HasPrefix(cred.Prefix, "/") {
		return fmt.Errorf("%w: prefix should not start with slash", errBadRequest)
	}

	return nil
}
