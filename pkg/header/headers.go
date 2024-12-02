package header

var ipHeaders = []string{
	"X-ProxyUser-Ip",
	"Client-IP",
	"X-Client-IP",
	"X-Originating-IP",
	"X-Real-IP",
	"X-Forwarded-For",
	"X-Remote-IP",
	"X-Remote-Addr",
	"Forwarded-For",
	"True-Client-IP",
	"X-Custom-IP-Authorization",
	"X-Forwarded",
	"X-Host",
}

var pathHeaders = []string{
	"X-Original-URL",
	"X-Rewrite-URL",
}
