package main

import "log"

func onUndefinedVariable(k string) {
	if flags.ErrorOnUndefined {
		log.Fatalf("undefined variable: %q", k)
	}
}

func onTokenizerError(err error) {
	if flags.ExitOnError {
		log.Fatalf("tokenizer: %v", err)
	}
}

func onExecuteError(command []string, err error) {
	nonzeroExit = true
	log.Printf("execute %v: %v", command, err)
	if flags.ExitOnError {
		log.Fatalf("exiting on error")
	}
}

func onExecutePre(command []string) {
	if flags.Verbose {
		log.Printf("+ %q", command)
	}
}

func onExecutePost(command []string) {
	if flags.Verbose {
		log.Printf("- %q", command)
	}
}
