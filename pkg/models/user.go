package models

type User struct {
	// Account 用户名 admin
	Account string `json:"account,omitempty"`
	// Password 密码 admin123
	Password string `json:"password,omitempty"`
	// BasicToken basic token -> [Basic base64(username:password)] "Basic YWRtaW46YWRtaW4xMjM="
	// 不要尝试去解密 数据库尽量存加密后的数据
	BasicToken string `json:"basicToken,omitempty"`
}
