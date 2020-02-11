package main

// Aiming to emulate the behaviour of typing 'gq' with a visual selection in
// vim, wrap lines to a specified column count and retain any leading prefix on
// the newly wrapped lines. This is primarily for markup / comments but it
// might prove useful for other things as well. Who knows!

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"
)

var (
	columns      = flag.Int("c", 80, "number of columns to wrap to")
	alphaNumeric = flag.Bool("a", false, "allow alphanumeric characters in the prefix")
)

// Determine the largest prefix that we can find for all of the lines that we
// have. By default, we stop at the first alphanumeric character encountered but
// this can be overriden using the '-a' flag.
func maximalPrefix(lines []string, allowAlphaNumeric bool) string {
	first := lines[0]
	prefix := ""

	k := len(first)
	if !allowAlphaNumeric {
		k = 0
		for _, c := range first {
			if unicode.IsLetter(c) || unicode.IsNumber(c) {
				break
			}
			k++
		}
	}

	for i := 1; i < k; i++ {
		s := first[:i]
		for _, l := range lines[1:] {
			if !strings.HasPrefix(l, s) {
				return strings.TrimSpace(prefix)
			}
		}
		prefix = s
	}

	return strings.TrimSpace(prefix)
}

func linesToWords(lines []string, prefix string) chan string {
	ch := make(chan string)

	go func() {
		for _, l := range lines {
			s := strings.Trim(strings.TrimPrefix(l, prefix), " \t")
			for _, w := range strings.Split(s, " ") {
				ch <- w
			}
		}
		close(ch)
	}()

	return ch
}

func wrapLinesWithPrefix(lines []string, prefix string, columns int) {
	var wrapped []string
	var l, m, w string

	l = prefix
	for w = range linesToWords(lines, prefix) {
		if w == "" {
			wrapped = append(wrapped, fmt.Sprintf("%s\n%s\n", l, prefix))
			l = prefix + " "
			continue
		}

		m = w
		if len(l) > 0 {
			m = fmt.Sprintf("%s %s", l, w)
		}

		if len(m) <= columns {
			l = m
		} else {
			wrapped = append(wrapped, l)
			l = fmt.Sprintf("%s %s", prefix, w)
		}
	}

	if len(l) > 0 {
		wrapped = append(wrapped, l)
	}

	fmt.Print(strings.Join(wrapped, "\n"))
}

func main() {
	var lines []string
	flag.Parse()

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		lines = append(lines, s.Text())
	}

	prefix := maximalPrefix(lines, *alphaNumeric)
	wrapLinesWithPrefix(lines, prefix, *columns)
}
