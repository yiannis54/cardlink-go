package digest

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
)

// Redirect computes Base64(SHA256(utf8(preimage + secret))).
func Redirect(preimage, secret string) string {
	h := sha256.Sum256([]byte(preimage + secret))
	return base64.StdEncoding.EncodeToString(h[:])
}

// RedirectSHA1 computes Base64(SHA1(utf8(preimage + secret))) for legacy recurring notifications without version=2.
func RedirectSHA1(preimage, secret string) string {
	h := sha1.Sum([]byte(preimage + secret))
	return base64.StdEncoding.EncodeToString(h[:])
}
