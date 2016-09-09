package errors

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
)

var (
	Todo = New("not implemented.")
)

// Equal
// compare err1 and err2 is same in memory index,error data or key code in engine/Err.
//
// Spec
// if Error Data is format with Err, compare with key code to the another error.
//
// Param
// err1 -- error one which want to compare.
// err2 -- error two which want to compare.
//
// Return
// return true is same, or return false.
func Equal(err1 error, err2 error) bool {
	if err1 == err2 {
		return true
	}
	if err1 == nil || err2 == nil {
		return false
	}

	// check Err type
	// if they are Err type,using errImpl compare.
	eImpl1, ok1 := err1.(*errImpl)
	eImpl2, ok2 := err2.(*errImpl)
	if ok1 && ok2 {
		return eImpl1.Key() == eImpl2.Key()
	}

	// if they are standar error,
	// compare the Message data.
	eMsg1 := err1.Error()
	eMsg2 := err2.Error()
	if eMsg1 == eMsg2 {
		return true
	}

	return ParseErr(err1).Key() == ParseErr(err2).Key()
}

type Err interface {
	Key() string
	Error() string
	Equal(err error) bool
	As(arg ...interface{}) Err
}

type errImpl struct {
	Code   string `json:"code"`
	Reason string `json:"reason"`
	Where  string `json:"where"`
}

// New
// create an Err implement error interface.
//
// Param
// code -- code or msg for the error struct,it will be a key.
//
// Return
// return a new Err interface
func New(code string) Err {
	return &errImpl{
		Code:  code,
		Where: caller(),
	}
}

// ParseErr
// Parse a standar error to Err interface.
// if the parameter is belong to Err, do a value copy an return a new Err.
// or parse string with error.Error(),
// if the string have a json struct with Err.Error(),return the origin struct with a new Err.
// or using error.Error() to create a new Err.
//
// Spec
// in the two case before, it will keep the key same as origin.
// the location is not change in parsing.
//
// Param
// src -- any error who implement error interface.
//
// Return
// return a new Err interface.
func ParseErr(src error) Err {
	if newErr, ok := src.(*errImpl); ok {
		return &errImpl{
			newErr.Code,
			newErr.Reason,
			newErr.Where,
		}
	}
	return parse(src.Error())
}

// As
// Parse the error, and fix with reason,it can make a replenishment for a same error.
//
// Spec
// because the value of error is change, so that location of Where is changed.
//
// Param
// err -- any error interface
// reason -- a array reason,it will be append to the reason of parameter.
//
// Return
// return a New Err,but with a same key with param error.
func As(err error, reason ...interface{}) Err {
	if err == nil {
		return nil
	}
	e := ParseErr(err).(*errImpl)
	return &errImpl{
		Code:   e.Code,
		Reason: fmt.Sprint(e.Reason, " ", reason),
		Where:  fmt.Sprint(e.Where, " ", caller()),
	}
}

func parse(src string) *errImpl {
	if len(src) == 0 || src[:1] != "{" {
		return New(src).(*errImpl)
	}

	eImpl := &errImpl{}
	if err := json.Unmarshal([]byte(src), eImpl); err != nil {
		return New(src).(*errImpl)
	}
	return eImpl
}

// call for domain
func caller() string {
	at := ""
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		at = "domain of caller is unknown"
	}
	me := runtime.FuncForPC(pc)
	if me == nil {
		at = "domain of call is unnamed"
	}

	fields := strings.Split(file, "/")

	filename := strings.Join(fields[len(fields)-2:], "/")

	at = me.Name() + "(file:" + filename + ",line:" + fmt.Sprintf("%d", line) + ")"
	return at
}

func (e *errImpl) Key() string {
	return e.Code
}

func (e *errImpl) Error() string {
	data, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		s := fmt.Sprintf("%v\n", *e)
		println(err.Error())
		return s
	}
	return string(data)
}

func (e *errImpl) Equal(l error) bool {
	if l == nil {
		return false
	}
	if e == l {
		return true
	}

	if t, ok := l.(*errImpl); ok {
		return e.Code == t.Code
	} else {
		return e.Code == parse(l.Error()).Code
	}
	return false
}

func (e *errImpl) As(arg ...interface{}) Err {
	newReason := fmt.Sprint(arg...)
	return &errImpl{
		Code:   e.Code,
		Reason: fmt.Sprint(e.Reason, " ", newReason),
		Where:  fmt.Sprint(e.Where, " ", caller()),
	}
}
