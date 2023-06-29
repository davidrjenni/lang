// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lexer // import "davidrjenni.io/lang/lexer"

import "fmt"

// Pos represents a position in the source code.
type Pos struct {
	Filename string // filename, optional
	Line     uint32 // line number, starting at 1
	Column   uint32 // column number within a line, starting at 1
}

// String returns a string representation of the position.
// An invalid position is indicated with "-".
func (p Pos) String() string {
	if p.Line > 0 && p.Column > 0 {
		s := p.Filename
		if s != "" {
			s += ":"
		}
		return fmt.Sprintf("%s%d:%d", s, p.Line, p.Column)
	}
	return "-"
}
