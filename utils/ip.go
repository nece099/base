package utils

import (
	"net/http"
	"strings"
)

func PeerIP(r *http.Request) string {
	fval := r.Header.Get("X-Forwarded-For")
	if len(fval) == 0 {
		rval := r.Header.Get("X-Real-IP")
		if len(rval) == 0 {
			raddr := r.RemoteAddr
			ip := strings.Split(raddr, ":")
			return ip[0]
		}
		addresses := strings.Split(rval, ",")
		address := strings.TrimSpace(addresses[0])
		return address
	}

	addresses := strings.Split(fval, ",")
	address := strings.TrimSpace(addresses[0])
	return address
}
