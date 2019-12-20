// Code generated by "stringer -type=RuleTypes"; DO NOT EDIT.

package esg

import (
	"errors"
	"strconv"
)

var _ = errors.New("dummy error")

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[UniformItems-0]
	_ = x[ProbItems-1]
	_ = x[CondItems-2]
	_ = x[SequentialItems-3]
	_ = x[PermutedItems-4]
	_ = x[RuleTypesN-5]
}

const _RuleTypes_name = "UniformItemsProbItemsCondItemsSequentialItemsPermutedItemsRuleTypesN"

var _RuleTypes_index = [...]uint8{0, 12, 21, 30, 45, 58, 68}

func (i RuleTypes) String() string {
	if i < 0 || i >= RuleTypes(len(_RuleTypes_index)-1) {
		return "RuleTypes(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _RuleTypes_name[_RuleTypes_index[i]:_RuleTypes_index[i+1]]
}

func (i *RuleTypes) FromString(s string) error {
	for j := 0; j < len(_RuleTypes_index)-1; j++ {
		if s == _RuleTypes_name[_RuleTypes_index[j]:_RuleTypes_index[j+1]] {
			*i = RuleTypes(j)
			return nil
		}
	}
	return errors.New("String: " + s + " is not a valid option for type: RuleTypes")
}