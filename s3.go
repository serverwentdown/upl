package main

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/minio/minio-go/v7/pkg/signer"
)

var httpClientS3 *http.Client

func setupS3() {
	httpClientS3 = &http.Client{
		Timeout: 10 * time.Second,
	}
}

/* signing */

func preSign(req *http.Request, cred credential) *http.Request {
	signedReq := signer.PreSignV4(
		*req,
		cred.AccessKey, cred.SecretKey, "",
		cred.Region,
		60*60, // seconds
	)
	return signedReq
}

func sign(req *http.Request, cred credential) *http.Request {
	req.Header.Set("X-Amz-Content-Sha256", "UNSIGNED-PAYLOAD")
	signedReq := signer.SignV4(
		*req,
		cred.AccessKey, cred.SecretKey, "",
		cred.Region,
	)
	return signedReq
}

/* helpers */

func stripETag(t string) string {
	return strings.TrimSuffix(strings.TrimPrefix(t, "\""), "\"")
}

/* initiateMultipartUpload */

type initiateMultipartUploadResult struct {
	Bucket   string
	Key      string
	UploadID string `xml:"UploadId"`
}

func initiateMultipartUpload(key string, cred credential) (initiateMultipartUploadResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	params := make(url.Values)
	params.Set("uploads", "")
	unsignedReq, err := http.NewRequestWithContext(ctx, http.MethodPost, cred.Endpoint+"/"+key+"?"+params.Encode(), nil)
	if err != nil {
		return initiateMultipartUploadResult{}, err
	}
	if cred.ACL != "" {
		unsignedReq.Header.Set("X-Amz-Acl", cred.ACL)
	}

	signedReq := sign(unsignedReq, cred)
	resp, err := httpClientS3.Do(signedReq)
	if err != nil {
		return initiateMultipartUploadResult{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return initiateMultipartUploadResult{}, fmt.Errorf("endpoint request failed: %d: %s", resp.StatusCode, body)
	}

	result := initiateMultipartUploadResult{}
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

/* listParts */

type part struct {
	XMLName    xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ Part" json:"-"`
	PartNumber uint16   `json:"PartNumber"`
	ETag       string   `json:"ETag"`
	Size       uint32   `json:"Size"`
}

type listPartsResult struct {
	XMLName  xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ ListPartsResult" json:"-"`
	Bucket   string
	Key      string
	UploadID string `xml:"UploadId"`
	// not implemented: Initiator
	// not implemented: Owner
	// not implemented: StorageClass
	PartNumberMarker     uint32
	NextPartNumberMarker uint32
	MaxParts             uint32
	IsTruncated          bool
	Parts                []part `xml:"Part"`
}

func listParts(key, uploadID string, cred credential, partNumberMarker uint32) (listPartsResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	params := make(url.Values)
	params.Set("max-parts", "1000")
	params.Set("part-number-marker", strconv.FormatUint(uint64(partNumberMarker), 10))
	params.Set("uploadId", uploadID)
	unsignedReq, err := http.NewRequestWithContext(ctx, http.MethodGet, cred.Endpoint+"/"+key+"?"+params.Encode(), nil)
	if err != nil {
		return listPartsResult{}, err
	}

	signedReq := sign(unsignedReq, cred)
	resp, err := httpClientS3.Do(signedReq)
	if err != nil {
		return listPartsResult{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return listPartsResult{}, fmt.Errorf("endpoint request failed: %d: %s", resp.StatusCode, body)
	}

	result := listPartsResult{}
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	if err != nil {
		return result, err
	}

	for i := range result.Parts {
		result.Parts[i].ETag = strings.TrimSuffix(strings.TrimPrefix(result.Parts[i].ETag, "\""), "\"")
	}

	return result, nil
}

/* completeMultipartUpload */

type completeMultipartUploadBody struct {
	XMLName xml.Name       `xml:"http://s3.amazonaws.com/doc/2006-03-01/ CompleteMultipartUpload" json:"-"`
	Parts   []completePart `xml:"Part"`
}

type completePart struct {
	XMLName    xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ Part" json:"-"`
	PartNumber uint16   `json:"PartNumber"`
	ETag       string   `json:"ETag"`
}

func (r completePart) validate() error {
	if r.PartNumber < 1 || r.PartNumber > 10000 {
		return errors.New("invalid part number")
	} else if r.ETag == "" {
		return errors.New("invalid etag")
	}
	return nil
}

type completeMultipartUploadResult struct {
	Location string
	Bucket   string
	Key      string
	ETag     string
}

func completeMultipartUpload(key, uploadID string, parts []completePart, cred credential) (completeMultipartUploadResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var body bytes.Buffer
	complete := completeMultipartUploadBody{Parts: parts}
	b := xml.NewEncoder(&body)
	err := b.Encode(complete)
	if err != nil {
		return completeMultipartUploadResult{}, err
	}

	params := make(url.Values)
	params.Set("uploadId", uploadID)
	unsignedReq, err := http.NewRequestWithContext(ctx, http.MethodPost, cred.Endpoint+"/"+key+"?"+params.Encode(), &body)
	if err != nil {
		return completeMultipartUploadResult{}, err
	}

	signedReq := sign(unsignedReq, cred)
	resp, err := httpClientS3.Do(signedReq)
	if err != nil {
		return completeMultipartUploadResult{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return completeMultipartUploadResult{}, fmt.Errorf("endpoint request failed: %d: %s", resp.StatusCode, body)
	}

	result := completeMultipartUploadResult{}
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	if err != nil {
		return result, err
	}

	result.ETag = stripETag(result.ETag)

	return result, nil
}

/* abortMultipartUpload */

func abortMultipartUpload(key, uploadID string, cred credential) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	params := make(url.Values)
	params.Set("uploadId", uploadID)
	unsignedReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, cred.Endpoint+"/"+key+"?"+params.Encode(), nil)
	if err != nil {
		return err
	}

	signedReq := sign(unsignedReq, cred)
	resp, err := httpClientS3.Do(signedReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("endpoint request failed: %d: %s", resp.StatusCode, body)
	}

	return nil
}
