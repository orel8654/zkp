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

func TestGet(t *testing.T) {
	var testCase = []struct {
		key  string
		val  int64
		want int64
	}{
		{"test1", 30, 30},
		{"test2", 40, 40},
		{"test3", 50, 50},
	}
	rT := MyNewRedis("localhost:6379")

	for _, elem := range testCase {
		name := fmt.Sprintf("case(%s,%d)", elem.key, elem.val)
		t.Run(name, func(t *testing.T) {
			got, _ := rT.GetVal(elem.key)
			if got != elem.want {
				t.Errorf("got %d, want %d", got, elem.want)
			}
		})
	}
}
