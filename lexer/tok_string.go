// Code generated by "stringer -type=Tok -linecomment"; DO NOT EDIT.

package lexer

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[EOF-0]
	_ = x[Illegal-1]
	_ = x[True-2]
	_ = x[False-3]
}

const _Tok_name = "EOFillegaltruefalse"

var _Tok_index = [...]uint8{0, 3, 10, 14, 19}

func (i Tok) String() string {
	if i < 0 || i >= Tok(len(_Tok_index)-1) {
		return "Tok(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Tok_name[_Tok_index[i]:_Tok_index[i+1]]
}
