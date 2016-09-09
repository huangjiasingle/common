package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	code := "new test"
	e := New(code)
	if e.Key() != code {
		t.Fatalf("want:%s,but:%s", code, e.Key())
		return
	}
	println(e.Error())

	eAs := As(e, "this is value copy for As, not same as before")
	if eAs.Key() != code {
		t.Fatalf("want:%s,but:%s", code, eAs.Key())
		return
	}
	if eAs == e {
		t.Fatalf("%s, same as:%s", eAs, e)
		return
	}
	println(eAs.Error())

	ePar := ParseErr(e)
	if ePar.Key() != code {
		t.Fatalf("want:%s,but:%s", code, ePar.Key())
		return
	}
	if ePar == e {
		t.Fatalf("%s, same as:%s", ePar, e)
		return
	}
	println(ePar.Error())
}

var equalTests = []struct {
	err1 Err
	err2 error
	out  bool
}{
	{New("New"), New("New"), true},
	{New("New"), ParseErr(New("New")), true},
	{New("New"), ParseErr(errors.New("New")), true},
	{New("New"), As(New("New")), true},
	{New("New"), As(errors.New("New")), true},
	{New("New"), As(New("New"), "reason"), true},
	{New("New"), As(errors.New("New"), "reason"), true},
	{ParseErr(New("ParseErr")), ParseErr(New("ParseErr")), true},
	{ParseErr(New("ParseErr")), ParseErr(errors.New("ParseErr")), true},
	{ParseErr(New("ParseErr")), As(New("ParseErr")), true},
	{ParseErr(New("ParseErr")), As(errors.New("ParseErr")), true},
	{ParseErr(New("ParseErr")), As(New("ParseErr"), "reason"), true},
	{ParseErr(New("ParseErr")), As(errors.New("ParseErr"), "reason"), true},
}

func TestEqual(t *testing.T) {
	for index, test := range equalTests {
		if test.err1.Equal(test.err2) != test.out {
			println(test.err1.Error())
			println(test.err2.Error())
			t.Fatalf("usercase %d,want:%s,but:%s", index, !test.out, test.out)
			return
		}
		if Equal(test.err1, test.err2) != test.out {
			println(test.err1.Error())
			println(test.err2.Error())
			t.Fatalf("usercase %d,want:%s,but:%s", index, !test.out, test.out)
			return
		}
	}
}

func TestAs(t *testing.T) {
	err1 := New("test")
	err2 := New("test")

	outErr1 := As(err1, "test")
	outErr2 := err2.As("test", t)
	fmt.Println(outErr1.Error())
	fmt.Println(outErr2.Error())
}
