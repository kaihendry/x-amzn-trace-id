package main

import (
	"context"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/apex/log"
	"github.com/aws/aws-xray-sdk-go/xray"
	"golang.org/x/net/context/ctxhttp"

	"github.com/apex/gateway/v2"
)

var Version string

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// Fetch url
		getAddress := "https://aaq9ogdzqd.execute-api.us-east-1.amazonaws.com/metrics"
		// trace request with Xray
		ctx, root := xray.BeginSegment(context.Background(), "fetch another microservice")

		resp, err := ctxhttp.Get(ctx, xray.Client(http.DefaultClient), getAddress)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if resp.StatusCode != 200 {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()
		root.Close(nil)

		t, err := template.New("").Parse(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />
<title>{{ .Name }}</title>
</head>
<body>
<pre>
{{ .Response }}
</pre>
</body>
</html>`)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = t.Execute(w, struct {
			Name     string
			Response string
		}{
			Name:     os.Getenv("AWS_LAMBDA_FUNCTION_NAME") + Version,
			Response: string(response),
		})

	})

	port := os.Getenv("_LAMBDA_SERVER_PORT")
	var err error
	if port == "" {
		err = http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	} else {
		err = gateway.ListenAndServe("", nil)
	}
	log.Fatalf("failed to start server: %v", err)
}
