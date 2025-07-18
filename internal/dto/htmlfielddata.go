package dto

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

type HtmlFieldData struct {
	Repr  string
	Value string
	X     int
	Y     int
	Z     int
}

func NewHtmlFieldData[T constraints.Integer](value string, x T, y T, z T) *HtmlFieldData {
	d := &HtmlFieldData{
		Value: value,
		X:     int(x) * 60,
		Y:     int(y) * 60,
		Z:     420 - (int(z) * 60),
	}

	d.Repr = d.Value +
		fmt.Sprintf("%X", x) +
		fmt.Sprintf("%X", y) +
		fmt.Sprintf("%X", z)
	return d
}
