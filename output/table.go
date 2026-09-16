package output

import (
	"strings"
	"fmt"
)

// Note Table represents a table with headers and rows.
type Table struct {
	offset int
	// NOTE header is treated like a row
	rows [][]string
}

func NewTable(headers []string, offset int) *Table {
	// Init a new table with given headers line
	return &Table{
		offset:  offset,
		rows:    [][]string{headers},
	}
}

func (t *Table) SetOffset(offset int) {
	t.offset = offset
}

func (t *Table) AddRow(row []string) error {
	// compare the row length with header line
	if len(row) != len(t.rows[0]) {
		return fmt.Errorf("row length %d does not match header length %d", len(row), len(t.rows[0]))
	}
	t.rows = append(t.rows, row)
	return nil
}

func (t *Table) calculateColumnWidth() []int {
    // Calculate maximum width for each column based on headers and rows

	// preallocate slice with the length of header row
	columnWidths := make([]int, len(t.rows[0]))

	// iterate over all rows
	for _, row := range t.rows {
		// iterate cell by cell
		for colIndex, cell := range row {		
			// calculate the cell size including offset
			cellWidth := len(cell) + t.offset
			if cellWidth > columnWidths[colIndex] {
				// set the maximum value inside column widths
				columnWidths[colIndex] = cellWidth
			}
		}
	}

	return columnWidths
}

// NOTE render the console table
func (t *Table) Render() {
	// calculate per column width, to determine the amount of whitespaces required
	column_widths := t.calculateColumnWidth()

	// prepare table
	var preparedTab [][]string

	for _, row := range t.rows {
		var currentRow []string
		for colIndex, cell := range row {
			// calculate the amount of whitspaces for cell	
			wsAmount := column_widths[colIndex] - len(cell)
			// prepare the proper amount of whitespace
			whitespaces := strings.Repeat(" ", wsAmount)
			// append the whitespaces to cell, and add to prepared row
			currentRow = append(currentRow, cell + whitespaces)
		}
		preparedTab = append(preparedTab, currentRow)
	}
		
	// print whole string ina single stdout print
	for _, line := range preparedTab {
		fmt.Println(strings.Join(line, ""))
	}
}
