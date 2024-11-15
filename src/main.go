package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/IJMacD/cows/logging"
	"github.com/IJMacD/cows/resources"
)

func main() {
	rm, err := getResourceMap("resources.json")

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	r := resources.NewResources(rm)

    fmt.Println("Listening on :8080")

	mux := http.NewServeMux()
	mux.Handle("/r/{resourceName}", r)

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