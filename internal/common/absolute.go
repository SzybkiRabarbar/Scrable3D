package common

import (
	"golang.org/x/exp/constraints"
)

func Abs[T constraints.Integer | constraints.Float](x T) T {
	if x < 0 {
		return -x
	} else {
		return x
	}
}
