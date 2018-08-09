package main

import (
	"fmt"
	"strings"

	"github.com/c-bata/go-prompt"
)

var Version, GitCommit, GitBranch, BuildDate string

const (
	// AppName app's name
	AppName = "HaxShell"
	// PromptText text to use for promt
	PromptText = "hax> "
)

func main() {
	fmt.Printf("%s v%s-%s (built %s)\n", AppName, Version, GitCommit, BuildDate)
	fmt.Printf("%s - (c) MikeModder\n", AppName)
	InitCommands()
	p := prompt.New(executor, completer,
		prompt.OptionPrefix(PromptText),
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
