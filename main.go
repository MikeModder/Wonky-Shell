package main

import (
	"fmt"
	"strings"

	"github.com/c-bata/go-prompt"
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
	AppName = "HaxShell"
	// PromptText text to use for promt
	PromptText = "hax> "
)

func main() {
	fmt.Printf("%s v%s-%s (built %s)\n", aurora.Green(AppName), Version, GitCommit, BuildDate)
	fmt.Printf("Type %s for a list of commands!\n", aurora.Bold("help"))
	InitCommands()
	p := prompt.New(executor, completer,
		prompt.OptionPrefix(PromptText),
		prompt.OptionTitle(fmt.Sprintf("%s v%s-%s", AppName, Version, GitCommit)))
	p.Run()
}

func executor(in string) {
	arr := strings.Split(in, " ")
	cmd := arr[0]
	args := arr[1:]
	CallCommand(cmd, args)
}

func completer(_ prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{}
}
