package sortstrings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSort(t *testing.T) {
	testCases := []struct {
		desc string
		app  appEnv
		data []string
		want []string
	}{
		{
			desc: "normal",
			app:  appEnv{},
			data: []string{
				"Standing on one's head at job interviews forms a lasting impression.",
				"The chic gangster liked to start the day with a pink scarf.",
				"He kept telling himself that one day it would all somehow make sense.",
				"Cats are good pets, for they are clean and are not noisy.",
				"I am my aunt's sister's daughter.",
			},
			want: []string{
				"Cats are good pets, for they are clean and are not noisy.",
				"He kept telling himself that one day it would all somehow make sense.",
				"I am my aunt's sister's daughter.",
				"Standing on one's head at job interviews forms a lasting impression.",
				"The chic gangster liked to start the day with a pink scarf.",
			},
		},
		{
			desc: "not numeric order",
			app:  appEnv{},
			data: []string{
				"1",
				"5",
				"13",
				"23",
				"11",
				"21",
				"31",
			},
			want: []string{
				"1",
				"11",
				"13",
				"21",
				"23",
				"31",
				"5",
			},
		},
		{
			desc: "reverse order",
			app: appEnv{
				isReverse: true,
			},
			data: []string{
				"1",
				"5",
				"13",
				"23",
				"11",
				"21",
				"31",
			},
			want: []string{
				"5",
				"31",
				"23",
				"21",
				"13",
				"11",
				"1",
			},
		},
		{
			desc: "delete duplicate",
			app: appEnv{
				deleteDuplicate: true,
			},
			data: []string{
				"1",
				"1",
				"5",
				"13",
				"23",
				"11",
				"11",
				"21",
				"31",
				"31",
				"31",
			},
			want: []string{
				"1",
				"11",
				"13",
				"21",
				"23",
				"31",
				"5",
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := tC.app.sort(tC.data)
			assert.Equal(t, tC.want, got)
		})
	}
}

func TestSortColumns(t *testing.T) {
	testCases := []struct {
		desc string
		app  appEnv
		data []string
		want []string
	}{
		{
			desc: "by 2nd column",
			app: appEnv{
				column: 2,
			},
			data: []string{
				"Standing on one's head at job interviews forms a lasting impression.",
				"The chic gangster liked to start the day with a pink scarf.",
				"He kept telling himself that one day it would all somehow make sense.",
				"Cats are good pets, for they are clean and are not noisy.",
				"I am my aunt's sister's daughter.",
			},
			want: []string{
				"I am my aunt's sister's daughter.",
				"Cats are good pets, for they are clean and are not noisy.",
				"The chic gangster liked to start the day with a pink scarf.",
				"He kept telling himself that one day it would all somehow make sense.",
				"Standing on one's head at job interviews forms a lasting impression.",
			},
		},
		{
			desc: "by column out of range",
			app: appEnv{
				column: 200,
			},
			data: []string{
				"Standing on one's head at job interviews forms a lasting impression.",
				"The chic gangster liked to start the day with a pink scarf.",
				"He kept telling himself that one day it would all somehow make sense.",
				"Cats are good pets, for they are clean and are not noisy.",
				"I am my aunt's sister's daughter.",
			},
			want: []string{
				"Cats are good pets, for they are clean and are not noisy.",
				"He kept telling himself that one day it would all somehow make sense.",
				"I am my aunt's sister's daughter.",
				"Standing on one's head at job interviews forms a lasting impression.",
				"The chic gangster liked to start the day with a pink scarf.",
			},
		},
		{
			desc: "numbers numeric order",
			app: appEnv{
				column:    1,
				isNumeric: true,
			},
			data: []string{
				"5",
				"23",
				"1",
				"21",
				"31",
				"13",
				"11",
			},
			want: []string{
				"1",
				"5",
				"11",
				"13",
				"21",
				"23",
				"31",
			},
		},
		{
			desc: "by 2nd column, in reverse",
			app: appEnv{
				column:    2,
				isReverse: true,
			},
			data: []string{
				"Standing on one's head at job interviews forms a lasting impression.",
				"The chic gangster liked to start the day with a pink scarf.",
				"He kept telling himself that one day it would all somehow make sense.",
				"Cats are good pets, for they are clean and are not noisy.",
				"I am my aunt's sister's daughter.",
			},
			want: []string{
				"Standing on one's head at job interviews forms a lasting impression.",
				"He kept telling himself that one day it would all somehow make sense.",
				"The chic gangster liked to start the day with a pink scarf.",
				"Cats are good pets, for they are clean and are not noisy.",
				"I am my aunt's sister's daughter.",
			},
		},
		{
			desc: "delete duplicate",
			app: appEnv{
				column:          1,
				deleteDuplicate: true,
			},
			data: []string{
				"1",
				"1",
				"5",
				"13",
				"23",
				"11",
				"11",
				"21",
				"31",
				"31",
				"31",
			},
			want: []string{
				"1",
				"11",
				"13",
				"21",
				"23",
				"31",
				"5",
			},
		},
		{
			desc: "numeric sort, but column starts with letter",
			app: appEnv{
				column:    1,
				isNumeric: true,
			},
			data: []string{
				"d1",
				"ad5",
				"asbv13",
				"sfg23",
				"fa11",
				"gh21",
				"31",
			},
			want: []string{
				"31",
				"ad5",
				"asbv13",
				"d1",
				"fa11",
				"gh21",
				"sfg23",
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := tC.app.sortColumns(tC.data)
			assert.Equal(t, tC.want, got)
		})
	}
}
