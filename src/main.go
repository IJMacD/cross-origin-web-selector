package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/IJMacD/cows/logging"
	"github.com/IJMacD/cows/resources"
)

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

	handler := logging.LoggingMiddleware(mux)

	err := http.ListenAndServe(":8080", handler)
	
  	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}