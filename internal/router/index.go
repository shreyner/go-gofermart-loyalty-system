package router

import "net/http"

func Index(rw http.ResponseWriter, _ *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Main page"))
}
