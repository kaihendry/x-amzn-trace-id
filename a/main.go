package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/apex/log"
	"github.com/aws/aws-xray-sdk-go/xray"
	"golang.org/x/net/context/ctxhttp"

	"github.com/apex/gateway/v2"
)

// Version is set during the build Makefile
var Version string

func main() {
	http.Handle("/", xray.Handler(xray.NewFixedSegmentNamer("aApp"), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Fetch url
		getAddress := "https://dx3kwbyjpd.execute-api.ap-southeast-1.amazonaws.com/"
		// trace request with Xray

		resp, err := ctxhttp.Get(r.Context(), xray.Client(nil), getAddress)
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

		t, err := template.New("").Parse(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />
<title>{{ .Name }}</title>
</head>
<body>
<h1>Trace ID</h1>
<pre>{{ .TraceID }}</pre>
<h1>Response</h1>
<pre>
{{ .Response }}
</pre>

<dl>
{{range $key, $value := .Env -}}
<dt>{{ $key }}</dt><dd>{{ $value }}</dd>
{{- end}}
</dl>

</body>
</html>`)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		err = t.Execute(w, struct {
			Name     string
			TraceID  string
			Response string
			Env      map[string]string
		}{
			Name:     os.Getenv("AWS_LAMBDA_FUNCTION_NAME") + Version,
			TraceID:  r.Header.Get("x-amzn-trace-id"),
			Response: string(response),
			Env:      envMap(),
		})

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

func envMap() map[string]string {
	envmap := make(map[string]string)
	for _, e := range os.Environ() {
		ep := strings.SplitN(e, "=", 2)
		// Skip potentially security sensitive AWS stuff
		if ep[0] == "AWS_SECRET_ACCESS_KEY" {
			continue
		}
		if ep[0] == "AWS_SESSION_TOKEN" {
			continue
		}
		envmap[ep[0]] = ep[1]
	}
	return envmap
}
