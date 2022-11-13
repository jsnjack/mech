package cmd

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"strconv"
	"time"
)

// List of actions to sync
type ResourceAction string

const ActionCreate ResourceAction = "create"
const ActionDelete ResourceAction = "delete"
const ActionUpate ResourceAction = "update"
const ActionOK ResourceAction = "ok"

const sonarRESTAPIBaseURL = "https://api.sonar.constellix.com/rest/api"

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
