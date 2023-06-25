// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors_test

import (
	"strings"
	"testing"

	"davidrjenni.io/lang/internal/errors"
)

func TestErrors(t *testing.T) {
	var errs errors.Errors

	if errs.Err() != nil {
		t.Errorf("got: %v, expected nil", errs.Err())
	}

	errs.Append("error %d", 23)
	errs.Append("error %d", 42)

	if errs.Err().Error() != errs.Error() {
		t.Errorf("got:\n%v\nexpected:\n%v", errs.Err(), errs)
	}

	const err = "error 23\nerror 42"
	if errs.Error() != err {
		t.Errorf("got:\n%v\nexpected:\n%v", errs.Error(), err)
	}

	for i := 0; i < 19; i++ {
		errs.Append("error %d", i)
	}

	lastLine := "and 1 more error"
	errMsg := errs.Error()
	actualLastLine := errMsg[strings.LastIndex(errMsg, "\n")+1:]
	if actualLastLine != lastLine {
		t.Errorf("got:\n%v\nexpected:\n%v", actualLastLine, lastLine)
	}

	errs.Append("error")

	lastLine = "and 2 more errors"
	errMsg = errs.Error()
	actualLastLine = errMsg[strings.LastIndex(errMsg, "\n")+1:]
	if actualLastLine != lastLine {
		t.Errorf("got:\n%v\nexpected:\n%v", actualLastLine, lastLine)
	}
}
