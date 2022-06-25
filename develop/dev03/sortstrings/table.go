package sortstrings

import (
	"strconv"
)

type stringTable struct {
	data      [][]string
	column    int
	isNumeric bool
}

func (t stringTable) Len() int {
	return len(t.data)
}

func (t stringTable) Less(i, j int) bool {
	if t.isNumeric {
		n1 := trimNonNumber(t.data[i][t.column])
		n2 := trimNonNumber(t.data[j][t.column])

		i1, err := strconv.Atoi(n1)
		if err != nil {
			return (t.data[i][t.column] < t.data[j][t.column])
		}
		j1, err := strconv.Atoi(n2)
		if err != nil {
			return (t.data[i][t.column] < t.data[j][t.column])
		}

		return i1 < j1
	}
	return (t.data[i][t.column] < t.data[j][t.column])
}

func (t stringTable) Swap(i, j int) {
	t.data[i], t.data[j] = t.data[j], t.data[i]
}
