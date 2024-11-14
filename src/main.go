package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/IJMacD/boc-prime-odata/resources"
	"github.com/urfave/negroni"
)

func loggingMiddleware (next http.Handler) http.Handler {
	return http.HandlerFunc(
		func (w http.ResponseWriter, r *http.Request) {
			// Wrapped writer to intercept and record response code
			ww := negroni.NewResponseWriter(w)

			next.ServeHTTP(ww, r)

			fmt.Printf("\"%s %s %s\" %d %d \"%s\" \"%s\"\n", r.Method, r.RequestURI, r.Proto, ww.Status(), ww.Size(), "-", strings.Join(r.Header["User-Agent"], " "))
		},
	)
}

func main() {
	resourceMap := make(map[string]resources.ResourceSpec)

	resourceMap["bocPrime"] = resources.ResourceSpec{
		URL: "https://www.bochk.com/whk/rates/hkDollarPrimeRate/hkDollarPrimeRate-enquiry.action?lang=en",
		QuerySelector: ".best-rate td:nth-child(2)",
	}

	r := resources.NewResources(resourceMap)

    fmt.Println("Listening on :8080")

	mux := http.NewServeMux()
	mux.Handle("/r/{resourceName}", r)

	handler := loggingMiddleware(mux)

	err := http.ListenAndServe(":8080", handler)
	
  	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}