package main

import (
	"fmt"
	"os"
	"strconv"
)

var builtins = map[string]func([]string) error{
	"cd": func(args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("exactly one argument required, got %d: %q", len(args), args)
		}
		if err := os.Chdir(args[0]); err != nil {
			return err
		}
		wd, _ := os.Getwd()
		os.Setenv("PWD", wd)
		return nil
	},
	"exit": func(args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("exactly one argument required, got %d: %q", len(args), args)
		}
		code, err := strconv.ParseInt(args[0], 32, 10)
		if err != nil {
			return err
		}
		os.Exit(int(code))
		return nil
	},
}
