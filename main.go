package main

import (
	"fmt"
	"strings"

	"github.com/c-bata/go-prompt"
)

var Version, GitCommit, BuildDate string
var AppName = "HaxShell"

func main() {
	fmt.Printf("%s v%s-%s (built %s)\n", AppName, Version, GitCommit, BuildDate)
	fmt.Printf("Setting up...\n")
	InitCommands()
	p := prompt.New(executor, completer,
		prompt.OptionPrefix("hax> "),
		prompt.OptionTitle(fmt.Sprintf("%s v%s-%s (built %s)", AppName, Version, GitCommit, BuildDate)))
	p.Run()
}

func executor(in string) {
	arr := strings.Split(in, " ")
	cmd := arr[0]
	args := arr[1:]
	CallCommand(cmd, args)
}

func completer(in prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{}
}
