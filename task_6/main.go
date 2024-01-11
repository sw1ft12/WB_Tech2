package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Flags struct {
	f []int
	d string
	s bool
}

func getFlags() (*Flags, error) {
	f := flag.String("f", "", "Select only these fields")
	d := flag.String("d", `\t`, "Use DELIM instead of TAB for field delimiter")
	s := flag.Bool("s", false, "Do not print lines not containing delimiters")
	flag.Parse()

	fragments := strings.Split(*f, " ")
	fields := make([]int, len(fragments))

	for i := range fragments {
		num, err := strconv.Atoi(fragments[i])
		if err != nil {
			return nil, err
		}
		fields[i] = num
	}

	return &Flags{
		f: fields,
		d: *d,
		s: *s,
	}, nil
}

func readFile(filename string) ([]string, error) {
	var lines []string

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func Cut() ([]string, error) {

	flags, err := getFlags()
	if err != nil {
		return nil, err
	}

	var result []string

	files := flag.Args()
	if len(files) == 0 {
		files = append(files, os.Stdin.Name())
	}

	for _, file := range files {
		lines, err := readFile(file)
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, line := range lines {
			if !flags.s {
				result = append(result, line)
			} else if strings.Contains(line, flags.d) {
				fragments := strings.Split(line, flags.d)
				builder := strings.Builder{}
				for _, val := range flags.f {
					if len(fragments) >= val {
						builder.WriteString(fragments[val-1])
						builder.WriteString(flags.d)
					}
				}
				result = append(result, strings.TrimSuffix(builder.String(), flags.d))
			}
		}
	}
	return result, nil
}

func main() {
	result, err := Cut()
	if err != nil {
		log.Fatal(err)
	}
	for _, line := range result {
		fmt.Println(line)
	}
}
