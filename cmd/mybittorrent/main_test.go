package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_decodeBencode(t *testing.T) {
	testData := map[string]interface{}{
		"9:pineapple":      "pineapple",
		"i18264343e":       18264343,
		"lli983e5:grapeee": []interface{}{983, "grape"},
		"l6:orangei463ee":  []interface{}{"orange", 463},
		"le":               []interface{}{},
	}

	for k, v := range testData {
		gotData, gotErr := decodeBencode(k)
		if gotErr != nil {
			assert.FailNow(t, "err is not nil", gotErr)
		}
		assert.Equal(t, gotData, v)
	}
}
