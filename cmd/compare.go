package cmd

import (
	"fmt"
	"reflect"
)

type IExpectedResource interface {
	GetDefinedStructFieldNames() []string
	GetCreateData() ([]byte, error)
	GetUpdateData() ([]byte, error)
	GetResource() interface{}
	GetUID() string
	SyncResourceCreate([]byte) error
}

type IActiveResource interface {
	GetResource() interface{}
	GetConstellixID() int
	GetUID() string
	SyncResourceUpdate([]byte) error
	SyncResourceDelete() error
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
