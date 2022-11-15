package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/juju/ansiterm"
)

func Sync(expectedCollection, activeCollection []ResourceMatcher, doit, remove bool) error {
	report := ansiterm.NewTabWriter(os.Stdout, 10, 0, 2, ' ', tabwriter.Debug)
	defer report.Flush()

	// Check if anything needs to be created / updated
	for _, r := range expectedCollection {
		expectedResource := r.(IExpectedResource)
		fmt.Printf("Inspecting %q...\n", expectedResource.GetUID())
		activeResource := getMatchingResource(expectedResource, activeCollection)
		action, details, err := Compare(expectedResource, activeResource)
		if err != nil {
			return err
		}
		fmt.Printf("  status: %s\n", action)
		fmt.Fprintf(report, "%s\t%s\t%s\n", colorAction(action), expectedResource.GetUID(), details)
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
			default:
				return fmt.Errorf("unhandled action %q", action)
			}
		}
	}

	// Check if anything needs to be deleted
OUTER:
	for _, a := range activeCollection {
		activeResource := a.(IActiveResource)
		fmt.Printf("Inspecting %q...\n", activeResource.GetUID())
		var expectedResource IExpectedResource
		for _, expectedResource := range expectedCollection {
			if expectedResource.GetUID() == activeResource.GetUID() {
				continue OUTER
			}
		}
		fmt.Printf("  status: %s\n", ActionDelete)
		fmt.Fprintf(
			report, "%s\t%s\t%s\n",
			colorAction(ActionDelete),
			activeResource.GetUID(),
			fmt.Sprintf("Resource ID %d", activeResource.GetConstellixID()),
		)
		if doit && remove {
			err := expectedResource.SyncResourceDelete(activeResource.GetConstellixID())
			if err != nil {
				return err
			}
		} else {
			fmt.Printf("  pass --remove flag to remove %q\n", activeResource.GetUID())
		}

	}
	return nil
}
