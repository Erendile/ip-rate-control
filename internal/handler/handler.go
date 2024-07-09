package handler

import (
	"fmt"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	response := fmt.Sprintf("Your IP address is %s", clientIP)
	w.Write([]byte(response))
}
