package cmd

import (
	"fmt"
	"reflect"
	"strconv"

	"golang.org/x/exp/slices"
)

// IExpectedResource implements expected resource which is extracted from the
// local configuration
type IExpectedResource interface {
	// List of fields defined in configuration file. Only this fields will be
	// synchronized with the remote
	GetDefinedStructFieldNames() []string
	// Constellix will reject requests with changes in immutable fields. Delete
	// and create a new resource in this case
	GetImmutableStructFields() []string
	// Returns the resource itself
	GetResource() interface{}
	// Returns struct field name with which resources will be matched
	GetResourceID() string
	// Create a new active resource via API from expected resource
	SyncResourceCreate() error
	// Update active resource with all fields prsented in configuration but
	// immutable ones
	SyncResourceUpdate(int) error
}

// IActiveResource implements remote / active resource
type IActiveResource interface {
	// Returns ID of the resource in Constellix
	GetConstellixID() int
	GetResource() interface{}
	GetResourceID() string
	// Delete active resource
	SyncResourceDelete(int) error
}

// ResourceMatcher implements resources to compare
type ResourceMatcher interface {
	GetResourceID() string
}

type FieldDiff struct {
	FieldName string
	OldValue  string
	NewValue  string
}

func (f *FieldDiff) String() string {
	return fmt.Sprintf(
		"%s: %s%s",
		f.FieldName,
		Red+Crossed+f.OldValue+Reset,
		Green+f.NewValue+Reset,
	)
}

func getFieldDiff(diffs []*FieldDiff, fieldName string) *FieldDiff {
	for _, diff := range diffs {
		if diff.FieldName == fieldName {
			return diff
		}
	}
	return nil
}

// Compare compares expected resource with active resource
func Compare(expected IExpectedResource, active IActiveResource) (ResourceAction, []*FieldDiff, error) {
	var action ResourceAction
	diffs := make([]*FieldDiff, 0)
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
				diffs = append(diffs, &FieldDiff{
					FieldName: structFieldName,
					OldValue:  valueToString(fieldActive),
					NewValue:  valueToString(fieldExpected),
				})
				if slices.Contains(expected.GetImmutableStructFields(), structFieldName) {
					return ActionError, make([]*FieldDiff, 0), fmt.Errorf("found change in immutable field: %s", structFieldName)
				}
			}
		}
	}

	switch action {
	case ActionOK, ActionCreate, ActionUpate:
		return action, diffs, nil
	}
	return "", make([]*FieldDiff, 0), fmt.Errorf("unexpected action %q", action)
}

func toResourceMatcher(collection interface{}) []ResourceMatcher {
	v := reflect.ValueOf(collection)
	// No check here, just panic!
	new := make([]ResourceMatcher, v.Len())
	for i := 0; i < v.Len(); i++ {
		new[i] = v.Index(i).Interface().(ResourceMatcher)
	}
	return new
}

// valueToString returns a textual representation of the reflection value val.
// Based on function from reflect library
// https://cs.opensource.google/go/go/+/master:src/reflect/tostring_test.go
func valueToString(val reflect.Value) string {
	var str string
	if !val.IsValid() {
		return "<zero Value>"
	}
	typ := val.Type()
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(val.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(val.Float(), 'g', -1, 64)
	case reflect.Complex64, reflect.Complex128:
		c := val.Complex()
		return strconv.FormatFloat(real(c), 'g', -1, 64) + "+" + strconv.FormatFloat(imag(c), 'g', -1, 64) + "i"
	case reflect.String:
		return val.String()
	case reflect.Bool:
		if val.Bool() {
			return "true"
		} else {
			return "false"
		}
	case reflect.Pointer:
		v := val
		str = typ.String() + "("
		if v.IsNil() {
			str += "0"
		} else {
			str += "&" + valueToString(v.Elem())
		}
		str += ")"
		return str
	case reflect.Array, reflect.Slice:
		v := val
		str += "["
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				str += ", "
			}
			str += valueToString(v.Index(i))
		}
		str += "]"
		return str
	case reflect.Map:
		t := typ
		str = t.String()
		str += "{"
		str += "<can't iterate on maps>"
		str += "}"
		return str
	case reflect.Chan:
		str = typ.String()
		return str
	case reflect.Struct:
		t := typ
		v := val
		str += t.String()
		str += "{"
		for i, n := 0, v.NumField(); i < n; i++ {
			if i > 0 {
				str += ", "
			}
			str += valueToString(v.Field(i))
		}
		str += "}"
		return str
	case reflect.Interface:
		return typ.String() + "(" + valueToString(val.Elem()) + ")"
	case reflect.Func:
		v := val
		return typ.String() + "(" + strconv.FormatUint(uint64(v.Pointer()), 10) + ")"
	default:
		return fmt.Sprintf("%s", val.Interface())
	}
}
