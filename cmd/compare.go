package cmd

import (
	"fmt"
	"reflect"
)

type IExpectedResource interface {
	GetCreateData() ([]byte, error)
	GetDefinedStructFieldNames() []string
	GetResource() interface{}
	GetUID() string
	GetUpdateData() ([]byte, error)
	SyncResourceCreate([]byte) error
	SyncResourceDelete(int) error
	SyncResourceUpdate([]byte, int) error
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
func Compare(expected IExpectedResource, active IActiveResource) (ResourceAction, []byte, error) {
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

		if action == ActionOK {
			return action, nil, nil
		}
	}

	switch action {
	case ActionOK:
		return action, nil, nil
	case ActionCreate:
		data, err := expected.GetCreateData()
		if err != nil {
			return "", nil, err
		}
		return action, data, nil
	case ActionUpate:
		data, err := expected.GetUpdateData()
		if err != nil {
			return "", nil, err
		}
		return action, data, nil
	}
	return "", nil, fmt.Errorf("unexpected action %q", action)
}
