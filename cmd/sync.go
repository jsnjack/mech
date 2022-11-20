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
	defer func() {
		if report.Length() > 0 {
			report.Render()
		} else {
			fmt.Println("  nothing to do")
		}
	}()

	report.SetOutputMirror(os.Stdout)
	report.AppendHeader(table.Row{"Action", "Resource", "Details"})

	// Check if anything needs to be created / updated
	for _, r := range expectedCollection {
		expectedResource := r.(IExpectedResource)
		if rootVerbose {
			fmt.Printf("Inspecting %q...\n", expectedResource.GetResourceID())
		}
		activeResource := getMatchingResource(expectedResource, activeCollection)
		action, diffs, err := Compare(expectedResource, activeResource)
		if err != nil {
			return err
		}
		if rootVerbose {
			fmt.Printf("  status: %s\n", action)
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
		if rootVerbose {
			fmt.Printf("Inspecting %q...\n", activeResource.GetResourceID())
		}
		matched := getMatchingResource(activeResource, expectedCollection)
		if matched == nil {
			if rootVerbose {
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

func generatePayload(obj interface{}, definedFieldsJSON, immutableFieldsJSON []string) ([]byte, error) {
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
		if slices.Contains(definedFieldsJSON, key) && !slices.Contains(immutableFieldsJSON, key) {
			dataOut[key] = value
		}
	}

	dataOutBytes, err := json.Marshal(dataOut)
	if err != nil {
		return nil, err
	}
	return dataOutBytes, nil
}
