package odata

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/IJMacD/boc-prime-odata/boc"
)

//go:embed metadata.xml
var metadata string

type Message struct {
	ODataContext string `json:"@odata.context"`
	Prime float32
}

func GetBank(w http.ResponseWriter, r *http.Request) {
	prime, err := boc.GetPrime()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}

	scheme := "http"
	if (len(r.Header["X-Forwarded-Scheme"]) > 0) {
		scheme = r.Header["X-Forwarded-Scheme"][0]
	}

	base := fmt.Sprintf("%s://%s%s", scheme, r.Host, "/V4")
	
	msg := Message{
		ODataContext: base + "/$metadata#Banks/$entity",
		Prime: prime,
	}
	
	data, err := json.Marshal(msg)
	
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "application/json; odata.metadata=minimal")
	w.Header().Set("OData-Version", "4.0")
	w.Write(data)
}

func GetMetaData (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")

	io.WriteString(w, metadata)
}