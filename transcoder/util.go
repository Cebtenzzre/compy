package transcoder

import (
	"net/http"
	"strings"
)

func SupportsWebP(headers http.Header) bool {
	accept := headers.Get("Accept")
	if accept == "" {
		return true
	}
	for _, v := range strings.Split(accept, ",") {
		mimeType := strings.SplitN(v, ";", 2)[0]
		if mimeType == "*/*" || mimeType == "image/*" || mimeType == "image/webp" {
			return true
		}
	}
	return false
}
