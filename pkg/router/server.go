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
	"github.com/clarechu/docker-proxy/pkg/nexus"
	"github.com/clarechu/docker-proxy/pkg/proxy"
	"github.com/clarechu/docker-proxy/pkg/utils/base64"
	"github.com/clarechu/docker-proxy/pkg/utils/queue"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	log "k8s.io/klog/v2"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Server struct {
	sv    *http.Server
	queue queue.Instance
	stop  chan struct{}
}

type NexusServer struct {
	servers []*http.Server
	queue   queue.Instance
	stop    chan struct{}
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

	router.Path("/v2/").HandlerFunc(app.VersionHandler).Methods(http.MethodGet)

	router.PathPrefix("/v2/").HandlerFunc(app.OtherHandler).Methods(http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut)
	router.PathPrefix("/v2/").HandlerFunc(app.HeadHandler).Methods(http.MethodHead)

}

var validate = validator.New()

func NewServer(root *models.Root) *Server {
	stop := make(chan struct{}, 0)
	//文件浏览
	r := mux.NewRouter()
	err := validate.Struct(root.App)
	if err != nil {
		panic(err)
	}
	app := buildApp(root.App, stop)
	AddRouter(r, app)
	addHTTPMiddleware(r, app)
	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%d", root.Port),
	}
	return &Server{
		stop:  stop,
		queue: app.Queue,
		sv:    srv,
	}
}

func NewNexusServer(nexus *models.NexusApp) *NexusServer {
	stop := make(chan struct{}, 0)
	//文件浏览
	err := validate.Struct(nexus)
	if err != nil {
		panic(err)
	}
	instance := queue.NewQueue(5 * time.Second)
	apps := buildNexusApp(nexus, instance, stop)
	servers := make([]*http.Server, 0)
	for _, app := range apps {
		r := mux.NewRouter()
		AddRouter(r, app)
		addHTTPMiddleware(r, app)
		srv := &http.Server{
			Handler: r,
			Addr:    fmt.Sprintf(":%d", app.Port),
		}
		servers = append(servers, srv)
	}
	return &NexusServer{
		stop:    stop,
		queue:   instance,
		servers: servers,
	}
}

func addHTTPMiddleware(router *mux.Router, app *proxy.App) {
	router.Use(CORSMethodMiddleware(router))
	router.Use(LogMiddleware(router))
	router.Use(app.LoggingHandlerFunc(router))
}

func (s *Server) Run() {
	log.V(0).Info("Available on:")
	log.V(0).Infof("   http://127.0.0.1%s", s.sv.Addr)
	log.V(0).Infof("Hit CTRL-C to stop the server")
	go s.queue.Run(s.stop)
	log.Fatal(s.sv.ListenAndServe())
}

func (s *NexusServer) Run() {
	go s.queue.Run(s.stop)
	for _, sv := range s.servers {
		log.V(0).Info("Available on:")
		log.V(0).Infof("   http://127.0.0.1%s", sv.Addr)
		log.V(0).Infof("Hit CTRL-C to stop the server")
		go log.Fatal(sv.ListenAndServe())
	}
	<-s.stop
}

func (s *NexusServer) Close() {
	defer close(s.stop)
	for _, server := range s.servers {
		err := server.Close()
		if err != nil {
			log.Fatalf("close server error:%v", err)
		}
	}
}

func (s *Server) Close() {
	defer close(s.stop)
	s.sv.Close()
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

func buildApp(a *models.App, stop chan struct{}) *proxy.App {
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
		Stop:                    stop,
		LoggingHandler:          a.LoggingHandler,
		Queue:                   queue.NewQueue(5 * time.Second),
		Schema:                  a.Schema.SchemaToString(),
		OAuth2EventHandlerFuncs: a.OAuth2EventHandlerFuncs,
	}
}

func buildNexusApp(a *models.NexusApp, queue queue.Instance, stop chan struct{}) []*proxy.App {
	token := fmt.Sprintf("Basic %s", base64.EncodeToString(fmt.Sprintf("%s:%s", a.Username, a.Password)))
	apps := make([]*proxy.App, 0)
	nexusClient := nexus.NewRepository(*a)
	ports, err := nexusClient.GetPortByDocker()
	if err != nil {
		log.Fatalf("nexus get port error: %v", err)
	}
	for _, port := range ports {
		nexusUrl, err := url.Parse(a.URL)
		if err != nil {
			log.Fatalf("parse url error:%s", err.Error())
		}
		log.V(2).Infof("nexus ip :%s", getIP(nexusUrl.Host))
		app := &proxy.App{
			Host:                    fmt.Sprintf("%s:%d", getIP(nexusUrl.Host), port),
			Port:                    port,
			Token:                   token,
			Stop:                    stop,
			LoggingHandler:          a.LoggingHandler,
			Queue:                   queue,
			Schema:                  a.Schema.SchemaToString(),
			OAuth2EventHandlerFuncs: a.OAuth2EventHandlerFuncs,
		}
		apps = append(apps, app)
	}
	if a.Ports != nil {
		for _, port := range a.Ports {
			for _, app := range apps {
				targetPort := int(port.TargetPort)
				if targetPort == app.Port {
					app.Port = int(port.Port)
				}
			}
		}
	}
	return apps
}

func getIP(host string) string {
	return strings.Split(host, ":")[0]
}
