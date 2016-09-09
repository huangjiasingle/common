package rest

import (
	"testing"
	"time"
)

func TestRest(t *testing.T) {
	for i := 0; i < 10; i++ {
		if _, err := Get("http://localhost:8000/console/about").End(); err != nil {
			t.Error(err)
		}
		time.Sleep(3*time.Second)
	}
}
