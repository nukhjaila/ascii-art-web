package renderer

import (
	"fmt"
	"io"
)

func Render(w io.Writer, matrix [][]rune, bannerMap map[rune][]string) {
	for _, row := range matrix {
		if len(row) == 0 {
			fmt.Fprintln(w)
			continue
		}

		for lineIndex := 0; lineIndex < 8; lineIndex++ {
			for _, ch := range row {
				glyphLines := bannerMap[ch]
				fmt.Fprint(w, glyphLines[lineIndex])
			}
			fmt.Fprintln(w)
		}
	}
}
