package redirect

import (
	"net/http"
	"strings"
)

const modirumVPOSUserAgent = "Modirum VPOS"

// IsServerToServerCallback reports whether the request is likely the delayed background
// confirmation POST (User-Agent Modirum VPOS) described in Cardlink redirect documentation.
func IsServerToServerCallback(r *http.Request) bool {
	if r == nil {
		return false
	}
	ua := r.Header.Get("User-Agent")
	return strings.Contains(ua, modirumVPOSUserAgent)
}
