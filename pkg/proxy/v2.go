package proxy

import (
	"fmt"
	"github.com/clarechu/docker-proxy/pkg/models"
	log "k8s.io/klog/v2"
	"net/http"
)

func (a *App) VersionHandler(w http.ResponseWriter, r *http.Request) {
	log.V(2).Infof("Received request %s %s %s", r.Method, r.Host, r.RemoteAddr)
	host := r.Host
	authURL := fmt.Sprintf("%s://%s/v2/token", models.HttpSchema, host)
	/*	if r.Header.Get(models.AuthorizationKey) == "Bearer eyJhbGciOiJFUzI1NiIsInR5" {
		w.WriteHeader(http.StatusOK)
		return
	}*/
	ok := a.OAuth2EventHandlerFuncs.CheckTokenFunc(r.Header.Get(models.AuthorizationKey))
	if ok == true {
		w.WriteHeader(http.StatusOK)
		return
	}
	authorizationValue := fmt.Sprintf("Bearer realm=\"%s\",service=\"%s\"", authURL, authURL)
	w.Header().Set(models.WwwAuthorizationKey, authorizationValue)
	w.WriteHeader(http.StatusUnauthorized)
}
