package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/clarechu/docker-proxy/pkg/models"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	log "k8s.io/klog/v2"
	"net/http"
	"time"
)

func (a *App) LoggingHandlerFunc(r *mux.Router) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			logging := &models.Logging{}
			blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: w}
			headers, _ := json.Marshal(req.Header)
			logging.Header = string(headers)
			logging.ClientIP = req.Header.Get("X-Forwarded-For")
			logging.Method = req.Method
			logging.URI = req.RequestURI
			reqBody, _ := io.ReadAll(req.Body)
			req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
			logging.RequestBody = string(reqBody)
			startTime := time.Now()
			next.ServeHTTP(blw, req)
			logging.HttpStatusCode = fmt.Sprintf("%d", blw.statusCode)
			logging.ResponseBody = blw.body.String()
			logging.ReturnTime = time.Since(startTime).String()
			log.V(4).Infof("logging -> %+v", logging)
			// push logging
			a.Queue.Push(func() error {
				return a.LoggingHandler(logging)
			})
		})
	}

}

type bodyLogWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)

	return w.ResponseWriter.Write(b)
}

func (w *bodyLogWriter) WriteHeader(statusCode int) {

	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
