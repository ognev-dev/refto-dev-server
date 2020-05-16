package githubwebhook

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
)

const SigPrefix = "sha1="

// HashMAC generates hmac of given body and key
// https://developer.github.com/webhooks/event-payloads/#delivery-headers
// The HMAC hex digest of the response body.
// This header will be sent if the webhook is configured with a secret.
// The HMAC hex digest is generated using the sha1 hash function and the secret as the HMAC key.
// Note: github also adds a prefix "sha1=" to the hash
func HashMAC(body, key string) string {
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(body))

	return SigPrefix + fmt.Sprintf("%+x", h.Sum(nil))
}

func ValidMAC(body []byte, hash, key string) bool {
	h := hmac.New(sha1.New, []byte(key))
	h.Write(body)
	expected := fmt.Sprintf("%x", h.Sum(nil))

	return hash == (SigPrefix + expected)
}
