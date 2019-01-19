# 0sh

Very minimal (sub-POSIX) shell for scripting.

## Contents

- [Contents](#contents)
- [Features](#features)
- [Get it](#get-it)
  - [Using `go get`](#using-go-get)
  - [Pre-built binary](#pre-built-binary)
- [Use it](#use-it)

## Features

- `$ENV` variables (read usage only, no assignment)
- double-quoted strings (expansion, escapes) and single-quoted strings (neither)
- built-ins: `cd` and `exit`
- no branching of any kind (conditionals or loops)
- no pipes (_NOTE_: this might change)

## Get it

### Using `go get`

```sh
go get -u github.com/sgreben/0sh
```

### Pre-built binary

Or [download a binary](https://github.com/sgreben/0sh/releases/latest) from the releases page, or from the shell:

```sh
# Linux
curl -L https://github.com/sgreben/0sh/releases/download/0.0.1/0sh_0.0.1_linux_x86_64.tar.gz | tar xz

# OS X
curl -L https://github.com/sgreben/0sh/releases/download/0.0.1/0sh_0.0.1_osx_x86_64.tar.gz | tar xz

# Windows
curl -LO https://github.com/sgreben/0sh/releases/download/0.0.1/0sh_0.0.1_windows_x86_64.zip
unzip 0sh_0.0.1_windows_x86_64.zip
```

## Use it

```text
0sh [OPTIONS]
```

```text
Usage of 0sh:
  -c string
    	run only the given COMMAND
  -e	exit on error
  -n	dry-run
  -u	error on undefined
  -v	verbose
  -version
    	print version and exit
```
