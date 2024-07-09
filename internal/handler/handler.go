package handler

import (
	"fmt"
	"ip-rate-control/pkg/ip"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	clientIP := ip.GetClientIP(r)
	response := fmt.Sprintf("Your IP address is %s", clientIP)
	w.Write([]byte(response))
}
