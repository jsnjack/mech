package cmd

import (
	"fmt"
	"reflect"
)

type IExpectedResource interface {
	GetDefinedStructFieldNames() []string
	GetResource() interface{}
	GetUID() string
	SyncResourceCreate() error
	SyncResourceDelete(int) error
	SyncResourceUpdate(int) (string, error)
}

type IActiveResource interface {
	GetConstellixID() int
	GetResource() interface{}
	GetUID() string
}

type ResourceMatcher interface {
	GetUID() string
}

// Compare compares expected resource with active resource
func Compare(expected IExpectedResource, active IActiveResource) (ResourceAction, error) {
	var action ResourceAction
	if active == nil {
		action = ActionCreate
	} else {
		action = ActionOK
		// reflect.Indirect deals with pointer
		expectedValue := reflect.Indirect(reflect.ValueOf(expected.GetResource()))
		activeValue := reflect.Indirect(reflect.ValueOf(active.GetResource()))
		for _, structFieldName := range expected.GetDefinedStructFieldNames() {
			fieldExpected := expectedValue.FieldByName(structFieldName)
			fieldActive := activeValue.FieldByName(structFieldName)
			// Compare field values
			if !reflect.DeepEqual(fieldExpected.Interface(), fieldActive.Interface()) {
				action = ActionUpate
				break
			}
		}
	}

	switch action {
	case ActionOK, ActionCreate, ActionUpate:
		return action, nil
	}
	return "", fmt.Errorf("unexpected action %q", action)
}
