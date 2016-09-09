package httpx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"common/errors"
	l "common/log"
	"github.com/gorilla/mux"
)

const (
	HTTP_GET    = "GET"
	HTTP_POST   = "POST"
	HTTP_PUT    = "PUT"
	HTTP_DELETE = "DELETE"
)

const (
	OK               = "OK"
	JSON_EMPTY_ARRAY = "[]"
	JSON_EMPTY_OBJ   = "{}"
)

// -----------------------
var log = l.New("")

// -----------------------
// RegisterHttpHandler

type HttpHandler func(*http.Request) (string, interface{})

func RegisterHttpHandler(router *mux.Router, path, method string, handler HttpHandler) {
	h := func(w http.ResponseWriter, r *http.Request) {
		// parseForm
		if err := r.ParseForm(); err != nil {
			log.Warn(errors.As(err))
		}

		// dump
		bytes, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Warn(err)
		} else {
			log.Debug(string(bytes))
		}

		dump := dumpHttpRequest(r)
		log.Debug(dump)

		t := time.Now()
		status, body := handler(r)
		writeHttpResp(w, dump, status, body, t)
	}
	router.HandleFunc(path, h).Methods(method)
}

func dumpHttpRequest(r *http.Request) string {
	if r.Method == "GET" {
		return fmt.Sprintf("%s %s", r.Method, r.URL.RequestURI())
	}

	if r.Method == "POST" {
		return fmt.Sprintf("%s %s", r.Method, r.URL.RequestURI())
	}

	return fmt.Sprintf("%s %s %s", r.Method, r.URL.RequestURI(), r.Form)
}

// --------------------------------
// response

const httpJsonRespFmt = `{
  "api": "1.0",
  "status": "%s",
  "err": %s,
  "data": %s
}
`

func writeHttpResp(w http.ResponseWriter, dump string, status string, body interface{}, t time.Time) {
	w.Header().Set("Content-Type", "application/json")

	sub := time.Now().Sub(t)
	// empty array
	if body == JSON_EMPTY_ARRAY {
		log.Info(dump, status, sub)
		fmt.Fprintf(w, httpJsonRespFmt, status, `""`, body)
		return
	}

	if body == JSON_EMPTY_OBJ {
		log.Info(dump, status, sub)
		fmt.Fprintf(w, httpJsonRespFmt, status, `""`, body)
		return
	}

	errStr, data := "", JSON_EMPTY_OBJ
	res, err := json.MarshalIndent(body, " ", "    ")
	if err != nil {
		errStr = `"` + err.Error() + `"`
		log.Debug(dump, status, errStr, data, sub)
		fmt.Fprintf(w, httpJsonRespFmt, status, errStr, data)
		return
	}

	// error
	if status != StatusOK && status != StatusCreated && status != StatusNoContent {
		errStr = `"` + string(res) + `"`
		log.Debug(dump, status, errStr, data, sub)
		fmt.Fprintf(w, httpJsonRespFmt, status, errStr, data)
		return
	}

	errStr = `"` + OK + `"`
	data = string(res)

	log.Debug(dump, status, sub)
	fmt.Fprintf(w, httpJsonRespFmt, status, errStr, data)
}

// -----------------------------------------------------------------------------
// HttpListen

func HttpServe(addr string) error {
	// root
	root := mux.NewRouter()

	// api
	api := root.PathPrefix("/api/v1").Subrouter()
	regHttpApiHandles(api)

	// log, admin,
	// ...

	// root
	http.Handle("/", root)

	// listen
	return http.ListenAndServe(addr, nil)
}

func regHttpApiHandles(router *mux.Router) {
	RegisterBuildHandler(router)
}
