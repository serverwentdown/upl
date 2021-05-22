package main

import (
	"net/http"

	"github.com/minio/minio-go/v7/pkg/signer"
)

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
