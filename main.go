package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/sgreben/0sh/pkg/shenv"
	"github.com/sgreben/0sh/pkg/shlex"
)

var flags struct {
	Command             string
	DryRun              bool
	ExitOnError         bool
	ErrorOnUndefined    bool
	Verbose             bool
	PrintVersionAndExit bool
}

const appName = "0sh"

var (
	version        = "SNAPSHOT"
	nonzeroExit    bool
	commandReader  io.Reader
	postionalArg0  string
	positionalArgs []string
)

func init() {
	log.SetOutput(os.Stderr)
	log.SetPrefix(fmt.Sprintf("[%s] ", appName))
	log.SetFlags(0)
	flag.BoolVar(&flags.PrintVersionAndExit, "version", false, "print version and exit")
	flag.BoolVar(&flags.Verbose, "v", false, "verbose")
	flag.BoolVar(&flags.ExitOnError, "e", false, "exit on error")
	flag.BoolVar(&flags.ErrorOnUndefined, "u", false, "error on undefined")
	flag.StringVar(&flags.Command, "c", "", "run only the given COMMAND")
	flag.BoolVar(&flags.DryRun, "n", false, "dry-run")
	flag.Parse()
}

func main() {
	if flags.PrintVersionAndExit {
		fmt.Println(version)
		os.Exit(0)
	}
	commandReader = os.Stdin
	postionalArg0, _ = os.Executable()
	switch {
	case flags.Command != "":
		commandReader = strings.NewReader(flags.Command)
		positionalArgs = flag.Args()
	case flag.NArg() > 0:
		path := flag.Arg(0)
		f, err := os.Open(path)
		if err != nil {
			log.Fatalf("open %q: %v", path, err)
		}
		if flag.NArg() > 1 {
			positionalArgs = flag.Args()[1:]
		}
		commandReader = f
	}
	executeFromReader(commandReader)
	if nonzeroExit {
		os.Exit(1)
	}
}

func substEnvVar(varName string) (out string) {
	if v, err := strconv.ParseUint(varName, 10, 32); err == nil {
		if v == 0 {
			return postionalArg0
		}
		if i := int(v); i <= len(positionalArgs) {
			return positionalArgs[i-1]
		}
	}
	if v, ok := os.LookupEnv(varName); ok {
		return v
	}
	onUndefinedVariable(varName)
	return
}

func substEnv(s string) string {
	return shenv.Expand(s, substEnvVar)
}

func execute(command []string) {
	defer onExecutePost(command)
	onExecutePre(command)
	if len(command) == 0 {
		return
	}
	name := command[0]
	var args []string
	if len(command) > 1 {
		args = command[1:]
	}
	if flags.DryRun {
		return
	}
	if builtin, ok := builtins[name]; ok {
		if err := builtin(args); err != nil {
			onExecuteError(command, fmt.Errorf("builtin %q: %v", name, err))
		}
		return
	}
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		onExecuteError(command, err)
	}
}

func executeFromReader(r io.Reader) {
	tokens := shlex.NewTokenizer(r)
	var command []string
	run := true
	for run {
		token, err := tokens.Next(substEnv)
		switch {
		case err == io.EOF:
			run = false
			token = &shlex.Token{Type: shlex.TokenTypeNewline}
		case err != nil:
			onTokenizerError(err)
		}
		switch token.Type {
		case shlex.TokenTypeWordThenNewline,
			shlex.TokenTypeWordThenSemicolon:
			command = append(command, token.Value)
			fallthrough
		case shlex.TokenTypeNewline,
			shlex.TokenTypeSemicolon:
			execute(command)
			command = nil
		case shlex.TokenTypeWord:
			command = append(command, token.Value)
		}
	}

}
