package digest

import (
	"crypto/sha256"
	"encoding/base64"
)

// VPOS21 computes Base64(SHA256(c14nMessage || secret)) per VPOS XML 2.1 shared-secret digest rules.
func VPOS21(c14nMessage []byte, secret string) string {
	h := sha256.Sum256(append(c14nMessage, []byte(secret)...))
	return base64.StdEncoding.EncodeToString(h[:])
}
