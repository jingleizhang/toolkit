package hash

import (
	"testing"
)

func TestTime33(t *testing.T) {
	str := "hello"
	count := 5

	k := Time33(str, uint32(count))

	if k != 3 {
		t.Error()
	}
}
