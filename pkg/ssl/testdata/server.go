package testdata

import (
	"crypto/tls"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	log "k8s.io/klog/v2"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello World")
	})
	srv := &http.Server{
		Handler:   r,
		TLSConfig: &tls.Config{ServerName: "localhost"},
		Addr:      fmt.Sprintf(":443"),
	}
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	dir = filepath.Join(dir, "pkg/ssl/testdata")
	log.Fatal(srv.ListenAndServeTLS(filepath.Join(dir, "server.crt"), filepath.Join(dir, "server.key")))
}
