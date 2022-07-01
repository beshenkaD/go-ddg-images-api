package goddgimagesapi

import (
	"testing"
)

func Test(t *testing.T) {
	_, err := token("duck")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	res, err := Do(Query{
		Keywords: "duck",
		Moderate: false,
	})
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if len(res.Results) == 0 {
		t.Log("no results!")
		t.Fail()
	}
}
