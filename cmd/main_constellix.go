package cmd

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"strconv"
	"time"
)

// buildSecurityToken returns security token which is used when authenticating
// Constellix REST API requests
func buildSecurityToken() string {
	millis := time.Now().UnixNano() / 1000000
	timestamp := strconv.FormatInt(millis, 10)
	mac := hmac.New(sha1.New, []byte(constellixSecretKey))
	mac.Write([]byte(timestamp))
	hmacstr := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return constellixAPIKey + ":" + hmacstr + ":" + timestamp
}
