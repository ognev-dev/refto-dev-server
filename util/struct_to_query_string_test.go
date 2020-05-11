package util_test

import (
	"fmt"
	"testing"

	"github.com/refto/server/test/assert"
	"github.com/refto/server/util"
)

func TestStructToQueryString(t *testing.T) {
	type A struct {
		B string `json:"b"`
		C int    `json:"c"`
	}
	st := A{
		B: "C",
		C: 1,
	}

	result, err := util.StructToQueryString(st)
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Equals(t, fmt.Sprintf(`b=%s&c=%d`, st.B, st.C), result)
}
