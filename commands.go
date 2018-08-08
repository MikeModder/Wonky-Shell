package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Command command
type Command struct {
	Help     string
	Function func([]string)
	Args     []Arg
	ReqArgs  []string
}

// Arg command arg
type Arg struct {
	Name string
	Type string
	Desc string
}

// Commands map of commands
var Commands = make(map[string]*Command)

// InitCommands init all commands
func InitCommands() {
	// Setup all the commands here
	// No arguments commands
	Commands["test"] = &Command{Help: "test", Function: testCmd}
	Commands["version"] = &Command{Help: "Display program versopm", Function: versionCmd}

	// 1+ argument commands
	Commands["help"] = &Command{
		Help:     "Help!",
		Function: helpCmd,
		Args: []Arg{
			{Name: "command", Type: "string", Desc: "option command to show help for"},
		},
	}
	Commands["echo"] = &Command{
		Help:     "Repeat after me!",
		Function: echoCmd,
		ReqArgs:  []string{"text"},
		Args: []Arg{
			{Name: "text", Type: "string", Desc: "What do I repeat"},
		},
	}
	Commands["exit"] = &Command{
		Help:     "Exit the shell",
		Function: exitCmd,
		Args: []Arg{
			{Name: "code", Type: "int", Desc: "optional code to exit with"},
		},
	}
}

// CallCommand call a command
func CallCommand(name string, args []string) {
	if c, exist := Commands[name]; exist {
		if len(args) >= len(c.ReqArgs) {
			c.Function(args)
			return
		}
		fmt.Printf("Not enough parameters for %s: %d required, %d passed\n", name, len(c.ReqArgs), len(args))
		printCommandUsage(name)
		return
	}
	fmt.Printf("%s is not a recognized command!\n", name)
	return
}

func printCommandUsage(name string) {
	c := Commands[name]
	fmt.Printf("Usage for %s:\n %s %s\n  Parameters:\n", name, name, strings.Join(c.ReqArgs, " "))
	for i := 0; i < len(c.Args); i++ {
		fmt.Printf("  %s (%s) - %s\n", c.Args[i].Name, c.Args[i].Type, c.Args[i].Desc)
	}
	return
}

func testCmd(_ []string) {
	fmt.Println("Test command ok!")
}

func exitCmd(args []string) {
	var e error
	code := 0
	if len(args) > 0 {
		code, e = strconv.Atoi(args[0])
		if e != nil {
			code = 0
			fmt.Printf("%s is not a valid int, exiting with %d\n", args[0], code)
		}
	}
	fmt.Printf("Exiting with code %d, bye!\n", code)
	os.Exit(code)
}

func echoCmd(args []string) {
	fmt.Println(strings.Join(args, " "))
}

func helpCmd(args []string) {
	if len(args) >= 1 {
		cmd := args[0]
		if _, exist := Commands[cmd]; exist {
			printCommandUsage(cmd)
			return
		}
		fmt.Printf("Connot get help for non existant command (%s)\n", cmd)
		return
	}
	var commandNames []string
	for cKey := range Commands {
		commandNames = append(commandNames, cKey)
	}
	sort.Strings(commandNames)
	for _, commandName := range commandNames {
		c := Commands[commandName]
		fmt.Printf("%s - %s\n", commandName, c.Help)
	}
	return
}

func versionCmd(_ []string) {
	fmt.Printf("%s v%s (git %s, built %s)\n", AppName, Version, GitCommit, BuildDate)
	return
}
