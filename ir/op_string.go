// Code generated by "stringer -type=Op -linecomment"; DO NOT EDIT.

package ir

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Push-0]
	_ = x[Pop-1]
	_ = x[Neg-2]
	_ = x[Add-3]
	_ = x[Sub-4]
	_ = x[Mul-5]
	_ = x[Div-6]
	_ = x[Cmp-7]
	_ = x[And-8]
	_ = x[Or-9]
	_ = x[Setl-10]
	_ = x[Setle-11]
	_ = x[Sete-12]
	_ = x[Setne-13]
	_ = x[Setg-14]
	_ = x[Setge-15]
}

const _Op_name = "pushpopnegaddsubmuldivcmpandorsetlsetlesetesetnesetgsetge"

var _Op_index = [...]uint8{0, 4, 7, 10, 13, 16, 19, 22, 25, 28, 30, 34, 39, 43, 48, 52, 57}

func (i Op) String() string {
	if i < 0 || i >= Op(len(_Op_index)-1) {
		return "Op(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Op_name[_Op_index[i]:_Op_index[i+1]]
}