// Code generated by "stringer -type=ServerType"; DO NOT EDIT.

package pixieapi

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ST_COMPUTE_G1-0]
	_ = x[ST_COMPUTE_G2-1]
	_ = x[DEFAULT-2]
	_ = x[IGNORE-3]
}

const _ServerType_name = "ST_COMPUTE_G1ST_COMPUTE_G2DEFAULTIGNORE"

var _ServerType_index = [...]uint8{0, 13, 26, 33, 39}

func (i ServerType) String() string {
	if i < 0 || i >= ServerType(len(_ServerType_index)-1) {
		return "ServerType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ServerType_name[_ServerType_index[i]:_ServerType_index[i+1]]
}
