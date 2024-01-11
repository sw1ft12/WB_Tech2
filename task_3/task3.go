package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Flags struct {
	k int
	n bool
	r bool
	u bool
	m bool
	b bool
	c bool
	h bool
}

func getFlags() *Flags {
	k := flag.Int("k", 1, "sort by kth column")
	n := flag.Bool("n", false, "compare according to string numerical value")
	r := flag.Bool("r", false, "sort in descending order")
	u := flag.Bool("u", false, "outputs only unique strings")
	m := flag.Bool("M", false, "sort by month")
	b := flag.Bool("b", false, "ignore leading blanks")
	c := flag.Bool("c", false, "heck for sorted input")
	h := flag.Bool("h", false, "compare human readable numbers (e.g., 2K 1G)")
	flag.Parse()
	return &Flags{
		k: *k,
		n: *n,
		r: *r,
		u: *u,
		m: *m,
		b: *b,
		c: *c,
		h: *h,
	}
}

func readFiles() []string {
	var lines []string
	for _, fileName := range flag.Args() {
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			continue
		}
		r := bufio.NewReader(file)
		for {
			line, err := r.ReadString('\n')
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
		file.Close()
	}
	return lines
}

func getColumn(lines []string, col int) []string {
	words := make([]string, 0)
	for i := range lines {
		v := strings.Fields(lines[i])
		if len(v) < col {
			continue
		}
		words = append(words, v[col-1])
	}
	return words
}

func siftUp(lines []string, col int) []string {
	start := 0
	for i := range lines {
		v := strings.Fields(lines[i])
		if len(v) < col {
			lines[start], lines[i] = lines[i], lines[start]
			start++
		}
	}
	sort.Strings(lines[:start])
	return lines[start:]
}

func sortByColumn(lines []string, cmp func(i, j int) bool) {
	sort.Slice(lines, cmp)
}

func sortByNumber(lines []string, column []string) {
	sortByColumn(lines, func(i, j int) bool {
		v1, _ := strconv.Atoi(column[i])
		v2, _ := strconv.Atoi(column[j])
		if v1 < v2 {
			column[i], column[j] = column[j], column[i]
			return true
		}
		return false
	})
}

func parseTime(date string) (time.Time, error) {
	t, err := time.Parse(`Jan`, date)
	if err == nil {
		return t, nil
	}
	t, err = time.Parse(`January`, date)
	if err == nil {
		return t, nil
	}
	return time.Time{}, err
}

func sortByMonth(lines []string, column []string) {
	sortByColumn(lines, func(i, j int) bool {
		t1, _ := parseTime(column[i])
		t2, _ := parseTime(column[j])
		if t1.Before(t2) {
			column[i], column[j] = column[j], column[i]
			return true
		}
		return false
	})
}

func sortByNumberWithSuff(lines []string, column []string) {
	numericSuff := map[byte]float64{
		'K': 1024,
		'M': 1024 * 1024,
		'G': 1024 * 1024 * 1024,
		'T': 1024 * 1024 * 1024 * 1024,
	}
	for i := range column {
		var s string
		for _, c := range column[i] {
			if c == ',' {
				continue
			}
			s += string(c)
		}
		column[i] = s
	}

	sortByColumn(lines, func(i, j int) bool {
		v1, _ := strconv.ParseFloat(column[i][:len(column[i])], 64)
		v2, _ := strconv.ParseFloat(column[j][:len(column[j])], 64)
		if s, ok := numericSuff[column[i][len(column[i])-1]]; ok {
			v1, _ = strconv.ParseFloat(column[i][:len(column[i])-1], 64)
			v1 *= s
		}
		if s, ok := numericSuff[column[j][len(column[j])-1]]; ok {
			v2, _ = strconv.ParseFloat(column[j][:len(column[j])-1], 64)
			v2 *= s
		}
		if v1 < v2 {
			column[i], column[j] = column[j], column[i]
			return true
		}
		return false
	})
}

func reverse(lines []string) []string {
	sz := len(lines)
	for i := 0; i < sz/2; i++ {
		lines[i], lines[sz-i-1] = lines[sz-i-1], lines[i]
	}
	return lines
}

func makeUnique(lines []string) []string {
	m := make(map[string]struct{})
	for _, line := range lines {
		m[line] = struct{}{}
	}
	res := make([]string, 0, len(m))
	for line := range m {
		res = append(res, line)
	}
	return res
}

func deleteBlanks(lines []string) {
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i]) + "\n"
	}
}

func Sort() string {
	flags := getFlags()
	lines := readFiles()
	if flags.u {
		lines = makeUnique(lines)
	}
	if flags.b {
		deleteBlanks(lines)
	}
	lines_ := siftUp(lines, flags.k)
	column := getColumn(lines_, flags.k)
	if flags.m {
		sortByMonth(lines_, column)
	} else if flags.n {
		sortByNumber(lines_, column)
	} else if flags.h {
		sortByNumberWithSuff(lines, column)
	} else {
		sortByColumn(lines_, func(i, j int) bool {
			if column[i] < column[j] {
				column[i], column[j] = column[j], column[i]
				return true
			}
			return false
		})
	}
	if flags.r {
		reverse(lines)
	}
	if flags.c {
		c := readFiles()
		for i := range lines {
			if c[i] != lines[i] {
				return fmt.Sprintf("Disorder: %d\n", i+1)
			}
		}
		return ""
	}
	return strings.Join(lines, "")
}

func main() {
	fmt.Print(Sort())
}
