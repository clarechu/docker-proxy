package models

// Logging 请求日志记录
type Logging struct {
	// ClientIP 客户端ip
	ClientIP string `json:"clientIP"`
	// 请求url
	URI string `json:"uri"`
	// 请求方法
	Method string `json:"method"`
	// Header 请求头
	Header string `json:"header"`
	// RequestBody 请求 body
	RequestBody string `json:"requestBody"`
	// ReturnTime 请求接口的时长
	ReturnTime string `json:"returnTime"`
	// HttpStatusCode 响应code
	HttpStatusCode string `json:"httpStatusCode"`
	// ResponseBody 响应 body
	ResponseBody string `json:"responseBody"`
}
