```
$ go install github.com/avamsi/ash@latest
```

```
$ ash --help

ash is a simple command runner (think sh -c) that supports defining flags
dynamically and executing (Go) templated commands (referring to said flags).

Usage:
  ash (--name [value [usage]])... <command ...{{.name}}> [options]
```
