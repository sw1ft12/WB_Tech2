package main

import (
	"errors"
	"fmt"
)

func isDigit(char rune) bool {
	if char >= '0' && char <= '9' {
		return true
	}
	return false
}

func isSlash(char rune) bool {
	return char == '\\'
}

func checkString(s string) bool {
	if isDigit(rune(s[0])) {
		return false
	}
	runes := []rune(s)
	slashFlag, digitFlag := false, false
	for i := range runes {
		if isDigit(runes[i]) {
			if digitFlag && !slashFlag {
				return false
			}
			digitFlag = true
		} else if isSlash(runes[i]) {
			if slashFlag && !digitFlag {
				slashFlag = false
				continue
			}
			slashFlag, digitFlag = true, false
		} else {
			slashFlag, digitFlag = false, false
		}
	}
	return true
}

func FormatString(s string) (string, error) {

	if len(s) == 0 {
		return s, nil
	}

	if !checkString(s) {
		return "", errors.New("invalid string")
	}

	runes := []rune(s)
	res := make([]rune, 0)
	slashFlag := false
	for i, val := range runes {
		if !isSlash(val) && !isDigit(val) {
			res = append(res, val)
		} else if slashFlag {
			res = append(res, val)
			slashFlag = false
		} else if isSlash(val) {
			slashFlag = true
		} else if isDigit(val) {
			for j := 1; j < int(val-'0'); j++ {
				res = append(res, runes[i-1])
			}
		}
	}
	return string(res), nil
}

func main() {
	{
		str, _ := FormatString(`abcd`)
		fmt.Print(str)
	}
}
