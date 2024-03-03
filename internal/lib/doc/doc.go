package doc

import (
	"github.com/mattn/go-runewidth"
	"strings"
)

// getMaxWidths calculates the maximum width for each column considering Unicode character widths
func getMaxWidths(table [][]string) []int {
	maxWidths := make([]int, len(table[0]))
	for _, row := range table {
		for i, cell := range row {
			cellDisplayWidth := runewidth.StringWidth(cell)
			if cellDisplayWidth > maxWidths[i] {
				maxWidths[i] = cellDisplayWidth
			}
		}
	}

	return maxWidths
}

// padRight returns a new string of a specified length in which the end of the original string is padded with spaces, considering Unicode width
func padRight(str string, length int) string {
	paddingNeeded := length - runewidth.StringWidth(str)
	return str + strings.Repeat(" ", paddingNeeded)
}

// GetAsMDTable generates a markdown table from a 2D string array.
func GetAsMDTable(table [][]string) string {
	if len(table) == 0 || len(table[0]) == 0 {
		return ""
	}

	maxWidths := getMaxWidths(table)
	var mdTable strings.Builder

	// Construct the header with padding
	for i, header := range table[0] {
		mdTable.WriteString("| " + padRight(header, maxWidths[i]+1))
	}
	mdTable.WriteString("|\n")

	// Construct the separator
	for _, width := range maxWidths {
		mdTable.WriteString("|" + strings.Repeat("-", width+2))
	}
	mdTable.WriteString("|\n")

	for i, row := range table[1:] {
		for i, cell := range row {
			mdTable.WriteString("| " + padRight(cell, maxWidths[i]+1))
		}

		mdTable.WriteString("|")

		if i+1 != len(table[1:]) {
			mdTable.WriteString("\n")
		}
	}

	return mdTable.String()
}
