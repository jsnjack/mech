package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/sergi/go-diff/diffmatchpatch"
)

const rateLimitWaitTime = 5
const textSpacer = 15

const Reset = "\033[0m"
const Red = "\033[31m"
const Green = "\033[32m"
const Yellow = "\033[33m"
const Blue = "\033[34m"
const Purple = "\033[35m"
const Cyan = "\033[36m"
const Gray = "\033[37m"
const White = "\033[97m"
const Crossed = "\033[9m"

func colorAction(action ResourceAction) string {
	var start string
	switch action {
	case ActionCreate:
		start = Green
	case ActionDelete:
		start = Red
	case ActionUpate:
		start = Yellow
	case ActionOK:
		start = White
	case ActionError:
		start = Purple
	}
	return start + string(action) + Reset
}

func colorStatus(status ResourceRuntimeStatus) string {
	var start string
	switch status {
	case StatusUp:
		start = Green
	case StatusDown:
		start = Red
	}
	return start + string(status) + Reset
}

var colorRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripBashColors(s string) string {
	return colorRe.ReplaceAllString(s, "")
}

// getFieldNamesMap returns struct field names from their tags
// (yaml/JSON keys to struct keys)
func getFieldNamesMap(obj interface{}, tagType string, tags ...string) map[string]string {
	res := make(map[string]string)
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	fields := reflect.VisibleFields(t)
OUTER:
	for _, tag := range tags {
		for _, f := range fields {
			val, ok := f.Tag.Lookup(tagType)
			if ok {
				tagName := strings.Split(val, ",")[0]
				if tagName == tag {
					res[tag] = f.Name
					continue OUTER
				}
			}
		}
	}
	return res
}

// getTag returns the tag value for a given field name
func getTag(obj interface{}, fieldName string, tagType string) string {
	t, ok := obj.(reflect.Type)
	if !ok {
		return ""
	}
	fields := reflect.VisibleFields(t)
	for _, f := range fields {
		val, ok := f.Tag.Lookup(tagType)
		if ok {
			if f.Name == fieldName {
				return strings.Split(val, ",")[0]
			}
		}
	}
	return ""
}

// makeSimpleAPIRequest makes a simple API request, normally to the Sonar API as
// it doesn't support pagination
func makeSimpleAPIRequest(method string, url string, payload io.Reader, expectedStatusCode int) (respBody []byte, err error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return nil, err
	}
	req.Header.Add("x-cns-security-token", buildSecurityToken())
	req.Header.Add("Content-Type", "application/json")
	if rootVerbose {
		logger.Printf("  requesting %s %s ...\n", method, url)
		if payload != nil {
			payloadBytes, err := io.ReadAll(payload)
			if err != nil {
				return nil, err
			}
			logger.Println("  payload: " + string(payloadBytes))
		} else {
			logger.Println("  no payload")
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		logger.Printf("Rate limit exceeded, waiting %d second...\n", rateLimitWaitTime)
		time.Sleep(time.Duration(rateLimitWaitTime) * time.Second)
		return makeSimpleAPIRequest(method, url, payload, expectedStatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != expectedStatusCode {
		return body, fmt.Errorf("unexpected status code %d, want %d", resp.StatusCode, expectedStatusCode)
	}
	if rootDebug {
		logger.Println(method, url, resp.StatusCode)
		logger.Println(string(body))
	}
	return body, nil
}

// makev4APIRequest makes a request to the v4 API, which supports pagination.
// It returns a slice of response bodies, one for each page.
func makev4APIRequest(method string, url string, payload io.Reader, expectedStatusCode int) (respBodys [][]byte, err error) {
	next := true
	for next {
		data, err := makeSimpleAPIRequest(method, url, payload, expectedStatusCode)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve resource: %s", err)
		}
		if len(data) != 0 {
			resp := DNSv4Response{}
			err = json.Unmarshal(data, &resp)
			if err != nil {
				return nil, err
			}
			respBodys = append(respBodys, resp.Data)
			if resp.Meta.Links.Next != "" {
				url = resp.Meta.Links.Next
			} else {
				next = false
			}
		} else {
			next = false
		}
	}
	return respBodys, nil
}

func getMatchingResource(item ResourceMatcher, collection []ResourceMatcher) interface{} {
	for _, el := range collection {
		if item.GetResourceID() == el.GetResourceID() {
			return el
		}
	}
	return nil
}

func DiffPrettyText(diffs []diffmatchpatch.Diff) string {
	var buff bytes.Buffer
	for idx, diff := range diffs {
		text := diff.Text

		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			_, _ = buff.WriteString(Green)
			_, _ = buff.WriteString(text)
			_, _ = buff.WriteString(Reset)
		case diffmatchpatch.DiffDelete:
			_, _ = buff.WriteString(Red + Crossed)
			_, _ = buff.WriteString(text)
			_, _ = buff.WriteString(Reset)
		case diffmatchpatch.DiffEqual:
			if idx == 0 {
				// the beginning is the same, we can skip it
				textRune := []rune(text)
				if len(textRune) > textSpacer {
					text = "..." + string(textRune[len(textRune)-textSpacer:])
				}
				_, _ = buff.WriteString(text)
			} else if idx == len(diffs)-1 {
				// the end is the same, we can skip it
				textRune := []rune(text)
				if len(textRune) > textSpacer {
					text = string(textRune[0:textSpacer]) + "..."
				}
				_, _ = buff.WriteString(text)
			} else {
				_, _ = buff.WriteString(text)
			}
		}

	}

	return buff.String()
}
