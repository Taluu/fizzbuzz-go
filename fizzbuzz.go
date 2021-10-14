package main

import (
	"strconv"
	"strings"
)

// Fizzbuzz prints all the ints between 1 and limit, replacing multiples of int1 with str1,
// multiples of int2 with str2, multiples of int1 and int2 with str1str2.
func Fizzbuzz(int1 uint, int2 uint, limit uint, str1 string, str2 string) string {
	result := make([]string, 0, limit)

	for i := uint(1); i <= limit; i++ {
		var j string

		if i%int1 == 0 {
			j += str1
		}

		if i%int2 == 0 {
			j += str2
		}

		if i%int1 != 0 && i%int2 != 0 {
			j = strconv.FormatUint(uint64(i), 10)
		}

		result = append(result, j)
	}

	return strings.Join(result, ",")
}
