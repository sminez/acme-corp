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
// have. By default, we stop at the first alphanumeric character encountered
// but this can be overriden using the '-a' flag.
func maximalPrefix(lines []string, allowAlphaNumeric bool) string {
	var k, n int
	first := lines[0]
	prefix := ""

	if allowAlphaNumeric {
		k = len(first)
	} else {
		k = 0
		for _, c := range first {
			if unicode.IsLetter(c) || unicode.IsNumber(c) {
				break
			}
			k++
		}
	}

	for n = 1; n < k; n++ {
		s := first[:n]
		for _, l := range lines[1:] {
			if !strings.HasPrefix(l, s) {
				return strings.TrimSpace(prefix)
			}
		}
		prefix = s
	}

	return strings.TrimSpace(prefix)
}

// runs in a goroutine
func linesToWords(lines []string, prefix string, ch chan string) {
	for _, l := range lines {
		s := strings.Trim(strings.TrimPrefix(l, prefix), " \t")
		words := strings.Split(s, " ")
		for _, w := range words {
			ch <- w
		}
	}

	close(ch)
}

func wrapLinesWithPrefix(lines []string, prefix string, columns int) {
	var l, m string
	ch := make(chan string)
	l = prefix

	go linesToWords(lines, prefix, ch)

	for w := range ch {
		if w == "" {
			fmt.Printf("%s\n%s\n", l, prefix)
			l = prefix
			continue
		}

		m = fmt.Sprintf("%s %s", l, w)
		if len(m) <= columns {
			l = m
		} else {
			fmt.Println(l)
			l = fmt.Sprintf("%s %s", prefix, w)
		}
	}

	if len(l) > 0 {
		fmt.Println(l)
	}
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
