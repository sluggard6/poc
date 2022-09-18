package util

import (
	"testing"
)

func TestUUID(t *testing.T) {
	t.Log(UUID())
}

func TestMd5(t *testing.T) {
	var (
		in       = "abc123"
		expected = "e99a18c428cb38d5f260853678922e03"
	)
	actual := Md5String(in)
	if actual != expected {
		t.Errorf("Md5String(%s) = %s; expected %s", in, actual, expected)
	}
}
