package parser

import (
	"errors"
	"strings"
)

func Parse(input string) ([][]rune, error) {
	parts := strings.Split(input, "\\n")
	var res [][]rune

	for _, part := range parts {
		var row []rune
		for _, ch := range part {
			if ch < 32 || ch > 126 {
				return nil, errors.New("not printable character")
			}
			row = append(row, ch)
		}
		res = append(res, row)
	}

	return res, nil
}
