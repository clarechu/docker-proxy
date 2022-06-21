package proxy

import (
	"github.com/clarechu/docker-proxy/pkg/models"
	"io"
	log "k8s.io/klog/v2"
	"net/http"
)

// HeadHandler 所有其余的路由全部走这个地方
func (a *App) HeadHandler(w http.ResponseWriter, r *http.Request) {
	log.V(2).Infof("Received request %s %s %s\n", r.Method, r.Host, r.RemoteAddr)
	transport := http.DefaultTransport
	outReq := new(http.Request)
	*outReq = *r // this only does shallow copies of maps
	// outReq.Host = "localhost:9001"
	outReq.Host = a.Host
	log.V(2).Infof("proxy url --> %s", outReq.Host)
	outReq.URL.Host = outReq.Host
	// 设置权限头
	outReq.Header = map[string][]string{}
	outReq.URL.Scheme = models.HttpSchema.SchemaToString()
	res, err := transport.RoundTrip(outReq)
	if err != nil {
		log.Errorf("RoundTrip error :%v", err)
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(err.Error()))
		return
	}
	for key, value := range res.Header {
		for _, v := range value {
			w.Header().Add(key, v)
		}
	}
	res.Request.URL.Host = r.Host
	res.Request.Host = r.Host
	if res.StatusCode != http.StatusOK {
		w.WriteHeader(res.StatusCode)
		return
	}
	defer res.Body.Close()
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
}
