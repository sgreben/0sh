# 0sh

Very simple (sub-POSIX) shell for scripting.

## Contents

- [Contents](#contents)
- [Features](#features)
  - [Example](#example)
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
- `-e` (errexit) and `-u` (noundef) enabled by default

### Example

```sh
#!/usr/bin/env 0sh
echo "$ABC" '$ABC' # a comment
exit 1
```

## Get it

### Using `go get`

```sh
go get -u github.com/sgreben/0sh
```

### Pre-built binary

Or [download a binary](https://github.com/sgreben/0sh/releases/latest) from the releases page, or from the shell:

```sh
# Linux
curl -L https://github.com/sgreben/0sh/releases/download/0.0.2/0sh_0.0.2_linux_x86_64.tar.gz | tar xz

# OS X
curl -L https://github.com/sgreben/0sh/releases/download/0.0.2/0sh_0.0.2_osx_x86_64.tar.gz | tar xz

# Windows
curl -LO https://github.com/sgreben/0sh/releases/download/0.0.2/0sh_0.0.2_windows_x86_64.zip
unzip 0sh_0.0.2_windows_x86_64.zip
```

## Use it

```text
0sh [OPTIONS]
```

```text
Usage of 0sh:
  -c string
    	run only the given COMMAND
  -e	exit on error (default true)
  -n	dry-run
  -u	error on undefined (default true)
  -v	verbose
  -version
    	print version and exit
```
