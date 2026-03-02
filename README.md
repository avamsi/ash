```
$ go install github.com/avamsi/ash@latest
```

```
$ ash --help

ash is a hybrid between getopts and 'sh -c'.

It lets you define CLI flags dynamically and run shell commands rendered from
Go templates. Flags and positional arguments can be referenced within the
command template -- flags by name, and the remaining arguments via {{.args}}.

Usage:
  ash [-n, --noop] [-q, --quiet] [-<name> [<value> [usage]]]... <template>

Example:
  alias rdiff='ash "-b main branch" "--remote origin" \
        "git fetch {{.remote}} && git diff {{.remote}}/{{.b}}"'

  $ rdiff --help

  Usage:
    -b string
          branch (default "main")
    -remote string
           (default "origin")

See https://pkg.go.dev/text/template for more information on Go templates.
```
