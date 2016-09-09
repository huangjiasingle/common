package log

import (
	"testing"
)

func TestLog(t *testing.T) {
	l := New("name")
	l.Warn("a", "b")
	l.Info("c", "d")
}
