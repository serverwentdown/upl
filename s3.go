package main

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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
		1*60, // seconds
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

type errEndpoint struct {
	err    error
	status string
	body   []byte
}

func (e errEndpoint) Unwrap() error {
	return e.err
}

func (e errEndpoint) Error() string {
	body := bytes.ReplaceAll(e.body, []byte("\n"), []byte(""))
	if e.err != nil {
		return fmt.Sprintf("endpoint responded with %v: %s", e.err, body)
	}
	return fmt.Sprintf("endpoint responded with %s: %s", e.status, body)
}

func endpointReturnedError(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, _ := ioutil.ReadAll(resp.Body)
		return errEndpoint{responseToError(resp), resp.Status, body}
	}
	return nil
}

/* initiateMultipartUpload */

type initiateMultipartUploadResult struct {
	Bucket   string
	Key      string
	UploadID string `xml:"UploadId"`
}

func initiateMultipartUpload(
	key string,
	cred credential,
) (initiateMultipartUploadResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	params := make(url.Values)
	params.Set("uploads", "")
	unsignedReq, err := http.NewRequestWithContext(ctx, http.MethodPost, cred.Endpoint+"/"+key+"?"+params.Encode(), nil)
	if err != nil {
		log.Printf("failure creating request: %v", err)
		return initiateMultipartUploadResult{}, err
	}
	if cred.ACL != "" {
		unsignedReq.Header.Set("X-Amz-Acl", cred.ACL)
	}

	signedReq := sign(unsignedReq, cred)
	resp, err := httpClientS3.Do(signedReq)
	if err != nil {
		log.Printf("failure connecting to endpoint: %v", err)
		return initiateMultipartUploadResult{}, err
	}
	defer resp.Body.Close()
	err = endpointReturnedError(resp)
	if err != nil {
		log.Printf("endpoint responded negatively: %v", err)
		return initiateMultipartUploadResult{}, err
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

func listParts(
	key, uploadID string,
	cred credential,
	partNumberMarker uint32,
) (listPartsResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	params := make(url.Values)
	params.Set("max-parts", "1000")
	params.Set("part-number-marker", strconv.FormatUint(uint64(partNumberMarker), 10))
	params.Set("uploadId", uploadID)
	unsignedReq, err := http.NewRequestWithContext(ctx, http.MethodGet, cred.Endpoint+"/"+key+"?"+params.Encode(), nil)
	if err != nil {
		log.Printf("failure creating request: %v", err)
		return listPartsResult{}, err
	}

	signedReq := sign(unsignedReq, cred)
	resp, err := httpClientS3.Do(signedReq)
	if err != nil {
		log.Printf("failure connecting to endpoint: %v", err)
		return listPartsResult{}, err
	}
	defer resp.Body.Close()
	err = endpointReturnedError(resp)
	if err != nil {
		log.Printf("endpoint responded negatively: %v", err)
		return listPartsResult{}, err
	}

	result := listPartsResult{}
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	if err != nil {
		return result, err
	}

	for i := range result.Parts {
		result.Parts[i].ETag = stripETag(result.Parts[i].ETag)
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

func completeMultipartUpload(
	key, uploadID string,
	parts []completePart,
	cred credential,
) (completeMultipartUploadResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
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
		log.Printf("failure creating request: %v", err)
		return completeMultipartUploadResult{}, err
	}

	signedReq := sign(unsignedReq, cred)
	resp, err := httpClientS3.Do(signedReq)
	if err != nil {
		log.Printf("failure connecting to endpoint: %v", err)
		return completeMultipartUploadResult{}, err
	}
	defer resp.Body.Close()
	err = endpointReturnedError(resp)
	if err != nil {
		log.Printf("endpoint responded negatively: %v", err)
		return completeMultipartUploadResult{}, err
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

func abortMultipartUpload(
	key, uploadID string,
	cred credential,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	params := make(url.Values)
	params.Set("uploadId", uploadID)
	unsignedReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, cred.Endpoint+"/"+key+"?"+params.Encode(), nil)
	if err != nil {
		log.Printf("failure creating request: %v", err)
		return err
	}

	signedReq := sign(unsignedReq, cred)
	resp, err := httpClientS3.Do(signedReq)
	if err != nil {
		log.Printf("failure connecting to endpoint: %v", err)
		return err
	}
	defer resp.Body.Close()
	err = endpointReturnedError(resp)
	if err != nil {
		log.Printf("endpoint responded negatively: %v", err)
		return err
	}

	return nil
}
