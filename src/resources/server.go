package resources

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func NewResources (resources map[string]ResourceSpec) Resources {
	r := Resources{registeredResources: resources}
	return r
}

func (res Resources) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v := r.PathValue("resourceName")

	rs, ok := res.registeredResources[v]

	if !ok {
		w.WriteHeader(404)
		return
	}

	var shouldOutputJSON bool

	if len(r.Header["Accept"]) > 0 {
		shouldOutputJSON = strings.Contains(r.Header["Accept"][0], "application/json")
	}

	var data []byte
	var err error

	switch{
	case rs.QuerySelector != "":
		val, e := GetScalar(rs.URL, rs.QuerySelector)
		
		if e != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, e.Error())
			return
		}

		if shouldOutputJSON {
			msg := ScalarMessage{
				Value: val,
			}
			
			data, err = json.Marshal(msg)
		} else {
			data = []byte(val)
		}

	case rs.QuerySelectorAll != "":
		vals, e := GetVector(rs.URL, rs.QuerySelectorAll)
		
		if e != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, e.Error())
			return
		}

		if shouldOutputJSON {
			msg := VectorMessage{
				Values: vals,
			}
			
			data, err = json.Marshal(msg)
		} else {
			data = []byte(strings.Join(vals, "\n"))
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	if shouldOutputJSON {
		w.Header().Set("Content-Type", "application/json")
	} else {
		w.Header().Set("Content-Type", "text/plain")
	}
	
	w.Write(data)
}