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

	if rs.QuerySelector != "" {
		val, err := GetScalar(rs.URL, rs.QuerySelector)
		
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, err.Error())
			return
		}

		msg := ScalarMessage{
			Value: val,
		}
		
		data, err := json.Marshal(msg)
		
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)

		return
	}

	if rs.QuerySelectorAll != "" {
		vals, err := GetVector(rs.URL, rs.QuerySelectorAll)
		
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, err.Error())
			return
		}

		msg := VectorMessage{
			Values: vals,
		}
		
		data, err := json.Marshal(msg)
		
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)

		return
	}
}