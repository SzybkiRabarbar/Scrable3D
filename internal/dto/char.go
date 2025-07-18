package dto

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Char struct {
	HtmlIdentifier string `json:"id"`
	Value          string `json:"val"`
	// Position on 2D grid
	//  index 0: Y axis
	//  index 1: X axis
	Position [2]int `json:"pos"`
}

func (d *Char) ParseID() (int64, error) {
	if !strings.HasPrefix(d.HtmlIdentifier, "char-") {
		return 0, fmt.Errorf("char id must start with 'char-'")
	}

	re := regexp.MustCompile(`^char-([A-Z])(\d+)$`)
	matches := re.FindStringSubmatch(d.HtmlIdentifier)
	if len(matches) != 3 {
		return 0, fmt.Errorf("invalid char id format")
	}

	letter := matches[1]
	if letter != d.Value {
		return 0, fmt.Errorf("letter in char id (%s) doesn't match its value (%s)", letter, d.Value)
	}

	id, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, err
	}

	return int64(id), nil
}

func (d *Char) Validate() error {
	if d.HtmlIdentifier == "" {
		return errors.New("required field `id` is missing")
	}
	if d.Value == "" {
		return errors.New("required field `val` is missing")
	}
	return nil
}
