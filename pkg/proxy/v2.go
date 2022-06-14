package proxy

import (
	"fmt"
	log "k8s.io/klog/v2"
	"net/http"
)

var AuthURL = "http://172.20.218.10:7777/v2/token"

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	log.Infof("Received request %s %s %s\n", r.Method, r.Host, r.RemoteAddr)
	if r.Header.Get("Authorization") == "Bearer eyJhbGciOiJFUzI1NiIsInR5" {
		w.WriteHeader(http.StatusOK)
		return
	}
	AuthValue := fmt.Sprintf("Bearer realm=\"%s\",service=\"%s\"", AuthURL, AuthURL)
	w.Header().Set("WWW-Authenticate", AuthValue)
	w.WriteHeader(http.StatusUnauthorized)
}
