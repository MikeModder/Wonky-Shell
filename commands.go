package main

/* I must give credit where credit is due:
 * a big thanks goes out to JoshuaDoes
 * (https://github.com/JoshuaDoes) and
 * his Discord bot Clinet for providing
 * me with the inspiration for the command
 * system being used here.
 *
 * Files used as refrence:
 * https://github.com/JoshuaDoes/clinet-discord/blob/master/commands.go
 * https://github.com/JoshuaDoes/clinet-discord/blob/master/cmd-moderation.go
 */

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/logrusorgru/aurora"
)

var (
	history []RanCmd
)

// Command command
type Command struct {
	Help     string
	Function func([]string) int
	Args     []Arg
	ReqArgs  []string
}

// Arg command arg
type Arg struct {
	Name string
	Type string
	Desc string
}

// RanCmd Ran command
type RanCmd struct {
	Name string
	Args []string
	Code int
}

// Commands map of commands
var Commands = make(map[string]*Command)

// InitCommands init all commands
func InitCommands() {
	// Setup all the commands here
	// No arguments commands
	Commands["about"] = &Command{Help: "Display program information", Function: versionCmd}
	Commands["os"] = &Command{Help: "Display host OS and arch", Function: osCmd}
	Commands["pwd"] = &Command{Help: "Print current directory", Function: pwdCmd}
	//Commands["update"] = &Command{Help: "Check for updates", Function: updateCmd}
	Commands["last"] = &Command{Help: "Check return code of last run command", Function: lastRanCmd}
	Commands["history"] = &Command{Help: "Get command history for this session", Function: historyCmd}

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
	Commands["download"] = &Command{
		Help:     "Download a file",
		Function: downloadCmd,
		ReqArgs:  []string{"url", "filename"},
		Args: []Arg{
			{Name: "url", Type: "url", Desc: "URL of file to download"},
			{Name: "filename", Type: "string", Desc: "Optional name of saved file"},
		},
	}
}

// CallCommand call a command
func CallCommand(name string, args []string) {
	if c, exist := Commands[name]; exist {
		if len(args) >= len(c.ReqArgs) {
			code := c.Function(args)
			if len(args) == 0 {
				args = []string{"N/A"}
			}
			history = append(history, RanCmd{Name: name, Code: code, Args: args})
			return
		}
		fmt.Printf("Not enough parameters for %s: %d required, %d passed\n", aurora.Bold(name), aurora.Bold(len(c.ReqArgs)), aurora.Bold(len(args)))
		printCommandUsage(name)
		return
	}
	fmt.Printf("%s is not a recognized command!\n", aurora.Bold(name))
	return
}

func printCommandUsage(name string) {
	c := Commands[name]
	fmt.Printf("Usage for %s:\n %s %s\n  Parameters:\n", aurora.Bold(name), name, strings.Join(c.ReqArgs, " "))
	for i := 0; i < len(c.Args); i++ {
		fmt.Printf("  %s (%s) - %s\n", c.Args[i].Name, c.Args[i].Type, c.Args[i].Desc)
	}
	return
}

func exitCmd(args []string) int {
	var e error
	code := 0
	if len(args) > 0 {
		code, e = strconv.Atoi(args[0])
		if e != nil {
			code = 0
			fmt.Printf("%s is not a valid exit code, exiting with %d\n", args[0], aurora.Cyan(code))
		}
	}
	fmt.Printf("Exiting with code %d, bye!\n", aurora.Cyan(code))
	os.Exit(code)
	return 2
}

func echoCmd(args []string) int {
	fmt.Println(strings.Join(args, " "))
	return 0
}

func helpCmd(args []string) int {
	if len(args) >= 1 {
		cmd := args[0]
		if _, exist := Commands[cmd]; exist {
			printCommandUsage(cmd)
			return 0
		}
		fmt.Printf("Connot get help for non existant command (%s)\n", aurora.Gray(cmd))
		return 1
	}
	var commandNames []string
	for cKey := range Commands {
		commandNames = append(commandNames, cKey)
	}
	sort.Strings(commandNames)
	for _, commandName := range commandNames {
		c := Commands[commandName]
		fmt.Printf("%s\t - %s\n", commandName, c.Help)
	}
	return 0
}

func versionCmd(_ []string) int {
	fmt.Printf("About %s:\n", aurora.Green(AppName))
	fmt.Printf(" Version %s\n", Version)
	fmt.Printf(" Branch %s, Commit %s\n", GitBranch, GitCommit)
	fmt.Printf(" Built %s\n", BuildDate)
	fmt.Printf(" Go runtime: %s\n", runtime.Version())
	fmt.Printf("\n(c) MikeModder 2018-present\nSpecial thanks to %s\n", aurora.Cyan("JoshuaDoes"))
	return 0
}

func osCmd(_ []string) int {
	fmt.Printf("You are using %s/%s\n", runtime.GOOS, runtime.GOARCH)
	return 0
}

func pwdCmd(_ []string) int {
	dir, e := os.Getwd()
	if e != nil {
		fmt.Printf("Error getting working directory: %s\n", e.Error())
		return 1
	}
	fmt.Printf("Current directory: %s\n", dir)
	return 0
}

func cdCmd(args []string) int {
	if len(args) >= 1 {
		dir := args[0]
		e := os.Chdir(dir)
		if e != nil {
			fmt.Printf("Error changing dir to %s, error: %s\n", dir, e.Error())
			return 1
		}
		fmt.Printf("Changed directory to %s\n", dir)
		return 0
	}
	return 0
}

func lsCmd(args []string) int {
	dir := "./"
	if len(args) >= 1 {
		dir = args[0]
	}
	files, e := ioutil.ReadDir(dir)
	if e != nil {
		fmt.Printf("Error getting directory list for %s, error: %s\n", dir, e.Error())
		return 1
	}
	fmt.Printf("Contents of %s:\n", dir)
	for _, f := range files {
		fmt.Printf(" %s - %d bytes - dir: %t\n", f.Name(), f.Size(), f.IsDir())
	}
	return 0
}

func execCmd(args []string) int {
	cmd := exec.Command(args[0], args[1:len(args)]...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Env = os.Environ()
	cmd.Start()
	cmd.Wait()
	return 0
}

func downloadCmd(args []string) int {
	url := args[0]
	filename := args[1]

	file, e := os.Create(filename)
	if e != nil {
		fmt.Printf("Error creating %s for download! Error: %s\n", filename, e.Error())
		return 1
	}

	resp, e := http.Get(url)
	if e != nil {
		fmt.Printf("Error downloading %s! Error: %s\n", url, e.Error())
		return 1
	}
	defer resp.Body.Close()

	_, e = io.Copy(file, resp.Body)
	if e != nil {
		fmt.Printf("Error writing file! Error: %s\n", e.Error())
		return 1
	}
	e = file.Close()
	if e != nil {
		fmt.Printf("Error closing file! Error: %s\n", e.Error())
		return 1
	}
	fmt.Printf("Downloaded %s successfully!\n", filename)
	return 0
}

func lastRanCmd(_ []string) int {
	last := history[len(history)-1]
	fmt.Printf("Last ran command (%s) returned with code %d\n", last.Name, last.Code)
	return 0
}

func historyCmd(args []string) int {
	if len(args) >= 1 {
		if args[0] == "clear" {
			history = make([]RanCmd, 0)
			fmt.Println("Cleared command history!")
			return 0
		}
	}
	fmt.Println("Command history since begining of session:\nTIP: Use `history clear` to clear!\n| name | args | returned |")
	if len(history) == 0 {
		fmt.Println("No commands ran yet!")
		return 0
	}
	for _, item := range history {
		fmt.Printf("| %s | %s | %d |\n", item.Name, strings.Join(item.Args, " "), item.Code)
	}
	return 0
}
