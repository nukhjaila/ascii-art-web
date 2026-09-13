package banner

import (
	"os"
	"strings"
)

func Load(filePath string) (map[rune][]string, error) {

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	bannerMap := make(map[rune][]string)

	for i := 32; i <= 126; i++ {
		startLine := (i-32)*9 + 1
		endLine := startLine + 8

		if endLine > len(lines) {
			break
		}

		characterGlyph := lines[startLine:endLine]
		ch := rune(i)
		bannerMap[ch] = characterGlyph
	}

	return bannerMap, nil
}
