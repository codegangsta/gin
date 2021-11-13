package gin_test

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

func getCaller() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "Unrecoverable location"
	}
	return fmt.Sprintf("%s:%d", file, line)
}

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v) at %s", b, reflect.TypeOf(b), a, reflect.TypeOf(a), getCaller())
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v) at %s", b, reflect.TypeOf(b), a, reflect.TypeOf(a), getCaller())
	}
}
