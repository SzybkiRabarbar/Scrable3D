package dto

import "errors"

type PlayData struct {
	SideInt int    `json:"side"`
	Chars   []Char `json:"chars"`
}

func (d *PlayData) Validate() error {
	if d.Chars == nil {
		return errors.New("required field `chars` is missing")
	}
	for _, char := range d.Chars {
		if err := char.Validate(); err != nil {
			return err
		}
	}
	return nil
}
