package dto

import "errors"

type ActionData struct {
	Type string `json:"actionType"`
}

func (d *ActionData) Validate() error {
	if d.Type == "" {
		return errors.New("required field `actionType` is missing")
	}
	return nil
}
