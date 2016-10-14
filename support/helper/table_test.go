package helper

import (
	"os"
)

func ExampleTable() {
	tb := NewTable(os.Stdout)
	tb.SetHeader([]string{"AAA", "BBB", "CCC"})
	tb.Append([]string{"aaa", "bbb", "ccc"})
	tb.Append([]string{"aaa", "bbbbbb", "ccc"})
	tb.Append([]string{"aaa", "bbbbbbbbb", "ccc"})
	tb.Render()
	// Output:
	//+-----+-----------+-----+
	//| AAA | BBB       | CCC |
	//+-----+-----------+-----+
	//| aaa | bbb       | ccc |
	//| aaa | bbbbbb    | ccc |
	//| aaa | bbbbbbbbb | ccc |
	//+-----+-----------+-----+
}

func ExamplePlainTable() {
	tb := NewPlainTable(os.Stdout)
	tb.SetHeader([]string{"AAA", "BBB", "CCC"})
	tb.Append([]string{"aaa", "bbb", "ccc"})
	tb.Append([]string{"aaa", "bbbbbb", "ccc"})
	tb.Append([]string{"aaa", "bbbbbbbbb", "ccc"})
	tb.Render()
	// Output:
	//AAA        BBB              CCC
	//aaa        bbb              ccc
	//aaa        bbbbbb           ccc
	//aaa        bbbbbbbbb        ccc
}
