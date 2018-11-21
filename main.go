package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/logrusorgru/aurora"
)

var (
	//Version Contents of ./VERSION
	Version string
	//GitBranch Current Git branch
	GitBranch string
	//GitCommit Current Git commit
	GitCommit string
	//BuildDate Build date
	BuildDate string
	//GitState State of the local repo
	GitState string
)

const (
	// AppName app's name
	AppName = "Wonky-Shell"
	// PromptText text to use for promt
	PromptText = "hax> "
)

func main() {
	fmt.Printf("%s v%s-%s (built %s)\n", aurora.Green(AppName), Version, GitCommit, BuildDate)
	fmt.Printf("Type %s for a list of commands!\n", aurora.Bold("help"))
	InitCommands()

	scanner := bufio.NewScanner(os.Stdin)
	loop := true

	for loop {
		fmt.Print(PromptText)
		if !scanner.Scan() {
			fmt.Printf("[err] failed to scan: %v\n", scanner.Err())
			os.Exit(1)
		}

		input := scanner.Text()
		executor(input)
	}
}

func executor(in string) {
	arr := strings.Split(in, " ")
	cmd := arr[0]
	args := arr[1:]
	CallCommand(cmd, args)
}
