package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func readyz(w http.ResponseWriter, req *http.Request) {
	status := http.StatusOK
	var msg bytes.Buffer

	storeErr := globalStore.ping()
	if storeErr != nil {
		status = http.StatusServiceUnavailable
		fmt.Fprintf(&msg, "store: %v\n", storeErr)
	} else {
		fmt.Fprintf(&msg, "store: ok\n")
	}

	w.WriteHeader(status)
	w.Write(msg.Bytes())
}
