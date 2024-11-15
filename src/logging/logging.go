package logging

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/urfave/negroni"
)

func LoggingMiddleware (next http.Handler) http.Handler {
	return http.HandlerFunc(
		func (w http.ResponseWriter, r *http.Request) {
			// Wrapped writer to intercept and record response code and bytes written
			ww := negroni.NewResponseWriter(w)

			next.ServeHTTP(ww, r)

			remote := r.RemoteAddr
			// Implicit trust
			switch {
			case len(r.Header["X-Real-Ip"]) > 0:
				remote = r.Header["X-Real-Ip"][0]
			case len(r.Header["X-Forwarded-For"]) > 0:
				f := strings.Join(r.Header["X-Forwarded-For"], ",")
				s := strings.Split(f, ",")
				remote = s[0]
			}

			status := ww.Status()

			user := "-"
			authUser, _, ok := getBasicAuthUser(r)
			if status >= 200 && status < 400 && ok {
				user = authUser
			}

			date := time.Now().Format(time.RFC3339)

			referer := strings.Join(r.Header["Referer"], " ")
			if len(referer) == 0 {
				referer = "-"
			}

			ua := strings.Join(r.Header["User-Agent"], " ")
			if len(ua) == 0 {
				ua = "-"
			}

			fmt.Printf("%s - %s [%s] \"%s %s %s\" %d %d \"%s\" \"%s\"\n", remote, user, date, r.Method, r.RequestURI, r.Proto, status, ww.Size(), referer, ua)
		},
	)
}

func getBasicAuthUser (r *http.Request) (username string, password string, ok bool) {
	if len(r.Header["Authorization"]) == 0 {
		return "", "", false
	}

	auth := r.Header["Authorization"][0]

	s, ok := strings.CutPrefix(auth, "Basic ")

	if !ok {
		return "", "", false
	}

	b, err := base64.StdEncoding.DecodeString(s)

	if err != nil {
		return 
	}

	p := strings.SplitN(string(b), ":", 2)

	if len(p) != 2 {
		return "", "", false
	} 

	return p[0], p[1], true
}