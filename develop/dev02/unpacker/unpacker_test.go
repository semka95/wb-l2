package unpacker

import "testing"

func TestUnpack(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		want  string
	}{
		{
			desc:  "normal",
			input: "a4df5cvs1",
			want:  "aaaadfffffcvs",
		},
		{
			desc:  "no numbers, only letters",
			input: "abcd",
			want:  "abcd",
		},
		{
			desc:  "only numbers",
			input: "45",
			want:  "",
		},
		{
			desc:  "empty string",
			input: "",
			want:  "",
		},
		{
			desc:  "multiple digit number",
			input: "a10df5cvs1",
			want:  "aaaaaaaaaadfffffcvs",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := unpack(tC.input)
			if got != tC.want {
				t.Errorf("got: %s, want: %s", got, tC.want)
			}
		})
	}
}
