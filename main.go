package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/c-bata/go-prompt"
)

var Version, GitCommit, BuildDate string
var AppName = "HaxShell"

/*var commands []prompt.Suggest = []prompt.Suggest{
{Text: "exit", Description: "exit the program"},
{Text: "echo", Description: "repeat text given"},
{Text: "whoami", Description: "what user am I?"},
{Text: "run", Description: "directly run a file"}}*/

func main() {
	fmt.Printf("%s v%s-%s (built %s)\n", AppName, Version, GitCommit, BuildDate)
	p := prompt.New(executor, completer,
		prompt.OptionPrefix("hax> "),
		prompt.OptionTitle(fmt.Sprintf("%s v%s-%s (built %s)", AppName, Version, GitCommit, BuildDate)))
	p.Run()
}

func executor(in string) {
	arr := strings.Split(in, " ")
	cmd := strings.ToLower(arr[0])
	args := arr[1:]
	handleCommand(cmd, args)
}

func completer(in prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{}
}

func handleCommand(c string, args []string) {
	switch c {
	case "exit":
		var e error
		code := 0
		if len(args) > 0 {
			code, e = strconv.Atoi(args[0])
			if e != nil {
				code = 0
			}
		}
		fmt.Printf("Exiting with code %d, bye!\n", code)
		os.Exit(code)
	case "echo":
		fmt.Println(strings.Join(args, " "))
	case "whoami":
		fmt.Printf("You are %s running HaxShell\n", os.Getenv("USER"))
	case "exist":
		f, e := os.Stat(args[0])
		if e != nil {
			fmt.Printf("%s does not exist or you don't have permission to see it,\nstat error: %s\n", args[0], e.Error())
			return
		}
		fmt.Printf("%s exists, directory: %t\nsize: %d\n", f.Name(), f.IsDir(), f.Size())
	case "os":
		fmt.Printf("You are running %s (arch: %s, %s)\n", runtime.GOOS, runtime.GOARCH, runtime.Version())
	case "version":
		fmt.Printf("%s v%s (git: %s, built %s)\n", AppName, Version, GitCommit, BuildDate)
	case "help":
		fmt.Printf("exit <code>\necho <text>\nwhomai\nexist <path>\nos\nversion\nhelp\n")
	default:
		fmt.Printf("%s is not a valid command!\n", c)
	}
}
