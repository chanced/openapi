package openapi

import "strings"

var (
	tokenDecodeReplacer = strings.NewReplacer("~1", "/", "~0", "~")
	tokenEncodeReplacer = strings.NewReplacer("/", "~1", "~", "~0")
)

func decodeJSONPtr(token string) string {
	return tokenDecodeReplacer.Replace(token)
}

func encodeJSONPtr(token string) string {
	return tokenEncodeReplacer.Replace(token)
}
