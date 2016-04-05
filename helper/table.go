package helper

// table helper refers to https://github.com/olekukonko/tablewriter

//
// Copyright (C) 2014 by Oleku Konko
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
//

import (
	"fmt"
	"github.com/mattn/go-runewidth"
	"io"
	"math"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	MAX_ROW_WIDTH = 30
)

const (
	CENTRE = "+"
	ROW    = "-"
	COLUMN = "|"
	SPACE  = " "
	EMPTY  = ""
)

const (
	ALIGN_DEFAULT = iota
	ALIGN_CENTRE
	ALIGN_RIGHT
	ALIGN_LEFT
)

var (
	decimal = regexp.MustCompile(`^\d*\.?\d*$`)
	percent = regexp.MustCompile(`^\d*\.?\d*$%$`)
)

type Table struct {
	out          io.Writer
	rows         [][]string
	lines        [][][]string
	cs           map[int]int
	rs           map[int]int
	headers      []string
	noHeaderLine bool
	autoFmt      bool
	autoWrap     bool
	mW           int
	pCenter      string
	pRow         string
	pColumn      string
	tColumn      int
	tRow         int
	align        int
	rowLine      bool
	border       bool
	colSize      int
}

// Start New Table
// Take io.Writer Directly
func NewTable(writer io.Writer) *Table {
	t := &Table{
		out:          writer,
		rows:         [][]string{},
		lines:        [][][]string{},
		cs:           make(map[int]int),
		rs:           make(map[int]int),
		headers:      []string{},
		noHeaderLine: false,
		autoFmt:      true,
		autoWrap:     true,
		// zero is unlimited
		mW:      MAX_ROW_WIDTH,
		pCenter: CENTRE,
		pRow:    ROW,
		pColumn: COLUMN,
		tColumn: -1,
		tRow:    -1,
		align:   ALIGN_DEFAULT,
		rowLine: false,
		border:  true,
		colSize: -1}
	return t
}

func NewPlainTable(writer io.Writer) *Table {
	t := NewTable(writer)
	t.SetBorder(false)
	t.SetRowLine(false)
	t.SetColWidth(0)
	t.SetRowSeparator("")
	t.SetNoHeaderLine(true)
	t.SetCenterSeparator("")
	t.SetColumnSeparator("    ")

	return t
}

// Render table output
func (t Table) Render() {
	if t.border {
		t.printLine(true)
	}
	t.printHeading()
	t.printRows()

	if !t.rowLine && t.border {
		t.printLine(true)
	}
	// t.printFooter()

}

// Set table header
func (t *Table) SetHeader(keys []string) {
	t.colSize = len(keys)
	for i, v := range keys {
		t.parseDimension(v, i, -1)
		t.headers = append(t.headers, v)
	}
}
func (t *Table) SetNoHeaderLine(b bool) {
	t.noHeaderLine = b
}

// Turn header autoformatting on/off. Default is on (true).
func (t *Table) SetAutoFormatHeaders(auto bool) {
	t.autoFmt = auto
}

// Turn automatic multiline text adjustment on/off. Default is on (true).
func (t *Table) SetAutoWrapText(auto bool) {
	t.autoWrap = auto
}

// Set the Default column width
func (t *Table) SetColWidth(width int) {
	t.mW = width
}

// Set the Column Separator
func (t *Table) SetColumnSeparator(sep string) {
	t.pColumn = sep
}

// Set the Row Separator
func (t *Table) SetRowSeparator(sep string) {
	t.pRow = sep
}

// Set the center Separator
func (t *Table) SetCenterSeparator(sep string) {
	t.pCenter = sep
}

// Set Table Alignment
func (t *Table) SetAlignment(align int) {
	t.align = align
}

// Set Row Line
// This would enable / disable a line on each row of the table
func (t *Table) SetRowLine(line bool) {
	t.rowLine = line
}

// Set Table Border
// This would enable / disable line around the table
func (t *Table) SetBorder(border bool) {
	t.border = border
}

// Append row to table
func (t *Table) Append(row []string) error {
	rowSize := len(t.headers)
	if rowSize > t.colSize {
		t.colSize = rowSize
	}

	n := len(t.lines)
	line := [][]string{}
	for i, v := range row {

		// Detect string  width
		// Detect String height
		// Break strings into words
		out := t.parseDimension(v, i, n)

		// Append broken words
		line = append(line, out)
	}
	t.lines = append(t.lines, line)
	return nil
}

// Allow Support for Bulk Append
// Eliminates repeated for loops
func (t *Table) AppendBulk(rows [][]string) (err error) {
	for _, row := range rows {
		err = t.Append(row)
		if err != nil {
			return err
		}
	}
	return nil
}

// Print line based on row width
func (t Table) printLine(nl bool) {
	fmt.Fprint(t.out, t.pCenter)
	for i := 0; i < len(t.cs); i++ {
		v := t.cs[i]
		fmt.Fprintf(t.out, "%s%s%s%s",
			t.pRow,
			strings.Repeat(string(t.pRow), v),
			t.pRow,
			t.pCenter)
	}
	if nl {
		fmt.Fprintln(t.out)
	}
}

// Print heading information
func (t Table) printHeading() {
	// Check if headers is available
	if len(t.headers) < 1 {
		return
	}

	// Check if border is set
	// Replace with space if not set
	fmt.Fprint(t.out, ConditionString(t.border, t.pColumn, EMPTY))

	// Identify last column
	end := len(t.cs) - 1

	// Print Heading column
	for i := 0; i <= end; i++ {
		v := t.cs[i]
		h := t.headers[i]
		if t.autoFmt {
			h = Title(h)
		}
		pad := ConditionString((i == end && !t.border), EMPTY, t.pColumn)

		if !t.border {
			fmt.Fprintf(t.out, "%s%s",
				// Pad(h, SPACE, v),
				PadRight(h, SPACE, v),
				pad)
		} else {
			fmt.Fprintf(t.out, " %s %s",
				// Pad(h, SPACE, v),
				PadRight(h, SPACE, v),
				pad)

		}
	}
	// Next line
	fmt.Fprintln(t.out)
	if !t.noHeaderLine {
		t.printLine(true)
	}
}

func (t Table) printRows() {
	for i, lines := range t.lines {
		t.printRow(lines, i)
	}

}

// Print Row Information
// Adjust column alignment based on type

func (t Table) printRow(columns [][]string, colKey int) {
	// Get Maximum Height
	max := t.rs[colKey]
	total := len(columns)

	// TODO Fix uneven col size
	// if total < t.colSize {
	//	for n := t.colSize - total; n < t.colSize ; n++ {
	//		columns = append(columns, []string{SPACE})
	//		t.cs[n] = t.mW
	//	}
	//}

	// Pad Each Height
	// pads := []int{}
	pads := []int{}

	for i, line := range columns {
		length := len(line)
		pad := max - length
		pads = append(pads, pad)
		for n := 0; n < pad; n++ {
			columns[i] = append(columns[i], "  ")
		}
	}
	//fmt.Println(max, "\n")
	for x := 0; x < max; x++ {
		for y := 0; y < total; y++ {

			// Check if border is set
			fmt.Fprint(t.out, ConditionString((!t.border && y == 0), EMPTY, t.pColumn))
			fmt.Fprint(t.out, ConditionString(!t.border, EMPTY, SPACE))

			str := columns[y][x]

			// This would print alignment
			// Default alignment  would use multiple configuration
			switch t.align {
			case ALIGN_CENTRE: //
				fmt.Fprintf(t.out, "%s", Pad(str, SPACE, t.cs[y]))
			case ALIGN_RIGHT:
				fmt.Fprintf(t.out, "%s", PadLeft(str, SPACE, t.cs[y]))
			case ALIGN_LEFT:
				fmt.Fprintf(t.out, "%s", PadRight(str, SPACE, t.cs[y]))
			default:
				if decimal.MatchString(strings.TrimSpace(str)) || percent.MatchString(strings.TrimSpace(str)) {
					fmt.Fprintf(t.out, "%s", PadLeft(str, SPACE, t.cs[y]))
				} else {
					fmt.Fprintf(t.out, "%s", PadRight(str, SPACE, t.cs[y]))

					// TODO Custom alignment per column
					//if max == 1 || pads[y] > 0 {
					//	fmt.Fprintf(t.out, "%s", Pad(str, SPACE, t.cs[y]))
					//} else {
					//	fmt.Fprintf(t.out, "%s", PadRight(str, SPACE, t.cs[y]))
					//}

				}
			}
			fmt.Fprintf(t.out, ConditionString(!t.border, EMPTY, SPACE))
		}
		// Check if border is set
		// Replace with space if not set
		fmt.Fprint(t.out, ConditionString(t.border, t.pColumn, EMPTY))
		fmt.Fprintln(t.out)
	}

	if t.rowLine {
		t.printLine(true)
	}

}

func (t *Table) parseDimension(str string, colKey, rowKey int) []string {
	var (
		raw []string
		max int
	)
	w := DisplayWidth(str)
	// Calculate Width
	// Check if with is grater than maximum width
	if w > t.mW && t.mW != 0 {
		w = t.mW
	}

	// Check if width exists
	v, ok := t.cs[colKey]
	if !ok || v < w || v == 0 {
		t.cs[colKey] = w
	}

	if rowKey == -1 {
		return raw
	}
	// Calculate Height
	if t.autoWrap {
		raw, _ = WrapString(str, t.cs[colKey])
	} else {
		raw = getLines(str)
	}

	for _, line := range raw {
		if w := DisplayWidth(line); w > max {
			max = w
		}
	}

	// Make sure the with is the same length as maximum word
	// Important for cases where the width is smaller than maxu word
	if max > t.cs[colKey] {
		t.cs[colKey] = max
	}

	h := len(raw)
	v, ok = t.rs[rowKey]

	if !ok || v < h || v == 0 {
		t.rs[rowKey] = h
	}
	//fmt.Printf("Raw %+v %d\n", raw, len(raw))
	return raw
}

var ansi = regexp.MustCompile("\033\\[(?:[0-9]{1,3}(?:;[0-9]{1,3})*)?[m|K]")

func DisplayWidth(str string) int {
	return runewidth.StringWidth(ansi.ReplaceAllLiteralString(str, ""))
}

// Simple Condition for string
// Returns value based on condition
func ConditionString(cond bool, valid, inValid string) string {
	if cond {
		return valid
	}
	return inValid
}

// Format Table Header
// Replace _ , . and spaces
func Title(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Replace(name, ".", " ", -1)
	name = strings.TrimSpace(name)
	return strings.ToUpper(name)
}

// Pad String
// Attempts to play string in the center
func Pad(s, pad string, width int) string {
	gap := width - DisplayWidth(s)
	if gap > 0 {
		gapLeft := int(math.Ceil(float64(gap / 2)))
		gapRight := gap - gapLeft
		return strings.Repeat(string(pad), gapLeft) + s + strings.Repeat(string(pad), gapRight)
	}
	return s
}

// Pad String Right position
// This would pace string at the left side fo the screen
func PadRight(s, pad string, width int) string {
	gap := width - DisplayWidth(s)
	if gap > 0 {
		return s + strings.Repeat(string(pad), gap)
	}
	return s
}

// Pad String Left position
// This would pace string at the right side fo the screen
func PadLeft(s, pad string, width int) string {
	gap := width - DisplayWidth(s)
	if gap > 0 {
		return strings.Repeat(string(pad), gap) + s
	}
	return s
}

var (
	nl = "\n"
	sp = " "
)

const defaultPenalty = 1e5

// Wrap wraps s into a paragraph of lines of length lim, with minimal
// raggedness.
func WrapString(s string, lim int) ([]string, int) {
	words := strings.Split(strings.Replace(strings.TrimSpace(s), nl, sp, -1), sp)
	var lines []string
	max := 0
	for _, v := range words {
		max = len(v)
		if max > lim {
			lim = max
		}
	}
	for _, line := range WrapWords(words, 1, lim, defaultPenalty) {
		lines = append(lines, strings.Join(line, sp))
	}
	return lines, lim
}

// WrapWords is the low-level line-breaking algorithm, useful if you need more
// control over the details of the text wrapping process. For most uses,
// WrapString will be sufficient and more convenient.
//
// WrapWords splits a list of words into lines with minimal "raggedness",
// treating each rune as one unit, accounting for spc units between adjacent
// words on each line, and attempting to limit lines to lim units. Raggedness
// is the total error over all lines, where error is the square of the
// difference of the length of the line and lim. Too-long lines (which only
// happen when a single word is longer than lim units) have pen penalty units
// added to the error.
func WrapWords(words []string, spc, lim, pen int) [][]string {
	n := len(words)

	length := make([][]int, n)
	for i := 0; i < n; i++ {
		length[i] = make([]int, n)
		length[i][i] = utf8.RuneCountInString(words[i])
		for j := i + 1; j < n; j++ {
			length[i][j] = length[i][j-1] + spc + utf8.RuneCountInString(words[j])
		}
	}
	nbrk := make([]int, n)
	cost := make([]int, n)
	for i := range cost {
		cost[i] = math.MaxInt32
	}
	for i := n - 1; i >= 0; i-- {
		if length[i][n-1] <= lim {
			cost[i] = 0
			nbrk[i] = n
		} else {
			for j := i + 1; j < n; j++ {
				d := lim - length[i][j-1]
				c := d*d + cost[j]
				if length[i][j-1] > lim {
					c += pen // too-long lines get a worse penalty
				}
				if c < cost[i] {
					cost[i] = c
					nbrk[i] = j
				}
			}
		}
	}
	var lines [][]string
	i := 0
	for i < n {
		lines = append(lines, words[i:nbrk[i]])
		i = nbrk[i]
	}
	return lines
}

// getLines decomposes a multiline string into a slice of strings.
func getLines(s string) []string {
	var lines []string

	for _, line := range strings.Split(strings.TrimSpace(s), nl) {
		lines = append(lines, line)
	}
	return lines
}
