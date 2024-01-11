package main

import (
	"fmt"
	"sort"
	"strings"
)

const alphSz int = 33

func FindAnagrams(data []string) map[string][]string {
	m := make(map[[alphSz]int][]string)
	set := make(map[string]struct{})
	for i := range data {
		s := strings.ToLower(data[i])
		if _, ok := set[s]; ok {
			continue
		}
		var letterCnt [alphSz]int
		for _, c := range s {
			letterCnt[c-'а']++
		}
		m[letterCnt] = append(m[letterCnt], s)
		set[s] = struct{}{}
	}
	res := make(map[string][]string)
	for i := range m {
		if len(m[i]) == 1 {
			continue
		}
		sort.Strings(m[i])
		res[m[i][0]] = m[i]
	}
	return res
}

func main() {
	data := []string{"Пятак", "пятка", "пятка", "тяпка", "листок", "слиток", "столик"}
	fmt.Print(FindAnagrams(data))
}
