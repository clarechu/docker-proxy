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
func OtherHandler(w http.ResponseWriter, r *http.Request) {
	log.Infof("Received request %s %s %s\n", r.Method, r.Host, r.RemoteAddr)
	transport := http.DefaultTransport
	outReq := new(http.Request)
	*outReq = *r // this only does shallow copies of maps
	if clientIP, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		if prior, ok := outReq.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		outReq.Header.Set("X-Forwarded-For", clientIP)
	}
	outReq.Host = "localhost:9001"
	log.Infof("proxy url --> %s", outReq.Host)
	outReq.URL.Host = outReq.Host
	// 设置权限头
	outReq.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	outReq.URL.Scheme = models.HttpSchema
	res, err := transport.RoundTrip(outReq)
	if err != nil {
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
