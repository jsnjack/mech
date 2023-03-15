package cmd

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"time"
)

// List of actions to sync
type ResourceAction string

const ActionCreate ResourceAction = "create"
const ActionDelete ResourceAction = "delete"
const ActionUpate ResourceAction = "update"
const ActionOK ResourceAction = "ok"
const ActionError ResourceAction = "error"

var sonarRESTAPIBaseURL string = "https://api.sonar.constellix.com/rest/api"
var dnsRESTAPIBaseURL string = "https://api.dns.constellix.com/v4"

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

// Runtime status of a resource
type ResourceRuntimeStatus string

const StatusUp ResourceRuntimeStatus = "UP"
const StatusDown ResourceRuntimeStatus = "DOWN"

// RuntimeStatus is a struct that represents the runtime status of a resource
type RuntimeStatus struct {
	Status ResourceRuntimeStatus `json:"status"`
}

type DNSv4Response struct {
	Data json.RawMessage `json:"data"`
	Meta v4ResponseMeta  `json:"meta"`
}

type v4ResponseMeta struct {
	Pagination v4MetaPagination `json:"pagination"`
	Links      v4MetaLinks      `json:"links"`
}

type v4MetaPagination struct {
	Total       int `json:"total"`
	Count       int `json:"count"`
	PerPage     int `json:"perPage"`
	CurrentPage int `json:"currentPage"`
	TotalPages  int `json:"totalPages"`
}

type v4MetaLinks struct {
	Self     string `json:"self"`
	First    string `json:"first"`
	Last     string `json:"last"`
	Previous string `json:"previous"`
	Next     string `json:"next"`
}
