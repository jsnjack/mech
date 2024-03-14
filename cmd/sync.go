package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"golang.org/x/exp/slices"
)

func Sync(expectedCollection, activeCollection []ResourceMatcher, doit, remove bool, title string) error {
	report := table.NewWriter()

	toDelete := []IActiveResource{}
	toUpdate := map[IExpectedResource]int{}
	toCreate := []IExpectedResource{}

	if reportToTestBuffer {
		// Skip header in tests
		report.SetOutputMirror(testBuffer)
	} else {
		report.SetOutputMirror(os.Stdout)
		if title != "" {
			report.SetTitle(title)
		}
		report.AppendHeader(table.Row{"Action", "Resource", "Details"})
	}

	// Check if anything needs to be deleted first
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
			toDelete = append(toDelete, activeResource)
		}
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
		switch action {
		case ActionOK:
		case ActionUpate:
			toUpdate[expectedResource] = activeResource.GetConstellixID()
		case ActionCreate:
			toCreate = append(toCreate, expectedResource)
		case ActionError:
			report.Render()
			os.Exit(1)
		default:
			return fmt.Errorf("unhandled action %q", action)
		}
	}

	printReport(report)
	logger.Printf("SUMMARY: %d to delete, %d to update, %d to create\n", len(toDelete), len(toUpdate), len(toCreate))

	if doit {
		if !remove && len(toDelete) > 0 {
			return fmt.Errorf("resource deletion is not allowed. Use --remove flag to allow it")
		}
		logger.Println("Syncing changes...")
		err := syncChanges(toDelete, toUpdate, toCreate)
		if err != nil {
			return err
		}
	}
	return nil
}

func printReport(report table.Writer) {
	if report.Length() > 0 {
		if reportToTestBuffer {
			// Skip header in tests to simplify testing
			report.RenderCSV()
		} else {
			report.Render()
		}
	}
}

func syncChanges(toDelete []IActiveResource, toUpdate map[IExpectedResource]int, toCreate []IExpectedResource) error {
	// First, we delete resources
	// If the resource is DNSRecord, we must first remove the ones with geoproximities.
	// They will not end with " 0)"
	// Note: the sorting is very simple and might affect other resources
	sort.Slice(toDelete, func(i, j int) bool {
		iEndsWithZero := strings.HasSuffix(toDelete[i].GetResourceID(), " 0)")
		jEndsWithZero := strings.HasSuffix(toDelete[j].GetResourceID(), " 0)")

		// If both end with " 0)" or both don't, sort by the resource ID
		if iEndsWithZero == jEndsWithZero {
			return toDelete[i].GetResourceID() < toDelete[j].GetResourceID()
		}

		// If i ends with " 0)" and j doesn't, i should come after j
		if iEndsWithZero && !jEndsWithZero {
			return false
		}

		// If j ends with " 0)" and i doesn't, i should come before j
		return true
	})
	for _, resource := range toDelete {
		err := resource.SyncResourceDelete(resource.GetConstellixID())
		if err != nil {
			return err
		}
	}

	// Then, we update resources
	for resource, constellixID := range toUpdate {
		err := resource.SyncResourceUpdate(constellixID)
		if err != nil {
			return err
		}
	}

	// Finally, we create resources
	for _, resource := range toCreate {
		err := resource.SyncResourceCreate()
		if err != nil {
			return err
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
			dataOut[key] = value
		}
	}

	dataOutBytes, err := json.Marshal(dataOut)
	if err != nil {
		return nil, err
	}
	return dataOutBytes, nil
}
