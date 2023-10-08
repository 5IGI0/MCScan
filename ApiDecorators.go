package main

import (
	"encoding/json"
	"net/http"
)

func ApiJsonDecorator(fn func(w http.ResponseWriter, r *http.Request) (error, interface{})) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err, data := fn(w, r)
		var err_str *string

		w.Header().Add("Content-Type", "application/json")
		if err != nil {
			w.WriteHeader(500)
			a := err.Error()
			err_str = &a
		} else {
			w.WriteHeader(200)
		}

		enc, _ := json.Marshal(
			ApiResponseBase{
				Error: err_str,
				Data:  data,
			})
		w.Write(enc)
	}
}
