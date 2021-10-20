package main

import "testing"

func TestFizzbuzz(t *testing.T) {
	type input struct {
		int1  uint
		int2  uint
		limit uint
		str1  string
		str2  string
	}

	type test struct {
		name  string
		input input
		want  []string
	}

	tests := []test{
		{name: "classic fizzbuzz", input: input{int1: 2, int2: 3, limit: 7, str1: "fizz", str2: "buzz"}, want: []string{"1", "fizz", "buzz", "fizz", "5", "fizzbuzz", "7"}},
		{name: "classic foobar with different terms", input: input{int1: 2, int2: 3, limit: 7, str1: "foo", str2: "bar"}, want: []string{"1", "foo", "bar", "foo", "5", "foobar", "7"}},
		{name: "same term str1 and str2", input: input{int1: 2, int2: 3, limit: 7, str1: "fizz", str2: "fizz"}, want: []string{"1", "fizz", "fizz", "fizz", "5", "fizzfizz", "7"}},
		{name: "same multiples int1 and int2", input: input{int1: 2, int2: 2, limit: 5, str1: "fizz", str2: "buzz"}, want: []string{"1", "fizzbuzz", "3", "fizzbuzz", "5"}},
		{name: "empty strings", input: input{int1: 2, int2: 3, limit: 7, str1: "fizz", str2: ""}, want: []string{"1", "fizz", "", "fizz", "5", "fizz", "7"}},
	}

	for _, tc := range tests {
		got := Fizzbuzz(tc.input.int1, tc.input.int2, tc.input.limit, tc.input.str1, tc.input.str2)

		if len(got) != len(tc.want) {
			t.Fatalf("%s :: Expected a result of size %d, had %d", tc.name, len(tc.want), len(got))
		}

		for k, v := range got {
			if v != tc.want[k] {
				t.Fatalf("%s :: Expected \"%s\", got \"%s\"", tc.name, tc.want, got)
			}
		}
	}
}

func BenchmarkFizzbuzz(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Fizzbuzz(3, 5, 100, "fizz", "buzz")
	}
}
