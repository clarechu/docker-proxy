package proxy

import (
	"encoding/json"
	"github.com/clarechu/docker-proxy/pkg/models"
	"github.com/clarechu/docker-proxy/pkg/utils"
	log "k8s.io/klog/v2"
	"net/http"
)

// TokenHandler 使用账号密码 来验证当前的用户是否合法
func (a *App) TokenHandler(w http.ResponseWriter, r *http.Request) {
	log.V(2).Infof("Received request %s %s %s\n", r.Method, r.Host, r.RemoteAddr)
	account := utils.GetParams(r, "account")
	auth := r.Header.Get(models.AuthorizationKey)
	user := &models.User{
		Account:    account,
		BasicToken: auth,
	}
	oauth2, err := a.OAuth2EventHandlerFuncs.LoginFunc(user)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	b, err := json.Marshal(oauth2)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Write(b)
	w.WriteHeader(http.StatusOK)
}

// PostTokenHandler 获取token令牌 和 刷新令牌
func (a *App) PostTokenHandler(w http.ResponseWriter, r *http.Request) {
	log.V(2).Infof("Received request %s %s %s\n", r.Method, r.Host, r.RemoteAddr)
	token := &models.Token{}
	err := utils.GetBodyByForm(r, token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	oauth2, err := a.OAuth2EventHandlerFuncs.PostTokenFunc(token)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	b, err := json.Marshal(oauth2)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(b)
	w.WriteHeader(http.StatusOK)
}
