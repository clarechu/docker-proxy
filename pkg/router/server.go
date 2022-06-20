// Copyright (c) 2021 The static-server Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package router

import (
	"fmt"
	"github.com/clarechu/docker-proxy/pkg/models"
	"github.com/clarechu/docker-proxy/pkg/proxy"
	"github.com/clarechu/docker-proxy/pkg/utils/base64"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	log "k8s.io/klog/v2"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Server struct {
	sv *http.Server
}

type HandlerFunc func(http.ResponseWriter, *http.Request)

type Router struct {
	handler HandlerFunc
	path    string
	methods []string
}

func AddRouter(router *mux.Router, app *proxy.App) {

	router.PathPrefix("/v2/token").HandlerFunc(app.TokenHandler).Methods(http.MethodGet)
	router.PathPrefix("/v2/token").HandlerFunc(app.PostTokenHandler).Methods(http.MethodPost)

	router.PathPrefix("/v2").HandlerFunc(app.VersionHandler).Methods(http.MethodGet)

	router.PathPrefix("/").HandlerFunc(app.OtherHandler).Methods(http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPatch, http.MethodPut)

}

var validate = validator.New()

func NewServer(root *models.Root) *Server {
	//文件浏览
	r := mux.NewRouter()
	addHTTPMiddleware(r)
	err := validate.Struct(root.App)
	if err != nil {
		panic(err)
	}
	AddRouter(r, buildApp(root.App))
	srv := &http.Server{
		// Handler: handlers.LoggingHandler(os.Stdout, r),
		Handler: r,
		Addr:    fmt.Sprintf(":%d", root.Port),
	}
	return &Server{
		sv: srv,
	}
}

func addHTTPMiddleware(router *mux.Router) {
	router.Use(CORSMethodMiddleware(router))
	router.Use(LogMiddleware(router))
}

func (s *Server) Run() {
	log.V(0).Info("Available on:")
	log.V(0).Infof("   http://127.0.0.1%s", s.sv.Addr)
	log.V(0).Infof("Hit CTRL-C to stop the server")
	log.Fatal(s.sv.ListenAndServe())
}

// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type spaHandler struct {
	rootPath   string
	staticPath string
	indexPath  string
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	path = strings.Replace(path, h.rootPath, "", 1)
	r.URL.Path = path
	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		path = filepath.Join(h.staticPath, h.indexPath)
		http.ServeFile(w, r, path)
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func buildApp(a *models.App) *proxy.App {
	token := ""
	switch a.RegistryType {
	case models.DockerRegistry:
	case models.NexusRegistry:
		token = fmt.Sprintf("Basic %s", base64.EncodeToString(fmt.Sprintf("%s:%s", a.Nexus.Username, a.Nexus.Password)))
	case models.HarborRegistry:
		token = fmt.Sprintf("Basic %s", base64.EncodeToString(fmt.Sprintf("%s:%s", a.Nexus.Username, a.Nexus.Password)))
	}
	return &proxy.App{
		Host:                    a.DockerRegistryHost,
		Token:                   token,
		Schema:                  a.Schema.SchemaToString(),
		OAuth2EventHandlerFuncs: a.OAuth2EventHandlerFuncs,
	}
}
