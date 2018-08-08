package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
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
	Commands["version"] = &Command{Help: "Display program versopm", Function: versionCmd}
	Commands["os"] = &Command{Help: "Display host OS and arch", Function: osCmd}
	Commands["pwd"] = &Command{Help: "Print current directory", Function: pwdCmd}
	Commands["update"] = &Command{Help: "Check for updates", Function: updateCmd}

	// 1+ argument commands
	Commands["help"] = &Command{
		Help:     "Help!",
		Function: helpCmd,
		Args: []Arg{
			{Name: "command", Type: "string", Desc: "option command to show help for"},
		},
	}
	Commands["cd"] = &Command{
		Help:     "Change directory",
		Function: cdCmd,
		ReqArgs:  []string{"path"},
		Args: []Arg{
			{Name: "path", Type: "path", Desc: "Path to enter"},
		},
	}
	Commands["ls"] = &Command{
		Help:     "Show files/folders in a dir",
		Function: lsCmd,
		Args: []Arg{
			{Name: "path", Type: "path", Desc: "Path to list contents of"},
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
	Commands["exec"] = &Command{
		Help:     "Execute a file",
		Function: execCmd,
		ReqArgs:  []string{"file"},
		Args: []Arg{
			{Name: "file", Type: "string", Desc: "Name of file/program to execute"},
			{Name: "args", Type: "[]string", Desc: "arguments to pass"},
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

func osCmd(_ []string) {
	fmt.Printf("You are using %s/%s\n", runtime.GOOS, runtime.GOARCH)
	return
}

func pwdCmd(_ []string) {
	dir, e := os.Getwd()
	if e != nil {
		fmt.Printf("Error getting working directory: %s\n", e.Error())
		return
	}
	fmt.Printf("Current directory: %s\n", dir)
	return
}

func cdCmd(args []string) {
	if len(args) >= 1 {
		dir := args[0]
		e := os.Chdir(dir)
		if e != nil {
			fmt.Printf("Error changing dir to %s, error: %s\n", dir, e.Error())
			return
		}
		fmt.Printf("Changed directory to %s\n", dir)
		return
	}
}

func lsCmd(args []string) {
	dir := "./"
	if len(args) >= 1 {
		dir = args[0]
	}
	files, e := ioutil.ReadDir(dir)
	if e != nil {
		fmt.Printf("Error getting directory list for %s, error: %s\n", dir, e.Error())
		return
	}
	fmt.Printf("Contents of %s:\n", dir)
	for _, f := range files {
		fmt.Printf(" %s - %d bytes - dir: %t\n", f.Name(), f.Size(), f.IsDir())
	}
}

func execCmd(args []string) {
	cmd := exec.Command(args[0], args[1:len(args)]...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Env = os.Environ()
	cmd.Start()
	cmd.Wait()
	return
}

func updateCmd(_ []string) {
	update, e, hash := CheckUpdate()
	if e {
		fmt.Println("Error while checking for update!")
		return
	}
	if update {
		fmt.Printf("Latest commit is %s while local is %s, there may be an update.\n", hash, GitCommit)
	}
	fmt.Println("It looks like you're all up-to-date!")
	return
}
