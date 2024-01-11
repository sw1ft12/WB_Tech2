package main

import (
	"fmt"
	"github.com/beevik/ntp"
	"os"
	"time"
)

func GetCurrentTime() (time.Time, error) {
	t, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func main() {
	t, err := GetCurrentTime()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Print(t)
}
