package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"
)

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

var colorRe = regexp.MustCompile(`\033\[[0-9;]*m`)

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

// Remove color codes from string, used in tests
func stripColor(s string) string {
	return colorRe.ReplaceAllString(s, "")
}

// getFieldNamesMap returns struct field names from their tags
// JSON to struct
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

func makeAPIRequest(method string, url string, payload io.Reader, expectedStatusCode int) (respBody []byte, err error) {
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
		fmt.Printf("  requesting %s %s ...\n", method, url)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != expectedStatusCode {
		return body, fmt.Errorf("unexpected status code %d, want %d", res.StatusCode, expectedStatusCode)
	}
	return body, nil
}

func getMatchingResource(item ResourceMatcher, collection []ResourceMatcher) IActiveResource {
	for _, el := range collection {
		if item.GetResourceID() == el.GetResourceID() {
			ar, ok := el.(IActiveResource)
			if !ok {
				fmt.Printf("failed type assertion %q\n", ar)
				os.Exit(1)
			}
			return ar
		}
	}
	return nil
}
