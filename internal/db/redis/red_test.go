package redis

import (
	"fmt"
	"testing"
)

func TestSave(t *testing.T) {
	var testCase = []struct {
		key  string
		val  int64
		want error
	}{
		{"test1", 30, nil},
		{"test2", 40, nil},
		{"test3", 50, nil},
	}
	rT := MyNewRedis("localhost:6379")

	for _, elem := range testCase {
		name := fmt.Sprintf("case(%s,%d)", elem.key, elem.val)
		t.Run(name, func(t *testing.T) {
			got := rT.SaveVal(elem.key, elem.val)
			if got != elem.want {
				t.Errorf("got %s, want %d", got, elem.want)
			}
		})
	}
}
