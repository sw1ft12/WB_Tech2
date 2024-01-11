package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type Flags struct {
	A int
	B int
	C int
	c bool
	i bool
	v bool
	F bool
	n bool
}

func getFlags() *Flags {
	A := flag.Int("f", 0, "Print NUM lines of trailing context after matching lines")
	B := flag.Int("d", 0, "Print NUM lines of leading context before matching lines")
	C := flag.Int("f", 0, "Print NUM lines of output context")
	c := flag.Bool("c", false, "Print a count of matching lines")
	i := flag.Bool("i", false, "Ignore case distinctions in patterns and input data")
	v := flag.Bool("v", false, "Invert the sense of matching")
	F := flag.Bool("F", false, "Interpret PATTERNS as fixed strings")
	n := flag.Bool("n", false, "Print line number")
	flag.Parse()
	return &Flags{
		A: *A,
		B: *B,
		C: *C,
		c: *c,
		i: *i,
		v: *v,
		F: *F,
		n: *n,
	}
}

func readFiles() []string {
	var lines []string
	files := flag.Args()[1:]
	for _, fileName := range files {
		file, err := os.Open(fileName)
		if err != nil {
			log.Println(err)
			continue
		}
		reader := bufio.NewReader(file)
		for {
			line, err := reader.ReadString('\n')
			if err == nil || len(line) > 0 {
				if err != nil {
					line += string('\n')
				}
				lines = append(lines, line)
			}
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Fprint(os.Stderr, err)
			}
		}
	}
	return lines
}

func getIndexes(lines []string, patterns []string, flags *Flags) []int {
	var indexes []int
	for i, v := range lines {
		if flags.i {
			v = strings.ToLower(v)
		}
		for _, pattern := range patterns {
			if flags.F {
				if strings.Contains(v, pattern) {
					if flags.n {
						lines[i] = strconv.Itoa(i+1) + ":" + lines[i]
					}
					indexes = append(indexes, i)
					break
				}
			} else {
				matched, err := regexp.MatchString(pattern, v)
				if err != nil {
					continue
				}
				if matched {
					if flags.n {
						lines[i] = strconv.Itoa(i+1) + ":" + lines[i]
					}
					indexes = append(indexes, i)
					break
				}
			}
		}
	}
	return indexes
}

func getInvertedIndexes(lines []string, patterns []string, flags *Flags) []int {
	var indexes []int
	for i, v := range lines {
		if flags.i {
			v = strings.ToLower(v)
		}
		var f bool
		for _, pattern := range patterns {
			if flags.F {
				if strings.Contains(v, pattern) {
					f = true
					break
				}
			} else {
				matched, err := regexp.MatchString(pattern, v)
				if err != nil {
					continue
				}
				if matched {
					f = true
					break
				}
			}
		}
		if !f {
			indexes = append(indexes, i)
			if flags.n {
				lines[i] = strconv.Itoa(i+1) + ":" + lines[i]
			}
		}
	}
	return indexes
}

func addBefore(indexes []int, lines []string, flags *Flags) []int {
	var result []int
	for i := 0; i < len(indexes)-1; i++ {
		for j := indexes[i] + 1; j < min(indexes[i]+flags.A+1, indexes[i+1]); j++ {
			result = append(result, j)
			if flags.n {
				lines[j] = strconv.Itoa(j+1) + "-" + lines[j]
			}
		}
	}
	last := indexes[len(indexes)-1] + 1
	for j := last; j < min(last+flags.A+1, len(lines)); j++ {
		result = append(result, j)
		if flags.n {
			lines[j] = strconv.Itoa(j+1) + "-" + lines[j]
		}
	}
	return result
}

func addAfter(indexes []int, lines []string, flags *Flags) []int {
	var result []int
	first := indexes[0] - 1
	for j := first; j > max(first-flags.B-1, -1); j-- {
		result = append(result, j)
		if flags.n {
			lines[j] = strconv.Itoa(j+1) + "-" + lines[j]
		}
	}
	for i := 1; i < len(indexes); i++ {
		for j := indexes[i] - 1; j > max(indexes[i]-flags.B-1, indexes[i-1]+flags.A); j-- {
			result = append(result, j)
			if flags.n {
				lines[j] = strconv.Itoa(j+1) + "-" + lines[j]
			}
		}
	}
	return result
}

func Grep() []string {
	flags := getFlags()
	lines := readFiles()
	patterns := strings.Split(flag.Arg(0), "|")

	if flags.i {
		for i := range patterns {
			patterns[i] = strings.ToLower(patterns[i])
		}
	}

	var indexes []int
	if flags.v {
		indexes = getInvertedIndexes(lines, patterns, flags)
	} else {
		indexes = getIndexes(lines, patterns, flags)
	}

	if flags.c {
		return []string{strconv.Itoa(len(indexes))}
	}

	if flags.C > 0 {
		flags.A = flags.C
		flags.B = flags.C
	}

	var resultInd []int
	if flags.A != 0 {
		resultInd = append(resultInd, addBefore(indexes, lines, flags)...)
	}
	if flags.B != 0 {
		resultInd = append(resultInd, addAfter(indexes, lines, flags)...)
	}

	resultInd = append(resultInd, indexes...)
	slices.Sort(resultInd)

	var result []string
	for _, s := range resultInd {
		result = append(result, lines[s])
	}
	return result
}

func main() {
	result := Grep()
	for i := range result {
		fmt.Print(result[i])
	}
}
