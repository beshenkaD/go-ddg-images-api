package goddgimagesapi

import (
	"testing"
)

func Test(t *testing.T) {
	r, err := Do("duck", false)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if len(r.Results) == 0 {
		t.Fail()
	}
}
