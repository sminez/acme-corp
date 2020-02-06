package main

// Aiming to emulate the behaviour of typing 'gq' with a visual
// selection in vim: wrap lines to a specified column count and
// retain any leading prefix on the newly wrapped lines.
//   This is primarily for markup / comments but it might prove
//   useful for other things as well. Who knows!

import (
	"bufio"
	"flag"
	"os"
	"strings"
	"unicode"
)

var (
	columns    = flag.Int("c", 80, "number of columns to wrap to")
	alphaNumer = flag.Bool("a", false, "allow alphanumeric characters in the prefix")
)

// Determine the largest prefix that we can find for all of the
// lines that we have. By default, we stop at the first
// alphanumeric character encountered but this can be overriden
// using the '-a' flag.
func maximalPrefix(lines []string, allowAlphaNumeric bool) (string, int) {
	first := lines[0]
	k := len(first)
	prefix := ""

	if !allowAlphaNumeric {
		k = 0
		for _, c := range first {
			if unicode.IsLetter(c) || unicode.IsNumber(c) {
				break
			}
			k++
		}
	}

	for n := 1; n < k; n++ {
		s := first[:n]
		for _, l := range lines[1:] {
			if !strings.HasPrefix(l, s) {
				return prefix
			}
			prefix = s
		}
	}

	return prefix, int
}

func main() {
	var lines []string
	flag.Parse()

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		lines = append(l, s.Text())
	}

	prefix := maximalPrefix(lines, *alphaNUmeric)
}
