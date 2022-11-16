package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/juju/ansiterm"
	"golang.org/x/exp/slices"
)

func Sync(expectedCollection, activeCollection []ResourceMatcher, doit, remove bool) error {
	report := ansiterm.NewTabWriter(os.Stdout, 10, 0, 2, ' ', tabwriter.Debug)
	defer report.Flush()

	// Check if anything needs to be created / updated
	for _, r := range expectedCollection {
		expectedResource := r.(IExpectedResource)
		if rootVerbose {
			fmt.Printf("Inspecting %q...\n", expectedResource.GetResourceID())
		}
		activeResource := getMatchingResource(expectedResource, activeCollection)
		action, details, err := Compare(expectedResource, activeResource)
		if err != nil {
			return err
		}
		if rootVerbose {
			fmt.Printf("  status: %s\n", action)
		}
		fmt.Fprintf(report, "%s\t%s\t%s\n", colorAction(action), expectedResource.GetResourceID(), details)
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
				report.Flush()
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
			fmt.Fprintf(
				report, "%s\t%s\t%s\n",
				colorAction(ActionDelete),
				activeResource.GetResourceID(),
				fmt.Sprintf("Resource ID %d", activeResource.GetConstellixID()),
			)
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
