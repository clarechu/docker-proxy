package utils

import (
	"encoding/json"
	"github.com/clarechu/docker-proxy/pkg/utils/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"io/ioutil"
	"net/http"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func init() {
	validate = validator.New()
}

type ErrorResponse struct {
	Code    int32       `json:"code"`
	Message string      `json:"message" `
	Data    interface{} `json:"data" schema:"zxzxzx" defaults:"111"`
}

func GetBody(r *http.Request, i interface{}) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, i)
	if err != nil {
		return err
	}
	if err := defaults.Set(i); err != nil {
		return err
	}
	return validate.Struct(i)
}

var decoder = schema.NewDecoder()

func GetBodyByForm(r *http.Request, i interface{}) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	err = decoder.Decode(i, r.PostForm)
	if err != nil {
		return err
	}
	if err := defaults.Set(i); err != nil {
		return err
	}
	return validate.Struct(i)
}

//GetParams /api/v2.0/cluster?name=xxx
// key = name
// return xxx
func GetParams(r *http.Request, key string) string {
	url := r.URL
	values := url.Query()
	return values.Get(key)
}

//GetPathVar
// /xs/sd/ew/foo/{foo}/bar/{bar1}
// key foo =x
func GetPathVar(r *http.Request, key string) string {
	values := mux.Vars(r)
	return values[key]
}

//GetHeader
func GetHeader(r *http.Request, key string) string {
	return r.Header.Get(key)
}
