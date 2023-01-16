package _func

import (
	"reflect"
	"testing"
)

func TestFunction(t *testing.T) {
	result := calculate(1, 2)

	if !reflect.DeepEqual(result, 3) {
		t.Errorf("Error")
		return
	}
}
