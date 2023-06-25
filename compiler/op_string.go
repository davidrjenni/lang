// Code generated by "stringer -type=Op -linecomment"; DO NOT EDIT.

package compiler

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Movq-0]
	_ = x[Movb-1]
	_ = x[Jump-2]
	_ = x[CJump-3]
	_ = x[Cmpq-4]
	_ = x[Cmpb-5]
	_ = x[Setne-6]
	_ = x[Call-7]
}

const _Op_name = "movqmovbjmpjecmpqcmpbsetnecall"

var _Op_index = [...]uint8{0, 4, 8, 11, 13, 17, 21, 26, 30}

func (i Op) String() string {
	if i < 0 || i >= Op(len(_Op_index)-1) {
		return "Op(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Op_name[_Op_index[i]:_Op_index[i+1]]
}