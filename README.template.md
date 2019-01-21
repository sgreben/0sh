# ${APP}

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
go get -u github.com/sgreben/${APP}
```

### Pre-built binary

Or [download a binary](https://github.com/sgreben/${APP}/releases/latest) from the releases page, or from the shell:

```sh
# Linux
curl -L https://github.com/sgreben/${APP}/releases/download/${VERSION}/${APP}_${VERSION}_linux_x86_64.tar.gz | tar xz

# OS X
curl -L https://github.com/sgreben/${APP}/releases/download/${VERSION}/${APP}_${VERSION}_osx_x86_64.tar.gz | tar xz

# Windows
curl -LO https://github.com/sgreben/${APP}/releases/download/${VERSION}/${APP}_${VERSION}_windows_x86_64.zip
unzip ${APP}_${VERSION}_windows_x86_64.zip
```

## Use it

```text
${APP} [OPTIONS]
```

```text
${USAGE}
```
