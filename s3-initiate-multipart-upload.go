package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"
)

type initiateMultipartUploadResult struct {
	Bucket   string
	Key      string
	UploadID string `xml:"UploadId"`
}

func createMultipartUpload(key string, cred credential) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	unsignedReq, err := http.NewRequestWithContext(ctx, http.MethodPost, cred.Endpoint+"/"+key+"?uploads", nil)
	if err != nil {
		return "", err
	}
	if cred.ACL != "" {
		unsignedReq.Header.Set("X-Amz-Acl", cred.ACL)
	}

	signedReq := sign(unsignedReq, cred)
	resp, err := httpClientS3.Do(signedReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("endpoint request failed: %d", resp.StatusCode)
	}

	initiateMultipartUploadResult := initiateMultipartUploadResult{}
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&initiateMultipartUploadResult)
	if err != nil {
		return "", err
	}

	return initiateMultipartUploadResult.UploadID, nil
}
