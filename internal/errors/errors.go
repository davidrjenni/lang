// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors // import "davidrjenni.io/lang/internal/errors"

import (
	"bytes"
	"fmt"

	"davidrjenni.io/lang/lexer"
)

// Errors is a slice of errors.
type Errors []error

// Append appends an error, starting with the given position and
// formatted according to the format string and the arguments given.
func (e *Errors) Append(pos lexer.Pos, format string, args ...interface{}) {
	*e = append(*e, fmt.Errorf(pos.String()+": "+format, args...))
}

// Err returns itself or nil, if there are none.
func (e Errors) Err() error {
	if len(e) == 0 {
		return nil
	}
	return e
}

// Error returns a string containing the first 20 errors.
func (e Errors) Error() string {
	const n = 20

	var b bytes.Buffer
	for i, err := range e {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(err.Error())
		if i == n-1 {
			if len(e) == n+1 {
				fmt.Fprintf(&b, "\nand 1 more error")
			} else {
				fmt.Fprintf(&b, "\nand %d more errors", len(e)-n)
			}
			break
		}
	}
	return b.String()
}
