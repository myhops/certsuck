package keytool

import (
	"bufio"
	"testing"
)

func TestYes(t *testing.T) {
	r := NewYesReader("")
	s := bufio.NewScanner(r)
	for i ,ok := 0, s.Scan(); i < 10 && ok; i, ok = i+1, s.Scan() {
		t.Logf("line %d: %s",i, s.Text())
	}
}
