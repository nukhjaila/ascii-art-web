package renderer

import (
	"fmt"
	"strings"
)

func Render(matrix [][]rune, bannerMap map[rune][]string) []byte {
	var sb strings.Builder

	for _, row := range matrix {
		if len(row) == 0 {
			fmt.Fprintln(&sb)
			continue
		}

		for lineIndex := 0; lineIndex < 8; lineIndex++ {
			for _, ch := range row {
				glyphLines := bannerMap[ch]
				fmt.Fprint(&sb, glyphLines[lineIndex])
			}
			fmt.Fprintln(&sb)
		}
	}
	return []byte(sb.String())
}
