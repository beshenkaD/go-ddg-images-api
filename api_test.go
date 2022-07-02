package goddgimagesapi

import (
	"net/http"
	"testing"
)

func Test(t *testing.T) {
	client := NewClient(http.DefaultClient)

	result, err := client.Do(Query{
		Keywords: "duck",
		Moderate: true,
	})

	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if len(result.Results) == 0 {
		t.Fail()
	}
}
