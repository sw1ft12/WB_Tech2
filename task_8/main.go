package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/shirou/gopsutil/v3/process"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func readLine() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.TrimSuffix(input, "\n")
	return input, nil
}

func ChangeDirectory(path []string) error {
	if len(path) > 1 {
		return errors.New("cd: too many arguments")
	}
	if len(path) == 0 {
		homeDir := os.Getenv("HOME")
		path = append(path, homeDir)
	}
	err := os.Chdir(path[0])

	if err != nil {
		err = errors.New(strings.Replace(err.Error(), "chdir", "cd", 1))
		return err
	}
	return nil
}

func Echo(args []string) string {
	return strings.Join(args, " ") + "\n"
}

func PrintWorkingDirectory() (string, error) {
	return os.Getwd()
}

func Kill(args []string) error {
	if len(args) < 1 {
		return errors.New("not enough arguments")
	}
	processes, err := process.Processes()

	if err != nil {
		return err
	}

	var errs []error
	for _, pid := range args {
		pi, err := strconv.Atoi(pid)
		if err != nil {
			errs = append(errs, fmt.Errorf("invalid pid: %d", pi))
		}
		var killed bool
		for _, p := range processes {
			if p.Pid == int32(pi) {
				err = p.Kill()
				if err != nil {
					errs = append(errs, err)
				}
				killed = true
			}
		}
		if !killed {
			errs = append(errs, fmt.Errorf("%d: no such process", pi))
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func PrintProcesses() (string, error) {
	processes, err := process.Processes()
	if err != nil {
		return "", err
	}
	result := strings.Builder{}
	result.WriteString("USER\tPID\tCOMMAND\n")

	for _, p := range processes {
		name, _ := p.Name()
		user, _ := p.Username()
		result.WriteString(fmt.Sprintf("%s\t%d\t%s\n", user, p.Pid, name))
	}
	return result.String(), nil
}

var (
	pipe     bytes.Buffer
	pipeFlag bool
)

func handleCmd(cmd string) (string, error) {
	var err error
	tokens := strings.Fields(cmd)
	utility := tokens[0]
	var args []string
	if len(tokens) > 1 {
		args = tokens[1:]
	}
	var result string
	switch utility {
	case "cd":
		err = ChangeDirectory(args)
	case "pwd":
		result, err = PrintWorkingDirectory()
	case "echo":
		result = Echo(args)
	case "kill":
		err = Kill(args)
	case "ps":
		result, err = PrintProcesses()
	default:
		cmd := exec.Command(utility, args...)
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		if pipeFlag {
			cmd.Stdout = &pipe
		}

		if pipe.Len() > 0 {
			cmd.Stdin = &pipe
		}

		err := cmd.Run()
		if err != nil {
			break
		}

		result = pipe.String()
	}

	if err != nil {
		return "", err
	}

	if pipeFlag {
		pipe.Reset()
		pipe.WriteString(result)
	}

	return result, nil
}

func Run() error {
	for {
		input, err := readLine()
		if err != nil {
			return err
		}
		if len(input) < 1 {
			continue
		}
		if input == "quit" {
			return nil
		}
		commands := strings.Split(input, "|")
		pipeFlag = true
		var result string
		for i, cmd := range commands {
			if i == len(commands) {
				pipeFlag = false
			}
			result, err = handleCmd(cmd)
			if err != nil {
				return err
			}
		}
		fmt.Print(result)
	}
}

func main() {
	err := Run()
	if err != nil {
		fmt.Println(err)
	}
}
