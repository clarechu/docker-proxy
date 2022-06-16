package proxy

import (
	"github.com/clarechu/docker-proxy/pkg/models"
	"io"
	log "k8s.io/klog/v2"
	"net"
	"net/http"
	"strings"
)

// OtherHandler 所有其余的路由全部走这个地方
func (a *App) OtherHandler(w http.ResponseWriter, r *http.Request) {
	log.V(2).Infof("Received request %s %s %s\n", r.Method, r.Host, r.RemoteAddr)
	transport := http.DefaultTransport
	outReq := new(http.Request)
	*outReq = *r // this only does shallow copies of maps
	if clientIP, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		if prior, ok := outReq.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		outReq.Header.Set("X-Forwarded-For", clientIP)
	}
	// outReq.Host = "localhost:9001"
	outReq.Host = a.Host
	log.V(2).Infof("proxy url --> %s", outReq.Host)
	outReq.URL.Host = outReq.Host
	// 设置权限头
	// outReq.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	outReq.Header.Set(models.AuthorizationKey, a.Token)
	outReq.URL.Scheme = models.HttpSchema.SchemaToString()
	res, err := transport.RoundTrip(outReq)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer io.Copy(w, res.Body)
	defer res.Body.Close()
	for key, value := range res.Header {
		for _, v := range value {
			w.Header().Add(key, v)
		}
	}
	w.WriteHeader(res.StatusCode)
}
