package util

import (
	"testing"
)

func TestAnyEmpty(t *testing.T) {
	params := []struct {
		input  []string
		is_err bool
	}{
		{
			input:  []string{"not empty", "not empty"},
			is_err: false,
		},
		{
			input:  []string{"not empty", ""},
			is_err: true,
		},
	}

	for _, param := range params {
		err := AnyEmpty(param.input)
		if param.is_err && err == nil {
			t.Errorf("wanted err but got %v for %v", err, param.input)
		}
	}
}
