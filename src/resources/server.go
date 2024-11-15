package resources

import (
	"encoding/json"
	"io"
	"net/http"
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

		msg := ScalarMessage{
			Value: val,
		}
		
		data, err = json.Marshal(msg)

	case rs.QuerySelectorAll != "":
		vals, e := GetVector(rs.URL, rs.QuerySelectorAll)
		
		if e != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, e.Error())
			return
		}

		msg := VectorMessage{
			Values: vals,
		}
		
		data, err = json.Marshal(msg)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(data)
}