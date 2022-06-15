package proxy

import (
	"fmt"
	"github.com/clarechu/docker-proxy/pkg/models"
	"github.com/clarechu/docker-proxy/pkg/utils"
	"github.com/clarechu/docker-proxy/pkg/utils/base64"
	log "k8s.io/klog/v2"
	"net/http"
)

// TokenHandler 使用账号密码 来验证当前的用户是否合法
func TokenHandler(w http.ResponseWriter, r *http.Request) {
	log.V(2).Infof("Received request %s %s %s\n", r.Method, r.Host, r.RemoteAddr)
	account := utils.GetParams(r, "account")
	auth := r.Header.Get(models.AuthorizationKey)
	password := fmt.Sprintf("Basic %s", base64.EncodeToString("admin:xxx"))
	if account == "admin" && auth == password {
		w.Write([]byte("{\"refresh_token\":\"kas9Da81Dfa8\",\"access_token\":\"eyJhbGciOiJFUzI1NiIsInR5\",\"expires_in\":900,\"scope\":\"\"}"))
		w.WriteHeader(http.StatusOK)
	}
	w.WriteHeader(http.StatusUnauthorized)
}

// PostTokenHandler 获取token令牌 和 刷新令牌
func PostTokenHandler(w http.ResponseWriter, r *http.Request) {
	log.V(2).Infof("Received request %s %s %s\n", r.Method, r.Host, r.RemoteAddr)
	token := &models.Token{}
	err := utils.GetBodyByForm(r, token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("{\"refresh_token\":\"kas9Da81Dfa8\",\"access_token\":\"eyJhbGciOiJFUzI1NiIsInR5\",\"expires_in\":900,\"scope\":\"\"}"))
	w.WriteHeader(http.StatusOK)
}
