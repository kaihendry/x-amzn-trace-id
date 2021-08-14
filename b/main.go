package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/apex/log"
	"github.com/aws/aws-xray-sdk-go/xray"

	"github.com/apex/gateway/v2"
)

var Version string

func main() {
	http.Handle("/", xray.Handler(xray.NewFixedSegmentNamer("bApp"), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get("x-amzn-trace-id")
		log.Info(traceID)
		fmt.Fprintf(w, "b "+traceID)
	})))

	port := os.Getenv("_LAMBDA_SERVER_PORT")
	var err error
	if port == "" {
		err = http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	} else {
		err = gateway.ListenAndServe("", nil)
	}
	log.Fatalf("failed to start server: %v", err)
}
