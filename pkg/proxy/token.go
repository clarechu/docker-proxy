package proxy

import (
	log "k8s.io/klog/v2"
	"net/http"
)

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	log.Infof("Received request %s %s %s\n", r.Method, r.Host, r.RemoteAddr)
	/*	username := utils.GetParams(r, "account")
		password := r.Header.Get("")
		if username == "admin" && password == "admin123" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}*/
	w.Write([]byte("{\"refresh_token\":\"kas9Da81Dfa8\",\"access_token\":\"eyJhbGciOiJFUzI1NiIsInR5\",\"expires_in\":900,\"scope\":\"\"}"))
	w.WriteHeader(http.StatusOK)
}
