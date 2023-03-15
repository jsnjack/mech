package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"golang.org/x/exp/slices"
)

func Sync(expectedCollection, activeCollection []ResourceMatcher, doit, remove bool) error {
	report := table.NewWriter()

	// For tests, render data in csv format
	defer func() {
		if report.Length() > 0 {
			if reportToTestBuffer {
				report.RenderCSV()
			} else {
				report.Render()
			}
		} else {
			logger.Println("  nothing to do")
		}
	}()

	if reportToTestBuffer {
		// Skip header in tests
		report.SetOutputMirror(testBuffer)
	} else {
		report.SetOutputMirror(os.Stdout)
		report.AppendHeader(table.Row{"Action", "Resource", "Details"})
	}

	// Check if anything needs to be created / updated
	for _, r := range expectedCollection {
		expectedResource := r.(IExpectedResource)
		if logLevel > 0 {
			logger.Printf("Inspecting %q...\n", expectedResource.GetResourceID())
		}

		matchedResource := getMatchingResource(expectedResource, activeCollection)
		var activeResource IActiveResource
		if matchedResource != nil {
			activeResource = matchedResource.(IActiveResource)
		}

		action, diffs, err := Compare(expectedResource, activeResource)
		if err != nil {
			return err
		}
		if logLevel > 0 {
			logger.Printf("  status: %s\n", action)
		}
		if len(diffs) == 0 {
			report.AppendRow(table.Row{
				colorAction(action), expectedResource.GetResourceID(), "",
			})
		} else {
			for idx, diff := range diffs {
				if idx == 0 {
					report.AppendRow(table.Row{
						colorAction(action), expectedResource.GetResourceID(), diff.String(),
					})
				} else {
					report.AppendRow(table.Row{
						"", "", diff.String(),
					})
				}
			}
		}
		report.AppendSeparator()
		if doit {
			switch action {
			case ActionOK:
				break
			case ActionUpate:
				err := expectedResource.SyncResourceUpdate(activeResource.GetConstellixID())
				if err != nil {
					return err
				}
			case ActionCreate:
				err = expectedResource.SyncResourceCreate()
				if err != nil {
					return err
				}
			case ActionError:
				report.Render()
				os.Exit(1)
			default:
				return fmt.Errorf("unhandled action %q", action)
			}
		}
	}

	// Check if anything needs to be deleted
	for _, a := range activeCollection {
		activeResource := a.(IActiveResource)
		if logLevel > 0 {
			fmt.Printf("Inspecting %q...\n", activeResource.GetResourceID())
		}
		matched := getMatchingResource(activeResource, expectedCollection)
		if matched == nil {
			if logLevel > 0 {
				fmt.Printf("  status: %s\n", ActionDelete)
			}
			report.AppendRow(table.Row{
				colorAction(ActionDelete),
				activeResource.GetResourceID(),
				fmt.Sprintf("Resource ID %d", activeResource.GetConstellixID()),
			})
			report.AppendSeparator()
			if doit && remove {
				err := activeResource.SyncResourceDelete(activeResource.GetConstellixID())
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// generatePayload generates a JSON payload for a given Expected* resource
// which is send to constellix API endpoint
// Note: Costellix API is inconsistent. Sometimes it forces the inclusion of immutable fields
// in the payload and sometimes it refuses to process the request if one of the immutable fields
// is present in payload. To overcome it, excludedFieldsJSON is used.
func generatePayload(obj interface{}, definedFieldsJSON []string, excludedFieldsJSON []string) ([]byte, error) {
	objBytes, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	// Convert obj to map to simplify iteration
	dataIn := map[string]interface{}{}
	json.Unmarshal(objBytes, &dataIn)

	dataOut := map[string]interface{}{}

	// Create a new data obj which contains only fields which need to be included (JSON)
	for key, value := range dataIn {
		if slices.Contains(definedFieldsJSON, key) && !slices.Contains(excludedFieldsJSON, key) {
			switch key {
			case "ipfilter":
				// ipfilter is configured as DNSIPFilter, but sent as int
				// Note: when typecasting, pointer is map[string]interface {}
				record, ok := obj.(*ExpectedDNSRecord)
				if !ok {
					return nil, fmt.Errorf("expected *ExpectedDNSRecord, got %T", value)
				}
				if record.IPFilter != nil {
					dataOut[key] = record.IPFilter.ID
				} else {
					dataOut[key] = nil
				}
			default:
				dataOut[key] = value
			}
		}
	}

	dataOutBytes, err := json.Marshal(dataOut)
	if err != nil {
		return nil, err
	}
	return dataOutBytes, nil
}
