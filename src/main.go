package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/IJMacD/cows/logging"
	"github.com/IJMacD/cows/resources"
)

func main() {
	f := "resources.json"

	if len(os.Args) > 1 {
		switch {
		case strings.HasPrefix(os.Args[1], "--resources="):
			f = strings.TrimPrefix(os.Args[1], "--resources=")
		case os.Args[1] == "--resources" && len(os.Args) > 2:
			f = os.Args[2]
		}
	}

	rm, err := getResourceMap(f)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	r := resources.NewResources(rm)

    fmt.Println("Listening on :8080")

	mux := http.NewServeMux()
	mux.Handle("/r/{resourceName}", r)

	mux.Handle("/", http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		h := "<h1>Cross Origin Web Selector</h1><p>Pre-defined named selectors are accessible at <code>/r/{resourceName}</code></p>"
		w.Write([]byte(h))
	}))

	handler := logging.LoggingMiddleware(mux)

	err = http.ListenAndServe(":8080", handler)

  	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func getResourceMap (fileName string) (resources.ResourceMap, error) {
	configFile, err := os.Open(fileName)

	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(configFile)

	if err != nil {
		return nil, err
	}

	var resourceMap resources.ResourceMap

	json.Unmarshal(bytes, &resourceMap)

	return resourceMap, nil
}