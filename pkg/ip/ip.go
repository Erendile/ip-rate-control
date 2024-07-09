package ip

import (
	"net"
	"net/http"
	"strings"
)

var requestHeaders = []string{
	"X-Client-IP",
	"X-Forwarded-For",
	"Cf-Connecting-IP",
	"Fastly-Client-IP",
	"True-Client-IP",
	"X-Real-IP",
	"X-Cluster-Client-IP",
	"X-Forwarded",
	"Forwarded-For",
	"Forwarded",
}

func GetClientIP(r *http.Request) string {
	for _, header := range requestHeaders {
		switch header {
		case "X-Forwarded-For":
			if host, correctIP := getClientIPFromXForwardedFor(r.Header.Get(header)); correctIP {
				return host
			}
		default:
			if host := r.Header.Get(header); isCorrectIP(host) {
				return host
			}
		}
	}

	host, _, splitHostPortError := net.SplitHostPort(r.RemoteAddr)
	if splitHostPortError == nil && isCorrectIP(host) {
		return host
	}
	return ""
}

func getClientIPFromXForwardedFor(headers string) (string, bool) {
	if headers == "" {
		return "", false
	}

	forwardedIps := strings.Split(headers, ",")
	for _, ip := range forwardedIps {
		ip = strings.TrimSpace(ip)
		if isCorrectIP(ip) {
			return ip, true
		}
	}
	return "", false
}

func isCorrectIP(ip string) bool {
	return net.ParseIP(ip) != nil
}
