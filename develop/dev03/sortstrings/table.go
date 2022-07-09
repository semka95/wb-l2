package sortstrings

import (
	"strconv"
)

// stringTable represents slice of slices of strings
type stringTable struct {
	data      [][]string
	column    int
	isNumeric bool
}

func (t stringTable) Len() int {
	return len(t.data)
}

func (t stringTable) Less(i, j int) bool {
	col := t.column
	if col > len(t.data[i])-1 || col > len(t.data[j]) {
		col = 0
	}

	if t.isNumeric {
		n1 := trimNonNumber(t.data[i][col])
		n2 := trimNonNumber(t.data[j][col])

		i1, err := strconv.Atoi(n1)
		if err != nil {
			return (t.data[i][col] < t.data[j][col])
		}
		j1, err := strconv.Atoi(n2)
		if err != nil {
			return (t.data[i][col] < t.data[j][col])
		}

		return i1 < j1
	}
	return (t.data[i][col] < t.data[j][col])
}

func (t stringTable) Swap(i, j int) {
	t.data[i], t.data[j] = t.data[j], t.data[i]
}
