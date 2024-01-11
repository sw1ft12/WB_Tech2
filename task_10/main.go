package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type args struct {
	host    string
	port    string
	timeout time.Duration
}

func getArgs() (*args, error) {
	if len(os.Args) < 3 {
		return nil, errors.New("not enough arguments")
	}
	var (
		timeout time.Duration
		host    string
		port    string
	)

	if strings.Contains(os.Args[1], "--timeout=") {
		modif := os.Args[1][len(os.Args[1])-1]
		if modif != 's' {
			return nil, errors.New("you need to specify time unit: e.g.: 10s")
		}

		index := strings.Index(os.Args[1], "=")
		num, err := strconv.Atoi(os.Args[1][index+1 : len(os.Args[1])-1])
		if err != nil || num < 1 {
			return nil, err
		}

		timeout = time.Duration(num) * time.Second
		host = os.Args[2]
		port = os.Args[3]
	} else {
		host = os.Args[1]
		port = os.Args[2]
		timeout = time.Second * 10
	}

	return &args{
		host:    host,
		port:    port,
		timeout: timeout,
	}, nil
}

func readFromSocket(conn net.Conn, errChan chan error) {
	input := make([]byte, 1024)
	for {
		n, err := conn.Read(input)
		if err != nil {
			errChan <- fmt.Errorf("remoute server stopped: %v", err)
			return
		}
		fmt.Println(string(input[:n]))
	}
}

func writeToSocket(conn net.Conn, errChan chan error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadBytes('\n')
		if err != nil {
			errChan <- err
			return
		}
		text = text[:len(text)-1]

		_, err = conn.Write(text)
		if err != nil {
			errChan <- err
			return
		}
	}
}

func telnet() error {
	args, err := getArgs()
	if err != nil {
		return err
	}
	address := fmt.Sprintf("%s:%s", args.host, args.port)
	conn, err := net.DialTimeout("tcp", address, args.timeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	errChan := make(chan error)

	go readFromSocket(conn, errChan)
	go writeToSocket(conn, errChan)

	select {
	case s := <-sigs:
		fmt.Println("\nConnection stopped by signal:", s)
	case e := <-errChan:
		fmt.Println("Connection stopped by", e)
	}
	return nil

}

func main() {
	err := telnet()
	if err != nil {
		log.Fatal(err)
	}
}
