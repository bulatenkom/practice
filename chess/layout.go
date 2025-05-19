package main

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Box may contain single(multi)line content with all lines expanded to the width of the longest line
type Box struct {
	// width, height
	w, h    int
	content []string
}

func (b *Box) String() string {
	return strings.Join(b.content, "")
}

func NewBox(s string) Box {
	lines := strings.Split(s, "\n")
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	width := 0
	for _, line := range lines {
		width = max(width, CountRunesIgnoringAnsiEscapes(line))
	}
	content := []string{}
	for _, line := range lines {
		fline := line + strings.Repeat(" ", width-CountRunesIgnoringAnsiEscapes(line)) + "\n"
		content = append(content, fline)
	}

	return Box{
		w:       width,
		h:       len(lines),
		content: content,
	}
}

var escapeSeqMatcher *regexp.Regexp = regexp.MustCompile("(\033\\[0m|\033\\[31m|\033\\[33m|\033\\[1m|\033\\[38;2;0;0;0m|\033\\[48;2;209;139;71m|\033\\[48;2;255;206;158m)")

func CountRunesIgnoringAnsiEscapes(s string) int {
	res := escapeSeqMatcher.ReplaceAllLiteralString(s, "")
	return utf8.RuneCountInString(res)
}

type Orientation byte

const (
	Horizontal Orientation = iota
	Vertical
)

func Layout(boxes []Box, orientation Orientation) Box {
	if orientation == Vertical {
		content := ""
		for _, b := range boxes {
			content += b.String()
		}
		return NewBox(content)
	}
	if orientation == Horizontal {
		w := 0
		h := 0

		for _, b := range boxes {
			w += b.w
			h = max(h, b.h)
		}

		content := ""
		for r := 0; r < h; r++ {
			line := ""
			for _, b := range boxes {
				if r+1 > b.h {
					line += strings.Repeat(" ", b.w)
					line += "|"
					continue
				}
				s := b.content[r]
				line += s[:len(s)-1]
				line += "|"
			}
			line += "\n"
			content += line
		}

		return NewBox(content)
	}
	return Box{}
}

func RenderBoard(board board) string {
	vMarkers := []rune{'8', '7', '6', '5', '4', '3', '2', '1'}
	hMarkers := []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}

	output := joinHMarkers(hMarkers)

	for i := range board {
		output += string(vMarkers[i]) + " "
		for j := range board[0] {
			cell := " " + string(board[i][j]) + " "
			if (i+j)%2 == 0 {
				output += bg_light_peach + fg_black + cell + c_reset
			} else {
				output += bg_med_brown + fg_black + cell + c_reset
			}
		}
		output += " " + string(vMarkers[i]) + "\n"
	}
	output += joinHMarkers(hMarkers)
	return output
}

func joinHMarkers(markers []rune) string {
	output := "  "
	for _, v := range markers {
		output += " " + string(v) + " "
	}
	output += "\n"
	return output
}

func ClearScreen() (int, error) {
	return fmt.Print("\x1b[H\x1b[2J")
}
